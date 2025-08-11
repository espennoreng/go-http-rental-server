package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type CreateOrganizationUserParams struct {
	OrgID  string `json:"org_id"`
	UserID string `json:"user_id"`
	Role   models.Role   `json:"role"`
}

type OrganizationUserRepository interface {
	Create(ctx context.Context, input *CreateOrganizationUserParams) (*models.OrganizationUser, error)
	GetByID(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error)
	GetUsersByOrganizationID(ctx context.Context, orgID string) ([]*models.UserWithRole, error)
	Delete(ctx context.Context, orgID string, userID string) error
	UpdateRole(ctx context.Context, orgID string, userID string, newRole models.Role) error
}
