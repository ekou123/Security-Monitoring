package db

import (
    "database/sql"
    "log"
	"fmt"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB 

func InitializeDB() {
	
	var err error
	DB, err = sql.Open("sqlite3", "./security_monitoring.db")
	if err != nil {
		fmt.Println("Deez")
		log.Fatal(err)
	}

	if _, err = DB.Exec(`DROP TABLE IF EXISTS files;`); err != nil {
		log.Println("Failed to drop table 'files':", err)
	}
	if _, err = DB.Exec(`DROP TABLE IF EXISTS file_changes;`); err != nil {
		log.Println("Failed to drop table 'file_changes':", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS files (
		file_id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT,
		file_hash TEXT,
		file_size INTEGER,
		modified DATETIME
	);`)
	if (err != nil) {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS file_changes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT,
		change_type TEXT CHECK (change_type IN ('new', 'modified', 'deleted', 'unchanged')),
		hash TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
}