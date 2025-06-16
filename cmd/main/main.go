package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"forum/internal/db"
	"forum/internal/handlers"
)

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/posts", handlers.PostsHandler)
	http.HandleFunc("/creatpost", handlers.CreatePostHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)
	http.HandleFunc("/logout", handlers.LogoutHabndler)
	http.HandleFunc("/likedislike", handlers.LikeDislikeHandler)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {

		filePath := r.URL.Path[len("/static/"):] // Strip the /static/ prefix
		fullPath := filepath.Join("static", filepath.Clean(filePath))

		// Ensure it doesn't point to a directory or nested malformed file path
		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			handlers.ErrorHandler(w, http.StatusNotFound, "path not exist", "")
			return
		}

		// Serve the file safely
		http.ServeFile(w, r, fullPath)
		
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {

			
			http.Redirect(w, r, "/posts", http.StatusSeeOther)
			return
		}
		handlers.ErrorHandler(w, http.StatusNotFound, "path not exist", "")
	})

	fmt.Println("server started at :8080\nVisit http://localhost:8080 to access the forum.")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
