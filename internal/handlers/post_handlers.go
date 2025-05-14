package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"forum/internal/db"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

type PostWithUser struct {
	Post_id    string
	Username   string
	Content    string
	CreatedAt  string
	Commenters []DataComment
}

//	type Comments struct {
//		Comment_id int
//		Content    string
//	}
type DataComment struct {
	Contentcomment string
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch posts from the database
	rows, err := db.DB.Query(`
		SELECT 
			p.post_id, 
			u.username, 
			p.content, 
			p.created_at 
		FROM 
			posts p
		JOIN 
			users u ON p.user_id = u.user_id
	`)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Error fetching post data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostWithUser

	// Loop through posts
	for rows.Next() {
		var p PostWithUser
		err = rows.Scan(&p.Post_id, &p.Username, &p.Content, &p.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue // Skip rows with errors
		}
		fmt.Println(p.Post_id,"&&&&")

		// Fetch comments for this post
		rows2, err := db.DB.Query(`SELECT content FROM comments WHERE post_id = ?`, p.Post_id)
		if err != nil {
			log.Printf("Error querying comments: %v", err)
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}
		defer rows2.Close()

		var comments []DataComment
		for rows2.Next() {
			var c DataComment
			err = rows2.Scan(&c.Contentcomment)
			if err != nil {
				log.Printf("Error scanning comment: %v", err)
				continue
			}
			comments = append(comments, c)
		}
		if err = rows2.Err(); err != nil {
			log.Printf("Error iterating comments: %v", err)
		}

		// Assign comments to the post
		p.Commenters = comments
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating posts: %v", err)
	}

	// Render the template
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
}
