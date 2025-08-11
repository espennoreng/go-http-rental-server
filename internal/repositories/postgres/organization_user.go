package postgres

import (
	"context"
	"time"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationUserRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationUserRepository(db *pgxpool.Pool) *OrganizationUserRepository {
	return &OrganizationUserRepository{
		db: db,
	}
}

var _ repositories.OrganizationUserRepository = (*OrganizationUserRepository)(nil)

func (r *OrganizationUserRepository) Create(ctx context.Context, params *repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	query := `
		INSERT INTO organization_users (organization_id, user_id, created_at, role)
		VALUES ($1, $2, $3, $4)
		RETURNING organization_id, user_id, created_at, role
	`

	var orgUser models.OrganizationUser
	err := r.db.QueryRow(ctx, query, params.OrgID, params.UserID, time.Now(), params.Role).Scan(&orgUser.OrgID, &orgUser.UserID, &orgUser.CreatedAt, &orgUser.Role)
	if err != nil {
		return nil, err
	}

	return &orgUser, nil
}

func (r *OrganizationUserRepository) GetByID(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
	query := `
		SELECT organization_id, user_id, created_at, role
		FROM organization_users
		WHERE organization_id = $1 AND user_id = $2
	`
	var orgUser models.OrganizationUser
	err := r.db.QueryRow(ctx, query, orgID, userID).Scan(&orgUser.OrgID, &orgUser.UserID, &orgUser.CreatedAt, &orgUser.Role)
	if err != nil {
		return nil, err
	}
	return &orgUser, nil
}

func (r *OrganizationUserRepository) GetUsersByOrganizationID(ctx context.Context, orgID string) ([]*models.UserWithRole, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.updated_at, ou.role
		FROM users u
		JOIN organization_users ou ON ou.user_id = u.id
		WHERE ou.organization_id = $1
	`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orgUsersWithRole := make([]*models.UserWithRole, 0)
	for rows.Next() {
		var orgUser models.UserWithRole
		if err := rows.Scan(&orgUser.User.ID, &orgUser.User.Username, &orgUser.User.Email, &orgUser.User.CreatedAt, &orgUser.User.UpdatedAt, &orgUser.Role); err != nil {
			return nil, err
		}
		orgUsersWithRole = append(orgUsersWithRole, &orgUser)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orgUsersWithRole, nil
}

func (r *OrganizationUserRepository) Delete(ctx context.Context, orgID string, userID string) error {
	query := `
		DELETE FROM organization_users
		WHERE organization_id = $1 AND user_id = $2
	`

	_, err := r.db.Exec(ctx, query, orgID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrganizationUserRepository) UpdateRole(ctx context.Context, orgID string, userID string, role models.Role) error {
	query := `
		UPDATE organization_users
		SET role = $1
		WHERE organization_id = $2 AND user_id = $3
	`

	_, err := r.db.Exec(ctx, query, role, orgID, userID)
	if err != nil {
		return err
	}

	return nil
}
