package handlers

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strings"
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

type DataComment struct {
	Contentcomment string
}

type PostWithUser struct {
	Post_id         string
	Username        string
	Title1          string
	Content         string
	CreatedAt       string
	Commenters      []DataComment 
	CommnetString   string       // Show All Commnter For Every Post
	Status          string       // Comma-separated categories (e.g. "Technical,Soft Skills")
	LeftCategories  []string     // Left-side categories
	RightCategories []string     // Right-side categories
	LikeDislike     string
	Colorlike       string
	ColorDislike    string
	ColorValue      int
	Bool            int
	CountUserlike    int
	CountUserDislike int
	UserProfil      string
	Headerhtml      string
}

type Alldata struct {
	Posts    []PostWithUser
	Username string 
}

// Helper to split categories into two slices
func splitCategories(categories []string) (left, right []string) {
	mid := (len(categories) + 1) / 2 // left gets the extra if odd
	left = categories[:mid]
	right = categories[mid:]
	return
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	user_id, _ := auth.CheckSession(w, r)

	db.DB.QueryRow("SELECT username FROM users WHERE user_id = ? ", user_id).Scan(&user.Usernameprofil)

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
		log.Printf("Error querying database: %v", err)
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

		db.DB.QueryRow("SELECT likedislike  FROM likedislike WHERE post_id = ? AND user_id = ?", p.Post_id, user_id).Scan(&p.LikeDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  likedislike = 'true' ", p.Post_id).Scan(&p.CountUserlike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  likedislike = 'false' ", p.Post_id).Scan(&p.CountUserDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?  and  user_id = ? ", p.Post_id, user_id).Scan(&p.ColorValue)

		if p.LikeDislike == "true" {
			p.Colorlike = "blue"
		}
		if p.LikeDislike == "false" {
			p.Colorlike = "red"
		}

		p.Commenters = comments

		// --- CATEGORY SPLITTING LOGIC HERE ---
		// If p.Status is a comma-separated string of categories
		statusList := []string{}
		for _, s := range strings.Split(p.Status, ",") {
			cat := strings.TrimSpace(s)
			if cat != "" {
				statusList = append(statusList, cat)
			}
		}
		left, right := splitCategories(statusList)
		p.LeftCategories = left
		p.RightCategories = right

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

	data := Alldata{
		Posts:    posts,
		Username: user.Usernameprofil,
	}

	fmt.Println("Username : ", data.Username)

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}