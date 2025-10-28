package cmd

import (
	"database/sql"
	"fmt"
	"example.com/ekou123/db"

	_ "github.com/mattn/go-sqlite3"
)

type FileChange struct {
	changeID   string
	fileID     string
	scanID     string
	filepath   string
	changeType string
	hash       string
	timestamp  string
}

func DiffHandler(args []string) {
	var scanID int
	var err error

	if len(args) == 0 {
		// Get the latest scan ID if none provided
		scanID, err = getLatestScanID()
		if err != nil {
			fmt.Println("Error getting latest scan ID:", err)
			return
		}
		fmt.Println("Using latest scan ID:", scanID)
	} else {
		// Try to parse scanID from argument (if user types `diff 2`)
		_, err = fmt.Sscanf(args[0], "%d", &scanID)
		if err != nil {
			fmt.Println("Invalid scan ID format. Use an integer, e.g. 'diff 2'")
			return
		}
	}

	if err := performDiff(scanID); err != nil {
		fmt.Println("Error performing diff:", err)
	}
}

func getLatestScanID() (int, error) {
	var latestID int
	query := "SELECT scan_id FROM scans ORDER BY scan_id DESC LIMIT 1"

	err := db.DB.QueryRow(query).Scan(&latestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no scans found in database")
		}
		return 0, err
	}
	return latestID, nil
}

func performDiff(scanID int) error {
	query := `
        SELECT change_id, file_id, scan_id, file_path, change_type, hash, timestamp
        FROM file_changes
        WHERE scan_id = ?
    `
	rows, err := db.DB.Query(query, scanID)
	if err != nil {
		return err
	}
	defer rows.Close()

	num := 0
	for rows.Next() {
		fc := &FileChange{}
		err := rows.Scan(
			&fc.changeID, &fc.fileID, &fc.scanID, &fc.filepath,
			&fc.changeType, &fc.hash, &fc.timestamp,
		)
		if err != nil {
			return err
		}

		num++
		fmt.Println(fc.changeID, fc.fileID, fc.scanID, fc.filepath, fc.changeType, fc.hash, fc.timestamp)
	}

	if num == 0 {
		fmt.Println("No differences detected within files.\n")
	}

	return nil
}
