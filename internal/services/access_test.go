package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestIsAdmin_Success asserts that the check passes for a valid admin.
func TestIsAdmin_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New().String()
	adminID := uuid.New().String()

	mockRepo := &mockOrganizationUserRepository{
		GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{Role: models.RoleAdmin}, nil
		},
	}
	s := services.NewAccessService(mockRepo)

	// Act
	err := s.IsAdmin(ctx, services.OrgAccessParams{OrgID: orgID, UserID: adminID})

	// Assert
	require.NoError(t, err)
}

// TestIsAdmin_FailsForNonAdmin asserts that the check fails for a user who is not an admin.
func TestIsAdmin_FailsForNonAdmin(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New().String()
	memberID := uuid.New().String()

	mockRepo := &mockOrganizationUserRepository{
		GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{Role: models.RoleMember}, nil
		},
	}
	s := services.NewAccessService(mockRepo)

	// Act
	err := s.IsAdmin(ctx, services.OrgAccessParams{OrgID: orgID, UserID: memberID})

	// Assert
	require.ErrorIs(t, err, services.ErrUnauthorized)
}

// TestIsAdmin_FailsForUserNotInOrg asserts that the check fails for a user not in the organization.
func TestIsAdmin_FailsForUserNotInOrg(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New().String()
	userID := uuid.New().String()

	mockRepo := &mockOrganizationUserRepository{
		GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
			return nil, nil // Simulate user not found
		},
	}
	s := services.NewAccessService(mockRepo)

	// Act
	err := s.IsAdmin(ctx, services.OrgAccessParams{OrgID: orgID, UserID: userID})

	// Assert
	require.ErrorIs(t, err, services.ErrUserNotPartOfOrganization)
}

// TestIsAdmin_FailsOnRepositoryError asserts that the check fails when the repository returns an error.
func TestIsAdmin_FailsOnRepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	orgID := uuid.New().String()
	userID := uuid.New().String()

	mockRepo := &mockOrganizationUserRepository{
		GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
			return nil, errors.New("database connection failed")
		},
	}
	s := services.NewAccessService(mockRepo)

	// Act
	err := s.IsAdmin(ctx, services.OrgAccessParams{OrgID: orgID, UserID: userID})

	// Assert
	// Note: A repository error should likely map to an internal server error, not an input error.
	require.ErrorIs(t, err, services.ErrInternalServer)
}

// TestIsAdmin_FailsForInvalidUUID asserts that the check fails for a malformed user ID.
func TestIsAdmin_FailsForInvalidUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	s := services.NewAccessService(&mockOrganizationUserRepository{}) // Mock can be empty

	// Act
	err := s.IsAdmin(ctx, services.OrgAccessParams{OrgID: uuid.New().String(), UserID: "invalid-uuid"})

	// Assert
	require.ErrorIs(t, err, services.ErrInvalidInput)
}
