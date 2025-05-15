package auth

import (
	"fmt"
	"net/http"
	"time"

	"forum/internal/db"

	"github.com/gofrs/uuid"
)

func CreateSession(w http.ResponseWriter, userID int) error {
	sessionID, err := uuid.NewV4()
	fmt.Println("UUid when Add Data base",sessionID)
	if err != nil {
		return err
	}
	expiration := time.Now().Add(24 * time.Hour)
	_, err = db.DB.Exec("INSERT INTO sessions (uuid, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID.String(), userID, expiration)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID.String(),
		Expires: expiration,
	})
	return nil
}

func CheckSession(w http.ResponseWriter, r *http.Request) (int, error) {
	// Retrieve the session ID from the cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			// No session cookie found
			return 0, fmt.Errorf("no session cookie found")
		}
		return 0, err
	}

	// Query the database to check if the session exists and is valid
	var userID int
	var expiresAt time.Time

	err = db.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE uuid = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		if err == nil {
			// Session not found
			return 0, fmt.Errorf("session not found ")
		}
		return 0, err
	}

	// Check if the session has expired
	if expiresAt.Before(time.Now()) {
		// Session expired
		return 0, fmt.Errorf("session expired")
	}

	// Session is valid
	return userID, nil
}

func DeletCoockies(w http.ResponseWriter, r *http.Request) error {
	// Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			// No session to delete
			return nil
		}
		return err
	}

	// Delete the session from the database
	_, err = db.DB.Exec("DELETE FROM sessions WHERE uuid = ?", cookie.Value)
	if err != nil {
		return err
	}

	// Clear the cookie by setting its MaxAge to -1
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",        // Make sure this matches the Path used when setting the cookie
		MaxAge:   -1,         // Instructs browser to delete the cookie
		HttpOnly: true,       // Optional: helps prevent XSS
	})

	return nil
}
