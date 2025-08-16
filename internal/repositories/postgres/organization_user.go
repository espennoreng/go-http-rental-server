package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationUserRepository struct {
	db *pgxpool.Pool
	log *slog.Logger
}

func NewOrganizationUserRepository(db *pgxpool.Pool, log *slog.Logger) *OrganizationUserRepository {
	return &OrganizationUserRepository{
		db: db,
		log: log.With("component", "organization_user_repository"),
	}
}

var _ repositories.OrganizationUserRepository = (*OrganizationUserRepository)(nil)

func (r *OrganizationUserRepository) Create(ctx context.Context, params *repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	query := `
		INSERT INTO organization_users (organization_id, user_id, created_at, role)
		VALUES ($1, $2, $3, $4)
		RETURNING organization_id, user_id, created_at, role
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.Any("params", params))

	var orgUser models.OrganizationUser
	err := r.db.QueryRow(ctx, query, params.OrgID, params.UserID, time.Now(), params.Role).Scan(&orgUser.OrgID, &orgUser.UserID, &orgUser.CreatedAt, &orgUser.Role)
	if err != nil {
		r.log.Error("Failed to create organization user", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("Organization user created successfully", slog.String("org_id", orgUser.OrgID), slog.String("user_id", orgUser.UserID))

	return &orgUser, nil
}

func (r *OrganizationUserRepository) GetByID(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
	query := `
		SELECT organization_id, user_id, created_at, role
		FROM organization_users
		WHERE organization_id = $1 AND user_id = $2
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.String("org_id", orgID), slog.String("user_id", userID))

	var orgUser models.OrganizationUser
	err := r.db.QueryRow(ctx, query, orgID, userID).Scan(&orgUser.OrgID, &orgUser.UserID, &orgUser.CreatedAt, &orgUser.Role)
	if err != nil {
		r.log.Error("Failed to retrieve organization user by ID", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("Organization user retrieved successfully", slog.String("org_id", orgUser.OrgID), slog.String("user_id", orgUser.UserID))

	return &orgUser, nil
}

func (r *OrganizationUserRepository) GetUsersByOrganizationID(ctx context.Context, orgID string) ([]*models.UserWithRole, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at, ou.role
		FROM users u
		JOIN organization_users ou ON ou.user_id = u.id
		WHERE ou.organization_id = $1
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.String("org_id", orgID))

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		r.log.Error("Failed to retrieve users by organization ID", slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	orgUsersWithRole := make([]*models.UserWithRole, 0)
	for rows.Next() {
		var orgUser models.UserWithRole
		if err := rows.Scan(&orgUser.User.ID, &orgUser.User.Username, &orgUser.User.Email, &orgUser.User.CreatedAt, &orgUser.User.UpdatedAt, &orgUser.Role); err != nil {
			r.log.Error("Failed to scan organization user row", slog.Any("error", err))
			return nil, err
		}
		orgUsersWithRole = append(orgUsersWithRole, &orgUser)
	}
	if err := rows.Err(); err != nil {
		r.log.Error("Error occurred while iterating over organization users", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("Users retrieved successfully for organization", slog.String("org_id", orgID), slog.Int("user_count", len(orgUsersWithRole)))
	return orgUsersWithRole, nil
}

func (r *OrganizationUserRepository) Delete(ctx context.Context, orgID string, userID string) error {
	query := `
		DELETE FROM organization_users
		WHERE organization_id = $1 AND user_id = $2
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.String("org_id", orgID), slog.String("user_id", userID))

	_, err := r.db.Exec(ctx, query, orgID, userID)
	if err != nil {
		r.log.Error("Failed to delete organization user", slog.Any("error", err))
		return err
	}

	r.log.Info("Organization user deleted successfully", slog.String("org_id", orgID), slog.String("user_id", userID))

	return nil
}

func (r *OrganizationUserRepository) UpdateRole(ctx context.Context, orgID string, userID string, role models.Role) error {
	query := `
		UPDATE organization_users
		SET role = $1
		WHERE organization_id = $2 AND user_id = $3
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.String("org_id", orgID), slog.String("user_id", userID), slog.String("role", string(role)))

	_, err := r.db.Exec(ctx, query, role, orgID, userID)
	if err != nil {
		r.log.Error("Failed to update organization user role", slog.Any("error", err))
		return err
	}

	r.log.Info("Organization user role updated successfully", slog.String("org_id", orgID), slog.String("user_id", userID), slog.String("role", string(role)))
	return nil
}
