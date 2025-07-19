package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pokecache "github.com/LouisRemes-95/PokedexCLI/internal/pokecache"
)

type HTTPError struct {
	StatusCode int
	URL        string
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP request to %s failed with status %d: %s", e.URL, e.StatusCode, e.Message)
}

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
		return LocationAreas{}, err
	}

	var locations LocationAreas
	if err := json.Unmarshal(res, &locations); err != nil {
		return LocationAreas{}, fmt.Errorf("error decoding location: %v", err)
	}

	return locations, nil
}

type Area struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetArea(url string) (Area, error) {
	res, err := APIRequest(url)
	if err != nil {
		return Area{}, err
	}

	var area Area
	if err := json.Unmarshal(res, &area); err != nil {
		return Area{}, fmt.Errorf("error decoding area: %v", err)
	}

	return area, nil
}

func APIRequest(url string) ([]byte, error) {

	data, ok := pokeCache.Get(url)

	if ok {
		return data, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, &HTTPError{StatusCode: 0, URL: url, Message: "GET failled"}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, &HTTPError{StatusCode: res.StatusCode, URL: url, Message: "Bad status code"}
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error failed to read res.body: %v", err)
	}

	pokeCache.Add(url, data)
	return data, nil
}
