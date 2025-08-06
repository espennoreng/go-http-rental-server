package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, input repositories.CreateUserParams) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type OrganizationService interface {
	CreateOrganization(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error)
}