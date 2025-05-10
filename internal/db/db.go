package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	// Create tables (from SQL provided earlier)
	createTables := `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		-- Add other CREATE TABLE statements here
	`
	_, err = DB.Exec(createTables)
	if err != nil {
		return err
	}

	return nil
}

	// IsUserExists checks if a username or email is already in use.
	func IsUserExists(username, email string) (bool, error) {
		var EXITs bool
		query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)`
		err := DB.QueryRow(query, username, email).Scan(&EXITs)
		if err != nil {
			return false, err
		}
		return EXITs, nil
	}

// CreateUser inserts a new user into the database.
	func CreateUser(username , email , hashedPassword string) error {
		query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
		_, err := DB.Exec(query, username, email, hashedPassword)
		if err != nil {
			return err
		}
		return err
	}

