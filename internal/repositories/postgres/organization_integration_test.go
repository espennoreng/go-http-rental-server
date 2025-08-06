package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/repositories/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPostgresOrganizationRepository(t *testing.T) {
	ctx := context.Background()

	_, err := dbpool.Exec(ctx, "TRUNCATE organizations RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	repo := postgres.NewOrganizationRepository(dbpool)

	t.Run("Create", func(t *testing.T) {
		testUserID := seedUser(t, ctx, "testuser")

		params := &repositories.CreateOrganizationParams{
			Name:      "Test Organization",
			CreatedBy: testUserID,
		}
		org, err := repo.Create(ctx, params)
		require.NoError(t, err)
		require.NotEmpty(t, org.ID)
	})

	t.Run("GetByID", func(t *testing.T) {
		org, err := repo.GetByID(ctx, "some-id")
		require.NoError(t, err)
		require.NotNil(t, org)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existent-id")
		require.Error(t, err)
		require.Equal(t, repositories.ErrNotFound, err)
	})
}

func seedUser(t *testing.T, ctx context.Context, username string) string {
	t.Helper()

	// In a real app, you might have a proper user repository.
	// For testing dependencies, a direct query is often simplest.
	userID := uuid.New().String()
	query := `INSERT INTO users (id, username, email) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`
	_, err := dbpool.Exec(ctx, query, userID, username, fmt.Sprintf("%s@example.com", username))
	require.NoError(t, err)

	return userID
}
