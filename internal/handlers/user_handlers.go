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
	_, err := auth.CheckSession(w, r)


	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Declar Struct Type Error From CSS
	Data := models.Data{
		ErrorColor: []models.ErrorRegister{
			{Error: "Password or email not correct", Color: "red"},
			{Error: "Registration successful ", Color: "green"},
			{Error: "", Color: ""},
		},
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
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

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := db.GetUserByEmail(email)
	if err != nil {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", "")
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", "")
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		

		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}
	// Creat  Session And Session Starting ..
	// THis is Session To stock a Value  --uuid
	db.DB.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
	err = auth.CreateSession(w, user.ID)

	if err != nil {
		/// Erro If Not Data Session Noty Working
		fmt.Println("Error Session Is Not Staritng")
		ErrorHandler(w, http.StatusInternalServerError, "Session error, Please try again later.", "")
		return
	}

	// Header Page "Home.html"
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

// handlink passwrd formats
