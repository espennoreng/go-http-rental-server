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

		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
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

		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)

		// Attempt to create the same organization-user relationship again
		orgUserDup, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
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

		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)

		orgUsers, err := th.orgUserRepo.GetUsersByOrganizationID(ctx, org.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUsers)
		// Check that the found organization users match the created ones
		require.Len(t, orgUsers, 1)
		require.Equal(t, user.ID, orgUsers[0].User.ID)
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
		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org2.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
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
		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
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
		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
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

	t.Run("GetByID", func(t *testing.T) {
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
		orgUser, err := th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user.ID,
			Role:   models.RoleMember,
		})
		require.NoError(t, err)
		require.NotNil(t, orgUser)
		require.Equal(t, org.ID, orgUser.OrgID)
		require.Equal(t, user.ID, orgUser.UserID)

		// Get the organization user by ID
		orgUserRetrieved, err := th.orgUserRepo.GetByID(ctx, org.ID, user.ID)
		require.NoError(t, err)
		require.NotNil(t, orgUserRetrieved)
		require.Equal(t, org.ID, orgUserRetrieved.OrgID)
		require.Equal(t, user.ID, orgUserRetrieved.UserID)
		require.Equal(t, models.RoleMember, orgUserRetrieved.Role)
	})
	t.Run("GetByID_NotFound", func(t *testing.T) {
		th.ResetDB(t)

		// Attempt to get an organization user that does not exist
		orgUserRetrieved, err := th.orgUserRepo.GetByID(ctx, "non-existent-org-id", "non-existent-user-id")
		require.Error(t, err)
		require.Nil(t, orgUserRetrieved, "should return nil for non-existent organization user")
	})

	t.Run("AreUsersInSameOrganization", func(t *testing.T) {
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
		require.NotNil(t, org)

		user1, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "User One",
			Email:    "user1@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user1)

		user2, err := th.userRepo.Create(ctx, &repositories.CreateUserParams{
			Username: "User Two",
			Email:    "user2@example.com",
		})
		require.NoError(t, err)
		require.NotNil(t, user2)

		_, err = th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user1.ID,
			Role:   models.RoleMember,
		})
		require.NoError(t, err)

		_, err = th.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
			OrgID:  org.ID,
			UserID: user2.ID,
			Role:   models.RoleMember,
		})
		require.NoError(t, err)

		OK, err := th.orgUserRepo.AreUsersInSameOrg(ctx, &repositories.AreUsersInSameOrgParams{
			UserID1: user1.ID,
			UserID2: user2.ID,
		})

		require.NoError(t, err)
		require.True(t, OK)
	})


}
