// ArtistAPI.go
package API

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
    cachedArtists   []Artist
    cachedLocations map[int][]string
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

func FetchAllLocations() (map[int][]string, error) {
    var err error
    locationOnce.Do(func() {
        resp, e := http.Get("https://groupietrackers.herokuapp.com/api/locations")
        if e != nil {
            err = e
            return
        }
        defer resp.Body.Close()

        var data AllLocationsData
        if e := json.NewDecoder(resp.Body).Decode(&data); e != nil {
            err = e
            return
        }

        cachedLocations = make(map[int][]string)
        for _, locEntry := range data.Index {
            cachedLocations[locEntry.ID] = locEntry.Locations
        }

        if len(cachedLocations) == 0 {
            log.Println("No locations data found.")
            err = fmt.Errorf("no locations data found")
            return
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
