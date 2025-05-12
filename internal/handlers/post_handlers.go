package handlers

import (
	//"fmt"
//	"log"
	"net/http"
	
	"time"

	//"forum/internal/db"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the template
	// //tmpl, err := template.ParseFiles("templates/home.html")
	// if err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// Declare variables to hold data
	

	// Query the database
	
}
