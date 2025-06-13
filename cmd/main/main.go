package main

import (
	"fmt"
	"forum/internal/db"
	"forum/internal/handlers"
	"forum/internal/models"
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
	//http.HandleFunc("/filter", handlers.FilterByCategoryHandler)
	//http.HandleFunc("/my-posts", handlers.MyPostsHandler)
	//http.HandleFunc("/my-likes", handlers.MyLikedPostsHandler)

	
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

	categories := []string{
		"Career Advice & Development",
		"Job Opportunities & Networking",
		"Cybersecurity",
		"Networking & Infrastructure",
		"Project Management",
		"Industry News & Updates",
		"Technical Discussions",
		"Mathematics & Data Science",
		"Soft Skills & Communication",
		"Leadership & Management",
	}

	dbConn, err := db.GetDBConnection() // Assuming GetDBConnection() returns the database connection
	if err != nil {
		fmt.Printf("Failed to get database connection: %v\n", err)
		return
	}

	for _, cat := range categories {
		err := models.AddCategory(dbConn, 1, 1, cat, cat) // post_id=1, user_id=1, status=cat, content=cat
		if err != nil {
			fmt.Println("Error adding category:", cat, err)
		}
	}

	fmt.Println("server started at :8080\nVisit http://localhost:8080 to access the forum.")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
