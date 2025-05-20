package handlers

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"time"

	"forum/internal/auth"
	"forum/internal/db"
)

type User struct {
	Usernameprofil string
}
type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

type PostWithUser struct {
	Post_id      string
	Username     string
	Title1       string
	Content      string
	CreatedAt    string
	Commenters   []DataComment 
	CommnetString   string       // Show All Commnter For Every Post
	Status       string
	LikeDislike  string
	Colorlike    string
	ColorDislike string
	ColorValue   int
	Bool         int
	// How much Like_Dislike Evrey posts like ANd Dislike
	CountUserlike    int
	CountUserDislike int
	// Your Profil Form Header

	UserProfil string
	Headerhtml string
}

type DataComment struct {
	Contentcomment string
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	//var html PostWithUser
	//	var addUser PostWithUser

	user_id, _ := auth.CheckSession(w, r)
	/* if errS != nil {
		html.Headerhtml = ` 
		
		<li><a href="/login">Log In</a></li>
    	<li><a href="/register">Sign Up </a></li>`
	} else {
		html.Headerhtml = `
				<li><a href="/logout">Log Out</a></li>
				<li>
                  <button id="openModalBtn" type="button">Create Post</button>
                </li>`
	} */

	db.DB.QueryRow("SELECT username FROM users WHERE user_id = ? ", user_id).Scan(&user.Usernameprofil)
	// mytmpl, _ := template.ParseFiles("templates/home.html")
	// mytmpl.Execute(w, user)

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

type Alldata struct {
	Posts []PostWithUser
	Username string 
}


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
						    FOREIGN KEY 	db.DB.Exec("SELECT username FROM WHERE user_id = ? ", user_id).Scan(&)
			(user_id) REFERENCES users(user_id)
						);
		*/
			// Just check Database count table likedislike  where my ip and evry post

		db.DB.QueryRow("SELECT likedislike  FROM likedislike WHERE post_id = ? AND user_id = ?", p.Post_id, user_id).Scan(&p.LikeDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  likedislike = 'true' ", p.Post_id).Scan(&p.CountUserlike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  likedislike = 'false' ", p.Post_id).Scan(&p.CountUserDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  user_id = ? ", p.Post_id, user_id).Scan(&p.ColorValue)

		//////////////////////////////////////////////////
		// Print Like AND Dislike every Post-id
		if p.LikeDislike == "true" {
			p.Colorlike = "blue"
		}
		if p.LikeDislike == "false" {
			p.Colorlike = "red"
		}

		p.Commenters = comments
		fmt.Println(p.Commenters)


		posts = append(posts, p)
	}
	// addUser.UserProfil = user.Usernameprofil
	// posts = append(posts, addUser)
	//posts = append(posts, html)



	if err = rows.Err(); err != nil {
		log.Printf("Error iterating posts: %v", err)
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}


	
data := Alldata{
	Posts: posts,
	Username: user.Usernameprofil ,
	
} 


fmt.Println("Username : ",data.Username)

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
