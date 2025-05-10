package db

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt string
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