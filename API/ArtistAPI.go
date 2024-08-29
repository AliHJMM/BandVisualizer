package API

import (
	"encoding/json"
	"net/http"
)

func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func GetArtistsLength() (int, error) {
	artists, err := FetchArtists()
	if err != nil {
		return 0, err
	}
	return len(artists), nil
}
