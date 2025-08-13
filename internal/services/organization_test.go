package services_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/stretchr/testify/assert"
)

type mockOrganizationRepository struct {
	CreateOrganizationFunc  func(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error)
	GetOrganizationByIDFunc func(ctx context.Context, id string) (*models.Organization, error)
}

func (m *mockOrganizationRepository) Create(ctx context.Context, params *repositories.CreateOrganizationParams) (*models.Organization, error) {
	return m.CreateOrganizationFunc(ctx, *params)
}

func (m *mockOrganizationRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	return m.GetOrganizationByIDFunc(ctx, id)
}

func TestOrganizationService_CreateOrganization(t *testing.T) {
	mockRepo := &mockOrganizationRepository{
		CreateOrganizationFunc: func(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error) {
			return &models.Organization{
				ID:   "1",
				Name: input.Name,
			}, nil
		},
	}

	service := services.NewOrganizationService(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		org, err := service.CreateOrganization(context.Background(), services.CreateOrganizationParams{
			Name:      "Test Organization",
			CreatedBy: "user-001",
		})
		assert.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, "Test Organization", org.Name)
	})
}
func TestOrganizationService_GetOrganizationByID(t *testing.T) {
	mockRepo := &mockOrganizationRepository{
		GetOrganizationByIDFunc: func(ctx context.Context, id string) (*models.Organization, error) {
			return &models.Organization{ID: "1", Name: "Test Organization"}, nil
		},
	}

	service := services.NewOrganizationService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		org, err := service.GetOrganizationByID(context.Background(), services.GetOrganizationByIDParams{ID: "1"})
		assert.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, "Test Organization", org.Name)
	})
}
