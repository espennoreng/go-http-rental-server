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

type OrganizationUserService interface {
	CreateOrganizationUser(ctx context.Context, userID string, input repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error)
	GetUsersByOrganizationID(ctx context.Context, orgID, userID string) ([]*models.UserWithRole, error)
	UpdateUserRole(ctx context.Context, orgID, userID string, newRole models.Role) error
	DeleteUserFromOrganization(ctx context.Context, orgID, userID,  userIDToDelete string) error
}

type AccessService interface {
	IsAdmin(ctx context.Context, orgID, userID string) (bool, error)
	IsMember(ctx context.Context, orgID, userID string) (bool, error)
}