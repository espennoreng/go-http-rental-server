package services

import (
	"errors"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserNotFound      = errors.New("user not found")
	ErrInternalServer    = errors.New("internal server error")
	ErrUserWithDuplicateDetailsExists = errors.New("user with duplicate details exists")
	ErrDuplicateInput    = errors.New("duplicate input")
)
