package api_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type mockOrganizationUserService struct {
	createOrganizationUserFunc     func(ctx context.Context, params services.CreateOrganizationUserParams) (*models.OrganizationUser, error)
	getUsersByOrganizationIDFunc   func(ctx context.Context, params services.GetUsersByOrganizationIDParams) ([]*models.UserWithRole, error)
	updateUserRoleFunc             func(ctx context.Context, params services.UpdateUserRoleParams) error
	deleteUserFromOrganizationFunc func(ctx context.Context, params services.DeleteOrganizationUserParams) error
}

func (m *mockOrganizationUserService) CreateOrganizationUser(ctx context.Context, params services.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	return m.createOrganizationUserFunc(ctx, params)
}

func (m *mockOrganizationUserService) GetUsersByOrganizationID(ctx context.Context, params services.GetUsersByOrganizationIDParams) ([]*models.UserWithRole, error) {
	return m.getUsersByOrganizationIDFunc(ctx, params)
}

func (m *mockOrganizationUserService) UpdateUserRole(ctx context.Context, params services.UpdateUserRoleParams) error {
	return m.updateUserRoleFunc(ctx, params)
}
func (m *mockOrganizationUserService) DeleteUserFromOrganization(ctx context.Context, params services.DeleteOrganizationUserParams) error {
	return m.deleteUserFromOrganizationFunc(ctx, params)
}

type mockAccessService struct {
	isAdminFunc  func(ctx context.Context, params services.OrgAccessParams) error
	isMemberFunc func(ctx context.Context, params services.OrgAccessParams) error
}

func (m *mockAccessService) IsAdmin(ctx context.Context, params services.OrgAccessParams) error {
	return m.isAdminFunc(ctx, params)
}
func (m *mockAccessService) IsMember(ctx context.Context, params services.OrgAccessParams) error {
	return m.isMemberFunc(ctx, params)
}

func TestOrganizationUserHandler_AddUserToOrganization(t *testing.T) {
	// Define the user performing the action
	const actingUserID = "admin-user-007"
	const userID = "member-user-001"
	const orgID = "org-001"
	const role = models.RoleMember

	mockService := &mockOrganizationUserService{
		createOrganizationUserFunc: func(ctx context.Context, params services.CreateOrganizationUserParams) (*models.OrganizationUser, error) {

			if params.OrgID != orgID {
				t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
			}
			if params.ActingUserID != actingUserID {
				t.Errorf("expected actingUserID %s, got %s", actingUserID, params.ActingUserID)
			}
			if params.UserID != userID {
				t.Errorf("expected userID %s, got %s", userID, params.UserID)
			}
			if params.Role != role {
				t.Errorf("expected role %s, got %s", role, params.Role)
			}
			// Simulating a successful creation of an organization user
			return &models.OrganizationUser{
				ID:     "org-user-001",
				OrgID:  params.OrgID,
				UserID: params.UserID,
				Role:   params.Role,
			}, nil
		},
	}

	mockAccessService := &mockAccessService{
		isAdminFunc: func(ctx context.Context, params services.OrgAccessParams) error {
			if params.UserID != actingUserID {
				t.Errorf("expected actingUserID %s, got %s", actingUserID, params.UserID)
			}
			if params.OrgID != orgID {
				t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
			}
			return nil // Simulating that the acting user is an admin
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationUserHandler(mockService)
	accessMiddleware := middleware.NewAccessMiddleware(mockAccessService)

	adminProtectedHandler := accessMiddleware.RequireAdmin(http.HandlerFunc(handler.AddUserToOrganization))
	authedHandler := middleware.NewTestAuthMiddleware(adminProtectedHandler, auth.Identity{UserID: actingUserID})

	r.Method(http.MethodPost, "/organizations/{orgID}/users", authedHandler)

	reqBody := fmt.Sprintf(`{"user_id": "%s", "role": "%s"}`, userID, role)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/organizations/%s/users", orgID), bytes.NewBufferString(reqBody))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.Code)
	}
}

func TestOrganizationUserHandler_GetUsersByOrganizationID(t *testing.T) {
	// Define the organization ID and user ID for the request
	const orgID = "org-001"
	const userID = "user-001"

	mockService := &mockOrganizationUserService{
		getUsersByOrganizationIDFunc: func(ctx context.Context, params services.GetUsersByOrganizationIDParams) ([]*models.UserWithRole, error) {
			if params.OrgID != orgID {
				t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
			}
			if params.ActingUserID != userID {
				t.Errorf("expected actingUserID %s, got %s", userID, params.ActingUserID)
			}
			// Simulating a successful response with one user
			return []*models.UserWithRole{
				{User: models.User{
					ID:        "user-001",
					Username:  "John Doe",
					Email:     "john.doe@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, Role: models.RoleMember},
			}, nil
		},
	}

	mockAccessService := &mockAccessService{
		isMemberFunc: func(ctx context.Context, params services.OrgAccessParams) error {
			if params.UserID != userID {
				t.Errorf("expected actingUserID %s, got %s", userID, params.UserID)
			}
			if params.OrgID != orgID {
				t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
			}
			return nil // Simulating that the user is a member of the organization
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationUserHandler(mockService)
	accessMiddleware := middleware.NewAccessMiddleware(mockAccessService)

	memberProtectedHandler := accessMiddleware.RequireMember(http.HandlerFunc(handler.GetUsersByOrganizationID))
	authedHandler := middleware.NewTestAuthMiddleware(memberProtectedHandler, auth.Identity{UserID: userID})

	r.Method(http.MethodGet, "/organizations/{orgID}/users", authedHandler)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/organizations/%s/users", orgID), nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestOrganizationUserHandler_DeleteOrganizationUser(t *testing.T) {
	// Define the user performing the action
	const actingUserID = "admin-user-007"
	const userToDeleteID = "member-user-001"
	const orgID = "org-001"

	mockService := &mockOrganizationUserService{
		deleteUserFromOrganizationFunc: func(ctx context.Context, params services.DeleteOrganizationUserParams) error {
			if params.OrgID != orgID {
				t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
			}
			if params.ActingUserID != actingUserID {
				t.Errorf("expected actingUserID %s, got %s", actingUserID, params.ActingUserID)
			}
			if params.UserIDToDelete != userToDeleteID {
				t.Errorf("expected userIDToDelete %s, got %s", userToDeleteID, params.UserIDToDelete)
			}
			// Simulating a successful deletion
			return nil
		},
	}

	mockAccessService := &mockAccessService{
		isAdminFunc: func(ctx context.Context, params services.OrgAccessParams) error {
			if params.UserID != actingUserID {
				t.Errorf("expected actingUserID %s, got %s", actingUserID, params.UserID)
			}
			if params.OrgID != orgID {
				t.Errorf("expected orgID %s, got %s", orgID, params.OrgID)
			}
			return nil
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationUserHandler(mockService)
	accessMiddleware := middleware.NewAccessMiddleware(mockAccessService)

	adminProtectedHandler := accessMiddleware.RequireAdmin(http.HandlerFunc(handler.DeleteUserFromOrganization))
	authedHandler := middleware.NewTestAuthMiddleware(adminProtectedHandler, auth.Identity{UserID: actingUserID})

	r.Method(http.MethodDelete, "/organizations/{orgID}/users/{userID}", authedHandler)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/organizations/%s/users/%s", orgID, userToDeleteID), nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, res.Code)
	}
}
