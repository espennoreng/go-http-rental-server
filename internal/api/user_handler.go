package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userService services.UserService
	log *slog.Logger
}

func NewUserHandler(userService services.UserService, log *slog.Logger) *userHandler {
	return &userHandler{
		userService: userService,
		log:          log.With(slog.String("component", "user_handler")),
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("Failed to decode request body", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.log.Warn("Validation failed for user creation", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	log := h.log.With(slog.String("username", input.Username), slog.String("email", input.Email))
	log.Info("Creating new user", slog.String("username", input.Username), slog.String("email", input.Email))

	user, err := h.userService.CreateUser(r.Context(), services.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			log.Warn("Invalid input for user creation", slog.Any("error", err))
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, services.ErrUserWithDuplicateDetailsExists) {
			log.Warn("User with duplicate details exists", slog.Any("error", err))
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		log.Error("Failed to create user", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	log.Info("User created successfully", slog.String("user_id", user.ID))

	response := NewUserResponse(user)

	respondJSON(w, http.StatusCreated, response)
}

func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	_, err := auth.FromContext(r.Context())
	if err != nil {
		h.log.Error("User identity not found in request context", slog.Any("error", err))
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.log.Warn("User ID is required for fetching user details", slog.String("id", id))
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	log := h.log.With(slog.String("user_id", id))
	log.Info("Fetching user by ID", slog.String("user_id", id))

	user, err := h.userService.GetUserByID(r.Context(), services.GetUserByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			log.Warn("User not found", slog.Any("error", err))
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Error("Internal error while fetching user by ID", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	log.Info("User retrieved successfully", slog.String("user_id", user.ID))

	respondJSON(w, http.StatusOK, user)
}
