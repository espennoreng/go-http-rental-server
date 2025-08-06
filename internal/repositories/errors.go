package repositories

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("resource not found")
	ErrUniqueConstraint = errors.New("unique constraint violation")
	ErrInternal  = errors.New("internal server error")
)
