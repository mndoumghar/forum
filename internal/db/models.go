package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
	Count     int
	CountAll  int
}

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

func CheckPostId(postId int) (*Post, error) {
	var p Post
	err := DB.QueryRow("SELECT post_id FROM posts WHERE post_id = ?", postId).
		Scan(&p.ID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT user_id, email, username, password, created_at FROM users WHERE email = ? OR username = ?", email, email).
		Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmailUsername(email string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT user_id, email, username, password, created_at FROM users WHERE email = ? OR  username = ?", email, email).
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

// GetUserReaction returns the current reaction of a user to a post
func GetUserReaction(user_id int, post_id string) (string, error) {
	var reaction string
	err := DB.QueryRow(
		"SELECT likedislike FROM likedislike WHERE user_id = ? AND post_id = ?",
		user_id, post_id,
	).Scan(&reaction)

	if err == sql.ErrNoRows {
		return "", nil
	}
	return reaction, err
}

// InsertUserReaction adds a new reaction
func InsertUserReaction(user_id int, post_id string, reaction string) error {
	_, err := DB.Exec(
		"INSERT INTO likedislike(user_id, post_id, likedislike) VALUES(?,?,?)",
		user_id, post_id, reaction,
	)
	return err
}

// UpdateUserReaction changes existing reaction
func UpdateUserReaction(user_id int, post_id string, newReaction string) error {
	_, err := DB.Exec(
		"UPDATE likedislike SET likedislike = ? WHERE user_id = ? AND post_id = ?",
		newReaction, user_id, post_id,
	)
	return err
}

// DeleteUserReaction removes a reaction
func DeleteUserReaction(user_id int, post_id string) error {
	_, err := DB.Exec(
		"DELETE FROM likedislike WHERE user_id = ? AND post_id = ?",
		user_id, post_id,
	)
	return err
}

// GetLikeCount returns the number of likes for a post
func GetLikeCount(post_id string) (int, error) {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM likedislike WHERE post_id = ? AND likedislike = 'true'",
		post_id,
	).Scan(&count)
	return count, err
}

// GetDislikeCount returns the number of dislikes for a post
func GetDislikeCount(post_id string) (int, error) {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM likedislike WHERE post_id = ? AND likedislike = 'false'",
		post_id,
	).Scan(&count)
	return count, err
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
