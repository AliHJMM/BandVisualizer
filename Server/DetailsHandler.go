// server/DetailsHandler.go
package server

import (
	"BandVisualizer/API"
	"encoding/json"
	"html/template"
	"net/http"
)

// DetailsTemplateData embeds the AllDetails data and adds pre‑marshaled JSON strings.
type DetailsTemplateData struct {
	API.AllDetails
	LocJSON template.JS // JSON representation of Location.Locations
	RelJSON template.JS // JSON representation of Relations.DatesLocation
}

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

	// Marshal location and relation data into JSON strings.
	locBytes, err := json.Marshal(data.Location.Locations)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	relBytes, err := json.Marshal(data.Relations.DatesLocation)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the template data.
	tplData := DetailsTemplateData{
		AllDetails: data,
		LocJSON:    template.JS(locBytes), // safe to output as JS
		RelJSON:    template.JS(relBytes),
	}

	tmpl, err := template.ParseFiles("static/html/details.html")
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, tplData)
	if err != nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
