# Home Page Development

## Run Program Usin this command 
-   `go run cmd/main/main.go`

## 🔧 Features Implemented T

- ✅ **Created Home Page**
- 🔗 **Performed a Join Between `posts` and `users` Tables**  
  - Selected and retrieved combined data
- 🔁 **Displayed All Posts**  
  - Iterated through the joined results on the Home Page
- 🎨 **Designed a Simple UI**  
  - Basic layout using HTML and CSS

---

### 📄 Description Page Home

This page fetches and displays all posts by combining data from both the `posts` and `users` tables using a SQL join. The front-end is built with a minimalist design using basic HTML and CSS to keep it lightweight and responsive 

  ## READ TODO list !!
 
















type DataComment struct {
	Contentcomment string
	Usercommnter   string
	TimeCommnter   time.Time
	TimePost       int
	TmieType       string
}


rows2, err := db.DB.Query(`SELECT c.content,s.username , c.created_at FROM comments c JOIN users s ON c.user_id = s.user_id WHERE post_id = ?`, p.Post_id)
		if err != nil {
			log.Printf("Error querying comments: %v", err)
			ErrorHandler(w, http.StatusInternalServerError, "Error fetching comments, Please try again later.", err)
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
				fmt.Println("Sent just now")
			}

			if err != nil {
				log.Printf("Error scanning comment: %v", err)
				continue
			}
			comments = append(comments, c)
		}