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
    lowerSearchTerm := strings.ToLower(searchTerm)

    for _, artist := range artists {
        matched := false
        switch searchType {
        case "artist/band":
            if strings.Contains(strings.ToLower(artist.Name), lowerSearchTerm) {
                matched = true
            }
        case "member":
            for _, member := range artist.Members {
                if strings.Contains(strings.ToLower(member), lowerSearchTerm) {
                    matched = true
                    break
                }
            }
        case "location":
            locations := locationsMap[artist.ID]
            for _, loc := range locations {
                if strings.Contains(strings.ToLower(loc), lowerSearchTerm) {
                    matched = true
                    break
                }
            }
        case "creation date":
            if strings.Contains(strconv.Itoa(artist.CreationDate), searchTerm) {
                matched = true
            }
        case "first album date":
            if strings.Contains(strings.ToLower(artist.FirstAlbum), lowerSearchTerm) {
                matched = true
            }
        default:
            // Search all fields if no specific type is provided
            if strings.Contains(strings.ToLower(artist.Name), lowerSearchTerm) ||
                strings.Contains(strings.ToLower(artist.FirstAlbum), lowerSearchTerm) ||
                strings.Contains(strconv.Itoa(artist.CreationDate), searchTerm) {
                matched = true
            }
            if !matched {
                for _, member := range artist.Members {
                    if strings.Contains(strings.ToLower(member), lowerSearchTerm) {
                        matched = true
                        break
                    }
                }
            }
            if !matched {
                locations := locationsMap[artist.ID]
                for _, loc := range locations {
                    if strings.Contains(strings.ToLower(loc), lowerSearchTerm) {
                        matched = true
                        break
                    }
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

    // Remove duplicate artists if needed
    filteredArtists = removeDuplicateArtists(filteredArtists)

    // Aggregate unique locations from locationsMap for the dropdown
    locationSet := make(map[string]struct{})
    for _, locs := range locationsMap {
        for _, loc := range locs {
            locationSet[loc] = struct{}{}
        }
    }
    var allLocations []string
    for loc := range locationSet {
        allLocations = append(allLocations, loc)
    }
    // Optionally sort the locations
    // sort.Strings(allLocations)

    // Wrap the data into a PageData struct as expected by the template
    data := PageData{
        Artists:   filteredArtists,
        Locations: allLocations,
    }

    tmpl, err := template.ParseFiles("static/html/index.html")
    if err != nil {
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func parseQuery(query string) (string, string) {
    parts := strings.Split(query, " - ")
    if len(parts) == 2 {
        return strings.TrimSpace(parts[0]), strings.ToLower(strings.TrimSpace(parts[1]))
    }
    return strings.TrimSpace(query), ""
}

func removeDuplicateArtists(artists []API.Artist) []API.Artist {
    uniqueArtists := make(map[int]API.Artist)
    for _, artist := range artists {
        uniqueArtists[artist.ID] = artist
    }
    var result []API.Artist
    for _, artist := range uniqueArtists {
        result = append(result, artist)
    }
    return result
}
