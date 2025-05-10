package handler

import (
	"fmt"
	"net/http"
	"forum/internal/models"
	"forum/internal/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if the email is already taken
		_, err := models.GetUserByEmail(email)
		if err == nil {
			http.Error(w, "Email is already taken", http.StatusBadRequest)
			return
		}

		// Create new user
		err = models.CreateUser(email, username, password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "static/register.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := models.GetUserByEmail(email)
		if err != nil || !utils.CheckPasswordHash(password, user.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set session (cookies)
		// TODO: Implement session management
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "static/login.html")
}
