package inmemory_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories/inmemory"
)

func TestUserCreate(t *testing.T) {
	repo := inmemory.NewUserRepository()

	user := models.User{
		ID:       "user-001",
		Username: "John Doe",
		Email:    "test@example.com",
	}

	err := repo.Create(context.Background(), &user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedUser, err := repo.GetByID(context.Background(), "user-001")
	if err != nil {
		t.Fatalf("expected to find user, got error: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Username != user.Username {
		t.Errorf("expected user name %s, got %s", user.Username, retrievedUser.Username)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("expected user email %s, got %s", user.Email, retrievedUser.Email)
	}
}
