package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
	"errors"
	"time"
	"encoding/json"
	"github.com/mikeheiberger/pokedexcli/internal/pokedexapi"
	"github.com/mikeheiberger/pokedexcli/internal/pokecache"
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

type LocationsResponse struct {
	Count	int			`json:"count"`
 	Next	string		`json:"next"`
	Prev	string		`json:"previous"`
	Results	[]Location	`json:"results"`
}

type Location struct {
	Name	string	`json:"name"`
	Url		string	`json:"url"`
}

var commands map[string]cliCommand

var cache *pokecache.Cache

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

	cacheDur, _ := time.ParseDuration("5s")
	cache = pokecache.NewCache(cacheDur)

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
	jsonData, ok := cache.Get(conf.nextUrl)
	if !ok {
		var err error
		jsonData, err = pokedexapi.QueryPokedexApi(conf.nextUrl)
		if err != nil {
			return err
		}

		cache.Add(conf.nextUrl, jsonData)
	}

	err := printLocationData(conf, jsonData)
	return err
}

func commandMapBack(conf *config) error {
	if len(conf.prevUrl) == 0 {
		return errors.New("you're on the first page")
	}

	jsonData, ok := cache.Get(conf.prevUrl)
	if !ok {
		var err error
		jsonData, err = pokedexapi.QueryPokedexApi(conf.prevUrl)
		if err != nil {
			return err
		}

		cache.Add(conf.prevUrl, jsonData)
	}

	err := printLocationData(conf, jsonData)
	return err
}

func printLocationData(conf *config, jsonData []byte) error {
	var locations LocationsResponse
	err := json.Unmarshal(jsonData, &locations)
	if err != nil {
		return fmt.Errorf("Unmarshal failed: %v", err)
	}

	conf.nextUrl = locations.Next
	conf.prevUrl = locations.Prev

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
