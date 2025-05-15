package handlers

import (
	"forum/internal/auth"
	"net/http"
)

func LogoutHabndler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := auth.DeletCoockies(w, r)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
