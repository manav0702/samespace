package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"samespace/db"
	"samespace/models"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// GetTodos retrieves all todo items, sorted by creation date
func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	m := map[string]interface{}{}

	iter := db.Session.Query("SELECT id, title, description, completed, created_at FROM todos").Iter()
	for iter.MapScan(m) {
		todo := models.Todo{
			ID:          m["id"].(gocql.UUID),
			Title:       m["title"].(string),
			Description: m["description"].(string),
			Completed:   m["completed"].(bool),
			CreatedAt:   m["created_at"].(time.Time),
		}
		todos = append(todos, todo)
		m = map[string]interface{}{}
	}

	if err := iter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.Before(todos[j].CreatedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("TODO LIST")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo creates a new todo item
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.ID = gocql.TimeUUID()
	todo.CreatedAt = time.Now()

	if err := db.Session.Query("INSERT INTO todos (id, title, description, completed, created_at) VALUES (?, ?, ?, ?, ?)",
		todo.ID, todo.Title, todo.Description, todo.Completed, todo.CreatedAt).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("NEW TASK ADDED")
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo updates an existing todo item
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Session.Query("UPDATE todos SET title = ?, description = ?, completed = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Completed, id).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("TASK UPDATED")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes an existing todo item
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Session.Query("DELETE FROM todos WHERE id = ?", id).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("TASK DELETED")
}
