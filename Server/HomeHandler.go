package server

import (
	"BandVisualizer/API"
	"net/http"
	"text/template"
)

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
	tmpl, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, artists)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
