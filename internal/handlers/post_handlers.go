package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/internal/auth"
	"forum/internal/db"
	"forum/internal/models"
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
	Usercommnter   string
	TimeCommnter   time.Time
	TimePost       int
	TmieType       string
}

type PostWithUser struct {
	Post_id          string
	Username         string
	Title1           string
	Content          string
	CreatedAt        time.Time
	Commenters       []DataComment
	CommnetString    string
	Status           string
	LeftCategories   []string
	RightCategories  []string
	LikeDislike      string
	Colorlike        string
	ColorDislike     string
	ColorValue       int
	Bool             int
	CountUserlike    int
	CountUserDislike int
	UserProfil       string
	Headerhtml       string
}

type Alldata struct {
	Posts            []PostWithUser
	Username         string
	AllCategories    []string
	SelectedCategory string
	Cate             string
	Cate2            string
}

// Helper to split categories into two slices
func splitCategories(categories []string) (left, right []string) {
	mid := (len(categories) + 1) / 2
	left = categories[:mid]
	right = categories[mid:]
	return
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	user_id, _ := auth.CheckSession(w, r)
	db.DB.QueryRow("SELECT username FROM users WHERE user_id = ? ", user_id).Scan(&user.Usernameprofil)

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed, Please use the correct HTTP method.", "")
		return
	}

	// Category filtering logic
	selectedCategory := r.URL.Query().Get("category")
	allCategories, err := models.GetalldistCat(db.DB)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		allCategories = []string{}
	}

	var rows *sql.Rows
	if selectedCategory != "" {
		// Query posts that have the selected category in their status
		rows, err = db.DB.Query(`
			SELECT
				p.post_id,
				u.username,
				p.title,
				p.content,
				p.status,
				p.created_at
			FROM
				posts p
			JOIN
				users u ON p.user_id = u.user_id
			WHERE
				p.status LIKE ?`, "%"+selectedCategory+"%")
	} else {
		// Query all posts if no category is selected
		rows, err = db.DB.Query(`
			SELECT 
				p.post_id, 
				u.username, 
				p.title,
				p.content,
				p.status, 
				p.created_at 
			FROM 
				posts p
			JOIN 
				users u ON p.user_id = u.user_id
			ORDER BY p.created_at DESC`)
	}
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Error fetching post data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []PostWithUser
	for rows.Next() {
		var p PostWithUser
		err = rows.Scan(&p.Post_id, &p.Username, &p.Title1, &p.Content, &p.Status, &p.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Fetch comments
		rows2, err := db.DB.Query(`SELECT c.content, s.username, c.created_at FROM comments c JOIN users s ON c.user_id = s.user_id WHERE post_id = ? ORDER BY c.created_at DESC`, p.Post_id)
		if err != nil {
			log.Printf("Error querying comments: %v", err)
			ErrorHandler(w, http.StatusInternalServerError, "Error fetching comments, Please try again later.", "")
			return
		}
		var comments []DataComment
		for rows2.Next() {
			var c DataComment
			err = rows2.Scan(&c.Contentcomment, &c.Usercommnter, &c.TimeCommnter)
			duration := time.Since(c.TimeCommnter)
			// Format the duration nicely
			if duration.Hours() >= 24 {
				days := int(duration.Hours()) / 24
				c.TimePost = days
				c.TmieType = "day(s) ago"
			} else if duration.Hours() >= 1 {
				hours := int(duration.Hours())
				c.TimePost = hours
				c.TmieType = " hour(s) ago"
			} else if duration.Minutes() >= 1 {
				minutes := int(duration.Minutes())
				c.TimePost = minutes
				c.TmieType = " minute(s) ago"
			} else {
				c.TimePost = 0
				c.TmieType = "just now"
			}

			if err != nil {
				log.Printf("Error scanning comment: %v", err)
				continue
			}
			comments = append(comments, c)
		}
		if err = rows2.Err(); err != nil {
			log.Printf("Error iterating comments: %v", err)
		}
		rows2.Close()

		// Get like/dislike information
		db.DB.QueryRow("SELECT likedislike FROM likedislike WHERE post_id = ? AND user_id = ?", p.Post_id, user_id).Scan(&p.LikeDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ? and likedislike = 'true'", p.Post_id).Scan(&p.CountUserlike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ? and likedislike = 'false'", p.Post_id).Scan(&p.CountUserDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ? and user_id = ?", p.Post_id, user_id).Scan(&p.ColorValue)

		if p.LikeDislike == "true" {
			p.Colorlike = "blue"
		} else if p.LikeDislike == "false" {
			p.ColorDislike = "red"
		}

		p.Commenters = comments

		// Split status into categories if comma-separated
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

	// Deduplicate allCategories
	uniqueCategories := make(map[string]struct{})
	dedupedCategories := []string{}
	for _, cat := range allCategories {
		if _, exists := uniqueCategories[cat]; !exists {
			uniqueCategories[cat] = struct{}{}
			dedupedCategories = append(dedupedCategories, cat)
		}
	}
	allCategories = dedupedCategories

	data := Alldata{
		Posts:            posts,
		Username:         user.Usernameprofil,
		AllCategories:    allCategories,
		SelectedCategory: selectedCategory,
		Cate:             selectedCategory,
		Cate2:            "Category",
	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error, Please try again later.", "")
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error, Please try again later.", "")
		return
	}
}