package api_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type mockOrganizationService struct {
	createOrganizationFunc  func(ctx context.Context, params services.CreateOrganizationParams) (*models.Organization, error)
	getOrganizationByIDFunc func(ctx context.Context, params services.GetOrganizationByIDParams) (*models.Organization, error)
}

func (m *mockOrganizationService) CreateOrganization(ctx context.Context, params services.CreateOrganizationParams) (*models.Organization, error) {
	return m.createOrganizationFunc(ctx, params)
}

func (m *mockOrganizationService) GetOrganizationByID(ctx context.Context, params services.GetOrganizationByIDParams) (*models.Organization, error) {
	return m.getOrganizationByIDFunc(ctx, params)
}

func TestOrganizationHandler_CreateOrganization(t *testing.T) {
	mockService := &mockOrganizationService{
		createOrganizationFunc: func(ctx context.Context, params services.CreateOrganizationParams) (*models.Organization, error) {
			return &models.Organization{
				ID:   "org-001",
				Name: params.Name,
			}, nil
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationHandler(mockService)
	r.Post("/organizations", handler.CreateOrganization)

	reqBody := `{"name": "New Organization"}`
	req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBufferString(reqBody))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.Code)
	}
}

func TestOrganizationHandler_GetOrganizationByID(t *testing.T) {
	mockService := &mockOrganizationService{
		getOrganizationByIDFunc: func(ctx context.Context, params services.GetOrganizationByIDParams) (*models.Organization, error) {
			if params.ID == "org-001" {
				return &models.Organization{
					ID:   "org-001",
					Name: "Existing Organization",
				}, nil
			}
			return nil, services.ErrOrganizationNotFound
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationHandler(mockService)
	r.Get("/organizations/{id}", handler.GetOrganizationByID)

	req := httptest.NewRequest(http.MethodGet, "/organizations/org-001", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}
