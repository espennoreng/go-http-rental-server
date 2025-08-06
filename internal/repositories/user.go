package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type CreateUserParams struct {
	Username string
	Email    string
}

type UserRepository interface {
	Create(ctx context.Context, params *CreateUserParams) error
	GetByID(ctx context.Context, id string) (*models.User, error)
}
