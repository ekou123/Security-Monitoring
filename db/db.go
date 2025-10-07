package db

import (
    "database/sql"
    "log"
	"fmt"

    _ "github.com/mattn/go-sqlite3"
)

func InitializeDB() {
	db, err := sql.Open("sqlite3", "./security_monitoring.db")
	if err != nil {
		fmt.Println("Deez")
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS file (
		file_id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT,
		file_hash TEXT,
		file_size INTEGER
	);`)
	if (err != nil) {
		log.Fatal(err)
	}
}