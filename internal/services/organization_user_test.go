package services_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockOrganizationUserRepository struct {
	CreateOrganizationUserFunc   func(ctx context.Context, input repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error)
	GetUsersByOrganizationIDFunc func(ctx context.Context, orgID string) ([]*models.UserWithRole, error)
	UpdateUserRoleFunc           func(ctx context.Context, orgID, userID string, newRole models.Role) error
	DeleteOrganizationUserFunc   func(ctx context.Context, orgID, userID string) error
	GetByIDFunc                  func(ctx context.Context, orgID, userID string) (*models.OrganizationUser, error)
}

func (m *mockOrganizationUserRepository) Create(ctx context.Context, input *repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	return m.CreateOrganizationUserFunc(ctx, *input)
}

func (m *mockOrganizationUserRepository) GetUsersByOrganizationID(ctx context.Context, orgID string) ([]*models.UserWithRole, error) {
	return m.GetUsersByOrganizationIDFunc(ctx, orgID)
}

func (m *mockOrganizationUserRepository) UpdateRole(ctx context.Context, orgID, userID string, newRole models.Role) error {
	return m.UpdateUserRoleFunc(ctx, orgID, userID, newRole)
}

func (m *mockOrganizationUserRepository) Delete(ctx context.Context, orgID, userID string) error {
	return m.DeleteOrganizationUserFunc(ctx, orgID, userID)
}
func (m *mockOrganizationUserRepository) GetByID(ctx context.Context, orgID, userID string) (*models.OrganizationUser, error) {
	return m.GetByIDFunc(ctx, orgID, userID)
}

func TestOrganizationUserService_Create(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockOrganizationUserRepository{
		CreateOrganizationUserFunc: func(ctx context.Context, input repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				OrgID:  input.OrgID,
				UserID: input.UserID,
				Role:   input.Role,
			}, nil
		},
		GetByIDFunc: func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				OrgID:  orgID,
				UserID: userID,
				Role:   models.RoleAdmin, // Simulating that the user is an admin
			}, nil
		},
	}

	service := services.NewOrganizationUserService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		orgID := uuid.New().String()
		userID := uuid.New().String()
		orgUser, err := service.CreateOrganizationUser(ctx, services.CreateOrganizationUserParams{
			ActingUserID: userID,
			OrgID:        orgID,
			UserID:       userID,
			Role:         models.RoleMember,
		})
		assert.NoError(t, err)
		assert.NotNil(t, orgUser)
		assert.Equal(t, orgID, orgUser.OrgID)
		assert.Equal(t, userID, orgUser.UserID)
		assert.Equal(t, models.RoleMember, orgUser.Role)
	})

	t.Run("invalid organization ID", func(t *testing.T) {
		_, err := service.CreateOrganizationUser(ctx, services.CreateOrganizationUserParams{
			ActingUserID: uuid.New().String(),
			OrgID:        "invalid-id",
			UserID:       uuid.New().String(),
			Role:         models.RoleMember,
		})
		assert.Error(t, err)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		_, err := service.CreateOrganizationUser(ctx, services.CreateOrganizationUserParams{
			ActingUserID: uuid.New().String(),
			OrgID:        uuid.New().String(),
			UserID:       "invalid-id",
			Role:         models.RoleMember,
		})
		assert.Error(t, err)
	})

	t.Run("invalid role", func(t *testing.T) {
		_, err := service.CreateOrganizationUser(ctx, services.CreateOrganizationUserParams{
			ActingUserID: uuid.New().String(),
			OrgID:        uuid.New().String(),
			UserID:       uuid.New().String(),
			Role:         "invalid-role",
		})
		assert.Error(t, err)
	})

	t.Run("user is not organization admin", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				OrgID:  orgID,
				UserID: userID,
				Role:   models.RoleMember, // Simulating that the user is not an admin
			}, nil
		}
		orgUser, err := service.CreateOrganizationUser(ctx, services.CreateOrganizationUserParams{
			ActingUserID: uuid.New().String(),
			OrgID:        uuid.New().String(),
			UserID:       uuid.New().String(),
			Role:         models.RoleMember,
		})
		assert.Error(t, err)
		assert.Equal(t, services.ErrUnauthorized, err)
		assert.Nil(t, orgUser)
	})

	t.Run("user not part of organization", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return nil, nil // Simulating that the user is not part of the organization
		}
		orgUser, err := service.CreateOrganizationUser(ctx, services.CreateOrganizationUserParams{
			ActingUserID: uuid.New().String(),
			OrgID:        uuid.New().String(),
			UserID:       uuid.New().String(),
			Role:         models.RoleMember,
		})
		assert.Error(t, err)
		assert.Equal(t, services.ErrUserNotPartOfOrganization, err)
		assert.Nil(t, orgUser)
	})
}

func TestOrganizationUserService_GetUsersByOrganizationID(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockOrganizationUserRepository{
		GetUsersByOrganizationIDFunc: func(ctx context.Context, orgID string) ([]*models.UserWithRole, error) {
			return []*models.UserWithRole{
				{
					User: models.User{
						ID:       uuid.New().String(),
						Username: "john_doe",
						Email:    "john_doe@example.com",
					},
					Role: models.RoleMember,
				},
			}, nil
		},
		GetByIDFunc: func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				OrgID:  orgID,
				UserID: userID,
				Role:   models.RoleMember,
			}, nil
		},
	}
	service := services.NewOrganizationUserService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		users, err := service.GetUsersByOrganizationID(ctx, services.GetUsersByOrganizationIDParams{
			OrgID:        uuid.New().String(),
			ActingUserID: uuid.New().String(),
		})
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 1)
		assert.Equal(t, "john_doe", users[0].Username)
		assert.Equal(t, "john_doe@example.com", users[0].Email)
		assert.Equal(t, models.RoleMember, users[0].Role)
	})

	t.Run("invalid organization ID", func(t *testing.T) {
		_, err := service.GetUsersByOrganizationID(ctx, services.GetUsersByOrganizationIDParams{
			OrgID:        "invalid-id",
			ActingUserID: uuid.New().String(),
		})
		assert.Error(t, err)
	})
	t.Run("empty organization ID", func(t *testing.T) {
		_, err := service.GetUsersByOrganizationID(ctx, services.GetUsersByOrganizationIDParams{
			OrgID:        "",
			ActingUserID: uuid.New().String(),
		})
		assert.Error(t, err)
	})
	t.Run("user not part of organization", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return nil, nil // Simulating that the user is not part of the organization
		}
		_, err := service.GetUsersByOrganizationID(ctx, services.GetUsersByOrganizationIDParams{
			OrgID:        uuid.New().String(),
			ActingUserID: uuid.New().String(),
		})
		assert.Error(t, err)
		assert.Equal(t, services.ErrUserNotPartOfOrganization, err)
	})
}

func TestOrganizationUserService_UpdateUserRole(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockOrganizationUserRepository{
		UpdateUserRoleFunc: func(ctx context.Context, orgID string, userID string, newRole models.Role) error {
			return nil
		},
		GetByIDFunc: func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				OrgID:  orgID,
				UserID: userID,
				Role:   models.RoleAdmin,
			}, nil
		},
	}
	service := services.NewOrganizationUserService(mockRepo)

	t.Run("successful role update", func(t *testing.T) {
		err := service.UpdateUserRole(ctx, services.UpdateUserRoleParams{
			OrgID:        uuid.New().String(),
			ActingUserID: uuid.New().String(),
			NewRole:      models.RoleMember,
		})
		assert.NoError(t, err)
	})

	t.Run("invalid organization ID", func(t *testing.T) {
		err := service.UpdateUserRole(ctx, services.UpdateUserRoleParams{
			OrgID:        "invalid-id",
			ActingUserID: uuid.New().String(),
			NewRole:      models.RoleAdmin,
		})
		assert.Error(t, err)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		err := service.UpdateUserRole(ctx, services.UpdateUserRoleParams{
			OrgID:        uuid.New().String(),
			ActingUserID: "invalid-id",
			NewRole:      models.RoleAdmin,
		})
		assert.Error(t, err)
	})

	t.Run("invalid role", func(t *testing.T) {
		err := service.UpdateUserRole(ctx, services.UpdateUserRoleParams{
			OrgID:        uuid.New().String(),
			ActingUserID: uuid.New().String(),
			NewRole:      "invalid-role",
		})
		assert.Error(t, err)
	})

	t.Run("user not part of organization", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return nil, nil // Simulating that the user is not part of the organization
		}
		err := service.UpdateUserRole(ctx, services.UpdateUserRoleParams{
			OrgID:        uuid.New().String(),
			ActingUserID: uuid.New().String(),
			NewRole:      models.RoleAdmin,
		})
		assert.Error(t, err)
		assert.Equal(t, services.ErrUserNotPartOfOrganization, err)
	})

	t.Run("current role is not admin", func(t *testing.T) {
		mockRepo.GetByIDFunc = func(ctx context.Context, orgID string, userID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				OrgID:  orgID,
				UserID: userID,
				Role:   models.RoleMember,
			}, nil
		}
		err := service.UpdateUserRole(ctx, services.UpdateUserRoleParams{
			OrgID:        uuid.New().String(),
			ActingUserID: uuid.New().String(),
			NewRole:      models.RoleAdmin,
		})
		assert.Error(t, err)
		assert.Equal(t, services.ErrUnauthorized, err)
	})

}
