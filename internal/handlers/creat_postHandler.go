package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"forum/internal/db"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET request: Render the create post form
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/creat_post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle POST request: Process form submission
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// if content == "" {
		// 	http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		// 	return
		// }

		// Temporary user ID (replace with actual session-based user ID)
		status := r.FormValue("status")
		content := r.FormValue("content")

		fmt.Println(status)

		// Insert post into database
		_, err = db.DB.Exec("INSERT INTO posts(user_id, content,status) VALUES(?, ?, ?)",1, content, status)
		if err != nil {
			http.Error(w, "Failed to create postss", http.StatusInternalServerError)
			return
		}

		// Redirect to posts page
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	// Handle other methods
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
