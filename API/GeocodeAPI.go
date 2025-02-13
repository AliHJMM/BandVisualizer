// API/GeocodeAPI.go
package API

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Geodata struct {
	Index []struct {
		CountryCoords map[string][]string `json:"countryCoords"`
	} `json:"index"`
}

var allGeodata Geodata

// GetGeocode reads the geocodes file and returns coordinates for the requested locations.
func GetGeocode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.FormValue("query")
	if query == "" {
		http.Error(w, "Bad Request: missing query", http.StatusBadRequest)
		return
	}

	// Open the geocodes file (assumed to be in static/data/)
	file, err := os.Open("static/data/geocodes.json")
	if err != nil {
		log.Println("Error opening geocodes.json:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Error reading geocodes.json:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(bytes, &allGeodata); err != nil {
		log.Println("Error unmarshalling geocodes.json:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// The query should be a comma-separated list of location names.
	queries := strings.Split(query, ",")
	type Country struct {
		Name   string   `json:"Name"`
		Coords []string `json:"Coords"`
	}
	var response []Country
	for _, q := range queries {
		// Normalize the query string to lowercase
		q = strings.ToLower(strings.TrimSpace(q))
		for _, entry := range allGeodata.Index {
			if coords, ok := entry.CountryCoords[q]; ok {
				response = append(response, Country{Name: q, Coords: coords})
				break
			}
		}
	}
	

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
