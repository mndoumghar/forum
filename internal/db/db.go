package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}

	// Create tables
	createTables := `
			CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT  UNIQUE NOT NULL ,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);


		CREATE TABLE IF NOT EXISTS sessions (
			uuid  TEXT  PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);

		CREATE TABLE IF NOT EXISTS posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT ,
			content TEXT,
			status TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);



		CREATE TABLE IF NOT EXISTS comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
CREATE TABLE IF NOT EXISTS comment_reactions (
    reaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_like BOOLEAN NOT NULL,  -- TRUE for like, FALSE for dislike
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    UNIQUE(comment_id, user_id)  -- Each user can only react once per comment
);


		



		CREATE TABLE IF NOT EXISTS likedislike (
    likedislike_id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    likedislike TEXT CHECK (likedislike IN ('true', 'false')),
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);


--CREATE TRIGGER  IF NOT EXISTS mohamed
--BEFORE INSERT ON comments
--FOR EACH ROW
--BEGIN
  --  SELECT
   -- CASE
     --   WHEN NEW.content LIKE '%<script>%' THEN
       --     RAISE(ABORT, ' this is script')
    --END;
--END;

		
	`
	_, err = DB.Exec(createTables)
	return err
}

// GetDBConnection returns the current DB connection.
func GetDBConnection() (*sql.DB, error) {
    if DB == nil {
        return nil, fmt.Errorf("database connection is not initialized")
    }
    return DB, nil
}
