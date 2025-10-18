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

    "example.com/ekou123/db"

    _ "github.com/mattn/go-sqlite3"
)

var scanPath string
// var baselinePath string

func ScanHandler(args []string) {
    path := "C:\\Users\\Ethan\\Desktop" // default path
    if len(args) > 0 {
        path = args[0]
    }

    fileSystem := os.DirFS(path)

    baselinePath = path
    
    _, err := os.Stat(path)
    if err != nil {
        fmt.Println("Cannot find path, It may be incorrect")
        return
    } else {
        displayMessage := fmt.Sprintf("Scanning %s now...", path)
        fmt.Println(displayMessage)
        fs.WalkDir(fileSystem, ".", performScan)
    }
    // monitor.Baseline(path, db)
}

func performScan(path string, d fs.DirEntry, err error) error {
    if db.DB == nil {
        fmt.Println("Database connection not initialized!")
        return nil
    }

    if err != nil {
        return err
    }

    if d.IsDir() {
        return nil
    }

    info, err := d.Info()
    if err != nil {
        return err
    }

    absPath := filepath.Join(baselinePath, path)

    f, err := os.Open(absPath)
    if err != nil {
        fmt.Println("Failed to open: ", absPath)
        return nil
    }
    defer f.Close()

    size := info.Size()
    modTime := info.ModTime().Format(time.RFC3339)

    var count int 
    err = db.DB.QueryRow(`SELECT COUNT(*) FROM files WHERE file_path = ?`, absPath).Scan(&count)
    if err != nil {
        return err
    }

    if count <= 0 {
        fmt.Println("")
    }

    hash := sha256.New()
    if _, err := io.Copy(hash, f); err != nil {
        return err
    }


    hashString := hex.EncodeToString(hash.Sum(nil))

    

    row := db.DB.QueryRow(
        `SELECT file_hash, file_size, modified FROM files WHERE file_path = ?`,
        absPath,
    )
    
    var dbHash string
    var dbSize string
    var dbModTime string

    err = row.Scan(&dbHash, &dbSize, &dbModTime)
    if err == sql.ErrNoRows {
        fmt.Printf("New file detected: %s | Adding to files", absPath)
        _, err = db.DB.Exec(`INSERT INTO files (file_path, file_hash, file_size, modified) VALUES (?, ?, ?, ?)`, absPath, hashString, size, modTime,)
        if err != nil {
            fmt.Println("DB Insert Error: ", err)
        }
        _, err = db.DB.Exec(`INSERT INTO file_changes (file_path, change_type, hash) VALUES (?, ?, ?)`, absPath, "new", hashString,)
        if err != nil {
            fmt.Println("DB Insert Error: ", err)
        }
        return nil
    }
    if err != nil {
        fmt.Println("Database Read Error for", absPath, ":", err)
        return nil
    }

    if dbHash != hashString {
        fmt.Println("Modified file detected:", absPath)
        db.DB.Exec(`INSERT INTO file_changes (file_path, change_type, hash) VALUES (?, ?, ?)`, absPath, "modified", hashString,)
    } else {
        fmt.Println("Unchanged:", absPath)
        db.DB.Exec(`INSERT INTO file_changes (file_path, change_type, hash) VALUES (?, ?, ?)`, absPath, "unchanged", hashString,)
    }

    
    return nil
}
