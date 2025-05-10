package main

import (
	"fmt"
	"forum/internal/handlers"
	"net/http"
)

func main() {
	
	// Register handlers
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	//http.HandleFunc("/posts", handlers.PostsHandler)
	//http.HandleFunc("/comment", handlers.CommentHandler)
	//http.HandleFunc("/like", handlers.LikeHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Add a root handler to serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the path is exactly "/", serve index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		// For any other unmatched paths, return 404
		http.NotFound(w, r)
	})

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}