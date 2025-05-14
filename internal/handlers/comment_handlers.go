package handlers

import (
	"fmt"
	"net/http"

	"forum/internal/db"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {

		contentCommenter := r.FormValue("comment")
		fmt.Println(contentCommenter)
		// test user id and post id 
		user_id:= 1
		post_id:= 1


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
