package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/api"
	"github.com/espennoreng/go-http-rental-server/internal/models"
)

type mockUserService struct {
	createUserFunc func(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	getUserByIDFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *mockUserService) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	return m.createUserFunc(ctx, input)
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("create user successfully", func(t *testing.T) {
		mockService := &mockUserService{
			createUserFunc: func(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
				return &models.User{
					ID:       "user-001",
					Username: input.Username,
					Email:    input.Email,
				}, nil
			},				
		}

		handler := api.NewUserHandler(mockService)

		reqBody := `{"username": "John Doe", "email": "john.doe@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		res := httptest.NewRecorder()

		handler.CreateUser(res, req)

		if res.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, res.Code)
		}

		var createdUser models.User
		if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if createdUser.Username != "John Doe" {
			t.Errorf("expected username 'John Doe', got '%s'", createdUser.Username)
		}

		if createdUser.Email != "john.doe@example.com" {
			t.Errorf("expected email 'john.doe@example.com', got '%s'", createdUser.Email)
		}
	})
}