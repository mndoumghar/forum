package handlers

import (
	"net/http"

	"forum/internal/auth"
	"forum/internal/db"
)

type LikeDislikeCount struct {
	Count     int
	CountAll  int
	Reaction  string
}

func LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed, Please use the correct HTTP method.", "")
		return
	}

	// Get form values
	likedislike := r.FormValue("likedislike")
	post_id := r.FormValue("post_id")

	// Check current reaction if any
	currentReaction, err := db.GetUserReaction(user_id, post_id)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Failed to get like/dislike info", "")
		return
	}

	// Handle the reaction
	if currentReaction != "" {
		if currentReaction == likedislike {
			// User clicked same button - remove reaction
			err = db.DeleteUserReaction(user_id, post_id)
		} else {
			// User clicked opposite button - update reaction
			err = db.UpdateUserReaction(user_id, post_id, likedislike)
		}
	} else {
		// User has no reaction - insert new one
		err = db.InsertUserReaction(user_id, post_id, likedislike)
	}

	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Failed to update reaction", "")
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}