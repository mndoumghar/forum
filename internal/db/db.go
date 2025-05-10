package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// OpenDB opens a connection to the SQLite database.
// It returns a pointer to the database connection.
func OpenDB() (*sql.DB, error) {
	// Open the SQLite database file
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Error opening the database: ", err)
		return nil, err
	}

	// Verify the connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
		return nil, err
	}

	return db, nil
}

// User represents a user in the database
type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt string
}

// GetUserByEmail retrieves a user by their email.
func GetUserByEmail(email string) (*User, error) {
	// Open the database connection
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close() // Ensure the DB connection is closed

	// Prepare the query to find a user by email
	var u User
	err = db.QueryRow("SELECT id, email, username, password, created_at FROM users WHERE email = ?", email).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
