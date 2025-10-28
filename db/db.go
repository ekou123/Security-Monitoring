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

	if _, err = DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
        log.Println("Warning: failed to enable foreign key enforcement:", err)
    }

	if _, err = DB.Exec(`DROP TABLE IF EXISTS files;`); err != nil {
		log.Println("Failed to drop table 'files':", err)
	}
	if _, err = DB.Exec(`DROP TABLE IF EXISTS file_changes;`); err != nil {
		log.Println("Failed to drop table 'file_changes':", err)
	}

	if _, err = DB.Exec(`DROP TABLE IF EXISTS scans;`); err != nil {
		log.Println("Failed to drop table 'file_changes':", err)
	}


	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS files (
		file_id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT,
		file_hash TEXT,
		file_size INTEGER,
		modified DATETIME,
		last_seen_scan INTEGER REFERENCES scans(id)
	);`)
	if (err != nil) {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS file_changes (
		change_id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_id INTEGER REFERENCES files(file_id),
		scan_id INTEGER REFERENCES scans(id),
		file_path TEXT,
		change_type TEXT CHECK (change_type IN ('new', 'modified', 'deleted', 'unchanged')),
		hash TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if (err != nil) {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time TEXT,
		scanned_path TEXT,
		total_files INTEGER,
		new_files INTEGER,
		modified_files INTEGER,
		deleted_files INTEGER
	)`)
	if err != nil {
		log.Fatal(err)
	}
}