package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

var _ repositories.UserRepository = (*userRepository)(nil)

type userRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

func NewUserRepository() *userRepository {
	return &userRepository{
		users: make(map[string]*models.User),
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}

	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}

	return user, nil
}
