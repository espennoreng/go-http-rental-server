package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type organizationService struct {
	organizationRepo repositories.OrganizationRepository
}

func NewOrganizationService(organizationRepo repositories.OrganizationRepository) *organizationService {
	return &organizationService{
		organizationRepo: organizationRepo,
	}
}

func (s *organizationService) CreateOrganization(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error) {
	if input.Name == "" {
		return nil, ErrInvalidInput
	}
	if input.CreatedBy == "" {
		return nil, ErrInvalidInput
	}

	newOrganization, err := s.organizationRepo.Create(ctx, &repositories.CreateOrganizationParams{
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

	organization, err := s.organizationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInternalServer
	}

	return organization, nil
}
