package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        string
	GoogleID  pgtype.Text
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserInput struct {
	Username string
	Email    string
}
