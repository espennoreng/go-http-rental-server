package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRepository interface {
	Create(ctx context.Context, params *CreateUserParams) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	FindOrCreateByGoogleID(ctx context.Context, googleID, email string) (*models.User, error)
}