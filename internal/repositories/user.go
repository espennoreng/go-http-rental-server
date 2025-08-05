package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
}
