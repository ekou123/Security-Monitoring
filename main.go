package main

import (
	"example.com/ekou123/repl"
	"example.com/ekou123/db"
)

func main() {
	db.InitializeDB()

	repl.StartREPL()
	
}
