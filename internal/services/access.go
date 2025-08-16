package services

import (
	"context"
	"log/slog"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
)

// AccessService provides reusable methods to check user access rights.
type accessService struct {
	orgUserRepo repositories.OrganizationUserRepository
	log         *slog.Logger
}

// NewAccessService creates a new instance of the access checker.
func NewAccessService(orgUserRepo repositories.OrganizationUserRepository, log *slog.Logger) *accessService {
	return &accessService{
		orgUserRepo: orgUserRepo,
		log:        log.With(slog.String("component", "access_service")),
	}
}

var _ AccessService = (*accessService)(nil)

// IsAdmin checks if a user has 'admin' privileges in an organization.
// It returns ErrInvalidInput if the UUIDs are malformed,
// ErrUnauthorized if the user is not an admin, or a database error.
func (s *accessService) IsAdmin(ctx context.Context, params OrgAccessParams) error {
	
	log := s.log.With(
		slog.String("org_id", params.OrgID),
		slog.String("user_id", params.UserID),
	)
	if err := uuid.Validate(params.OrgID); err != nil || params.OrgID == "" {
		log.Error("Invalid input: organization ID is required")
		return ErrInvalidInput
	}
	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
		log.Error("Invalid input: user ID is required")
		return ErrInvalidInput
	}

	log.Info("Checking if user is admin in organization")

	orgUser, err := s.orgUserRepo.GetByID(ctx, params.OrgID, params.UserID)
	if err != nil {
		// Log the underlying error for debugging purposes
		log.Error("Failed to retrieve organization user", slog.Any("error", err))
		return ErrInternalServer
	}
	if orgUser == nil {
		log.Warn("User is not part of the organization")
		return ErrUserNotPartOfOrganization
	}
	if orgUser.Role != models.RoleAdmin {
		log.Warn("User is not an admin", slog.String("role", string(orgUser.Role)))
		return ErrUnauthorized
	}
	log.Info("User is an admin in the organization")

	return nil
}

// IsMember checks if a user is a member of an organization.
// It returns ErrInvalidInput if the UUIDs are malformed,
// ErrUnauthorized if the user is not a member, or a database error.
func (s *accessService) IsMember(ctx context.Context, params OrgAccessParams) error {
	log := s.log.With(
		slog.String("org_id", params.OrgID),
		slog.String("user_id", params.UserID),
	)


	if err := uuid.Validate(params.OrgID); err != nil || params.OrgID == "" {
		log.Error("Invalid input: organization ID is required")
		return ErrInvalidInput
	}
	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
		log.Error("Invalid input: user ID is required")
		return ErrInvalidInput
	}

	log.Info("Checking if user is a member of the organization")

	orgUser, err := s.orgUserRepo.GetByID(ctx, params.OrgID, params.UserID)
	if err != nil {
		// Log the underlying error for debugging purposes
		log.Error("Failed to retrieve organization user", slog.Any("error", err))
		return ErrInternalServer
	}
	if orgUser == nil {
		log.Warn("User is not part of the organization")
		return ErrUserNotPartOfOrganization
	}

	log.Info("User is a member of the organization")

	return nil
}
