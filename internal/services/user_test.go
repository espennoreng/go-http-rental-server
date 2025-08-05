package services_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/services"
)

type mockUserRepository struct {
	createFunc  func(ctx context.Context, user *models.User) error
	getByIDFunc func(ctx context.Context, id string) (*models.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	return m.getByIDFunc(ctx, id)
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("create user successfully", func(t *testing.T) {
		repo := &mockUserRepository{
			createFunc: func(ctx context.Context, user *models.User) error {
				return nil // Simulate successful creation
			},
			getByIDFunc: func(ctx context.Context, id string) (*models.User, error) {
				return &models.User{ID: id, Username: "John Doe", Email: "john.doe@example.com"}, nil
			},
		}

		service := services.NewUserService(repo)

		user := models.CreateUserInput{Username: "John Doe", Email: "john.doe@example.com"}
		createdUser, err := service.CreateUser(context.Background(), user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if createdUser.ID == "" {
			t.Fatal("expected user ID to be set")
		}
	})
}
