package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationRepository struct {
	db *pgxpool.Pool
	log *slog.Logger
}

func NewOrganizationRepository(db *pgxpool.Pool, log *slog.Logger) *OrganizationRepository {
	return &OrganizationRepository{
		db: db,
		log: log.With("component", "organization_repository"),
	}
}

var _ repositories.OrganizationRepository = (*OrganizationRepository)(nil)

func (r *OrganizationRepository) Create(ctx context.Context, params *repositories.CreateOrganizationParams) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (name, created_by)
		VALUES ($1, $2)
		RETURNING id, name, created_by, created_at, updated_at
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.Any("params", params))

	var org models.Organization
	err := r.db.QueryRow(ctx, query, params.Name, params.CreatedBy).Scan(&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			r.log.Warn("Organization with the same name already exists", slog.Any("error", err))
			return nil, repositories.ErrConflict
		}
		r.log.Error("Failed to create organization", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("Organization created successfully", slog.String("org_id", org.ID))

	return &org, nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	query := `
		SELECT id, name, created_by, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.String("org_id", id))

	row := r.db.QueryRow(ctx, query, id)

	var org models.Organization
	if err := row.Scan(&org.ID, &org.Name, &org.CreatedBy, &org.CreatedAt, &org.UpdatedAt); err != nil {
		r.log.Error("Failed to retrieve organization by ID", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("Organization retrieved successfully", slog.String("org_id", org.ID))

	return &org, nil
}
