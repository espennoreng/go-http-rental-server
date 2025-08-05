package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}
