package main

import (
	"example.com/ekou123/repl"
	"example.com/ekou123/db"
)

func main() {
	db.InitializeDB()
	defer db.DB.Close()

	repl.StartREPL()
	
}
