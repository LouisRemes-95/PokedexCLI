package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pokecache "github.com/LouisRemes-95/PokedexCLI/internal/pokecache"
)

var pokeCache = pokecache.NewCache(5 * time.Second)

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(url string) (LocationAreas, error) {
	res, err := APIRequest(url)
	if err != nil {
		return LocationAreas{}, fmt.Errorf("error API request: %v", err)
	}

	var locations LocationAreas
	if err := json.Unmarshal(res, &locations); err != nil {
		return LocationAreas{}, fmt.Errorf("error decoding location: %v", err)
	}

	return locations, nil
}

func APIRequest(url string) ([]byte, error) {

	data, ok := pokeCache.Get(url)

	if ok {
		return data, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting locations: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error failed status: %d", res.StatusCode)
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error failed to read res.body: %v", err)
	}

	pokeCache.Add(url, data)
	return data, nil
}
