package models

import (
	"database/sql"
	"forum/internal/db"
	"forum/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the forum
type User struct {
	ID       int
	UUID     string
	Email    string
	Username string
	Password string
}

// CreateUser creates a new user in the database with the provided email, username, and password.
func CreateUser(email, username, password string) error {
	// Hash the password before storing it
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// Generate a UUID for the user
	userUUID := utils.GenerateUUID()

	// Open the database connection
	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (uuid, email, username, password) VALUES (?, ?, ?, ?)", userUUID, email, username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(email string) (*User, error) {
	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL query to find a user by email
	row := db.QueryRow("SELECT id, uuid, email, username, password FROM users WHERE email = ?", email)

	var user User
	err = row.Scan(&user.ID, &user.UUID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Other errors
	}
	return &user, nil
}

// ValidatePassword checks if the provided password matches the stored password (hashed).
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// GetUserByUUID retrieves a user by their UUID.
func GetUserByUUID(uuid string) (*User, error) {
	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL query to find a user by UUID
	row := db.QueryRow("SELECT id, uuid, email, username, password FROM users WHERE uuid = ?", uuid)

	var user User
	err = row.Scan(&user.ID, &user.UUID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err // Other errors
	}
	return &user, nil
}

// IsEmailRegistered checks if the email already exists in the database.
func IsEmailRegistered(email string) (bool, error) {
	db, err := db.OpenDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
