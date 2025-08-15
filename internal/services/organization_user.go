package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
)

type organizationUserService struct {
	orgUserRepo repositories.OrganizationUserRepository
	accessService AccessService
}

// NewOrganizationUserService initializes a new organizationUserService.
func NewOrganizationUserService(orgUserRepo repositories.OrganizationUserRepository, accessService AccessService) *organizationUserService {
	return &organizationUserService{
		orgUserRepo: orgUserRepo,
		accessService: accessService,
	}
}

var _ OrganizationUserService = (*organizationUserService)(nil)

// CreateOrganizationUser handles the creation of a new organization-user relationship.
func (s *organizationUserService) CreateOrganizationUser(ctx context.Context, params CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	err := s.accessService.IsAdmin(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})
	if err != nil {
		return nil, err
	}

	if !models.ValidRoles[params.Role] {
		return nil, ErrInvalidInput
	}

	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
		return nil, ErrInvalidInput
	}

	if err := uuid.Validate(params.OrgID); err != nil || params.OrgID == "" {
		return nil, ErrInvalidInput
	}

	newOrgUser, err := s.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
		OrgID:  params.OrgID,
		UserID: params.UserID,
		Role:   params.Role,
	})
	if err != nil {
		// Log the error to get more context
		return nil, ErrInternalServer
	}

	return newOrgUser, nil
}

// GetUsersByOrganizationID retrieves all users within an organization.
func (s *organizationUserService) GetUsersByOrganizationID(ctx context.Context, params GetUsersByOrganizationIDParams) ([]*models.UserWithRole, error) {
	err := s.accessService.IsMember(ctx, OrgAccessParams{
		OrgID: params.OrgID,
		UserID: params.ActingUserID,
	})

	if err != nil {
		return nil, err
	}

	users, err := s.orgUserRepo.GetUsersByOrganizationID(ctx, params.OrgID)
	if err != nil {
		// Log the error to get more context
		return nil, ErrInternalServer
	}

	return users, nil
}

// UpdateRole updates a user's role within an organization.
func (s *organizationUserService) UpdateUserRole(ctx context.Context, params UpdateUserRoleParams) error {
	err := s.accessService.IsAdmin(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})
	if err != nil {
		return err
	}

	if !models.ValidRoles[params.Role] {
		return ErrInvalidInput
	}

	err = s.orgUserRepo.UpdateRole(ctx, params.OrgID, params.ActingUserID, params.Role)
	if err != nil {
		// Log the error to get more context
		return ErrInternalServer
	}

	return nil
}

// DeleteUserFromOrganization removes a user from an organization.
func (s *organizationUserService) DeleteUserFromOrganization(ctx context.Context, params DeleteOrganizationUserParams) error {
	err := s.accessService.IsAdmin(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})
	if err != nil {
		return err
	}

	err = s.orgUserRepo.Delete(ctx, params.OrgID, params.UserIDToDelete)
	if err != nil {
		// Log the error to get more context
		return ErrInternalServer
	}

	return nil
}
