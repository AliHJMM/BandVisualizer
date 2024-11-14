package server

import (
    "BandVisualizer/API"
    "sort"
    "strconv"
    "strings"
)

// FilterOptions holds the possible values for each filter category.
type FilterOptions struct {
    MinCreationDate   int
    MaxCreationDate   int
    MinFirstAlbumYear int
    MaxFirstAlbumYear int
    NumberOfMembers   []int
    Locations         []string
}

// getFilterOptions extracts the possible filter values from the artist data.
func getFilterOptions(artists []API.Artist, locationsMap map[int][]string) FilterOptions {
    minCreationDate := int(^uint(0) >> 1) // Max int
    maxCreationDate := 0
    minFirstAlbumYear := int(^uint(0) >> 1)
    maxFirstAlbumYear := 0
    numberOfMembersSet := make(map[int]struct{})
    locationsSet := make(map[string]struct{})

    for _, artist := range artists {
        // Creation date
        if artist.CreationDate < minCreationDate {
            minCreationDate = artist.CreationDate
        }
        if artist.CreationDate > maxCreationDate {
            maxCreationDate = artist.CreationDate
        }

        // First album year
        if len(artist.FirstAlbum) >= 4 {
            firstAlbumYearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
            firstAlbumYear, err := strconv.Atoi(firstAlbumYearStr)
            if err == nil {
                if firstAlbumYear < minFirstAlbumYear {
                    minFirstAlbumYear = firstAlbumYear
                }
                if firstAlbumYear > maxFirstAlbumYear {
                    maxFirstAlbumYear = firstAlbumYear
                }
            }
        }

        // Number of members
        numMembers := len(artist.Members)
        numberOfMembersSet[numMembers] = struct{}{}

        // Locations
        locations, ok := locationsMap[artist.ID]
        if ok {
            for _, loc := range locations {
                // Extract country/state/city
                parts := strings.Split(loc, "-")
                for i := range parts {
                    locationPart := strings.Join(parts[i:], "-")
                    locationsSet[locationPart] = struct{}{}
                }
            }
        }
    }

    // Convert numberOfMembersSet to slice
    var numberOfMembers []int
    for num := range numberOfMembersSet {
        numberOfMembers = append(numberOfMembers, num)
    }
    sort.Ints(numberOfMembers)

    // Convert locationsSet to slice
    var locations []string
    for loc := range locationsSet {
        locations = append(locations, strings.ReplaceAll(loc, "-", ", "))
    }
    sort.Strings(locations)

    return FilterOptions{
        MinCreationDate:   minCreationDate,
        MaxCreationDate:   maxCreationDate,
        MinFirstAlbumYear: minFirstAlbumYear,
        MaxFirstAlbumYear: maxFirstAlbumYear,
        NumberOfMembers:   numberOfMembers,
        Locations:         locations,
    }
}

// removeDuplicateArtists removes duplicate artists based on their ID.
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
