package models

import "time"

type User struct {
	ID        int
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}
type ErrorRegister struct {
	Error string
	Color string
}

type Data struct {
	ErrorColor []ErrorRegister
}
