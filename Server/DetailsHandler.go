package server

import (
	"BandVisualizer/API"
	"html/template"
	"net/http"
)

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		renderErrorPage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	isValid := CheckId(id)

	if !isValid {
		renderErrorPage(w, "Not Found", http.StatusNotFound)
		return
	}

	data, err := API.GetAllDetails(id)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("static/html/details.html")
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
