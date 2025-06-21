package pokeApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/adis-abazovic/pokedexcli/internal/pokecache"
)

type PokeApiClient struct {
	cache pokecache.Cache
}

func NewPokeApiClient(cacheInterval time.Duration) PokeApiClient {
	return PokeApiClient{
		cache: pokecache.NewCache(cacheInterval),
	}
}

type LocationResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonInfo struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
}

const baseUrl string = "https://pokeapi.co/api/v2"

func (c *PokeApiClient) GetPokemonsAtLocation(location string) (PokemonResponse, error) {

	// read from cache
	url := fmt.Sprintf("%s/location-area/%s", baseUrl, location)
	val, ok := c.cache.Get(url)

	if ok {
		pokResp := PokemonResponse{}
		err := json.Unmarshal(val, &pokResp)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("error: unmarshaling response body failed")
		}

		return pokResp, nil
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("error: invalid response for URL: '%s'", url)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("error: invalid response body")
		}
		defer resp.Body.Close()

		// add to cache
		c.cache.Add(url, body)

		pokResp := PokemonResponse{}
		err = json.Unmarshal(body, &pokResp)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("error: unmarshaling response body failed")
		}

		return pokResp, nil
	}
}

func (c *PokeApiClient) GetLocation(url string) (LocationResponse, error) {

	// read from cache
	val, ok := c.cache.Get(url)

	if ok {
		locResp := LocationResponse{}
		err := json.Unmarshal(val, &locResp)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("error: unmarshaling response body failed")
		}

		return locResp, nil
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("error: invalid response for URL: '%s'", url)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("error: invalid response body")
		}
		defer resp.Body.Close()

		// add to cache
		c.cache.Add(url, body)

		locResp := LocationResponse{}
		err = json.Unmarshal(body, &locResp)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("error: unmarshaling response body failed")
		}

		return locResp, nil
	}
}

func (c *PokeApiClient) GetPokemon(name string) (PokemonInfo, error) {

	// read from cache
	url := fmt.Sprintf("%s/pokemon/%s", baseUrl, name)
	val, ok := c.cache.Get(url)

	if ok {
		pokResp := PokemonInfo{}
		err := json.Unmarshal(val, &pokResp)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("error: unmarshaling response body failed")
		}

		return pokResp, nil
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("error: invalid response for URL: '%s'", url)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("error: invalid response body")
		}
		defer resp.Body.Close()

		// add to cache
		c.cache.Add(url, body)

		pokResp := PokemonInfo{}
		err = json.Unmarshal(body, &pokResp)
		if err != nil {
			return PokemonInfo{}, fmt.Errorf("error: unmarshaling response body failed")
		}

		return pokResp, nil
	}
}
