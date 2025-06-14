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

type CountLikeAll struct {
	Nums int
}

func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed", "Please use the correct HTTP method.", nil)
		return
	}
	if r.Method == http.MethodGet {

		likedislike := r.FormValue("likedislike")

		// input Hidden Send post_id In page Home
		post_id := r.FormValue("post_id")

		// Checxk FRom Databnase How much line
		u, err := db.GetLikeDisle(user_id, post_id)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Failed to get like/dislike info", "Please try again later.", err)
			return
		}

		// if likedislike!= "" {

		// 		err = db.UpdateLikeDislike(user_id,post_id,likedislike)
		// 		if err != nil {
		// 			return
		// 		}
		// 	}

		if u.Count > 0 {
			err = db.DeleteIdUserikeDislike(user_id, post_id)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, "Failed to update like/dislike", "Please try again later.", err)
				return
			}
		} else {

			_, err = db.DB.Exec("INSERT INTO likedislike(user_id, post_id, likedislike) VALUES(?,?,?)", user_id, post_id, likedislike)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, "Like Or Dislike failed", "Please try again later.", err)
				return
			}

		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)

		return
	}
}
