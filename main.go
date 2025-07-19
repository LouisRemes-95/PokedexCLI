package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/LouisRemes-95/PokedexCLI/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var registerOfCommands = map[string]cliCommand{}

var locations = internal.LocationAreas{
	Next:     "https://pokeapi.co/api/v2/location-area",
	Previous: "",
}

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

	registerOfCommands["map"] = cliCommand{
		name:        "map",
		description: "Displays the next 20 Pokémon location areas in order. Use repeatedly to explore more locations.",
		callback:    commandMap,
	}

	registerOfCommands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 Pokémon location areas in order. Use repeatedly to explore more locations.",
		callback:    commandMapB,
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

func commandMap() error {
	if locations.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	var err error
	locations, err = internal.GetLocations(locations.Next)
	fmt.Println(locations.Next)
	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB() error {
	if locations.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	var err error
	locations, err = internal.GetLocations(locations.Previous)
	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
