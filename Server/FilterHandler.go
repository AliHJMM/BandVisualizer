// FilterHandler.go
package server

import (
	"BandVisualizer/API"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// FilterHandler applies filters based on query parameters and renders the filtered artist list.
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Allow only GET requests.
	if r.Method != "GET" {
		renderErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	// ----- Creation Date Range Filter -----
	// Expecting parameters "creation_min" and/or "creation_max"
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
	// Expecting parameters "album_min" and/or "album_max"
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

	// ----- Number of Members Filter (Checkbox Filter) -----
	// Users may select one or several member counts. For example, ?members=2&members=3
	memberFilterStrs := q["members"]
	var memberFilters []int
	for _, m := range memberFilterStrs {
		num, err := strconv.Atoi(m)
		if err == nil {
			memberFilters = append(memberFilters, num)
		} else {
			log.Println("Error converting members value:", m, err)
		}
	}

	// ----- Locations Filter (Checkbox Filter) -----
	// Users may select one or more location keywords.
	locationsFilters := q["locations"]
	// Normalize the location filter values to lower case.
	for i, loc := range locationsFilters {
		locationsFilters[i] = strings.ToLower(strings.TrimSpace(loc))
	}

	// ----- Fetch Data -----
	artists, err := API.FetchArtists()
	if err != nil {
		log.Println("Error fetching artists:", err)
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	locationsMap, err := API.FetchAllLocations()
	if err != nil {
		log.Println("Error fetching locations:", err)
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare a regex to find a 4-digit year.
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
		// Only apply if at least one album filter parameter is provided.
		if albumMinStr != "" || albumMaxStr != "" {
			yearMatch := reYear.FindString(artist.FirstAlbum)
			if yearMatch != "" {
				albumYear, err := strconv.Atoi(yearMatch)
				if err != nil {
					log.Printf("Error converting extracted album year for %s: %v", artist.Name, err)
					// If conversion fails, skip filtering by album for this artist.
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
				// No 4-digit year was found in the firstAlbum field.
				log.Printf("No 4-digit year found in album field for %s, ignoring album filter", artist.Name)
			}
		}

		// (c) Number of Members Filter.
		if len(memberFilters) > 0 {
			matched := false
			for _, count := range memberFilters {
				if len(artist.Members) == count {
					matched = true
					break
				}
			}
			if !matched {
				log.Printf("Skipping %s due to members count (%d not in %v)", artist.Name, len(artist.Members), memberFilters)
				continue
			}
		}

		// (d) Locations Filter.
		if len(locationsFilters) > 0 {
			locs, ok := locationsMap[artist.ID]
			if !ok {
				log.Printf("Skipping %s because no locations found", artist.Name)
				continue
			}
			locationMatched := false
			for _, filterLoc := range locationsFilters {
				for _, loc := range locs {
					if strings.Contains(strings.ToLower(loc), filterLoc) {
						locationMatched = true
						break
					}
				}
				if locationMatched {
					break
				}
			}
			if !locationMatched {
				log.Printf("Skipping %s because none of the locations (%v) matched filter: %v", artist.Name, locs, locationsFilters)
				continue
			}
		}

		// Artist passed all filters.
		log.Println("Including artist:", artist.Name)
		filteredArtists = append(filteredArtists, artist)
	}

	log.Printf("Total filtered artists: %d", len(filteredArtists))
	if len(filteredArtists) == 0 {
		renderErrorPage(w, "No results found for the applied filters.", http.StatusNotFound)
		return
	}

	// Render the results using the index template (artist cards).
	tmpl, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		log.Println("Template parsing error:", err)
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, filteredArtists); err != nil {
		log.Println("Template execution error:", err)
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
