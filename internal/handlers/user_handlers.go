package handlers

import (
	
	"net/http"
	"text/template"
	"time"

	"forum/internal/db"
	"forum/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}
type ErrorRegister struct {
	Error string
	Color string
}

type Data struct {
	ErrorColor []ErrorRegister
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	Data := Data{
		ErrorColor: []ErrorRegister{
			{Error: "Email already taken", Color: "red"},
			{Error: "Registration successful 🚀✨💪🏆 ", Color: "green"},
			{Error: "", Color: ""},
		},
	}
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if email exists
	if _, err := db.GetUserByEmail(email); err == nil {
		tmpl, _ := template.ParseFiles("templates/register.html")
		// error.Error =  "Email already taken"
		tmpl.Execute(w, Data.ErrorColor[0])
		return
		// http.Error(w, "Email already taken", http.StatusBadRequest)
		// return
	}
	//    transfer passwordd to Hash password
	hashedPw, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)",
		email, username, hashedPw)
	if err != nil {

		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("templates/register.html")

	tmpl.Execute(w, Data.ErrorColor[1])

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	Data := Data{
		ErrorColor: []ErrorRegister{
			{Error: "Password or email not correct", Color: "red"},
			{Error: "Registration successful 🚀✨💪🏆 ", Color: "green"},
			{Error: "", Color: ""},
		},
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := db.GetUserByEmail(email)

	if err != nil {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}

	http.Redirect(w, r, "/posts",http.StatusSeeOther)


	
}
