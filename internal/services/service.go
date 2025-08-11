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

type CreateOrganizationUserParams struct {
	ActingUserID string
	OrgID        string
	UserID       string
	Role         models.Role
}

type GetUsersByOrganizationIDParams struct {
	OrgID        string
	ActingUserID string
}

type UpdateUserRoleParams struct {
	OrgID        string
	ActingUserID string
	NewRole      models.Role
}

type DeleteOrganizationUserParams struct {
	OrgID          string
	ActingUserID   string
	UserIDToDelete string
}

type OrganizationUserService interface {
	CreateOrganizationUser(ctx context.Context, params CreateOrganizationUserParams) (*models.OrganizationUser, error)
	GetUsersByOrganizationID(ctx context.Context, params GetUsersByOrganizationIDParams) ([]*models.UserWithRole, error)
	UpdateUserRole(ctx context.Context, params UpdateUserRoleParams) error
	DeleteUserFromOrganization(ctx context.Context, params DeleteOrganizationUserParams) error
}

type AccessService interface {
	IsAdmin(ctx context.Context, orgID, userID string) (bool, error)
	IsMember(ctx context.Context, orgID, userID string) (bool, error)
}
