package DataBase

import ( "database/sql" ; "log" )

var DB *sql.DB // dont touch my global var "Sc4rlx" ; i guess we will need it for registrer login without create a new connections

func InitDB() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	//check database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	//create the users table
	qryyy := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(qryyy)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database initialized and users table created")
}