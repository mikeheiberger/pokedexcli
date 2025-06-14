package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
)

type cliCommand struct {
    name        string
    description string
    callback    func(*config) error
}

type config struct {
	nextUrl	string
	prevUrl	string
}

var commands map[string]cliCommand

func initCommands() {
    commands = map[string]cliCommand{
        "help" : {
            name:           "help",
            description:    "Displays a help message",
            callback:       commandHelp,
        },
        "exit" : {
            name:           "exit",
            description:    "Exit the pokedex",
            callback:       commandExit,
        },
		"map" : {
			name:			"map",
			description:	"Displays the next 20 locations",
			callback:		commandMap,
		},
		"mapb" : {
			name:			"mapb",
			description:	"Displays the previous 20 locations",
			callback:		commandMapBack,
		},
    }
}

func main() {
    initCommands()

	configuration := config{
		"",
		"https://pokeapi.co/api/v2/location-area/",
	}

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())

        if command, ok := commands[input[0]]; ok {
            err := command.callback(&configuration)
            if err != nil {
                fmt.Println(err.Error())
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}

func cleanInput(text string) []string {
    split := strings.Fields(strings.ToLower(text))
    return split
}

func commandExit(conf *config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(conf *config) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Print("Usage:\n\n")
    for _, value := range commands {
        fmt.Printf("%s: %s\n", value.name, value.description)
    }
    return nil
}

func commandMap(conf *config) error {
	// TODO: get request from next field

	return nil
}

func commandMapBack(conf *config) error {
	if config.prevUrl == "" {
		fmt.Println("you're on the first page")
	}

	return nil
}
