package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.CreateUser(r.Context(), services.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, services.ErrUserWithDuplicateDetailsExists) {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := NewUserResponse(user)

	respondJSON(w, http.StatusCreated, response)
}

func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	_, err := auth.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	id := chi.URLParam(r, "id")

	user, err := h.userService.GetUserByID(r.Context(), services.GetUserByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
