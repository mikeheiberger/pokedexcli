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
    callback    func() error
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
    }
}

func main() {
    initCommands()

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := cleanInput(scanner.Text())

        if command, ok := commands[input[0]]; ok {
            err := command.callback()
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

func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp() error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Print("Usage:\n\n")
    for _, value := range commands {
        fmt.Printf("%s: %s\n", value.name, value.description)
    }
    return nil
}
