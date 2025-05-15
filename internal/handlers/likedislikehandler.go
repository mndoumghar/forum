package handlers

import (
	
	"net/http"

	"forum/internal/auth"
	"forum/internal/db"
)

/*
CREATE TABLE IF NOT EXISTS likedislike (

		likedislike_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		likedislike  BOOLEAN NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(post_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);
*/
func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {

		likedislike := r.FormValue("likedislike")

		// input Hidden Send post_id In page Home
		post_id := r.FormValue("post_id")

		_, err := db.DB.Exec("INSERT INTO likedislike(user_id, post_id, likedislike) VALUES(?,?,?)", user_id, post_id, likedislike)
		if err != nil {

			http.Error(w, "Like Or Dislike  failed", http.StatusInternalServerError)
			return
		}
		

		http.Redirect(w, r, "/posts", http.StatusSeeOther)

		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
