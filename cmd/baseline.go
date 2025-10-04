package cmd

import (
	"fmt"
)

var baselinePath string

func BaselineHandler(args []string) {
    path := "/etc"
    if len(args) > 0 {
        path = args[0]
    }
    fmt.Println("Creating baseline for", path)
    // monitor.Baseline(path, db)
}
