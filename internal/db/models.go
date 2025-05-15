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
	Count     int
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

func GetLikeDisle(user_id int) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE user_id = ? ", user_id).
		Scan(&u.Count)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

/*
CREATE TABLE IF NOT EXISTS likedislike (
    likedislike_id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    likedislike TEXT CHECK (likedislike IN ('true', 'false')),
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);


*/
