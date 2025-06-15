package handlers

import (
	"net/http"

	"forum/internal/auth"
)

func LogoutHabndler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method Not Allowed, Please use the correct HTTP method.", nil)
		return
	}
	err := auth.DeletCoockies(w, r)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Logout failed, Please try again later.", err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
