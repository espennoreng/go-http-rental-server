package postgres

import (
	"context"
	"errors"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: db,
	}
}

var _ repositories.UserRepository = (*userRepository)(nil)

func (r *userRepository) Create(ctx context.Context, user *models.User) error {

	query := `
		INSERT INTO users (id, username, email)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.Username, user.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			return repositories.ErrUniqueConstraint
		}
		return repositories.ErrInternal
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {

	query := `
		SELECT id, username, email
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repositories.ErrNotFound
		}
		return nil, repositories.ErrInternal
	}

	return &user, nil
}
