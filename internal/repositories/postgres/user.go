package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db: db,
		log: log.With("component", "user_repository"),
	}
}

var _ repositories.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, user *repositories.CreateUserParams) (*models.User, error) {

	query := `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id, username, email, created_at, updated_at
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.Any("params", user))

	var newUser models.User
	err := r.db.QueryRow(ctx, query, user.Username, user.Email).Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.CreatedAt, &newUser.UpdatedAt)

	if err != nil {
		r.log.Error("Failed to create user", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("User created successfully", slog.String("user_id", newUser.ID))

	return &newUser, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {

	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	r.log.Debug("Executing database query", slog.String("query", query), slog.String("user_id", id))

	row := r.db.QueryRow(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		r.log.Error("Failed to retrieve user by ID", slog.Any("error", err))
		return nil, err
	}

	r.log.Info("User retrieved successfully", slog.String("user_id", user.ID))

	return &user, nil
}

func (r *UserRepository) FindOrCreateByGoogleID(ctx context.Context, googleID, email string) (*models.User, error) {
	var user models.User

	log := r.log.With(
		slog.String("google_id", googleID),
		slog.String("email", email),
	)

	// Start a new transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error("Failed to begin transaction", slog.Any("error", err))
		return nil, err
	}
	// Defer a rollback in case anything fails. The rollback will be ignored if the tx is committed.
	defer tx.Rollback(ctx)

	// 1. First, try to find the user by their Google ID. This is the most common case after the first login.
	queryByGoogleID := "SELECT id, google_id, email, username FROM users WHERE google_id = $1"
	err = tx.QueryRow(ctx, queryByGoogleID, googleID).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Username)
	if err == nil {
		log.Info("User found by Google ID", slog.String("user_id", user.ID))
		return &user, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		log.Error("Error querying user by Google ID", slog.Any("error", err))
		return nil, err
	}

	// 2. User not found by Google ID. Try to find by email to link the account.
	log.Debug("User with google_id not found. Checking for existing user with email", slog.String("google_id", googleID), slog.String("email", email))
	queryByEmail := "SELECT id, google_id, email, username FROM users WHERE email = $1"
	err = tx.QueryRow(ctx, queryByEmail, email).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Username)
	if err == nil {
		// User found by email! Link the Google ID to this account.
		log.Info("Found existing user by email", slog.String("user_id", user.ID))
		updateQuery := "UPDATE users SET google_id = $1 WHERE id = $2"
		_, updateErr := tx.Exec(ctx, updateQuery, googleID, user.ID)
		if updateErr != nil {
			log.Error("Failed to update user with Google ID", slog.Any("error", updateErr))
			return nil, updateErr
		}
		user.GoogleID.String = googleID
		user.GoogleID.Valid = true

		return &user, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		log.Error("Error querying user by email", slog.Any("error", err))
		// A real database error occurred.
		return nil, err
	}

	// 3. No user found by Google ID or email. Create a new user.
	log.Debug("No existing user found. Creating new user for email", slog.String("email", email))
	// NOTE: You may need to adjust the columns (e.g., 'username') based on your user model
	// and what you get from the Google token.
	insertQuery := "INSERT INTO users (google_id, email, username) VALUES ($1, $2, $3) RETURNING id, google_id, email, username"
	err = tx.QueryRow(ctx, insertQuery, googleID, email, "Default Username").Scan(&user.ID, &user.GoogleID, &user.Email, &user.Username)
	if err != nil {
		log.Error("Failed to create new user", slog.Any("error", err))
		return nil, err
	}

	log.Info("New user created successfully", slog.String("user_id", user.ID))
	return &user, tx.Commit(ctx)
}