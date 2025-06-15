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
	CountAll  int
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

func CheckPostId(postId int) (*Post, error) {
	var p Post
	err := DB.QueryRow("SELECT post_id FROM posts WHERE post_id = ?", postId).
		Scan(&p.ID)
	if err != nil {
		return nil, err
	}
	return &p, nil
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

func GetLikeDisle(user_id int, post_id string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE user_id = ? AND post_id = ?", user_id, post_id).
		Scan(&u.Count)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CountLikeEveryPost(post_id string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT COUNT(*) FROM likedislike WHERE post_id = ?", post_id).Scan(&u.CountAll)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// this function checki lina Ila Can 3ndna Ktr mn user_id F Table likeDislike katrimove Azero bdpt Fdak Id User
// exmple mli kandght f form 3la like or Dislke  browser kol mra kaystocki true or false Ftable like  dislike Whna Bghina ghir mra whda Istock value Dyalo
// ila wrkana 3awtani 3la buton like Katcheck Ila deja m stock fih true or false kayremove mn jdid ...

func UpdateLikeDislike(user_id int, post_id string, Like string) error {
	if Like == "true" {
		_, err := DB.Exec("UPDATE likedislike set likedislike == 'false'  WHERE user_id = ? AND post_id = ? ", Like, user_id, post_id)
		if err != nil {
			return err
		}
	}
	if Like == "false" {
		_, err := DB.Exec("UPDATE likedislike set likedislike == 'true'  WHERE user_id = ? AND post_id = ? ", Like, user_id, post_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteIdUserikeDislike(user_id int, post_id string) error {
	_, err := DB.Exec("DELETE FROM likedislike WHERE user_id = ? AND post_id = ?", user_id, post_id)
	if err != nil {
		return err
	}
	return err
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
