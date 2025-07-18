package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var registerOfCommands = map[string]cliCommand{}

func initCommands() {
	registerOfCommands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	registerOfCommands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
}

func main() {

	initCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			continue
		}

		CleanedInput := cleanInput(scanner.Text())
		if len(CleanedInput) == 0 {
			continue
		}

		command := strings.ToLower(CleanedInput[0])

		if cmd, ok := registerOfCommands[command]; ok {
			if err := cmd.callback(); err != nil {
				fmt.Println("Error: ", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(input string) []string {
	return strings.Fields(input)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range registerOfCommands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}
