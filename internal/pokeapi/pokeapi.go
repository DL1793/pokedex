package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/DL1793/pokedex/internal/pokecache"
)

type Client struct {
	Cache      *pokecache.Cache
	HttpClient http.Client
}

type LocationResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Locations struct {
	Results  []LocationResource `json:"results"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
}

type PokemonEncounters struct {
	Results []Pokemon `json:"pokemon_encounters"`
}

type Pokemon struct {
	Result LocationResource `json:"pokemon"`
}

func (c *Client) GetLocations(url string) (Locations, error) {

	var locations Locations
	if cachedBytes, ok := c.Cache.Get(url); ok {
		err := json.Unmarshal(cachedBytes, &locations)
		if err != nil {
			return Locations{}, err
		}
		return locations, nil
	} else {
		res, err := c.HttpClient.Get(url)
		if err != nil {
			return Locations{}, err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			log.Fatal("Response failed with status code ", res.StatusCode)
			return Locations{}, err
		}
		if err != nil {
			log.Fatal(err)
			return Locations{}, err
		}
		c.Cache.Add(url, body)
		err = json.Unmarshal(body, &locations)
		if err != nil {
			log.Fatal(err)
			return Locations{}, err
		}
		return locations, nil
	}
}

func (c *Client) GetPokemon(url string) (PokemonEncounters, error) {
	var encounters PokemonEncounters
	if cachedBytes, ok := c.Cache.Get(url); ok {
		err := json.Unmarshal(cachedBytes, &encounters)
		if err != nil {
			return PokemonEncounters{}, err
		}
		return encounters, nil
	}
	res, err := c.HttpClient.Get(url)
	if err != nil {
		return PokemonEncounters{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Println("Response failed with status code", res.StatusCode)
		return PokemonEncounters{}, err
	}
	if err != nil {
		log.Fatal(err)
		return PokemonEncounters{}, err
	}
	c.Cache.Add(url, body)
	err = json.Unmarshal(body, &encounters)
	if err != nil {
		log.Fatal(err)
		return PokemonEncounters{}, err
	}
	return encounters, nil
}
