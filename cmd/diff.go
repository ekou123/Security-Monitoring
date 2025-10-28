package cmd

import (
	"database/sql"
	"fmt"
	"example.com/ekou123/db"

	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
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

	if err := OpenDiffConsole(scanID); err != nil {
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

// func performDiff(scanID int) error {
// 	query := `
//         SELECT change_id, file_id, scan_id, file_path, change_type, hash, timestamp
//         FROM file_changes
//         WHERE scan_id = ?
//     `
// 	rows, err := db.DB.Query(query, scanID)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	fmt.Printf("\nDiff Results for Scan #%d\n", scanID)
// 	fmt.Println("──────────────────────────────────────────────────────────────")

// 	// Create a tabwriter for aligned columns
// 	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
// 	fmt.Fprintln(w, "Change ID\tFile ID\tType\tFile Path\tHash\tTimestamp")

// 	num := 0
// 	for rows.Next() {
// 		fc := &FileChange{}
// 		err := rows.Scan(
// 			&fc.changeID, &fc.fileID, &fc.scanID, &fc.filepath,
// 			&fc.changeType, &fc.hash, &fc.timestamp,
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		// Skip unchanged
// 		if fc.changeType == "unchanged" {
// 			continue
// 		}

// 		num++
// 		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
// 			fc.changeID, fc.fileID, fc.changeType, fc.filepath, fc.hash[:10]+"...", fc.timestamp)
// 	}

// 	w.Flush()

// 	if num == 0 {
// 		fmt.Println("\nNo differences detected within files.\n")
// 	} else {
// 		fmt.Println("──────────────────────────────────────────────────────────────")
// 		fmt.Printf("Total: %d changes detected.\n\n", num)
// 	}

// 	return nil
// }

func OpenDiffConsole(scanID int) error {
	query := `
        SELECT change_id, file_id, scan_id, file_path, change_type, hash, timestamp
        FROM file_changes
        WHERE scan_id = ?
        ORDER BY change_type
    `
	rows, err := db.DB.Query(query, scanID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var changes []FileChange
	for rows.Next() {
		var fc FileChange
		rows.Scan(&fc.changeID, &fc.fileID, &fc.scanID, &fc.filepath, &fc.changeType, &fc.hash, &fc.timestamp)
		if fc.changeType != "unchanged" {
			changes = append(changes, fc)
		}
	}

	if len(changes) == 0 {
		fmt.Println("\nNo differences detected.\n")
		return nil
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(false)
	table.SetBorder(true).SetTitle(fmt.Sprintf("📊 Diff Console – Scan #%d", scanID))

	// headers
	headers := []string{"Change ID", "File ID", "Type", "File Path", "Hash", "Timestamp"}
	for i, h := range headers {
		cell := tview.NewTableCell(fmt.Sprintf("[yellow::b]%s", h))
		table.SetCell(0, i, cell)
	}

	// data rows
	for r, fc := range changes {
		row := r + 1
		color := "[white]"
		switch fc.changeType {
		case "new":
			color = "[green]"
		case "modified":
			color = "[yellow]"
		case "deleted":
			color = "[red]"
		}
		table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%s%s", color, fc.changeID)))
		table.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%s%s", color, fc.fileID)))
		table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%s%s", color, fc.changeType)))
		table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%s%s", color, fc.filepath)))
		table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%s%s...", color, fc.hash[:10])))
		table.SetCell(row, 5, tview.NewTableCell(fmt.Sprintf("%s%s", color, fc.timestamp)))
	}

	table.Select(1, 0).SetFixed(1, 0).SetSelectable(true, false)

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
			app.Stop()
		}
	})

	table.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		fc := changes[row-1]
		details := fmt.Sprintf(
			"File: %s\nType: %s\nHash: %s\nTimestamp: %s",
			fc.filepath, fc.changeType, fc.hash, fc.timestamp,
		)
		modal := tview.NewModal().
			SetText(details).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.SetRoot(table, true).SetFocus(table)
			})
		app.SetRoot(modal, true).SetFocus(modal)
	})

	if err := app.SetRoot(table, true).Run(); err != nil {
		return err
	}
	return nil
}

