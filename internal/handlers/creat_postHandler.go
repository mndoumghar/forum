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
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error, Please try again later.", "")
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error, Please try again later.", "")
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, "Bad Request, Please check your form data and try again.", "")
			return
		}

		status := r.Form["status"]

		validCategories := map[string]bool{
			"sport":  true,
			"jobs":   true,
			"news":   true,
			"movies": true,
			"tech":   true,
		}

		seen := make(map[string]bool)

		for _, selectedCategory := range status {
			if !validCategories[selectedCategory] {

				ErrorHandler(w, http.StatusBadRequest, "Bad Request, Please check your form data and try again.", "")
				return
			}
			if seen[selectedCategory] {
				ErrorHandler(w, http.StatusBadRequest, "Bad Request, duplicate categories are not allowed.", "")
				return
			}
			
			seen[selectedCategory] = true
		}

		statusStr := strings.Join(status, " ")

		title := r.FormValue("title")
		content := r.FormValue("content")
		if title == "" || content == "" || len(status) == 0 {
			ErrorHandler(w, http.StatusBadRequest, "Bad Request, Please check your form data and try again.", "")
			return

		}

		// Insert post and get post_id
		_, err = db.DB.Exec(
			"INSERT INTO posts (user_id, title, content, status) VALUES (?, ?, ?, ?)",
			user_id, title, content, statusStr,
		)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error, Please try again later.", "")
			return
		}

		

		// Insert each category into the category table
		

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed, Please use the correct HTTP method.", "")
}
