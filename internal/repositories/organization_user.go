package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type OrganizationUserRepository interface {
	Create(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error)
	GetByOrganizationID(ctx context.Context, orgID string) ([]*models.User, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Organization, error)
	Delete(ctx context.Context, orgID string, userID string) error
}
