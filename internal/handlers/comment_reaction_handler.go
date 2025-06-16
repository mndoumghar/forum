package handlers


import (
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/auth"
	"forum/internal/db"
)

func CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	// Get current user ID from session
	userID, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, "Invalid form data", "")
		return
	}

	// Get comment ID from form
	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, "Invalid comment ID", "")
		return
	}

	// Get reaction type (like/dislike) from form
	isLikeStr := r.FormValue("is_like")
	isLike, err := strconv.ParseBool(isLikeStr)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, "Invalid reaction value", "")
		return
	}

	// Check if user already reacted to this comment
	var existingReaction bool
	err = db.DB.QueryRow(`
		SELECT is_like FROM comment_reactions 
		WHERE comment_id = ? AND user_id = ?`,
		commentID, userID).Scan(&existingReaction)

	if err == nil {
		// User already reacted - check if they're changing their reaction
		if existingReaction == isLike {
			// Same reaction - remove it
			_, err = db.DB.Exec(`
				DELETE FROM comment_reactions 
				WHERE comment_id = ? AND user_id = ?`,
				commentID, userID)
		} else {
			// Different reaction - update it
			_, err = db.DB.Exec(`
				UPDATE comment_reactions 
				SET is_like = ? 
				WHERE comment_id = ? AND user_id = ?`,
				isLike, commentID, userID)
		}
	} else {
		// No existing reaction - insert new one
		_, err = db.DB.Exec(`
			INSERT INTO comment_reactions (comment_id, user_id, is_like) 
			VALUES (?, ?, ?)`,
			commentID, userID, isLike)
	}

	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Failed to process reaction", "")
		return
	}

	// Get post ID for this comment to redirect back to the post
	var postID int
	err = db.DB.QueryRow(`
		SELECT post_id FROM comments 
		WHERE comment_id = ?`,
		commentID).Scan(&postID)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Failed to find post", "")
		return
	}

	// Redirect back to the post page
	http.Redirect(w, r, fmt.Sprintf("/posts#post-%d", postID), http.StatusSeeOther)
}