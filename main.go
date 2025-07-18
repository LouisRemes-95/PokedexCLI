package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			CleanedInput := cleanInput(scanner.Text())
			if len(CleanedInput) != 0 {
				fmt.Print("Your command was: " + strings.ToLower(CleanedInput[0]) + "\n")
			}
		}
	}
}

func cleanInput(input string) []string {
	return strings.Fields(input)
}
