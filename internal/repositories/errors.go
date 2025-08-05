package repositories

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("resource not found")
	ErrDuplicate = errors.New("duplicate resource")
	ErrInternal  = errors.New("internal server error")
)
