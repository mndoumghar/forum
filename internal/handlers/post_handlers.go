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
	CommentID      int
	Contentcomment string
	Usercommnter   string
	TimeCommnter   time.Time
	TimePost       int
	TmieType       string
	LikeCount      int
	DislikeCount   int
	UserReaction   int // 1 for like, 0 for dislike, -1 for no reaction
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
	PostFilter       string
}

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

	selectedCategory := r.URL.Query().Get("category")
	postFilter := r.URL.Query().Get("post")

	allCategories, err := models.GetalldistCat(db.DB)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		allCategories = []string{}
	}

	var rows *sql.Rows
	baseQuery := `
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
	`
switch {
case selectedCategory != "" && postFilter != "":
    // Handle combined category and post filter
    switch postFilter {
    case "my":
        rows, err = db.DB.Query(baseQuery+`
            WHERE p.status LIKE ? 
            AND p.user_id = ?
            ORDER BY p.created_at DESC`, "%"+selectedCategory+"%", user_id)
    case "liked":
        rows, err = db.DB.Query(baseQuery+`
            JOIN likedislike ld ON p.post_id = ld.post_id
            WHERE p.status LIKE ? 
            AND ld.user_id = ?
            AND ld.likedislike = 'true'
            ORDER BY p.created_at DESC`, "%"+selectedCategory+"%", user_id)
    case "disliked":
        rows, err = db.DB.Query(baseQuery+`
            JOIN likedislike ld ON p.post_id = ld.post_id
            WHERE p.status LIKE ? 
            AND ld.user_id = ?
            AND ld.likedislike = 'false'
            ORDER BY p.created_at DESC`, "%"+selectedCategory+"%", user_id)
    }

case postFilter == "my":
    // Filter only by my posts
    rows, err = db.DB.Query(baseQuery+`
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC`, user_id)

case postFilter == "liked":
    // Filter only by liked posts
    rows, err = db.DB.Query(baseQuery+`
        JOIN likedislike ld ON p.post_id = ld.post_id
        WHERE ld.user_id = ?
        AND ld.likedislike = 'true'
        ORDER BY p.created_at DESC`, user_id)

case postFilter == "disliked":
    // Filter only by disliked posts
    rows, err = db.DB.Query(baseQuery+`
        JOIN likedislike ld ON p.post_id = ld.post_id
        WHERE ld.user_id = ?
        AND ld.likedislike = 'false'
        ORDER BY p.created_at DESC`, user_id)

case selectedCategory != "":
    // Filter only by category
    rows, err = db.DB.Query(baseQuery+`
        WHERE p.status LIKE ?
        ORDER BY p.created_at DESC`, "%"+selectedCategory+"%")

default:
    // No filters - get all posts
    rows, err = db.DB.Query(baseQuery + ` ORDER BY p.created_at DESC`)
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

		// Fetch comments with reactions
		rows2, err := db.DB.Query(`
			SELECT 
				c.comment_id, 
				c.content, 
				u.username, 
				c.created_at,
				COALESCE(SUM(CASE WHEN cr.is_like = 1 THEN 1 ELSE 0 END), 0) AS like_count,
				COALESCE(SUM(CASE WHEN cr.is_like = 0 THEN 1 ELSE 0 END), 0) AS dislike_count,
				COALESCE((SELECT cr2.is_like FROM comment_reactions cr2 
					WHERE cr2.comment_id = c.comment_id AND cr2.user_id = ?), -1) AS user_reaction
			FROM comments c
			JOIN users u ON c.user_id = u.user_id
			LEFT JOIN comment_reactions cr ON c.comment_id = cr.comment_id
			WHERE c.post_id = ?
			GROUP BY c.comment_id
			ORDER BY c.created_at DESC`, user_id, p.Post_id)
		
		if err != nil {
			log.Printf("Error querying comments: %v", err)
			continue
		}

		var comments []DataComment
		for rows2.Next() {
			var c DataComment
			var userReaction sql.NullInt64
			err = rows2.Scan(
				&c.CommentID,
				&c.Contentcomment,
				&c.Usercommnter,
				&c.TimeCommnter,
				&c.LikeCount,
				&c.DislikeCount,
				&userReaction,
			)
			if err != nil {
				log.Printf("Error scanning comment: %v", err)
				continue
			}

			if userReaction.Valid {
				c.UserReaction = int(userReaction.Int64)
			} else {
				c.UserReaction = -1
			}

			duration := time.Since(c.TimeCommnter)
			if duration.Hours() >= 24 {
				c.TimePost = int(duration.Hours()) / 24
				c.TmieType = "day(s) ago"
			} else if duration.Hours() >= 1 {
				c.TimePost = int(duration.Hours())
				c.TmieType = "hour(s) ago"
			} else if duration.Minutes() >= 1 {
				c.TimePost = int(duration.Minutes())
				c.TmieType = "minute(s) ago"
			} else {
				c.TimePost = 0
				c.TmieType = "just now"
			}

			comments = append(comments, c)
		}
		rows2.Close()

		// Get post reactions
		db.DB.QueryRow("SELECT likedislike FROM likedislike WHERE post_id = ? AND user_id = ?", p.Post_id, user_id).Scan(&p.LikeDislike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ? AND likedislike = 'true'", p.Post_id).Scan(&p.CountUserlike)
		db.DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ? AND likedislike = 'false'", p.Post_id).Scan(&p.CountUserDislike)

		if p.LikeDislike == "true" {
			p.Colorlike = "blue"
		} else if p.LikeDislike == "false" {
			p.ColorDislike = "red"
		}

		p.Commenters = comments

		// Split categories
		statusList := []string{}
		for _, s := range strings.Split(p.Status, ",") {
			cat := strings.TrimSpace(s)
			if cat != "" {
				statusList = append(statusList, cat)
			}
		}
		p.LeftCategories, p.RightCategories = splitCategories(statusList)
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating posts: %v", err)
	}

	// Deduplicate categories
	uniqueCategories := make(map[string]struct{})
	dedupedCategories := []string{}
	for _, cat := range allCategories {
		if _, exists := uniqueCategories[cat]; !exists {
			uniqueCategories[cat] = struct{}{}
			dedupedCategories = append(dedupedCategories, cat)
		}
	}

	data := Alldata{
		Posts:            posts,
		Username:         user.Usernameprofil,
		AllCategories:    dedupedCategories,
		SelectedCategory: selectedCategory,
		Cate:             selectedCategory,
		Cate2:            "Category",
		PostFilter:       postFilter,
	}

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", "")
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, http.StatusInternalServerError, "Internal Server Error", "")
		return
	}
}