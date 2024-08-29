package server

import (
	"BandVisualizer/API"
	"strconv"
)

func CheckId(id string) bool {
	value, err := strconv.Atoi(id)
	if err != nil {
		return false
	}
	len, err := API.GetArtistsLength()
	if err != nil {
		return false
	}
	if !(value > 0 && value <= len) {
		return false
	}
	return true
}
