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
	t.Run("successful organization creation", func(t *testing.T) {

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
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockService := &mockOrganizationService{}
		r := chi.NewRouter()
		handler := api.NewOrganizationHandler(mockService)
		r.Post("/organizations", handler.CreateOrganization)

		reqBody := fmt.Sprintf(`{"name": "%s"}`, newOrgName)
		req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func TestOrganizationHandler_GetOrganizationByID(t *testing.T) {
	actingUserID := "test-acting-user-id"
	orgID := "org-001"
	orgName := "Test Organization"

	t.Run("successful organization retrieval", func(t *testing.T) {
		mockService := &mockOrganizationService{
			getOrganizationByIDFunc: func(ctx context.Context, params services.GetOrganizationByIDParams) (*models.Organization, error) {
				if params.ID == orgID {
					return &models.Organization{
						ID:   orgID,
						Name: orgName,
					}, nil
				}
				return nil, services.ErrOrganizationNotFound
			},
		}

		mockAccessService := &mockAccessService{
			isMemberFunc: func(ctx context.Context, params services.OrgAccessParams) error {
				if params.UserID != actingUserID {
					t.Errorf("expected actingUserID %s, got %s", actingUserID, params.UserID)
				}
				if params.OrgID != orgID {
					t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
				}
				return nil // Simulating that the acting user is a member
			},
		}

		r := chi.NewRouter()
		handler := api.NewOrganizationHandler(mockService)
		accessMiddleware := middleware.NewAccessMiddleware(mockAccessService)

		memberProtectedHandler := accessMiddleware.RequireMember(http.HandlerFunc(handler.GetOrganizationByID))
		authedHandler := middleware.TestAuthMiddleware(memberProtectedHandler, auth.Identity{UserID: actingUserID})

		r.Method(http.MethodGet, "/organizations/{orgID}", authedHandler)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/organizations/%s", orgID), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response api.OrganizationResponse
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, orgID, response.ID)
		assert.Equal(t, orgName, response.Name)
	})

	t.Run("Unauthorized access", func(t *testing.T) {
		mockService := &mockOrganizationService{}
		r := chi.NewRouter()
		handler := api.NewOrganizationHandler(mockService)

		r.Get("/organizations/{orgID}", handler.GetOrganizationByID)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/organizations/%s", orgID), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}
