package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
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
	userID := "test-user-id"
	newOrgID := "org-001"
	newOrgName := "New Organization"

	mockService := &mockOrganizationService{
		createOrganizationFunc: func(ctx context.Context, params services.CreateOrganizationParams) (*models.Organization, error) {
			return &models.Organization{
				ID:   newOrgID,
				Name: params.Name,
			}, nil
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationHandler(mockService)

	authedHandler := middleware.TestAuthMiddleware(http.HandlerFunc(handler.CreateOrganization), auth.Identity{UserID: userID})
	r.Method(http.MethodPost, "/organizations", authedHandler)

	reqBody := fmt.Sprintf(`{"name": "%s"}`, newOrgName)
	req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBufferString(reqBody))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
	var response api.OrganizationResponse
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, newOrgID, response.ID)
	assert.Equal(t, newOrgName, response.Name)
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
