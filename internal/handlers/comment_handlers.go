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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {

		contentCommenter := r.FormValue("comment")

		// input Hidden Send post_id In page Home
		post_id := r.FormValue("post_id")

		_, err := db.DB.Exec("INSERT INTO comments(user_id, post_id, content) VALUES(?,?,?)", user_id, post_id, contentCommenter)
		if err != nil {

			http.Error(w, "Registration fvailed", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)

		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
