package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/censys/scan-takehome/pkg/scanning"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

type tester struct {
	dbClient *sql.DB
}

func main() {
	dbClient, err := newClient()
	if err != nil {
		panic("Failed to connect to database:" + err.Error())
	}
	defer dbClient.Close()

	t := &tester{dbClient: dbClient}
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		uniqueScanCount, err := t.getUniqueScansCount()
		if err != nil {
			log.Error("Failed to get unique scans count:", err)
			continue
		}
		fmt.Println("Unique scan count: ", uniqueScanCount)
		fmt.Println("Logging last 10 scans of the database:")
		t.logLast10DbLines()
	}
	fmt.Println("test_db_reader exiting")
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

// Every 10 seconds, log the last 10 lines of the database to verify inserts are working
func (t *tester) logLast10DbLines() {
	rows, err := t.dbClient.Query(`SELECT ip, port, service, last_scanned, response FROM scans ORDER BY last_scanned DESC LIMIT 10`)
	if err != nil {
		log.Error("Failed to query last 10 scans:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var scan scanning.Scan
		var dataStr string
		if err := rows.Scan(&scan.Ip, &scan.Port, &scan.Service, &scan.Timestamp, &dataStr); err != nil {
			log.Error("Failed to scan row:", err)
			return
		}
		fmt.Println("Scan:", scan.Ip, scan.Port, scan.Service, scan.Timestamp, "Data:", dataStr)
	}
}

func (t *tester) getUniqueScansCount() (uint64, error) {
	rows, err := t.dbClient.Query(`SELECT DISTINCT ip, port, service FROM scans`)
	if err != nil {
		return 0, fmt.Errorf("failed to query unique scans: %w", err)
	}
	defer rows.Close()

	var count uint64
	for rows.Next() {
		count++
	}
	return count, nil
}
