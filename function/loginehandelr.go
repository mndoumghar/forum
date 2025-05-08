package function

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	 fmt.Println(username, "||||||||||||", password)

	var id int
	var hashedPassword string

	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &hashedPassword)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println(password, "|||||", hashedPassword)
	if err != nil {
		http.Error(w, "##Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate session ID
	sessionID := uuid.New().String()
	_, err = db.Exec("UPDATE users SET session = ? WHERE id = ?", sessionID, id)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
