package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
)

type organizationUserService struct {
	orgUserRepo repositories.OrganizationUserRepository
	access      AccessService
}

// NewOrganizationUserService initializes a new organizationUserService.
func NewOrganizationUserService(orgUserRepo repositories.OrganizationUserRepository, access AccessService) *organizationUserService {
	return &organizationUserService{
		orgUserRepo: orgUserRepo,
		access:      access,
	}
}

var _ OrganizationUserService = (*organizationUserService)(nil)

// CreateOrganizationUser handles the creation of a new organization-user relationship.
func (s *organizationUserService) CreateOrganizationUser(ctx context.Context, params CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	if _, err := s.access.IsAdmin(ctx, params.OrgID, params.ActingUserID); err != nil {
		return nil, err
	}

	if !models.ValidRoles[params.Role] {
		return nil, ErrInvalidInput
	}

	// orgID is validated in IsAdmin, so we can assume it's valid here
	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
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
	if _, err := s.access.IsMember(ctx, params.OrgID, params.ActingUserID); err != nil {
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
	if _, err := s.access.IsAdmin(ctx, params.OrgID, params.ActingUserID); err != nil {
		return err
	}

	if !models.ValidRoles[params.NewRole] {
		return ErrInvalidInput
	}

	err := s.orgUserRepo.UpdateRole(ctx, params.OrgID, params.ActingUserID, params.NewRole)
	if err != nil {
		// Log the error to get more context
		return ErrInternalServer
	}

	return nil
}

// DeleteUserFromOrganization removes a user from an organization.
func (s *organizationUserService) DeleteUserFromOrganization(ctx context.Context, params DeleteOrganizationUserParams) error {
	if _, err := s.access.IsAdmin(ctx, params.OrgID, params.ActingUserID); err != nil {
		return err
	}

	err := s.orgUserRepo.Delete(ctx, params.OrgID, params.UserIDToDelete)
	if err != nil {
		// Log the error to get more context
		return ErrInternalServer
	}

	return nil
}
