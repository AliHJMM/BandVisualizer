// server/GeocodeHandler.go
package server

import (
	"BandVisualizer/API"
	"net/http"
)

// GeocodeHandler wraps the API.GetGeocode function.
func GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	API.GetGeocode(w, r)
}
