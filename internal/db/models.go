package db

import (
	"time"
)

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT user_id, email, username, password, created_at FROM users WHERE email = ?", email).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
func GetPost() (*Post, *User, error) {
	var p Post
	var u User
	// Ensure you correct the column name `usernam` to `username`
	err := DB.QueryRow(`SELECT u.username, p.content FROM posts p JOIN users u ON p.user_id = u.user_id LIMIT 1`).
		Scan(&u.Username, &p.Content)
	if err != nil {
		return nil, nil, err // Return nil for both User and Post in case of an error
	}
	return &p, &u, nil // Return the Post and User objects if no error
}

/*


{{range .Posts}}
				<li>{{.Content}}</li>
			{{end}}




CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	username TEXT NOT NULL,
	password TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS sessions (
	session_id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS posts (
	post_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(user_id)
);
*/
