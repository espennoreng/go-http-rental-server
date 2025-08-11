package postgres_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPostgresUserRepository(t *testing.T) {
	th := SetupTestHelper(t)
	ctx := context.Background()


	t.Run("Create", func(t *testing.T) {
		th.ResetDB(t)

		params := &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		}
		newUser, err := th.userRepo.Create(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, newUser)
	})

	t.Run("Create_UniqueConstraint", func(t *testing.T) {
		th.ResetDB(t)

		params := &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		}
		newUser, err := th.userRepo.Create(ctx, params)
		require.NoError(t, err, "should not return an error on first creation")
		require.NotNil(t, newUser)

		params2 := &repositories.CreateUserParams{
			Username: "Jane Doe", // Same username as user1
			Email:    "jane@example.com",
		}
		newUser, err = th.userRepo.Create(ctx, params2)
		require.Error(t, err)
		require.Nil(t, newUser, "should return nil user on unique constraint violation")
	})

	t.Run("GetByID", func(t *testing.T) {
		th.ResetDB(t)

		newUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "johndoe@example.com",
		})
		require.NoError(t, err)

		user, err := th.userRepo.GetByID(ctx, newUser.ID)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, newUser.ID, user.ID)
		require.Equal(t, "John Doe", user.Username)
		require.Equal(t, "johndoe@example.com", user.Email)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		th.ResetDB(t)
		
		randomID := uuid.New().String()
		_, err := th.userRepo.GetByID(ctx, randomID)
		require.Error(t, err)
	})
}
