package cmd

import (
	"fmt"
)

var scanPath string

func ScanHandler(args []string) {
    path := "/etc"
    if len(args) > 0 {
        path = args[0]
    }
    fmt.Println("Scanning ", path)
    // monitor.Baseline(path, db)
}
