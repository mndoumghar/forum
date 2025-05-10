<<<<<<< HEAD:server/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/function"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "../users.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			username TEXT UNIQUE,
			password TEXT,
			session TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create default user (admin / 123456)
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	_, _ = db.Exec("INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)", "admin", password)

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/login", function.LoginHandler)
	http.HandleFunc("/dashboard", function.DashboardHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

/*

tmpl, err := template.ParseFiles("templates/result.html")
		if err != nil {
			RenderPageNotFound(w, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, asciiArt)
		if err != nil {
			RenderPageNotFound(w, http.StatusInternalServerError)
			return
		}


*/
=======
package main

import ( "log" ; "net/http" ;  "forum/server/rt")

func main() {
    rt.InitRoutes()

    log.Println("The server is running on port 8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Error starting server: ", err)
    }
}
>>>>>>> f5d7041ae4b99c153e912c23f0b06b8c2bc42bbe:server/server.go
