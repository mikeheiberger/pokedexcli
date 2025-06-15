package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
	"net/http"
	"io"
	"encoding/json"
	"errors"
)

type cliCommand struct {
    name        string
    description string
    callback    func(*config) error
}

type config struct {
	prevUrl	string
	nextUrl	string
}

type locationsResponse struct {
	Count	int			`json:"count"`
 	Next	string		`json:"next"`
	Prev	string		`json:"previous"`
	Results	[]location	`json:"results"`
}

type location struct {
	Name	string	`json:"name"`
	Url		string	`json:"url"`
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
	locations, err := queryPokedexApi(conf.nextUrl)
	if err != nil {
		return err
	}

	conf.nextUrl = locations.Next
	conf.prevUrl = locations.Prev

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapBack(conf *config) error {
	if len(conf.prevUrl) == 0 {
		return errors.New("you're on the first page")
	}

	locations, err := queryPokedexApi(conf.prevUrl)
	if err != nil {
		return err
	}

	conf.nextUrl = locations.Next
	conf.prevUrl = locations.Prev

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func queryPokedexApi(url string) (*locationsResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	var locations locationsResponse
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal failed: %v", err)
	}

	return &locations, nil
}
