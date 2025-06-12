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

    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Bad Request", http.StatusBadRequest)
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
            http.Error(w, "Failed to create post", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/posts", http.StatusSeeOther)
        return
    }

    http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}