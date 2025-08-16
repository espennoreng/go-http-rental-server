package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type organizationService struct {
	orgRepo repositories.OrganizationRepository
	log *slog.Logger
}

func NewOrganizationService(orgRepo repositories.OrganizationRepository, log *slog.Logger) *organizationService {
	return &organizationService{
		orgRepo: orgRepo,
		log: log.With(slog.String("component", "organization_service")),
	}
}

var _ OrganizationService = (*organizationService)(nil)

func (s *organizationService) CreateOrganization(ctx context.Context, params CreateOrganizationParams) (*models.Organization, error) {
	log := s.log.With(
		slog.String("created_by", params.CreatedBy),
		slog.String("org_name", params.Name),
	)

	if params.Name == "" {
		log.Error("Invalid input: organization name is required")
		return nil, ErrInvalidInput
	}
	if params.CreatedBy == "" {
		log.Error("Invalid input: created by is required")
		return nil, ErrInvalidInput
	}

	log.Info("Creating new organization")

	newOrganization, err := s.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
		Name:      params.Name,
		CreatedBy: params.CreatedBy,
	})

	if err != nil {
		if errors.Is(err, repositories.ErrConflict) {
			log.Warn("Organization already exists", slog.Any("error", err))
			return nil, ErrOrganizationWithDuplicateDetailsExists
		}
		log.Error("Failed to create organization", slog.Any("error", err))
		return nil, ErrInternalServer
	}

	log.Info("Organization created successfully", slog.String("org_id", newOrganization.ID))

	return newOrganization, nil
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, params GetOrganizationByIDParams) (*models.Organization, error) {
	log := s.log.With(slog.String("org_id", params.ID))

	if params.ID == "" {
		log.Error("Invalid input: organization ID is required")
		return nil, ErrInvalidInput
	}

	log.Info("Retrieving organization by ID")

	organization, err := s.orgRepo.GetByID(ctx, params.ID)
	if err != nil {
		log.Error("Failed to retrieve organization", slog.Any("error", err))
		return nil, ErrInternalServer
	}

	log.Info("Organization retrieved successfully", slog.String("org_id", organization.ID))

	return organization, nil
}
