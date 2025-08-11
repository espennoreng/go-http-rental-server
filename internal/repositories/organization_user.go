package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type OrganizationUserRepository interface {
	Create(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error)
	GetUsersByOrganizationID(ctx context.Context, orgID string) ([]*models.UserWithRole, error)
	GetOrganizationsByUserID(ctx context.Context, userID string) ([]*models.Organization, error)
	Delete(ctx context.Context, orgID string, userID string) error
	UpdateRole(ctx context.Context, orgID string, userID string, newRole models.Role) error
}
