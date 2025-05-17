package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"forum/internal/auth"
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
	Post_id  string
	Username string
	Title1    string
	Content     string
	CreatedAt   string
	Commenters  []DataComment // Show All Commnter For Every Post
	Status      string
	LikeDislike string
	Colorlike string
	ColorDislike string
	ColorValue int
	Bool int
	// How much Like_Dislike Evrey posts like ANd Dislike
	CountUserlike    int
	CountUserDislike int
}

type DataComment struct {
	Contentcomment string
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	user_id, _ := auth.CheckSession(w, r)

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
			p.status, 
			p.created_at 
		FROM 
			posts p
		JOIN 
			users u ON p.user_id = u.user_id
	`)
	if err != nil {
		log.Printf("Error querying databa  se: %v", err)
		http.Error(w, "Error fetching post data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostWithUser

	for rows.Next() {
		var p PostWithUser
		err = rows.Scan(&p.Post_id, &p.Username, &p.Content, &p.Status, &p.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		rows2, err := db.DB.Query(`SELECT content FROM comments WHERE post_id = ?`, p.Post_id)
		if err != nil {
			log.Printf("Error Content: erying comments: %v", err)
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}
		defer rows2.Close()
		var comments []DataComment
		for rows2.Next() {
			var c DataComment
			err = rows2.Scan(&c.Contentcomment)
			if err != nil {
				fmt.Println(comments)
				log.Printf("Error scanning comment: %v", err)
				continue
			}
			comments = append(comments, c)
		}

		if err = rows2.Err(); err != nil {
			log.Printf("Error iterating comments: %v", err)
		}

		/*
				CREATE TABLE IF NOT EXISTS likedislike (
			    likedislike_id INTEGER PRIMARY KEY AUTOINCREMENT,
			    post_id INTEGER NOT NULL,
			    user_id INTEGER NOT NULL,
			      TEXT CHECK (likedislike IN ('true', 'false')),
			    FOREIGN KEY (post_id) REFERENCES posts(post_id),
			    FOREIGN KEY (user_id) REFERENCES users(user_id)
			);
		*/

		db.DB.QueryRow("SELECT likedislike  FROM likedislike WHERE post_id = ? AND user_id = ?", p.Post_id, user_id).Scan(&p.LikeDislike)
		//fmt.Println(p.Title1)

		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  likedislike = 'true' ", p.Post_id).Scan(&p.CountUserlike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  likedislike = 'false' ", p.Post_id).Scan(&p.CountUserDislike)
				db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  user_id = ? ",p.Post_id,user_id).Scan(&p.ColorValue)


		//////////////////////////////////////////////////
		// Print Like AND Dislike every Post-id

		fmt.Println(" like :  ", p.CountUserlike, "  ____  Dislike : ", p.CountUserDislike, " Conter : ", p.LikeDislike)
		if p.LikeDislike == "true" {
			p.Colorlike = "blue"
		} 
		if p.LikeDislike == "false" {
			p.Colorlike = "red"
		} 
		

		p.Commenters = comments
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating posts: %v", err)
	}

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
