package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	pokeApi "github.com/adis-abazovic/pokedexcli/internal/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	nextUrl        string
	previousUrl    string
	pokeApiClient  pokeApi.PokeApiClient
	caughtPokemons map[string]pokeApi.PokemonInfo
}

func main() {

	pokeApiClient := pokeApi.NewPokeApiClient(10 * time.Second)
	cfg := config{
		nextUrl:        "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		previousUrl:    "",
		pokeApiClient:  pokeApiClient,
		caughtPokemons: make(map[string]pokeApi.PokemonInfo),
	}

	cmds := getCommands()

	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)

		fields := strings.Fields(input)
		cmdName := fields[0]
		var args []string
		for i := 1; i < len(fields); i++ {
			args = append(args, fields[i])
		}

		cmd, ok := cmds[cmdName]
		if !ok {
			fmt.Println("Unknown command")
		}

		err := cmd.callback(&cfg, args)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getCommands() map[string]cliCommand {

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the names of 20 location areas in the Pokemon world",
			callback:    commandMapForward,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations areas in the Pokemon world",
			callback:    commandMapBackward,
		},
		"explore": {
			name:        "explore",
			description: "list of all the Pokémon located at given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch the Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect the Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "print all the names of the Pokemon the user has caught",
			callback:    commandPokedex,
		},
	}
}

func commandPokedex(cfg *config, args []string) error {

	fmt.Println("Your Pokedex:")
	for _, p := range cfg.caughtPokemons {
		fmt.Printf("\n - %s", p.Name)
	}

	fmt.Println()
	return nil
}

func commandInspect(cfg *config, args []string) error {

	if len(args) == 0 {
		return fmt.Errorf("command 'inspect' requires argument (pokemon name)")
	}

	name := args[0]
	p, ok := cfg.caughtPokemons[name]
	if !ok {
		fmt.Printf("\nyou have not caught that pokemon\n")
		return nil
	}

	fmt.Printf("\nName: %s", p.Name)
	fmt.Printf("\nHeight: %d", p.Height)
	fmt.Printf("\nWeight: %d", p.Weight)

	fmt.Printf("\nStats:")
	for _, s := range p.Stats {
		fmt.Printf("\n  -%s: %d", s.Stat.Name, s.BaseStat)
	}

	fmt.Printf("\nTypes:")
	for _, t := range p.Types {
		fmt.Printf("\n  -%s", t.Type.Name)
	}

	fmt.Println()

	return nil
}

func commandCatch(cfg *config, args []string) error {

	if len(args) == 0 {
		return fmt.Errorf("command 'catch' requires argument (location)")
	}

	fmt.Printf("\nThrowing a Pokeball at %s...\n", args[0])

	pokResp, err := cfg.pokeApiClient.GetPokemon(args[0])
	if err != nil {
		return fmt.Errorf("error fetching pokemon info")
	}

	fmt.Println()

	chance := 300 - pokResp.BaseExperience
	roll := rand.Intn(300)
	if roll < chance {

		fmt.Printf("\n%s was caught!\n", pokResp.Name)

		cfg.caughtPokemons[pokResp.Name] = pokResp

	} else {
		fmt.Printf("\n%s escaped!\n", pokResp.Name)
	}

	return nil
}

func commandExplore(cfg *config, args []string) error {

	if len(args) == 0 {
		return fmt.Errorf("command 'explore' requires argument (location)")
	}

	pokResp, err := cfg.pokeApiClient.GetPokemonsAtLocation(args[0])
	if err != nil {
		return fmt.Errorf("error fetching pokemons")
	}

	for _, res := range pokResp.PokemonEncounters {
		fmt.Println(res.Pokemon.Name)
	}

	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Printf("\nWelcome to the Pokedex!")
	fmt.Printf("\nUsage:\n")

	cmds := getCommands()
	for k, v := range cmds {
		fmt.Printf("\n%s: %s", k, v.description)
	}

	fmt.Printf("\n")

	return nil
}

func commandMapForward(cfg *config, args []string) error {

	locResp, err := cfg.pokeApiClient.GetLocation(cfg.nextUrl)
	if err != nil {
		return fmt.Errorf("error fetching locations from '%s'", cfg.nextUrl)
	}

	for _, res := range locResp.Results {
		fmt.Println(res.Name)
	}

	cfg.nextUrl = locResp.Next
	if locResp.Previous != nil {
		cfg.previousUrl = *locResp.Previous
	}

	return nil
}

func commandMapBackward(cfg *config, args []string) error {

	if cfg.previousUrl == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locResp, err := cfg.pokeApiClient.GetLocation(cfg.previousUrl)
	if err != nil {
		return fmt.Errorf("error fetching locations from '%s'", cfg.nextUrl)
	}

	for _, res := range locResp.Results {
		fmt.Println(res.Name)
	}

	cfg.nextUrl = locResp.Next
	if locResp.Previous != nil {
		cfg.previousUrl = *locResp.Previous
	}

	return nil
}

func cleanInput(text string) []string {
	var s []string

	splitted := strings.Split(text, " ")

	for _, w := range splitted {
		t := strings.TrimSpace(w)
		if len(t) != 0 {
			lower := strings.ToLower(w)
			s = append(s, lower)
		}
	}

	return s
}
