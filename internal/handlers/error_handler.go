package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, msg string, indecator string) {
	w.WriteHeader(statusCode)

	tmpl, tmplErr := template.ParseFiles("templates/error.html")
	if tmplErr != nil {
		log.Println("Template parsing error:", tmplErr)
		w.Write([]byte("Internal Server Error"))
		return
	}
	data := struct {
		Error_message string
		Status        string
		statusCode    int
		indec         string
	}{
		Error_message: msg,
		Status:        http.StatusText(statusCode),
		statusCode:    statusCode,
		indec:         indecator,
	}
	tmpl.Execute(w, data)
}
