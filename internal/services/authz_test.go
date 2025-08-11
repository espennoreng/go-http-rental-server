package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Ensure the mock implements the interface
var _ repositories.OrganizationUserRepository = (*mockOrganizationUserRepository)(nil)

// TestAccessService_IsAdmin tests the IsAdmin method of the AccessService.
func TestAccessService_IsAdmin(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New().String()
	adminID := uuid.New().String()
	memberID := uuid.New().String()
	nonExistentID := uuid.New().String()

	tests := []struct {
		name       string
		orgID      string
		userID     string
		mock       *mockOrganizationUserRepository
		wantErr    error
		wantErrMsg string
	}{
		{
			name:   "Valid Admin",
			orgID:  orgID,
			userID: adminID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return &models.OrganizationUser{
						OrgID:  oID,
						UserID: uID,
						Role:   models.RoleAdmin,
					}, nil
				},
			},
			wantErr: nil,
		},
		{
			name:    "Invalid UserID UUID",
			orgID:   orgID,
			userID:  "invalid-uuid",
			mock:    &mockOrganizationUserRepository{},
			wantErr: services.ErrInvalidInput,
		},
		{
			name:   "User Is Not Admin",
			orgID:  orgID,
			userID: memberID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return &models.OrganizationUser{
						OrgID:  oID,
						UserID: uID,
						Role:   models.RoleMember,
					}, nil
				},
			},
			wantErr: services.ErrUnauthorized,
		},
		{
			name:   "User Not Part Of Organization",
			orgID:  orgID,
			userID: nonExistentID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return nil, nil // Return nil for the user to simulate not found
				},
			},
			wantErr: services.ErrUserNotPartOfOrganization,
		},
		{
			name:   "Repository Returns Error",
			orgID:  orgID,
			userID: adminID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return nil, errors.New("database connection failed")
				},
			},
			wantErr: services.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := services.NewAccessService(tt.mock)
			err := s.IsAdmin(ctx, tt.orgID, tt.userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestAccessService_IsMember tests the IsMember method of the AccessService.
func TestAccessService_IsMember(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New().String()
	memberID := uuid.New().String()
	nonExistentID := uuid.New().String()

	tests := []struct {
		name       string
		orgID      string
		userID     string
		mock       *mockOrganizationUserRepository
		wantErr    error
		wantErrMsg string
	}{
		{
			name:   "Valid Member",
			orgID:  orgID,
			userID: memberID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return &models.OrganizationUser{
						OrgID:  oID,
						UserID: uID,
						Role:   models.RoleMember,
					}, nil
				},
			},
			wantErr: nil,
		},
		{
			name:    "Invalid UserID UUID",
			orgID:   orgID,
			userID:  "invalid-uuid",
			mock:    &mockOrganizationUserRepository{},
			wantErr: services.ErrInvalidInput,
		},
		{
			name:   "User Not Part Of Organization",
			orgID:  orgID,
			userID: nonExistentID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return nil, nil
				},
			},
			wantErr: services.ErrUserNotPartOfOrganization,
		},
		{
			name:   "Repository Returns Error",
			orgID:  orgID,
			userID: memberID,
			mock: &mockOrganizationUserRepository{
				GetByIDFunc: func(ctx context.Context, oID, uID string) (*models.OrganizationUser, error) {
					return nil, errors.New("database connection failed")
				},
			},
			wantErr: services.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := services.NewAccessService(tt.mock)
			err := s.IsMember(ctx, tt.orgID, tt.userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
