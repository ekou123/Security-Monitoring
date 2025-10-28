package cmd

import (
	"fmt"
	"os"
)

// Command structure for REPL commands
type Command struct {
	Name        string
	Description string
	Usage 	 string
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
			Usage: "help | help [command]",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the application",
			Usage: "exit",
			Callback:    commandExit,
		},
		"baseline": {
			Name:        "baseline",
			Description: "Create a baseline for a file or directory",
			Usage: "baseline [path]",
			Callback:    BaselineHandler, // TODO: implement
		},
		"scan": {
			Name:        "scan",
			Description: "Scan a file for changes against its baseline",
			Usage: "scan [path]",
			Callback:    ScanHandler, // TODO: implement
		},
		"diff": {
			Name: "diff",
			Description: "Check differences in files",
			Usage: "diff",
			Callback: DiffHandler,
		},
	}
}

// Built-in commands
func commandHelp(args []string) {
	if (len(args) == 0) {
		fmt.Println("Available commands:")
		for _, cmd := range commands {
			fmt.Printf(" - %-10s %s\n", cmd.Name, cmd.Description)
		}
	} else {
		cmdName := args[0]
		cmd, exists := commands[cmdName]
		if (!exists) {
			fmt.Printf("Unknown command '%s'. Type 'help' to list all commands.\n", cmdName)
			return
		} 

		fmt.Println("\n===========================Command Details===========================")
		fmt.Printf("Command: %s\n", cmd.Name)
		fmt.Printf("Description: %s\n", cmd.Description)
		fmt.Printf("Usage: %s\n", cmd.Usage)
		fmt.Println("=====================================================================\n")
	
	}
}

func commandExit(args []string) {
	fmt.Println("Exiting...")
	os.Exit(0)
}
