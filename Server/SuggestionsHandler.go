package server

import (
	"BandVisualizer/API"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func SuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	artists, err := API.FetchArtists()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch all locations once to improve performance
	locationsMap, err := API.FetchAllLocations()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var suggestions []string
	lowerQuery := strings.ToLower(query)

	for _, artist := range artists {
		// Artist/Band Name
		if strings.Contains(strings.ToLower(artist.Name), lowerQuery) {
			suggestion := artist.Name + " - artist/band"
			suggestions = append(suggestions, suggestion)
		}

		// Members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), lowerQuery) {
				suggestion := member + " - member"
				suggestions = append(suggestions, suggestion)
			}
		}

		// Creation Date
		creationDateStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(creationDateStr, query) {
			suggestion := creationDateStr + " - creation date"
			suggestions = append(suggestions, suggestion)
		}

		// First Album Date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), lowerQuery) {
			suggestion := artist.FirstAlbum + " - first album date"
			suggestions = append(suggestions, suggestion)
		}

		// Locations
		locations := locationsMap[artist.ID]
		for _, loc := range locations.Locations {
			if strings.Contains(strings.ToLower(loc), lowerQuery) {
				suggestion := loc + " - location"
				suggestions = append(suggestions, suggestion)
			}
		}
	}

	// Remove duplicates
	suggestions = removeDuplicates(suggestions)

	// Convert to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func removeDuplicates(suggestions []string) []string {
	unique := make(map[string]bool)
	var result []string
	for _, suggestion := range suggestions {
		if !unique[suggestion] {
			unique[suggestion] = true
			result = append(result, suggestion)
		}
	}
	return result
}
