package cmd

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "time"

    "example.com/ekou123/db"

    _ "github.com/mattn/go-sqlite3"
)

var baselinePath string



func BaselineHandler(args []string) {
    path := "C:\\Users\\Ethan\\Desktop" // default path
    if len(args) > 0 {
        path = args[0]
    }

    baselinePath = path

    fileSystem := os.DirFS(path)

    

    fmt.Println("Creating baseline for", path)
    fs.WalkDir(fileSystem, ".", performBaseline)

    


    // monitor.Baseline(path, db)
}

func performBaseline(path string, d fs.DirEntry, err error) error {
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

    hash := sha256.New()
    if _, err := io.Copy(hash, f); err != nil {
        return err
    }

    hashString := hex.EncodeToString(hash.Sum(nil))

    size := info.Size()
    modTime := info.ModTime().Format(time.RFC3339)

    _, err = db.DB.Exec(
        `INSERT INTO files (file_path, file_hash, file_size, modified)
        VALUES (?, ?, ?, ?)`,
        path, hashString, size, modTime,
    )
    if err != nil {
        fmt.Println("DB Insert error for", path, ":", err)
        return nil
    }

    fmt.Println("Added:", path)
    
    return nil
}
