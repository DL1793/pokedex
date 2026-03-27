package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DL1793/pokedex/internal/pokeapi"
	"github.com/DL1793/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
}

type config struct {
	nextURL       *string
	prevURL       *string
	pokeapiClient pokeapi.Client
	arg           string
	pokedex       map[string]pokeapi.Pokemon
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Show area pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch pokemon",
			callback:    commandCatch,
		},
	}
}

func cleanInput(text string) []string {
	var output []string
	lowerText := strings.ToLower(text)
	output = strings.Fields(lowerText)
	return output
}

func commandMap(cfg *config) error {
	var currentURL string
	if cfg.nextURL == nil {
		currentURL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		currentURL = *cfg.nextURL
	}
	locations, err := cfg.pokeapiClient.GetLocations(currentURL)
	if err != nil {
		return err
	}

	cfg.nextURL = locations.Next
	cfg.prevURL = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	var currentURL string
	if cfg.prevURL == nil {
		fmt.Println("You're on the first page!")
		return nil
	} else {
		currentURL = *cfg.prevURL
	}
	locations, err := cfg.pokeapiClient.GetLocations(currentURL)
	if err != nil {
		return err
	}

	cfg.nextURL = locations.Next
	cfg.prevURL = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	validCommands := getCommands()
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, cmd := range validCommands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExplore(cfg *config) error {
	if cfg.arg == "" {
		fmt.Println("Usage: explore <location-name>")
		return errors.New("no location name provided")
	}

	encounters, err := cfg.pokeapiClient.GetPokemon("https://pokeapi.co/api/v2/location-area/" + cfg.arg)
	if len(encounters.Results) == 0 {
		return errors.New("location not found")
	}
	if err != nil {
		fmt.Println("Error getting location:", err)
		return err
	}
	fmt.Println("Exploring " + cfg.arg + "...")
	fmt.Println("Found Pokemon:")
	for _, encounter := range encounters.Results {
		fmt.Println("- " + encounter.Result.Name)
	}
	return nil
}

func commandCatch(cfg *config) error {
	if cfg.arg == "" {
		fmt.Println("Usage: catch <pokemon-name>")
		return errors.New("no pokemon name provided")
	}
	pokemon, err := cfg.pokeapiClient.CatchPokemon("https://pokeapi.co/api/v2/pokemon/" + cfg.arg)
	if err != nil {
		fmt.Println("Error getting pokemon:", err)
		return err
	}
	fmt.Println("Throwing a Pokeball at " + cfg.arg + "...")
	var maxDifficulty float32 = 365.0 //Highest base experience - Lowest base experience
	var catchDifficulty float32 = float32(pokemon.BaseExperience-25) / maxDifficulty
	catchChance := rand.Float32()
	if catchChance < catchDifficulty {
		fmt.Println(cfg.arg + " escaped!")
		return nil
	} else {
		fmt.Println(cfg.arg + " was caught!")
		cfg.pokedex[cfg.arg] = pokemon
		return nil
	}
}

func startRepl() {
	cache := pokecache.NewCache(5 * time.Minute)
	cfg := config{
		nextURL: nil,
		prevURL: nil,
		pokeapiClient: pokeapi.Client{
			cache,
			http.Client{},
		},
		arg:     "",
		pokedex: make(map[string]pokeapi.Pokemon),
	}
	validCommands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		fields := cleanInput(text)
		if len(fields) == 0 {
			continue
		} else {
			if len(fields) == 2 {
				cfg.arg = fields[1]
			}
			val, ok := validCommands[fields[0]]
			if !ok {
				fmt.Println("Unknown command")
			} else {

				err := val.callback(&cfg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

	}
}
