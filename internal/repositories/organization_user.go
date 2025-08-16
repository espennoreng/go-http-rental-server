package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type CreateOrganizationUserParams struct {
	OrgID  string      `json:"org_id"`
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
}

type UpdateUserRoleParams struct {
	NewRole models.Role `json:"new_role"`
}

type AreUsersInSameOrgParams struct {
	UserID1 string `json:"user_id_1"`
	UserID2 string `json:"user_id_2"`
}

type OrganizationUserRepository interface {
	Create(ctx context.Context, input *CreateOrganizationUserParams) (*models.OrganizationUser, error)
	GetByID(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error)
	GetUsersByOrganizationID(ctx context.Context, orgID string) ([]*models.UserWithRole, error)
	Delete(ctx context.Context, orgID string, userID string) error
	UpdateRole(ctx context.Context, orgID string, userID string, newRole models.Role) error
	AreUsersInSameOrg(ctx context.Context, params *AreUsersInSameOrgParams) (bool, error)
}
