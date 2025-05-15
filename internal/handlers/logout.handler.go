package handlers

import "net/http"

func LogoutHabndler(w http.ResponseWriter, r *http.Request) {

if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}


}