package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"forum/internal/auth"
	"forum/internal/db"
	"forum/internal/models"
	"forum/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery != "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	
	_, err := auth.CheckSession(w, r)

	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", "")
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed, Please use the correct HTTP method.", "")
		return
	}

	formErrors := &models.FormErrors{}
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	hasError := false

	if !utils.ValidateEmail(email) {
		formErrors.EmailError = "Invalid email format"
		hasError = true
	}

	if !utils.ValidateUsername(username) {
		formErrors.UsernameError = "Username must be between 4-20 characters and contain only letters, numbers, and underscores"
		hasError = true
	}

	if !utils.ValidatePassword(password) {
		formErrors.PasswordError = "Password must be between 4-20 characters"
		hasError = true
	}

	// Check if email/username already exists
	if _, err := db.GetUserByEmailUsername(email); err == nil {
		formErrors.EmailError = "Email already taken"
		hasError = true
	}

	if hasError {
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, formErrors)
		return
	}

	// Hash password
	hashedPw, err := utils.HashPassword(password)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", "")
		return
	}

	// Insert user into database
	_, err = db.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)",
		email, username, hashedPw)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Registration failed, Please try again later.", "")
		return
	}

	// Redirect to login page after successful registration
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login Page if Exist Your Information

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check for existing session
	_, err := auth.CheckSession(w, r)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Predefined error messages
	data := models.Data{
		ErrorColor: []models.ErrorRegister{
			{Error: "Password or email not correct", Color: "red"},
			{Error: "Registration successful", Color: "green"},
			{Error: "", Color: ""},
		},
	}

	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, please try again later.", "")
			return
		}
		tmpl.Execute(w, nil)
		return

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := db.GetUserByEmail(email)
		if err != nil {
			renderLoginWithError(w, data.ErrorColor[0])
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			renderLoginWithError(w, data.ErrorColor[0])
			return
		}

		// Clear old sessions and create a new one
		db.DB.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
		err = auth.CreateSession(w, user.ID)
		if err != nil {
			fmt.Println("Error: Session not starting")
			ErrorHandler(w, http.StatusInternalServerError, "Session error, please try again later.", "")
			return
		}

		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return

	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed. Please use the correct HTTP method.", "")
	}
}

// Helper function to render login template with error data
func renderLoginWithError(w http.ResponseWriter, errorData models.ErrorRegister) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Internal server error, please try again later.", "")
		return
	}
	tmpl.Execute(w, errorData)
}


// handlink passwrd formats
