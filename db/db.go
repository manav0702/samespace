package db

import (
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var Session *gocql.Session

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cluster := gocql.NewCluster(os.Getenv("SCYLLA_HOST"))
	cluster.Keyspace = os.Getenv("SCYLLA_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	log.Println("Connecting to Scylla at ", os.Getenv("SCYLLA_HOST"))
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error creating ScyllaDB session: %v", err)
	}

	log.Println("ScyllaDB session initialized")
}
