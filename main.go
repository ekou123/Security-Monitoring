package main

import (
	"example.com/ekou123/repl"
	"example.com/ekou123/db"
	"fmt"
)

func main() {
	fmt.Println("Testing")
	db.InitializeDB()

	repl.StartREPL()
	
}
