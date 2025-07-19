package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

	res, err := http.Get(url)
	if err != nil {
		return LocationAreas{}, fmt.Errorf("error getting locations: %v", err)
	}
	defer res.Body.Close()

	var locations LocationAreas
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locations); err != nil {
		return LocationAreas{}, fmt.Errorf("error decoding location: %v", err)
	}

	return locations, nil
}
