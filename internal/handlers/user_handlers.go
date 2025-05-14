package handlers

import (
	"net/http"
	"text/template"

	"forum/internal/db"
	"forum/internal/models"
	"forum/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	Data := models.Data{
		ErrorColor: []models.ErrorRegister{
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
	// Definier Variable
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
	// insert All my Information From Data
	_, err = db.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)",
		email, username, hashedPw)
	if err != nil {

		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	tmpl, _ := template.ParseFiles("templates/register.html")

	tmpl.Execute(w, Data.ErrorColor[1])
}

// Login Page if Exist Your Information

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Declar Struct Type Error From CSS   
	Data := models.Data{
		ErrorColor: []models.ErrorRegister{
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

					//  Tach---> oussama\\
	//   can u Check all Errors possible About email And Password
	//  If these Variables WAS empty  then Add in Struct error for Exmple 	{Error: "Registration successful 🚀✨💪🏆 ", Color: "green"}, 
	// *** "Email Or Password is Empty please Insert Your Information"
	// You Can Add THe most Error obligatiore like you should to Add all Input Maximum superier >6 charcater ...
	//  Add Another baakcgriund like this  BUt i need  like as Project Forum Or media onginral  and you u most to respect color background Shadow in  some Black image

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

	// Header Page "Home.html"
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
