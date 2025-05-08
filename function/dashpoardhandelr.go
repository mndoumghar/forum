package function

import (
	"fmt"
	"net/http"
	"text/template"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type Users struct {
		Username string
		password string
	}

	var username string
	var password string

	err = db.QueryRow("SELECT username password FROM users WHERE session = ?", cookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "Session invalid", http.StatusUnauthorized)
		return
	}
	err = db.QueryRow("SELECT password FROM users WHERE session = ?", cookie.Value).Scan(&password)
	if err != nil {
		http.Error(w, "Session invalid", http.StatusUnauthorized)
		return
	}

	users := &Users{
		Username: username,
	}

	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		fmt.Sprintln(w, "Error File")
		return

	}
	tmpl.Execute(w, users)

	// fmt.Println(cookie.Value)
	// fmt.Fprintf(w, "<h1>Welcome, %s ----- %s</h1>", username, password)
}
