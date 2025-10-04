package main

import (
	"example.com/ekou123/cmd"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func main() {
	fmt.Println("Welcome to your Threat Monitoring System")
	fmt.Println("Type 'help' to see available commands.")
	fmt.Println("Type 'exit' to quit the application.")

	reader := bufio.NewReader(os.Stdin)
	commands := cmd.GetCommands() // ✅ now works

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Split(input, " ")
		cmdName := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		chosenCommand, exists := commands[cmdName]
		if exists {
			chosenCommand.Callback(args)
		} else {
			fmt.Println("Unknown command:", cmdName)
			fmt.Println("Type 'help' to see available commands.")
		}
	}
}
