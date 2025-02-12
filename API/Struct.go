// Struct.go
package API

type LocationData struct {
    Locations []string `json:"locations"`
}

type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
    Locations    string   `json:"locations"`
    ConcertDates string   `json:"concertDates"`
    Relations    string   `json:"relations"`
}

type LocationEntry struct {
    ID        int      `json:"id"`
    Locations []string `json:"locations"`
    Dates     string   `json:"dates"`
}

type AllLocationsData struct {
    Index []LocationEntry `json:"index"`
}

// Existing structs remain unchanged
type DatesData struct {
    Dates []string `json:"dates"`
}

type RelationData struct {
    DatesLocation map[string][]string `json:"datesLocations"`
}

type Details struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
}

type AllDetails struct {
    Details   Details
    Dates     DatesData
    Relations RelationData
    Location  LocationData
}