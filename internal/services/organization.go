package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type organizationService struct {
	orgRepo repositories.OrganizationRepository
}

func NewOrganizationService(orgRepo repositories.OrganizationRepository) *organizationService {
	return &organizationService{
		orgRepo: orgRepo,
	}
}

var _ OrganizationService = (*organizationService)(nil)

func (s *organizationService) CreateOrganization(ctx context.Context, params CreateOrganizationParams) (*models.Organization, error) {
	if params.Name == "" {
		return nil, ErrInvalidInput
	}
	if params.CreatedBy == "" {
		return nil, ErrInvalidInput
	}

	newOrganization, err := s.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
		Name:      params.Name,
		CreatedBy: params.CreatedBy,
	})

	if err != nil {
		return nil, ErrInternalServer
	}

	return newOrganization, nil
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, params GetOrganizationByIDParams) (*models.Organization, error) {
	if params.ID == "" {
		return nil, ErrInvalidInput
	}

	organization, err := s.orgRepo.GetByID(ctx, params.ID)
	if err != nil {
		return nil, ErrInternalServer
	}

	return organization, nil
}
