package cmd

import (
    "fmt"
    "example.com/ekou123/db"

    _ "github.com/mattn/go-sqlite3"
)

type FileChange struct {
	changeID string 
	fileID string
	scanID string 
	filepath string
	changeType string 
	hash string 
	timestamp string
}

func DiffHandler(args []string) {

	performDiff(nil)
	


    
    // monitor.Baseline(path, db)
}

func performDiff(err error) error {

	sql := "SELECT * FROM file_changes"

    rows, err := db.DB.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		fc := &FileChange{}

		err := rows.Scan(&fc.changeID, &fc.fileID, &fc.scanID, &fc.filepath, &fc.changeType, &fc.hash, &fc.timestamp)
		if err != nil {
			return err
		}

		if ()		
		fmt.Println(fc.changeID, fc.fileID, fc.scanID, fc.filepath, fc.changeType, fc.hash, fc.timestamp)
	}

	return nil
}
