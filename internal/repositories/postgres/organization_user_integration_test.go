package postgres_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/stretchr/testify/require"
)

func TestPostgresOrganizationUserRepository(t *testing.T) {

	th := SetupTestHelper(t)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		th.ResetDB(t)

		createOrgUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, createOrgUser)

		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)

		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		})

		require.NoError(t, err)
		require.NotNil(t, user)

		orgUser, err := th.orgUserRepo.Create(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)
	})

	t.Run("Create_AlreadyExists", func(t *testing.T) {

		th.ResetDB(t)

		createOrgUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, createOrgUser)

		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)

		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		})

		require.NoError(t, err)
		require.NotNil(t, user)

		orgUser, err := th.orgUserRepo.Create(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)

		// Attempt to create the same organization-user relationship again
		orgUserDup, err := th.orgUserRepo.Create(ctx, org.ID, user.ID)
		require.Error(t, err)
		require.Nil(t, orgUserDup, "should return nil on duplicate creation")
		require.Contains(t, err.Error(), "duplicate key value violates unique constraint")
	})

	t.Run("GetByOrganizationID", func(t *testing.T) {
		th.ResetDB(t)

		createOrgUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, createOrgUser)
		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)
		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user)

		orgUser, err := th.orgUserRepo.Create(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)

		orgUsers, err := th.orgUserRepo.GetByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		// Check that the found organization users match the created ones
		require.Len(t, orgUsers, 1)
		require.Equal(t, user.ID, orgUsers[0].ID)
	})

	t.Run("GetByOrganizationID_InvalidID", func(t *testing.T) {
		th.ResetDB(t)

		randomID := "invalid-org-id"
		orgUsers, err := th.orgUserRepo.GetByOrganizationID(ctx, randomID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		require.Len(t, orgUsers, 0, "should return an empty slice for non-existent organization ID")
	})

	t.Run("GetByOrganizationID_NotPartOfOrg", func(t *testing.T) {
		th.ResetDB(t)

		createOrgUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, createOrgUser)
		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)
		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user)
		// Create a different organization-user relationship
		org2, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Another Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)
		orgUser, err := th.orgUserRepo.Create(ctx, org2.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org2.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)
		// Attempt to get users by the first organization ID
		orgUsers, err := th.orgUserRepo.GetByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		// Should return an empty slice since the user is not part of the first organization
		require.Len(t, orgUsers, 0)
	})

	t.Run("GetByUserID", func(t *testing.T) {
		th.ResetDB(t)

		createOrgUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, createOrgUser)
		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)
		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user)

		orgUser, err := th.orgUserRepo.Create(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)

		userOrgs, err := th.orgUserRepo.GetByUserID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, userOrgs)
		// Check that the found user organizations match the created ones
		require.Len(t, userOrgs, 1)
		require.Equal(t, org.ID, userOrgs[0].ID)
	})

	t.Run("GetByUserID_InvalidID", func(t *testing.T) {
		th.ResetDB(t)

		randomID := "invalid-user-id"
		userOrgs, err := th.orgUserRepo.GetByUserID(ctx, randomID)
		require.NoError(t, err)
		require.NotNil(t, userOrgs)
		require.Len(t, userOrgs, 0, "should return an empty slice for non-existent user ID")
	})

	t.Run("Delete", func(t *testing.T) {
		th.ResetDB(t)

		createOrgUser, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "Jane Doe",
			Email:    "jane@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, createOrgUser)
		org, err := th.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
			Name:      "Test Org",
			CreatedBy: createOrgUser.ID,
		})
		require.NoError(t, err)
		user, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "John Doe",
			Email:    "john@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user)
		orgUser, err := th.orgUserRepo.Create(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUser)

		// Delete the organization-user relationship
		err = th.orgUserRepo.Delete(ctx, org.ID, user.ID)
		require.NoError(t, err)
		// Verify that the relationship no longer exists
		orgUsers, err := th.orgUserRepo.GetByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		require.Len(t, orgUsers, 0, "should return an empty slice after deletion")
	})
}
