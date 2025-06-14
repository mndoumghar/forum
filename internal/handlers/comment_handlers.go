package handlers

import (
	"net/http"

	"forum/internal/auth"
	"forum/internal/db"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Failed to add comment, ", "Please try again later.", nil)
		return
	}
	if r.Method == http.MethodGet {

		contentCommenter := r.FormValue("comment")

		// input Hidden Send post_id In page Home
		post_id := r.FormValue("post_id")

		_, err := db.DB.Exec("INSERT INTO comments(user_id, post_id, content) VALUES(?,?,?)", user_id, post_id, contentCommenter)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Failed to add comment", "Please try again later.", err)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)

		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
