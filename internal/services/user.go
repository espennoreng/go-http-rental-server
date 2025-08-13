package services

import (
	"context"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, params CreateUserParams) (*models.User, error) {
	if params.Username == "" {
		return nil, ErrInvalidInput
	}
	if params.Email == "" {
		return nil, ErrInvalidInput
	}

	newUser, err := s.userRepo.Create(ctx, &repositories.CreateUserParams{
		Username: params.Username,
		Email:    params.Email,
	})

	if err != nil {
		return nil, ErrInternalServer
	}

	return newUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, params GetUserByIDParams) (*models.User, error) {
	if params.ID == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.userRepo.GetByID(ctx, params.ID)
	if err != nil {
		return nil, ErrInternalServer
	}

	return user, nil
}
