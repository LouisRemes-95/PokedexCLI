package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/LouisRemes-95/PokedexCLI/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

var registerOfCommands = map[string]cliCommand{}

var locations = pokeapi.LocationAreas{
	Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
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

	registerOfCommands["explore"] = cliCommand{
		name:        "explore",
		description: "Displays the pokemons found in the stated area",
		callback:    commandExplore,
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

		for i, str := range CleanedInput {
			CleanedInput[i] = strings.ToLower(str)
		}

		command := CleanedInput[0]

		if cmd, ok := registerOfCommands[command]; ok {
			if err := cmd.callback(CleanedInput[1:]); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(input string) []string {
	return strings.Fields(input)
}

func commandExit([]string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp([]string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range registerOfCommands {
		fmt.Println(command.name + ": " + command.description)
	}
	return nil
}

func commandMap([]string) error {
	if locations.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}

	var err error
	locations, err = pokeapi.GetLocations(locations.Next)
	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB([]string) error {
	if locations.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	var err error
	locations, err = pokeapi.GetLocations(locations.Previous)
	if err != nil {
		return err
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(input []string) error {
	if len(input) == 0 {
		fmt.Println("Missing location")
		return nil
	}

	area, err := pokeapi.GetArea("https://pokeapi.co/api/v2/location-area/" + input[0])
	if err != nil {
		var HTTPError *pokeapi.HTTPError
		switch {
		case errors.As(err, &HTTPError) && HTTPError.StatusCode > 500:
			fmt.Println("Server side issue")
			return nil
		case errors.As(err, &HTTPError) && HTTPError.StatusCode > 400:
			fmt.Println("Location not found")
			return nil
		default:
			return err
		}
	}

	fmt.Println("Exploring" + input[0] + "...")
	fmt.Println("Found Pokemon:")
	for _, encounter := range area.PokemonEncounters {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}

	return nil
}
