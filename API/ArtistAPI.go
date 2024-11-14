package API

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	cachedArtists   []Artist
	cachedLocations map[int]LocationData
	artistOnce      sync.Once
	locationOnce    sync.Once
)

func FetchArtists() ([]Artist, error) {
	var err error
	artistOnce.Do(func() {
		resp, e := http.Get("https://groupietrackers.herokuapp.com/api/artists")
		if e != nil {
			err = e
			return
		}
		defer resp.Body.Close()

		var artists []Artist
		if e := json.NewDecoder(resp.Body).Decode(&artists); e != nil {
			err = e
			return
		}
		cachedArtists = artists
	})
	return cachedArtists, err
}

func FetchAllLocations() (map[int]LocationData, error) {
	var err error
	locationOnce.Do(func() {
		resp, e := http.Get("https://groupietrackers.herokuapp.com/api/locations")
		if e != nil {
			err = e
			return
		}
		defer resp.Body.Close()

		var allLocations []struct {
			ID        int         `json:"id"`
			Locations LocationData `json:"locations"`
		}

		if e := json.NewDecoder(resp.Body).Decode(&allLocations); e != nil {
			err = e
			return
		}

		cachedLocations = make(map[int]LocationData)
		for _, loc := range allLocations {
			cachedLocations[loc.ID] = loc.Locations
		}
	})
	return cachedLocations, err
}

func GetArtistsLength() (int, error) {
	artists, err := FetchArtists()
	if err != nil {
		return 0, err
	}
	return len(artists), nil
}
