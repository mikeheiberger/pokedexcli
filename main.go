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
    callback    func(*config, string) error
}

type config struct {
	prevUrl	string
	nextUrl	string
}

type LocationsResponse struct {
	Count	int			`json:"count"`
 	Next	string		`json:"next"`
	Prev	string		`json:"previous"`
	Results	[]struct {
		Name	string	`json:"name"`
		Url		string	`json:"url"`
	}`json:"results"`
}

type ExploreResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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
		"explore" : {
			name:			"explore",
			description:	"Displays the pokemon at a location",
			callback:		commandExplore,
		},
    }
}


func main() {
	const interval = 5 * time.Minute

    initCommands()

	cache = pokecache.NewCache(interval)

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
			param := ""
			if len(input) > 1 {
				param = input[1]
			}

            err := command.callback(&configuration, param)
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

func commandExit(conf *config, param string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(conf *config, param string) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Print("Usage:\n\n")
    for _, value := range commands {
        fmt.Printf("%s: %s\n", value.name, value.description)
    }
    return nil
}

func commandMap(conf *config, param string) error {
	jsonData, err := getJsonFromCacheOrServer(conf.nextUrl)
	if err != nil {
		return err
	}

	err = printLocationData(conf, jsonData)
	return err
}

func commandMapBack(conf *config, param string) error {
	if len(conf.prevUrl) == 0 {
		return errors.New("you're on the first page")
	}

	jsonData, err := getJsonFromCacheOrServer(conf.prevUrl)
	if err != nil {
		return err
	}

	err = printLocationData(conf, jsonData)
	return err
}

func commandExplore(conf *config, param string) error {
	const baseUrl = "https://pokeapi.co/api/v2/location-area/"

	if len(param) == 0 {
		return errors.New("Must pass a location to the explore command")
	}

	fullUrl := baseUrl + param + "/"
	jsonData, err := getJsonFromCacheOrServer(fullUrl)
	if err != nil {
		return err
	}

	var explore ExploreResponse
	err = json.Unmarshal(jsonData, &explore)
	if err != nil {
		return fmt.Errorf("Unmarshal failed: %v", err)
	}

	if len(explore.PokemonEncounters) == 0 {
		fmt.Println("No pokemon in the area!")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range explore.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func getJsonFromCacheOrServer(url string) ([]byte, error) {
	jsonData, ok := cache.Get(url)
	if !ok {
		var err error
		jsonData, err = pokedexapi.QueryPokedexApi(url)
		if err != nil {
			return nil, err
		}

		cache.Add(url, jsonData)
	}

	return jsonData, nil
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
