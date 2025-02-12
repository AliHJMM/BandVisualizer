package server

import (
	"BandVisualizer/API"
	"html/template"
	"net/http"
)

// ErrorData represents the error message and status code passed to the error template.
type ErrorData struct {
    ErrorMessage string
    Code         string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        renderErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    if r.URL.Path != "/" {
        renderErrorPage(w, "Not Found", http.StatusNotFound)
        return
    }
    artists, err := API.FetchArtists()
    if err != nil {
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    locations, err := AggregateLocations()
    if err != nil {
        renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    data := PageData{
        Artists:   artists,
        Locations: locations,
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
