package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetUserByIDParams struct {
	ActingUserID string `json:"acting_user_id"`
	UserID       string `json:"user_id"`
}

type UserService interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*models.User, error)
	GetUserByID(ctx context.Context, params GetUserByIDParams) (*models.User, error)
	FindOrCreateByGoogleID(ctx context.Context, googleID, email string) (*models.User, error)
}

type CreateOrganizationParams struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

type GetOrganizationByIDParams struct {
	ID string `json:"id"`
}

type OrganizationService interface {
	CreateOrganization(ctx context.Context, params CreateOrganizationParams) (*models.Organization, error)
	GetOrganizationByID(ctx context.Context, params GetOrganizationByIDParams) (*models.Organization, error)
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
	UserID      string
	Role         models.Role
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

type OrgAccessParams struct {
	OrgID  string
	UserID string
}

type AccessService interface {
	IsAdmin(ctx context.Context, params OrgAccessParams) error
	IsMember(ctx context.Context, params OrgAccessParams) error
}
