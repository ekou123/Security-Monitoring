package cmd

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "database/sql"
    "time"
    "strings"

    "example.com/ekou123/db"

    _ "github.com/mattn/go-sqlite3"
)

var scanPath string
// var baselinePath string

type ScanCounters struct {
    Total     int
    New       int
    Modified  int
    Deleted   int
}

func ScanHandler(args []string) {
    path := "C:\\Users\\Ethan\\Desktop" // default path
    if len(args) > 0 {
        path = args[0]
    }

    _, err := os.Stat(path)
    if err != nil {
        fmt.Println("Cannot find path. It may be incorrect.")
        return
    }

    currentTime := time.Now().Format(time.RFC3339)
    res, err := db.DB.Exec(`INSERT INTO scans (timestamp, scanned_path, total_files, new_files, modified_files, deleted_files)
                            VALUES (?, ?, 0, 0, 0, 0)`, currentTime, path)
    if err != nil {
        fmt.Println("Failed to start scan session:", err)
        return
    }

    scanID, _ := res.LastInsertId()

    fmt.Printf("Starting scan on #%d on %s...\n", scanID, path)

    fileSystem := os.DirFS(path)

    baselinePath = path
    
    counters := &ScanCounters{}
    fs.WalkDir(fileSystem, ".", func(p string, d fs.DirEntry, e error) error {
        return performScan(scanID, p, d, e, counters)
    })

    // Step 3: Detect deleted files after the walk
    detectDeletedFiles(scanID, counters)

    // Step 4: Update scan summary
    _, err = db.DB.Exec(`UPDATE scans SET total_files=?, new_files=?, modified_files=?, deleted_files=? WHERE scan_id=?`,
        counters.Total, counters.New, counters.Modified, counters.Deleted, scanID)
    if err != nil {
        fmt.Println("Failed to update scan summary:", err)
    }

    fmt.Printf("Scan #%d complete. Total: %d | New: %d | Modified: %d | Deleted: %d\n",
        scanID, counters.Total, counters.New, counters.Modified, counters.Deleted)
}

func performScan(scanID int64, path string, d fs.DirEntry, err error, counters *ScanCounters) error {
    if db.DB == nil {
        fmt.Println("Database connection not initialized!")
        return nil
    }

    if err != nil || d.IsDir() {
        return err
    }

    absPath := filepath.Join(baselinePath, path)

    if strings.HasPrefix(filepath.Base(absPath), ".") {
        return nil
    }

    skipDirs := []string {
        `C:\Windows\Temp`,
        `C:\Windows\Logs`,
        `C:\Windows\SoftwareDistribution`,
        `C:\ProgramData\Microsoft\Windows\Caches`,
        `C:\Users\%USERNAME%\AppData\Local\Temp`,
    }

    for _, skip := range skipDirs {
        if strings.HasPrefix(strings.ToLower(absPath), strings.ToLower(skip)) {
            return nil
        }
    }

    watchExtensions := []string{
        ".exe", ".dll", ".sys", ".bat", ".ps1", ".vbs",
        ".ini", ".cfg", ".xml", ".reg", ".txt",
    }

    ext := strings.ToLower(filepath.Ext(absPath))
    shouldWatch := false

    for _, e := range watchExtensions {
        if e == ext {
            shouldWatch = true
            break
        }
    }

    if !shouldWatch {
        return nil
    }


    info, err := d.Info()
    if err != nil {
        return err
    }

    f, err := os.Open(absPath)
    if err != nil {
        fmt.Println("Failed to open:", absPath)
        return nil
    }
    defer f.Close()

    hash := sha256.New()
    io.Copy(hash, f)
    hashString := hex.EncodeToString(hash.Sum(nil))

    size := info.Size()
    modTime := info.ModTime().Format(time.RFC3339)
    counters.Total++

    // Check if file exists in baseline
    var fileID int
    var dbHash string
    err = db.DB.QueryRow(`SELECT file_id, file_hash FROM files WHERE file_path=?`, absPath).Scan(&fileID, &dbHash)

    if err == sql.ErrNoRows {
        // New file
        fmt.Println("🆕 New file:", absPath)
        res, _ := db.DB.Exec(`INSERT INTO files (file_path, file_hash, file_size, modified, last_seen_scan)
                              VALUES (?, ?, ?, ?, ?)`, absPath, hashString, size, modTime, scanID)
        fileID64, _ := res.LastInsertId()
        db.DB.Exec(`INSERT INTO file_changes (file_id, scan_id, file_path, change_type, hash)
                    VALUES (?, ?, ?, ?, ?)`, fileID64, scanID, absPath, "new", hashString)
        counters.New++
        return nil
    }

    if err != nil {
        fmt.Println("DB Read Error:", err)
        return nil
    }

    // Existing file: compare hash
    if dbHash != hashString {
        fmt.Println("✏️ Modified:", absPath)
        db.DB.Exec(`UPDATE files SET file_hash=?, file_size=?, modified=?, last_seen_scan=? WHERE file_id=?`,
            hashString, size, modTime, scanID, fileID)
        db.DB.Exec(`INSERT INTO file_changes (file_id, scan_id, file_path, change_type, hash)
                    VALUES (?, ?, ?, ?, ?)`, fileID, scanID, absPath, "modified", hashString)
        counters.Modified++
    } else {
        // Unchanged
        db.DB.Exec(`UPDATE files SET last_seen_scan=? WHERE file_id=?`, scanID, fileID)
        db.DB.Exec(`INSERT INTO file_changes (file_id, scan_id, file_path, change_type, hash)
                    VALUES (?, ?, ?, ?, ?)`, fileID, scanID, absPath, "unchanged", hashString)
    }

    return nil
}


func detectDeletedFiles(scanID int64, counters *ScanCounters) {
    rows, err := db.DB.Query(`SELECT file_id, file_path FROM files WHERE last_seen_scan < ?`, scanID)
    if err != nil {
        fmt.Println("Deletion check failed:", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var fileID int
        var filePath string
        rows.Scan(&fileID, &filePath)

        fmt.Println("🗑️ Deleted file:", filePath)
        db.DB.Exec(`INSERT INTO file_changes (file_id, scan_id, file_path, change_type, hash)
                    VALUES (?, ?, ?, ?, '')`, fileID, scanID, filePath, "deleted")
        counters.Deleted++
    }
}
