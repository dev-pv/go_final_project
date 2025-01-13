package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database and creates required tables.
func InitDB() (*sql.DB, error) {
	// Determine database file path
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Unable to open database: %v", err)
	}

	if install {
		// Create table if it doesn't exist
		tableQuery := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128),
			UNIQUE(date, title)
		);`

		if _, err := db.Exec(tableQuery); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}

		// Create index on date column
		indexQuery := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
		if _, err := db.Exec(indexQuery); err != nil {
			log.Fatalf("Failed to create index: %v", err)
		}
	}

	return db, nil
}
