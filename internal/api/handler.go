package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/middleware"
	"github.com/espennoreng/go-http-rental-server/internal/models"
	"github.com/espennoreng/go-http-rental-server/internal/repositories"
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
	var input repositories.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), input)
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

	respondJSON(w, http.StatusCreated, user)
}

func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.userService.GetUserByID(r.Context(), id)
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

type organizationHandler struct {
	organizationService services.OrganizationService
}

func NewOrganizationHandler(organizationService services.OrganizationService) *organizationHandler {
	return &organizationHandler{
		organizationService: organizationService,
	}
}

func (h *organizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var input repositories.CreateOrganizationParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	org, err := h.organizationService.CreateOrganization(r.Context(), input)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, services.ErrOrganizationWithDuplicateDetailsExists) {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, org)
}

func (h *organizationHandler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	org, err := h.organizationService.GetOrganizationByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, org)
}

type accessHandler struct {
	accessService services.AccessService
}

func NewAccessHandler(accessService services.AccessService) *accessHandler {
	return &accessHandler{
		accessService: accessService,
	}
}

func (h *accessHandler) IsAdmin(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	userID := chi.URLParam(r, "userID")

	isAdmin, err := h.accessService.IsAdmin(r.Context(), orgID, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"is_admin": isAdmin})
}

func (h *accessHandler) IsMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "orgID")
	userID := chi.URLParam(r, "userID")

	isMember, err := h.accessService.IsMember(r.Context(), orgID, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"is_member": isMember})
}

type organizationUserHandler struct {
	organizationUserService services.OrganizationUserService
}

func NewOrganizationUserHandler(organizationUserService services.OrganizationUserService) *organizationUserHandler {
	return &organizationUserHandler{
		organizationUserService: organizationUserService,
	}
}

func (h *organizationUserHandler) AddUserToOrganization(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	orgID := chi.URLParam(r, "orgID")

	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	var input struct {
		UserID string      `json:"user_id"`
		Role   models.Role `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newOrgUser, err := h.organizationUserService.CreateOrganizationUser(context.Background(), services.CreateOrganizationUserParams{
		ActingUserID: userID,
		OrgID:        orgID,
		UserID:       input.UserID,
		Role:         input.Role,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, newOrgUser)
}

func (h *organizationUserHandler) GetUsersByOrganizationID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}
	orgID := chi.URLParam(r, "orgID")

	users, err := h.organizationUserService.GetUsersByOrganizationID(context.Background(), services.GetUsersByOrganizationIDParams{
		OrgID:        orgID,
		ActingUserID: userID,
	})

	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *organizationUserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")

	var input repositories.UpdateUserRoleParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.organizationUserService.UpdateUserRole(context.Background(), services.UpdateUserRoleParams{
		OrgID:        orgID,
		ActingUserID: userID,
		NewRole:      input.NewRole,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *organizationUserHandler) DeleteUserFromOrganization(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}
	orgID := chi.URLParam(r, "orgID")
	userIDToDelete := chi.URLParam(r, "userID")

	err = h.organizationUserService.DeleteUserFromOrganization(context.Background(), services.DeleteOrganizationUserParams{
		ActingUserID:   userID,
		OrgID:          orgID,
		UserIDToDelete: userIDToDelete,
	})

	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
