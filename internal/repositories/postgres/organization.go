package postgres

import (
	"context"
	"errors"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5/pgconn"
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
	`

	_, err := r.db.Exec(ctx, query, params.Name, params.CreatedBy)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			return nil, repositories.ErrUniqueConstraint
		}
		return nil, repositories.ErrInternal
	}
	return nil, nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	return nil, nil // Implementation will go here
}
