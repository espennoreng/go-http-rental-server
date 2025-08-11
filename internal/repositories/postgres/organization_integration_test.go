package postgres_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPostgresOrganizationRepository(t *testing.T) {
	th := SetupTestHelper(t)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		th.ResetDB(t)

		newUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "johndoe@example.com",
		})
		require.NoError(t, err)

		params := &repositories.CreateOrganizationParams{
			Name:      "Test Organization",
			CreatedBy: newUser.ID,
		}
		org, err := th.orgRepo.Create(ctx, params)
		require.NoError(t, err)
		require.NotEmpty(t, org.ID)
	})

	t.Run("GetByID", func(t *testing.T) {
		th.ResetDB(t)

		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "johndoe@example.com",
		})
		require.NoError(t, err)

		newOrg, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Organization",
			CreatedBy: user.ID,
		})
		require.NoError(t, err)

		org, err := th.orgRepo.GetByID(ctx, newOrg.ID)
		require.NoError(t, err)
		require.NotNil(t, org)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		th.ResetDB(t)

		randomUUID := uuid.New().String()
		_, err := th.orgRepo.GetByID(ctx, randomUUID)
		require.Error(t, err)
	})

	t.Run("Create_UniqueConstraint", func(t *testing.T) {
		th.ResetDB(t)

		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "janedoe@example.com",
		})
		require.NoError(t, err)

		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Organization",
			CreatedBy: user.ID,
		})
		require.NoError(t, err)
		require.NotEmpty(t, org.ID)

		// Try to create the same organization again
		org2, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Organization",
			CreatedBy: user.ID,
		})
		require.Error(t, err)
		require.Nil(t, org2)
	})
}
