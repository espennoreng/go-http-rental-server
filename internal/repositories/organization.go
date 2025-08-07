package repositories

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type CreateOrganizationParams struct {
	Name      string  `json:"name"`
	CreatedBy string  `json:"created_by"`
}

type OrganizationRepository interface {
	Create(ctx context.Context, params *CreateOrganizationParams) (*models.Organization, error)
	GetByID(ctx context.Context, id string) (*models.Organization, error)
}
