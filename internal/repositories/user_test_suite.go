package repositories

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/stretchr/testify/require"
)

func UserRepositoryTestSuite(t *testing.T, repo UserRepository) {
	t.Helper()
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		user := &models.User{
			ID:       "user-001",
			Username: "John Doe",
			Email:    "john@example.com",
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	})

	t.Run("Create_UniqueConstraint", func(t *testing.T) {
		user1 := &models.User{
			ID:       "user-002",
			Username: "Jane Doe",
			Email:    "jane@example.com",
		}
		err := repo.Create(ctx, user1)
		require.NoError(t, err, "should not return an error on first creation")

		user2 := &models.User{
			ID:       "user-003",
			Username: "Jane Doe", // Same username as user1
			Email:    "jane@example.com",
		}
		err = repo.Create(ctx, user2)
		require.Error(t, err)
		require.Equal(t, ErrUniqueConstraint, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		userID := "user-001"
		user, err := repo.GetByID(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, userID, user.ID)
		require.Equal(t, "John Doe", user.Username)
		require.Equal(t, "john@example.com", user.Email)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existent-id")
		require.Error(t, err)
		require.Equal(t, ErrNotFound, err)
	})
}
