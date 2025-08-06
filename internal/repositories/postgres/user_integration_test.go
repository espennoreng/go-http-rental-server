package postgres_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/repositories/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)



func TestPostgresUserRepository(t *testing.T) {
	ctx := context.Background()

	_, err := dbpool.Exec(ctx, "TRUNCATE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	repo := postgres.NewUserRepository(dbpool)


	t.Run("Create", func(t *testing.T) {
		params := &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		}
		newUser, err := repo.Create(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, newUser)
	})

	t.Run("Create_UniqueConstraint", func(t *testing.T) {
		params := &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		}
		newUser, err := repo.Create(ctx, params)
		require.NoError(t, err, "should not return an error on first creation")
		require.NotNil(t, newUser)

		params2 := &repositories.CreateUserParams{
			Username: "Jane Doe", // Same username as user1
			Email:    "jane@example.com",
		}
		newUser, err = repo.Create(ctx, params2)
		require.Error(t, err)
		require.Nil(t, newUser, "should return nil user on unique constraint violation")
	})

	t.Run("GetByID", func(t *testing.T) {
		testUserID := seedUser(t, ctx, "testuser")

		user, err := repo.GetByID(ctx, testUserID)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, testUserID, user.ID)
		require.Equal(t, "testuser", user.Username)
		require.Equal(t, "testuser@example.com", user.Email)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		randomID := uuid.New().String()
		_, err := repo.GetByID(ctx, randomID)
		require.Error(t, err)
	})
}
