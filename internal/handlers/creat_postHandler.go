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

for _, selectedCategory := range status {
    if !validCategories[selectedCategory] {
        ErrorHandler(w, http.StatusBadRequest, "Bad Request, Please check your form data and try again.", "")
        return
    }
}

		
		statusStr := strings.Join(status, " ")

		title := r.FormValue("title")
		content := r.FormValue("content")
		if title == "" || content == ""  || len(status) == 0  {
			ErrorHandler(w, http.StatusBadRequest, "Bad Request, Please check your form data and try again.", "")
			return

		}

		// Insert post and get post_id
		result, err := db.DB.Exec(
			"INSERT INTO posts (user_id, title, content, status) VALUES (?, ?, ?, ?)",
			user_id, title, content, statusStr,
		)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error, Please try again later.", "")
			return
		}

		postID, err := result.LastInsertId()
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Failed to get post ID.", "")
			return
		}

		// Insert each category into the category table
		for _, cat := range status {
			cat = strings.TrimSpace(cat)
			if cat == "" {
				continue
			}
			_, err := db.DB.Exec(
				"INSERT INTO category (post_id, user_id, status, content) VALUES (?, ?, ?, ?)",
				postID, user_id, strings.ToLower(cat), content,
			)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, "Failed to insert category.", "")
				return
			}
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed, Please use the correct HTTP method.", "")
}
