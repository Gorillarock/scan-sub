package main

import (
	"context"
	"database/sql"
	"flag"
	"os"

	"github.com/censys/scan-takehome/pkg/database"
	"github.com/censys/scan-takehome/pkg/processing"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	topicId := flag.String("topic", "scan-topic", "GCP PubSub Topic ID")

	dbClient, err := newClient()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Info("Connected to database successfully")

	h := database.NewHandler(dbClient)
	defer h.Close()
	processor := processing.NewHandler(h, *projectId)

	err = processor.SetSubscription(context.Background(), topicId)
	if err != nil {
		log.Fatal("Failed to set subscription:", err)
	}

	for {
		err = processor.Receive(context.Background())
		if err != nil {
			log.Error("Error receiving messages:", err)
		}
	}
}

// NewClient opens a connection to the SQLite database.
// It reads the DB path from the DB_PATH environment variable.
func newClient() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/db/scans.db" // default path, matches docker-compose volume mount
	}
	return sql.Open("sqlite3", dbPath)
}
