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
func (s *organizationUserService) CreateOrganizationUser(ctx context.Context, userID string, input repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	if _, err := s.access.IsAdmin(ctx, input.OrgID, userID); err != nil {
		return nil, err
	}

	if !models.ValidRoles[input.Role] {
		return nil, ErrInvalidInput
	}

	// orgID is validated in IsAdmin, so we can assume it's valid here
	if err := uuid.Validate(input.UserID); err != nil || input.UserID == "" {
		return nil, ErrInvalidInput
	}

	newOrgUser, err := s.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
		OrgID:  input.OrgID,
		UserID: input.UserID,
		Role:   input.Role,
	})
	if err != nil {
		// Log the error to get more context
		return nil, ErrInternalServer
	}

	return newOrgUser, nil
}

// GetUsersByOrganizationID retrieves all users within an organization.
func (s *organizationUserService) GetUsersByOrganizationID(ctx context.Context, userID, orgID string) ([]*models.UserWithRole, error) {
	if _, err := s.access.IsMember(ctx, orgID, userID); err != nil {
		return nil, err
	}

	users, err := s.orgUserRepo.GetUsersByOrganizationID(ctx, orgID)
	if err != nil {
		// Log the error to get more context
		return nil, ErrInternalServer
	}

	return users, nil
}

// UpdateRole updates a user's role within an organization.
func (s *organizationUserService) UpdateUserRole(ctx context.Context, orgID, userID string, newRole models.Role) error {
	if _, err := s.access.IsAdmin(ctx, orgID, userID); err != nil {
		return err
	}

	if !models.ValidRoles[newRole] {
		return ErrInvalidInput
	}

	err := s.orgUserRepo.UpdateRole(ctx, orgID, userID, newRole)
	if err != nil {
		// Log the error to get more context
		return ErrInternalServer
	}

	return nil
}

// DeleteUserFromOrganization removes a user from an organization.
func (s *organizationUserService) DeleteUserFromOrganization(ctx context.Context, orgID, userID, userIDToDelete string) error {
	if _, err := s.access.IsAdmin(ctx, orgID, userID); err != nil {
		return err
	}

	err := s.orgUserRepo.Delete(ctx, orgID, userIDToDelete)
	if err != nil {
		// Log the error to get more context
		return ErrInternalServer
	}

	return nil
}