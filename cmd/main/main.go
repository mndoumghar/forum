package main

import (
	"fmt"
	"forum/internal/db"
	"forum/internal/handlers"
	"net/http"
)

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	// Register handlers All Function from Dossier Hanlres --->----
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/posts", handlers.PostsHandler)
	http.HandleFunc("/creatpost", handlers.CreatePostHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)

	//	http.HandleFunc("/posts/", handlers.PostHandler) // e.g., /posts/1
	//	http.HandleFunc("/comment", handlers.CommentHandler)
	//	http.HandleFunc("/like", handlers.LikeHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Add a root handler to serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}

		http.NotFound(w, r)
	})

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
