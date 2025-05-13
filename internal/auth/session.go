package auth

import (
	"net/http"
	"time"

	"forum/internal/db"
	"github.com/gofrs/uuid"
)

func CreateSession(w http.ResponseWriter, userID int) error {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	expiration := time.Now().Add(24 * time.Hour)
	_, err = db.DB.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)",
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
