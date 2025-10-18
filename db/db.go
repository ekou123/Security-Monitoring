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

	_, err = DB.Exec(`DROP TABLE IF EXISTS files`)

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
}