package models

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserInput struct {
	Username string
	Email    string
}
