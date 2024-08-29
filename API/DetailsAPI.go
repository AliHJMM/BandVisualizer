package API

import (
	"encoding/json"
	"net/http"
)

func FetchLocation(id string) (LocationData, error) {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + id)
	if err != nil {
		return LocationData{}, err
	}
	defer res.Body.Close()

	var data LocationData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return LocationData{}, err
	}

	return data, nil
}

func FetchDates(id string) (DatesData, error) {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + id)
	if err != nil {
		return DatesData{}, err
	}
	defer res.Body.Close()

	var data DatesData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return DatesData{}, err
	}

	return data, nil
}

func FetchRelations(id string) (RelationData, error) {

	res, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err != nil {
		return RelationData{}, err
	}
	defer res.Body.Close()

	var data RelationData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return RelationData{}, err
	}

	return data, nil
}

func FetchDetails(id string) (Details, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		return Details{}, err
	}
	defer resp.Body.Close()

	var details Details
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return Details{}, err
	}
	return details, nil
}

func GetAllDetails(id string) (AllDetails, error) {

	details, err := FetchDetails(id)
	if err != nil {
		return AllDetails{}, err
	}

	dates, err := FetchDates(id)
	if err != nil {
		return AllDetails{}, err
	}

	relation, err := FetchRelations(id)
	if err != nil {
		return AllDetails{}, err
	}

	location, err := FetchLocation(id)
	if err != nil {
		return AllDetails{}, err
	}

	allDetails := AllDetails{
		Details:   details,
		Dates:     dates,
		Relations: relation,
		Location:  location,
	}
	return allDetails, nil
}
