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

func (s *organizationService) CreateOrganization(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error) {
	if input.Name == "" {
		return nil, ErrInvalidInput
	}
	if input.CreatedBy == "" {
		return nil, ErrInvalidInput
	}

	newOrganization, err := s.orgRepo.Create(ctx, &repositories.CreateOrganizationParams{
		Name:      input.Name,
		CreatedBy: input.CreatedBy,
	})

	if err != nil {
		return nil, ErrInternalServer
	}

	return newOrganization, nil
}

func (s *organizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	organization, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInternalServer
	}

	return organization, nil
}
