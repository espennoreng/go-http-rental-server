package services

import (
	"context"
	"log/slog"

	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
	log      *slog.Logger
}

func NewUserService(userRepo repositories.UserRepository, log *slog.Logger) *userService {
	return &userService{
		userRepo: userRepo,
		log:      log.With("component", "user_service"),
	}
}

var _ UserService = (*userService)(nil)

func (s *userService) CreateUser(ctx context.Context, params CreateUserParams) (*models.User, error) {
	log := s.log.With(
		slog.String("username", params.Username),
		slog.String("email", params.Email),
	)

	if params.Username == "" {
		log.Warn("CreateUser called with empty username")
		return nil, ErrInvalidInput
	}
	if params.Email == "" {
		log.Warn("CreateUser called with empty email")
		return nil, ErrInvalidInput
	}

	log.Info("Creating user")

	newUser, err := s.userRepo.Create(ctx, &repositories.CreateUserParams{
		Username: params.Username,
		Email:    params.Email,
	})

	if err != nil {
		log.Error("Internal error while creating user", "error", err)
		return nil, ErrInternalServer
	}

	log.Info("User created successfully", "user_id", newUser.ID)

	return newUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, params GetUserByIDParams) (*models.User, error) {
	log := s.log.With("user_id", params.ID)

	if params.ID == "" {
		log.Warn("GetUserByID called with empty ID")
		return nil, ErrInvalidInput
	}

	log.Info("Retrieving user by ID")

	user, err := s.userRepo.GetByID(ctx, params.ID)
	if err != nil {
		log.Error("Error retrieving user by ID", "error", err)
		return nil, ErrInternalServer
	}

	log.Info("User retrieved successfully", "user_id", user.ID)

	return user, nil
}

func (s *userService) FindOrCreateByGoogleID(ctx context.Context, googleID, email string) (*models.User, error) {
	log := s.log.With("google_id", googleID, "email", email)

	if googleID == "" {
		log.Warn("FindOrCreateByGoogleID called with empty Google ID")
		return nil, ErrInvalidInput
	}

	if email == "" {
		log.Warn("FindOrCreateByGoogleID called with empty email")
		return nil, ErrInvalidInput
	}

	log.Info("Finding or creating user by Google ID")

	user, err := s.userRepo.FindOrCreateByGoogleID(ctx, googleID, email)
	if err != nil {
		log.Error("Error finding or creating user by Google ID", "error", err)
		return nil, ErrInternalServer
	}

	log.Info("User found or created successfully", "user_id", user.ID)
	
	return user, nil
}