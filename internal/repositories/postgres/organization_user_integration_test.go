package postgres_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
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

		orgUsers, err := th.orgUserRepo.GetUsersByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		// Check that the found organization users match the created ones
		require.Len(t, orgUsers, 1)
		require.Equal(t, user.ID, orgUsers[0].ID)
	})

	t.Run("GetUsersByOrganizationID_InvalidID", func(t *testing.T) {
		th.ResetDB(t)

		randomID := uuid.New().String()
		orgUsers, err := th.orgUserRepo.GetUsersByOrganizationID(ctx, randomID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		require.Len(t, orgUsers, 0, "should return an empty slice for non-existent organization ID")
	})

	t.Run("GetUsersByOrganizationID_NotPartOfOrg", func(t *testing.T) {
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
		orgUsers, err := th.orgUserRepo.GetUsersByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		// Should return an empty slice since the user is not part of the first organization
		require.Len(t, orgUsers, 0)
	})

	t.Run("GetOrganizationsByUserID", func(t *testing.T) {
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

		userOrgs, err := th.orgUserRepo.GetOrganizationsByUserID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, userOrgs)
		// Check that the found user organizations match the created ones
		require.Len(t, userOrgs, 1)
		require.Equal(t, org.ID, userOrgs[0].ID)
	})

	t.Run("GetOrganizationsByUserID_InvalidID", func(t *testing.T) {
		th.ResetDB(t)

		randomID := "invalid-user-id"
		userOrgs, err := th.orgUserRepo.GetOrganizationsByUserID(ctx, randomID)
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
		orgUsers, err := th.orgUserRepo.GetUsersByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		require.Len(t, orgUsers, 0, "should return an empty slice after deletion")
	})

	t.Run("Update organization user role", func(t *testing.T) {
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
		// Change the role of the organization user
		err = th.orgUserRepo.UpdateRole(ctx, org.ID, user.ID, models.RoleAdmin)
		require.NoError(t, err)
		// Verify that the role has been changed
		orgUsers, err := th.orgUserRepo.GetUsersByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.Len(t, orgUsers, 1)
		require.Equal(t, models.RoleAdmin, orgUsers[0].Role, "should have updated the role to admin")

	})
}
