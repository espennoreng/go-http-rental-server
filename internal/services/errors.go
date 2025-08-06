package services

import (
	"errors"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserNotFound      = errors.New("user not found")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationWithDuplicateDetailsExists = errors.New("organization with duplicate details exists")
	ErrInternalServer    = errors.New("internal server error")
	ErrUserWithDuplicateDetailsExists = errors.New("user with similar details already exists")
	ErrDuplicateInput    = errors.New("duplicate input")
)
