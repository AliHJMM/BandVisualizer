package server

import (
	"BandVisualizer/API"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		renderErrorPage(w, "Bad Request: Query parameter is missing", http.StatusBadRequest)
		return
	}

	// Parse the query to separate the search term and type if provided
	searchTerm, searchType := parseQuery(query)

	artists, err := API.FetchArtists()
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch all locations once to improve performance
	locationsMap, err := API.FetchAllLocations()
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var filteredArtists []API.Artist
	for _, artist := range artists {
		matched := false
		switch searchType {
		case "artist/band", "":
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchTerm)) {
				matched = true
			}
		case "member":
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(searchTerm)) {
					matched = true
					break
				}
			}
		case "location":
			locations := locationsMap[artist.ID]
			for _, loc := range locations.Locations {
				if strings.Contains(strings.ToLower(loc), strings.ToLower(searchTerm)) {
					matched = true
					break
				}
			}
		case "creation date":
			if strings.Contains(strconv.Itoa(artist.CreationDate), searchTerm) {
				matched = true
			}
		case "first album date":
			if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(searchTerm)) {
				matched = true
			}
		default:
			// Search all fields if no specific type is provided
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchTerm)) ||
				strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(searchTerm)) ||
				strings.Contains(strconv.Itoa(artist.CreationDate), searchTerm) {
				matched = true
			}
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(searchTerm)) {
					matched = true
					break
				}
			}
			locations := locationsMap[artist.ID]
			for _, loc := range locations.Locations {
				if strings.Contains(strings.ToLower(loc), strings.ToLower(searchTerm)) {
					matched = true
					break
				}
			}
		}
		if matched {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	if len(filteredArtists) == 0 {
		renderErrorPage(w, "No results found for your search.", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, filteredArtists)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func parseQuery(query string) (string, string) {
	// Check if query contains " - "
	parts := strings.Split(query, " - ")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.ToLower(strings.TrimSpace(parts[1]))
	}
	return strings.TrimSpace(query), ""
}
