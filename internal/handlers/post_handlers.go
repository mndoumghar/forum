package handlers

import (
	"fmt"
	"forum/internal/db"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

type PostWithUser struct {
	Username string
	Content  string
	CreatedAt string

}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// استرجاع البيانات من قاعدة البيانات
	rows, err := db.DB.Query(`SELECT u.username, p.content,  p.created_at FROM posts p JOIN users u ON p.user_id = u.user_id`)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Error fetching post data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostWithUser
	for rows.Next() {
		var p PostWithUser
		err = rows.Scan(&p.Username, &p.Content,&p.CreatedAt)
		//CreatedAt  created_at
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue // if i meet some line Erro its skip Error
		}
		posts = append(posts, p) // show allPost 
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
	}

	// تشيك الـ template وتنفذو
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println(posts) // للتحقق من البيانات فالـ console
}