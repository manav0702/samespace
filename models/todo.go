package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Todo struct {
	ID          gocql.UUID `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
}
