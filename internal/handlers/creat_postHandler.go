package handlers

import (
	// "fmt"
	"net/http"
	"strings"
	"text/template"

	"forum/internal/auth"
	"forum/internal/db"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/home.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", "Please try again later.", err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", "Please try again later.", err)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, "Bad Request", "Please check your form data and try again.", err)
			return
		}

		status := r.Form["status"]
		statusStr := strings.Join(status, " ")

		title := r.FormValue("title")
		content := r.FormValue("content")

		_, err = db.DB.Exec(
			"INSERT INTO posts (user_id, title, content, status) VALUES (?, ?, ?, ?)",
			user_id, title, content, statusStr,
		)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", "Please try again later.", err)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Please use the correct HTTP method.", nil)
}
