package models

import (
	"fmt"
	"forum/internal/db"
	"forum/internal/utils"
)

type User struct {
	ID       string
	Email    string
	Username string
	Password string
}

func CreateUser(email, username, password string) error {
	id := utils.NewUUID()
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec("INSERT INTO users (id, email, username, password_hash) VALUES (?, ?, ?, ?)", id, email, username, passwordHash)
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	return nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	row := db.DB.QueryRow("SELECT id, email, username, password_hash FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return user, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}
