package postgres

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type organizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *organizationRepository {
	return &organizationRepository{
		db: db,
	}
}

var _ repositories.OrganizationRepository = (*organizationRepository)(nil)

func (r *organizationRepository) Create(ctx context.Context, params *repositories.CreateOrganizationParams) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (name, created_by)
		VALUES ($1, $2)
		RETURNING id, name, created_by, created_at, updated_at
	`

	var org models.Organization
	err := r.db.QueryRow(ctx, query, params.Name, params.CreatedBy).Scan(&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	query := `
		SELECT id, name, created_by, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var org models.Organization
	if err := row.Scan(&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt, &org.UpdatedAt); err != nil {
		return nil, err
	}

	return &org, nil
}
