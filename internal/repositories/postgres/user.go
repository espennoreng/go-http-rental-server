package postgres

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

var _ repositories.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, user *repositories.CreateUserParams) (*models.User, error) {

	query := `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id, username, email, created_at, updated_at
	`

	var newUser models.User
	err := r.db.QueryRow(ctx, query, user.Username, user.Email).Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {

	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
