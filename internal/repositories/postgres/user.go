package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) FindOrCreateByGoogleID(ctx context.Context, googleID, email string) (*models.User, error) {
	var user models.User

	// Start a new transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails. The rollback will be ignored if the tx is committed.
	defer tx.Rollback(ctx)

	// 1. First, try to find the user by their Google ID. This is the most common case after the first login.
	queryByGoogleID := "SELECT id, google_id, email, username FROM users WHERE google_id = $1"
	err = tx.QueryRow(ctx, queryByGoogleID, googleID).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Username)
	if err == nil {
		// User found by Google ID, we are done.
		return &user, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		// A real database error occurred.
		return nil, err
	}

	// 2. User not found by Google ID. Try to find by email to link the account.
	log.Printf("User with google_id %s not found. Checking for existing user with email %s", googleID, email)
	queryByEmail := "SELECT id, google_id, email, username FROM users WHERE email = $1"
	err = tx.QueryRow(ctx, queryByEmail, email).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Username)
	if err == nil {
		// User found by email! Link the Google ID to this account.
		log.Printf("Found existing user with ID %s. Linking google_id.", user.ID)
		updateQuery := "UPDATE users SET google_id = $1 WHERE id = $2"
		_, updateErr := tx.Exec(ctx, updateQuery, googleID, user.ID)
		if updateErr != nil {
			return nil, updateErr
		}
		user.GoogleID.String = googleID
		user.GoogleID.Valid = true
		
		return &user, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Error querying user by email: %v", err)
		// A real database error occurred.
		return nil, err
	}

	// 3. No user found by Google ID or email. Create a new user.
	log.Printf("No existing user found. Creating new user for email %s", email)
	// NOTE: You may need to adjust the columns (e.g., 'username') based on your user model
	// and what you get from the Google token.
	insertQuery := "INSERT INTO users (google_id, email, username) VALUES ($1, $2, $3) RETURNING id, google_id, email, username"
	err = tx.QueryRow(ctx, insertQuery, googleID, email, "Default Username").Scan(&user.ID, &user.GoogleID, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	return &user, tx.Commit(ctx)
}