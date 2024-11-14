package main

import (
	"fmt"
	"log"
	"net/http"
	"BandVisualizer/server"
)

func main() {
	http.HandleFunc("/", server.HomeHandler)
	http.HandleFunc("/details", server.DetailsHandler)
	// Add the search and suggestions handlers
	http.HandleFunc("/search", server.SearchHandler)
	http.HandleFunc("/suggestions", server.SuggestionsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server running at http://localhost:2004")
	log.Fatal(http.ListenAndServe(":2004", nil))
}
