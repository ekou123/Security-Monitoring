package cmd

import (
	"fmt"
	"os"
)

// Command structure for REPL commands
type Command struct {
	Name        string
	Description string
	Callback    func(args []string)
}

// commands registry (declared but not filled yet)
var commands map[string]Command

// Exported so main.go can use it
func GetCommands() map[string]Command {
	return commands
}

// Initialize commands in init()
func init() {
	commands = map[string]Command{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the application",
			Callback:    commandExit,
		},
		"baseline": {
			Name:        "baseline",
			Description: "Create a baseline (usage: baseline [path])",
			Callback:    BaselineHandler, // TODO: implement
		},
		"scan": {
			Name:        "scan",
			Description: "Scan a file for changes (usage: scan [path])",
			Callback:    ScanHandler, // TODO: implement
		},
	}
}

// Built-in commands
func commandHelp(args []string) {
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf(" - %-10s %s\n", cmd.Name, cmd.Description)
	}
}

func commandExit(args []string) {
	fmt.Println("Exiting...")
	os.Exit(0)
}
