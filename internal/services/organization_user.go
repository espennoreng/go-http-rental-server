package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/google/uuid"
)

type organizationUserService struct {
	orgUserRepo   repositories.OrganizationUserRepository
	accessService AccessService
	log           *slog.Logger
}

// NewOrganizationUserService initializes a new organizationUserService.
func NewOrganizationUserService(orgUserRepo repositories.OrganizationUserRepository, accessService AccessService) *organizationUserService {
	return &organizationUserService{
		orgUserRepo:   orgUserRepo,
		accessService: accessService,
		log:           slog.With(slog.String("component", "organization_user_service")),
	}
}

var _ OrganizationUserService = (*organizationUserService)(nil)

// CreateOrganizationUser handles the creation of a new organization-user relationship.
func (s *organizationUserService) CreateOrganizationUser(ctx context.Context, params CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	log := s.log.With(
		slog.String("acting_user_id", params.ActingUserID),
		slog.String("org_id", params.OrgID),
		slog.String("user_id", params.UserID),
		slog.String("role", string(params.Role)),
	)
	log.Info("Creating new organization user")

	err := s.accessService.IsAdmin(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})

	if err != nil {
		log.Warn("Failed to create organization user, probably due to insufficient permissions", slog.Any("error", err))
		return nil, err
	}

	if !models.ValidRoles[params.Role] {
		log.Error("Invalid role provided for organization user")
		return nil, ErrInvalidInput
	}

	if err := uuid.Validate(params.UserID); err != nil || params.UserID == "" {
		log.Error("Invalid user ID provided for organization user")
		return nil, ErrInvalidInput
	}

	if err := uuid.Validate(params.OrgID); err != nil || params.OrgID == "" {
		log.Error("Invalid organization ID provided for organization user")
		return nil, ErrInvalidInput
	}

	log.Info("Creating new organization user")

	newOrgUser, err := s.orgUserRepo.Create(ctx, &repositories.CreateOrganizationUserParams{
		OrgID:  params.OrgID,
		UserID: params.UserID,
		Role:   params.Role,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			log.Warn("Organization user already exists", slog.Any("error", err))
			return nil, ErrUserAlreadyHasARoleInOrganization
		}
		log.Error("Failed to create organization user", slog.Any("error", err))
		return nil, ErrInternalServer
	}

	log.Info("Organization user created successfully", slog.String("org_user_id", newOrgUser.ID))

	return newOrgUser, nil
}

// GetUsersByOrganizationID retrieves all users within an organization.
func (s *organizationUserService) GetUsersByOrganizationID(ctx context.Context, params GetUsersByOrganizationIDParams) ([]*models.UserWithRole, error) {
	log  := s.log.With(slog.String("org_id", params.OrgID), slog.String("acting_user_id", params.ActingUserID))

	err := s.accessService.IsMember(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})

	if err != nil {
		log.Warn("Failed to retrieve users by organization ID, probably due to insufficient permissions", slog.Any("error", err))
		return nil, err
	}

	log.Info("Fetching users for organization")

	users, err := s.orgUserRepo.GetUsersByOrganizationID(ctx, params.OrgID)
	if err != nil {
		log.Error("Failed to fetch users for organization", slog.Any("error", err))
		return nil, ErrInternalServer
	}

	log.Info("Users retrieved successfully for organization", slog.Int("user_count", len(users)))
	return users, nil
}

// UpdateRole updates a user's role within an organization.
func (s *organizationUserService) UpdateUserRole(ctx context.Context, params UpdateUserRoleParams) error {
	log := s.log.With(
		slog.String("acting_user_id", params.ActingUserID),
		slog.String("org_id", params.OrgID),
		slog.String("user_id", params.UserID),
		slog.String("role", string(params.Role)),
	)

	log.Info("Updating user role in organization")

	err := s.accessService.IsAdmin(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})

	if err != nil {
		log.Warn("Failed to update user role, probably due to insufficient permissions", slog.Any("error", err))
		return err
	}

	if !models.ValidRoles[params.Role] {
		log.Warn("Invalid role provided for organization user")
		return ErrInvalidInput
	}

	log.Info("Updating user role in organization")

	err = s.orgUserRepo.UpdateRole(ctx, params.OrgID, params.UserID, params.Role)
	if err != nil {
		log.Error("Failed to update user role in organization", slog.Any("error", err))
		return ErrInternalServer
	}

	log.Info("User role updated successfully in organization")

	return nil
}

// DeleteUserFromOrganization removes a user from an organization.
func (s *organizationUserService) DeleteUserFromOrganization(ctx context.Context, params DeleteOrganizationUserParams) error {
	log := s.log.With(
		slog.String("acting_user_id", params.ActingUserID),
		slog.String("org_id", params.OrgID),
		slog.String("user_id_to_delete", params.UserIDToDelete),
	)

	log.Info("Deleting user from organization")

	err := s.accessService.IsAdmin(ctx, OrgAccessParams{
		OrgID:  params.OrgID,
		UserID: params.ActingUserID,
	})
	
	if err != nil {
		log.Warn("Failed to delete user from organization, probably due to insufficient permissions", slog.Any("error", err))
		return err
	}

	err = s.orgUserRepo.Delete(ctx, params.OrgID, params.UserIDToDelete)
	if err != nil {
		log.Error("Failed to delete user from organization", slog.Any("error", err))
		return ErrInternalServer
	}

	log.Info("User deleted successfully from organization", slog.String("user_id_deleted", params.UserIDToDelete))

	return nil
}