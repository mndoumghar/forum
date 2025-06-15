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
	Data := models.Data{
		ErrorColor: []models.ErrorRegister{
			{Error: "Email already taken", Color: "red"},
			{Error: "Username must be at least 5 characters", Color: "red"},
			{Error: "Password must be at least 6 characters", Color: "red"},
			{Error: "Email must be at least 6 characters", Color: "red"},
			{Error: "Registration successful", Color: "green"},
			{Error: "", Color: ""},
		},
	}
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", err)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed, Please use the correct HTTP method.", nil)
		return
	}
	// Definier Variable

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

		if len(email) < 6 || len(username) < 5 || len(password) < 6 {
		tmpl, _ := template.ParseFiles("templates/register.html")
		if len(email) < 6 { tmpl.Execute(w, Data.ErrorColor[3])} 
		if len(username) < 5 { tmpl.Execute(w, Data.ErrorColor[1]) }
		if len(password) < 6 { tmpl.Execute(w, Data.ErrorColor[2]) }
								   
		return
	}

	// Check if email exists
	if _, err := db.GetUserByEmailUsername(email); err == nil {
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
		ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", err)
		return
	}
	// insert All my Information From Data
	_, err = db.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)",
		email, username, hashedPw)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Registration failed, Please try again later.", err)
		return
	}

	tmpl, _ := template.ParseFiles("templates/register.html")

	tmpl.Execute(w, Data.ErrorColor[1])
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
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", err)
			return
		}

		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed, "Method not allowed, Please use the correct HTTP method.", nil)
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

	user, err := db.GetUserByEmailUsername(email)
	if err != nil {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", err)
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "Internal server error, Please try again later.", err)
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}
	// Creat  Session And Session Starting ..
	// THis is Session To stock a Value  --uuid
	err = auth.CreateSession(w, user.ID)
	if err != nil {
		/// Erro If Not Data Session Noty Working
		fmt.Println("Error Session Is Not Staritng")
		ErrorHandler(w, http.StatusInternalServerError, "Session error, Please try again later.", err)
		return
	}

	// Header Page "Home.html"
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

// handlink passwrd formats
