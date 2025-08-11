package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

func testAuthMiddleware(next http.Handler, userID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockUserService struct {
	createUserFunc func(ctx context.Context, input repositories.CreateUserParams) (*models.User, error)
	getUserByIDFunc func(ctx context.Context, id string) (*models.User, error)	
}

func (m *mockUserService) CreateUser(ctx context.Context, input repositories.CreateUserParams) (*models.User, error) {
	return m.createUserFunc(ctx, input)
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return m.getUserByIDFunc(ctx, id)
}


func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("successful user creation", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input repositories.CreateUserParams) (*models.User, error) {
				return &models.User{
					ID:       "user-001",
					Username: input.Username,
					Email:    input.Email,
				}, nil
			},				
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		api.AssertStatus(t, res, http.StatusCreated)
		api.AssertJSONContentType(t, res)

		var user models.User
		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if user.ID != "user-001" {
			t.Errorf("expected user ID 'user-001', got %s", user.ID)
		}

		if user.Username != "John Doe" {
			t.Errorf("expected username 'John Doe', got %s", user.Username)
		}

		if user.Email != "john.doe@example.com" {
			t.Errorf("expected email 'john.doe@example.com', got %s", user.Email)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := &mockUserService{}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		api.AssertStatus(t, res, http.StatusBadRequest)
		api.AssertJSONContentType(t, res)
		api.AssertJSONErrorBody(t, res, "Invalid request body")
	})

	t.Run("service returns validation error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input repositories.CreateUserParams) (*models.User, error) {
				return nil, services.ErrInvalidInput
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.Code)
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input repositories.CreateUserParams) (*models.User, error) {
				return nil, services.ErrInternalServer
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		api.AssertStatus(t, res, http.StatusInternalServerError)
		api.AssertJSONContentType(t, res)
		api.AssertJSONErrorBody(t, res, "internal server error")
	})

	t.Run("service returns user with duplicate details error", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input repositories.CreateUserParams) (*models.User, error) {
				return nil, services.ErrUserWithDuplicateDetailsExists
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Post("/users", handler.CreateUser)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		api.AssertStatus(t, res, http.StatusConflict)
		api.AssertJSONContentType(t, res)
		api.AssertJSONErrorBody(t, res, "user with similar details already exists")
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	t.Run("successful user retrieval", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				if id == "user-001" {
					return &models.User{
						ID:       "user-001",
						Username: "John Doe",
						Email:    "john.doe@example.com",
					}, nil
				}
				return nil, nil
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Get("/users/{id}", handler.GetUserByID)


		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("service returns not found error", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, services.ErrUserNotFound
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Get("/users/{id}", handler.GetUserByID)

		req := httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, res.Code)
		}
	})

	t.Run("service returns internal server error", func(t *testing.T) {
		mockService := &mockUserService{
			getUserByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return nil, services.ErrInternalServer
			},
		}

		r := chi.NewRouter()
		handler := api.NewUserHandler(mockService)
		r.Get("/users/{id}", handler.GetUserByID)

		req := httptest.NewRequest(http.MethodGet, "/users/user-001", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
		}
	})
}


type mockOrganizationService struct {
	createOrganizationFunc func(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error)
	getOrganizationByIDFunc func(ctx context.Context, id string) (*models.Organization, error)
}

func (m *mockOrganizationService) CreateOrganization(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error) {
	return m.createOrganizationFunc(ctx, input)
}

func (m *mockOrganizationService) GetOrganizationByID(ctx context.Context, id string) (*models.Organization, error) {
	return m.getOrganizationByIDFunc(ctx, id)
}

func TestOrganizationHandler_CreateOrganization(t *testing.T) {
	mockService := &mockOrganizationService{
		createOrganizationFunc: func(ctx context.Context, input repositories.CreateOrganizationParams) (*models.Organization, error) {
			return &models.Organization{
				ID:   "org-001",
				Name: input.Name,
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
		getOrganizationByIDFunc: func(ctx context.Context, id string) (*models.Organization, error) {
			if id == "org-001" {
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

type mockAccessChecker struct {
	IsAdminFunc  func(ctx context.Context, orgID, userID string) (bool, error)
	IsMemberFunc func(ctx context.Context, orgID, userID string) (bool, error)
}

func (m *mockAccessChecker) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	return m.IsAdminFunc(ctx, orgID, userID)
}

func (m *mockAccessChecker) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	return m.IsMemberFunc(ctx, orgID, userID)
}

func TestAccessHandler_IsAdmin(t *testing.T) {
	mockChecker := &mockAccessChecker{
		IsAdminFunc: func(ctx context.Context, orgID, userID string) (bool, error) {
			if orgID == "org-001" && userID == "user-001" {
				return true, nil // Simulating that the user is an admin
			}
			return false, services.ErrUnauthorized
		},
	}

	r := chi.NewRouter()
	handler := api.NewAccessHandler(mockChecker)
	r.Get("/access/is-admin/{orgID}/{userID}", handler.IsAdmin)

	req := httptest.NewRequest(http.MethodGet, "/access/is-admin/org-001/user-001", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestAccessHandler_IsMember(t *testing.T) {
	mockChecker := &mockAccessChecker{
		IsMemberFunc: func(ctx context.Context, orgID, userID string) (bool, error) {
			if orgID == "org-001" && userID == "user-001" {
				return true, nil // Simulating that the user is a member
			}
			return false, services.ErrUnauthorized
		},
	}

	r := chi.NewRouter()
	handler := api.NewAccessHandler(mockChecker)
	r.Get("/access/member/{orgID}/{userID}", handler.IsMember)

	req := httptest.NewRequest(http.MethodGet, "/access/member/org-001/user-001", nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

type mockOrganizationUserService struct {
	createOrganizationUserFunc func(ctx context.Context, orgID string, params repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error)
	getUsersByOrganizationIDFunc func(ctx context.Context, orgID, userID string) ([]*models.UserWithRole, error)
	updateUserRoleFunc                func(ctx context.Context, orgID, userID string, role models.Role) error
	deleteUserFromOrganizationFunc func(ctx context.Context, orgID, userID, userIDToDelete string) error
}

func (m *mockOrganizationUserService) CreateOrganizationUser(ctx context.Context, orgID string, params repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
	return m.createOrganizationUserFunc(ctx, orgID, params)
}

func (m *mockOrganizationUserService) GetUsersByOrganizationID(ctx context.Context, orgID, userID string) ([]*models.UserWithRole, error) {
	return m.getUsersByOrganizationIDFunc(ctx, orgID, userID)
}

func (m *mockOrganizationUserService) UpdateUserRole(ctx context.Context, orgID, userID string, role models.Role) error {
	return m.updateUserRoleFunc(ctx, orgID, userID, role)
}
func (m *mockOrganizationUserService) DeleteUserFromOrganization(ctx context.Context, orgID, userID, userIDToDelete string) error {
	return m.deleteUserFromOrganizationFunc(ctx, orgID, userID, userIDToDelete)
}

func TestOrganizationUserHandler_CreateOrganizationUser(t *testing.T) {
	mockService := &mockOrganizationUserService{
		createOrganizationUserFunc: func(ctx context.Context, orgID string, params repositories.CreateOrganizationUserParams) (*models.OrganizationUser, error) {
			return &models.OrganizationUser{
				ID:     "org-user-001",
				OrgID:  orgID,
				UserID: params.UserID,
				Role:   params.Role,
			}, nil
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationUserHandler(mockService)
	r.Post("/organizations/{orgID}/users", handler.AddUserToOrganization)

	reqBody := `{"user_id": "user-001", "role": "member"}`
	req := httptest.NewRequest(http.MethodPost, "/organizations/org-001/users", bytes.NewBufferString(reqBody))
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.Code)
	}
}

func TestOrganizationUserHandler_GetUsersByOrganizationID(t *testing.T) {
	mockService := &mockOrganizationUserService{
		getUsersByOrganizationIDFunc: func(ctx context.Context, orgID, userID string) ([]*models.UserWithRole, error) {
			return []*models.UserWithRole{
				{User: models.User{
					ID: "user-001", 
					Username: "John Doe", 
					Email: "john.doe@example.com", 
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, Role: models.RoleMember},
			}, nil
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationUserHandler(mockService)
	r.Get("/organizations/{orgID}/users", handler.GetUsersByOrganizationID)

	req := httptest.NewRequest(http.MethodGet, "/organizations/org-001/users", nil)
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
		deleteUserFromOrganizationFunc: func(ctx context.Context, reqActingUserID, reqOrgID, reqUserIDToDelete string) error {
			if orgID != reqOrgID {
				t.Errorf("expected orgID %s, got %s", orgID, reqOrgID)
			}
			if actingUserID != reqActingUserID {
				t.Errorf("expected userID %s, got %s", actingUserID, reqActingUserID)
			}
			if userToDeleteID != reqUserIDToDelete {
				t.Errorf("expected userIDToDelete %s, got %s", userToDeleteID, reqUserIDToDelete)
			}
			return nil
		},
	}

	r := chi.NewRouter()
	handler := api.NewOrganizationUserHandler(mockService)

	authedHandler := testAuthMiddleware(http.HandlerFunc(handler.DeleteUserFromOrganization), actingUserID)

	r.Method(http.MethodDelete, "/organizations/{orgID}/users/{userID}", authedHandler)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/organizations/%s/users/%s", orgID, userToDeleteID), nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, res.Code)
	}
}