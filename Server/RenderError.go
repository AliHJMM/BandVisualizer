package server

import (
	"html/template"
	"net/http"
	"strconv"
)

func renderErrorPage(w http.ResponseWriter, errMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	data := ErrorData{
		ErrorMessage: errMsg,
		Code:         strconv.Itoa(statusCode),
	}

	tmpl, err := template.ParseFiles("static/html/error.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
