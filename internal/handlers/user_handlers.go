package handlers

import (
	"fmt"
	"log"
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
	// error.Error = "Registration y____successful"

	tmpl.Execute(w, Data.ErrorColor[1])

	// w.WriteHeader(http.StatusCreated)

	// w.Write([]byte("Registration successful"))
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
		// http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		// return

		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		// return

		// data.Error = "Password or email not correct"
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, Data.ErrorColor[0])
		return
	}

	// Create session (to be implemented in auth/session.go)

	//          w.Write([]byte("Login successful"))

	var u User

	// db.DB.QueryRow("SELECT user_id, email, username, password, created_at FROM users").
	// 	Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)

	tmpl, _ := template.ParseFiles("templates/home.html")
	// tmpl.Execute(w, u)

	rows, err := db.DB.Query(`SELECT u.username, p.content FROM posts p JOIN users u ON p.user_id = u.user_id`)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error querying database: %v", err)
		// Send an error response
		http.Error(w, "Error fetching post data", http.StatusInternalServerError)
		return
	}

	var posts []Post
	var users []User
	var p Post
	for rows.Next() {

		err := rows.Scan(&u.Username, &p.Content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue // Skip this row
		}
		posts = append(posts, p)
		users = append(users, u)
	}
	fmt.Println(posts)
	err = tmpl.Execute(w, posts)

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
	}

	if err != nil {
		// Log the error if template execution fails
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
