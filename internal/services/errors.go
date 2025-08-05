package services

import (
	"errors"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserNotFound      = errors.New("user not found")
	ErrInternalServer    = errors.New("internal server error")
	ErrUserAlreadyExists = errors.New("user already exists")
)
