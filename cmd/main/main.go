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
	http.HandleFunc("/logout", handlers.LogoutHabndler)
	http.HandleFunc("/likedislike", handlers.LikeDislikeHandler)

	//	http.HandleFunc("/posts/", handlers.PostHandler) // e.g., /posts/1
	//	http.HandleFunc("/comment", handlers.CommentHandler)
	//	http.HandleFunc("/like", handlers.LikeHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Add a root handler to serve index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			//http.ServeFile(w, r, "static/index.html")
			// Page Home
			http.Redirect(w, r, "/posts", http.StatusSeeOther)

			return
		}

		http.NotFound(w, r)
	})

	fmt.Println("server started at :8080\nVisit http://localhost:8080 to access the forum.")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
