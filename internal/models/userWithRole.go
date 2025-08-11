package models

type UserWithRole struct {
	User `json:"user"`
	Role Role `json:"role"`
}
