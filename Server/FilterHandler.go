package server

import (
	"BandVisualizer/API"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// PageData is used to pass both the artist slice and the locations slice to the template.
type PageData struct {
    Artists   []API.Artist
    Locations []string
}


// AggregateLocations collects all distinct locations from the locations map.
func AggregateLocations() ([]string, error) {
    locMap, err := API.FetchAllLocations()
    if err != nil {
        return nil, err
    }
    locationSet := make(map[string]struct{})
    for _, locs := range locMap {
        for _, loc := range locs {
            // You can trim or modify the location string if needed.
            locationSet[loc] = struct{}{}
        }
    }
    var locations []string
    for loc := range locationSet {
        locations = append(locations, loc)
    }
    sort.Strings(locations)
    return locations, nil
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        renderErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    q := r.URL.Query()

    // ----- Creation Date Range Filter -----
    creationMinStr := q.Get("creation_min")
    creationMaxStr := q.Get("creation_max")
    var creationMin, creationMax int
    var err error
    if creationMinStr != "" {
        creationMin, err = strconv.Atoi(creationMinStr)
        if err != nil {
            log.Println("Error converting creation_min:", err)
            renderErrorPage(w, "Invalid creation_min parameter", http.StatusBadRequest)
            return
        }
    }
    if creationMaxStr != "" {
        creationMax, err = strconv.Atoi(creationMaxStr)
        if err != nil {
            log.Println("Error converting creation_max:", err)
            renderErrorPage(w, "Invalid creation_max parameter", http.StatusBadRequest)
            return
        }
    }

    // ----- First Album Year Range Filter -----
    albumMinStr := q.Get("album_min")
    albumMaxStr := q.Get("album_max")
    var albumMin, albumMax int
    if albumMinStr != "" {
        albumMin, err = strconv.Atoi(albumMinStr)
        if err != nil {
            log.Println("Error converting album_min:", err)
            renderErrorPage(w, "Invalid album_min parameter", http.StatusBadRequest)
            return
        }
    }
    if albumMaxStr != "" {
        albumMax, err = strconv.Atoi(albumMaxStr)
        if err != nil {
            log.Println("Error converting album_max:", err)
            renderErrorPage(w, "Invalid album_max parameter", http.StatusBadRequest)
            return
        }
    }

    // ----- Members Filter (Checkboxes) -----
    memberFilterStrs := q["members"]

    // ----- Location Filter (Dropdown) -----
    locationFilter := strings.TrimSpace(q.Get("locations"))
    locationFilter = strings.ToLower(locationFilter)

    // ----- Fetch Data -----
    artists, err := API.FetchArtists()
    if err != nil {
        log.Println("Error fetching artists:", err)
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    locMap, err := API.FetchAllLocations()
    if err != nil {
        log.Println("Error fetching locations:", err)
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Prepare a regex to extract a 4-digit year from FirstAlbum.
    reYear := regexp.MustCompile(`\d{4}`)

    var filteredArtists []API.Artist
    for _, artist := range artists {
        log.Println("Checking artist:", artist.Name)

        // (a) Creation Date Filter.
       	if creationMinStr != "" && artist.CreationDate < creationMin {
            log.Printf("Skipping %s due to creation date (%d < %d)", artist.Name, artist.CreationDate, creationMin)
            continue
        }
        if creationMaxStr != "" && artist.CreationDate > creationMax {
            log.Printf("Skipping %s due to creation date (%d > %d)", artist.Name, artist.CreationDate, creationMax)
            continue
        }

        // (b) First Album Year Filter.
        if albumMinStr != "" || albumMaxStr != "" {
            yearMatch := reYear.FindString(artist.FirstAlbum)
            if yearMatch != "" {
                albumYear, err := strconv.Atoi(yearMatch)
                if err != nil {
                    log.Printf("Error converting extracted album year for %s: %v", artist.Name, err)
                } else {
                    if albumMinStr != "" && albumYear < albumMin {
                        log.Printf("Skipping %s due to album year (%d < %d)", artist.Name, albumYear, albumMin)
                        continue
                    }
                    if albumMaxStr != "" && albumYear > albumMax {
                        log.Printf("Skipping %s due to album year (%d > %d)", artist.Name, albumYear, albumMax)
                        continue
                    }
                }
            } else {
                log.Printf("No 4-digit year found in album field for %s, ignoring album filter", artist.Name)
            }
        }

        // (c) Members Filter.
       	if len(memberFilterStrs) > 0 {
            matched := false
            for _, m := range memberFilterStrs {
                num, err := strconv.Atoi(m)
                if err == nil && len(artist.Members) == num {
                    matched = true
                    break
                }
            }
            if !matched {
                log.Printf("Skipping %s due to members count (%d not in %v)", artist.Name, len(artist.Members), memberFilterStrs)
                continue
            }
        }

        // (d) Location Filter.
        if locationFilter != "" {
            locs, ok := locMap[artist.ID]
            if !ok {
                log.Printf("Skipping %s because no locations found", artist.Name)
                continue
            }
            locationMatched := false
            for _, loc := range locs {
                if strings.ToLower(loc) == locationFilter {
                    locationMatched = true
                    break
                }
            }
            if !locationMatched {
                log.Printf("Skipping %s because none of the locations (%v) exactly match filter: %s", artist.Name, locs, locationFilter)
                continue
            }
        }

        log.Println("Including artist:", artist.Name)
        filteredArtists = append(filteredArtists, artist)
    }

    // Aggregate all locations for the dropdown.
    locationSet := make(map[string]struct{})
    for _, locs := range locMap {
        for _, loc := range locs {
            locationSet[loc] = struct{}{}
        }
    }
    var allLocations []string
    for loc := range locationSet {
        allLocations = append(allLocations, loc)
    }
    sort.Strings(allLocations)

    if len(filteredArtists) == 0 {
        renderErrorPage(w, "No results found for the applied filters.", http.StatusNotFound)
        return
    }
  
    data := PageData{
        Artists:   filteredArtists,
        Locations: allLocations,
    }
  
    tmpl, err := template.ParseFiles("static/html/index.html")
    if err != nil {
        log.Println("Template parsing error:", err)
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w, data); err != nil {
        log.Println("Template execution error:", err)
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
