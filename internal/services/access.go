package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
)

// AccessService provides reusable methods to check user access rights.
type accessService struct {
	orgUserRepo repositories.OrganizationUserRepository
}

// NewAccessService creates a new instance of the access checker.
func NewAccessService(orgUserRepo repositories.OrganizationUserRepository) *accessService {
	return &accessService{
		orgUserRepo: orgUserRepo,
	}
}

var _ AccessService = (*accessService)(nil)

// IsAdmin checks if a user has 'admin' privileges in an organization.
// It returns ErrInvalidInput if the UUIDs are malformed,
// ErrUnauthorized if the user is not an admin, or a database error.
func (s *accessService) IsAdmin(ctx context.Context, params OrgAccessParams) error {
	if err := uuid.Validate(params.OrgID); err != nil || params.OrgID == "" {
		return ErrInvalidInput
	}
	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
		return ErrInvalidInput
	}

	orgUser, err := s.orgUserRepo.GetByID(ctx, params.OrgID, params.UserID)
	if err != nil {
		// Log the underlying error for debugging purposes
		return ErrInternalServer
	}
	if orgUser == nil {
		return ErrUserNotPartOfOrganization
	}
	if orgUser.Role != models.RoleAdmin {
		return ErrUnauthorized
	}
	return nil
}

// IsMember checks if a user is a member of an organization.
// It returns ErrInvalidInput if the UUIDs are malformed,
// ErrUnauthorized if the user is not a member, or a database error.
func (s *accessService) IsMember(ctx context.Context, params OrgAccessParams) error {
	if err := uuid.Validate(params.OrgID); err != nil || params.OrgID == "" {
		return ErrInvalidInput
	}
	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
		return ErrInvalidInput
	}

	orgUser, err := s.orgUserRepo.GetByID(ctx, params.OrgID, params.UserID)
	if err != nil {
		// Log the underlying error for debugging purposes
		return ErrInternalServer
	}
	if orgUser == nil {
		return ErrUserNotPartOfOrganization
	}
	return nil
}
