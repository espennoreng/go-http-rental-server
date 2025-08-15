package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type organizationUserHandler struct {
	organizationUserService services.OrganizationUserService
}

func NewOrganizationUserHandler(organizationUserService services.OrganizationUserService) *organizationUserHandler {
	return &organizationUserHandler{
		organizationUserService: organizationUserService,
	}
}

func (h *organizationUserHandler) AddUserToOrganization(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}
	orgID := chi.URLParam(r, "orgID")

	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	var input AddUserToOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	newOrgUser, err := h.organizationUserService.CreateOrganizationUser(context.Background(), services.CreateOrganizationUserParams{
		ActingUserID: identity.UserID,
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

	response := toUserResponse(newOrgUser)

	respondJSON(w, http.StatusCreated, response)
}

func (h *organizationUserHandler) GetUsersByOrganizationID(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())

	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")

	if orgID == "" {
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	users, err := h.organizationUserService.GetUsersByOrganizationID(context.Background(), services.GetUsersByOrganizationIDParams{
		OrgID:        orgID,
		ActingUserID: identity.UserID,
	})

	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := toUsersResponse(users)

	respondJSON(w, http.StatusOK, response)
}

func (h *organizationUserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")

	if orgID == "" {
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	var input UpdateUserRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.organizationUserService.UpdateUserRole(context.Background(), services.UpdateUserRoleParams{
		OrgID:        orgID,
		ActingUserID: identity.UserID,
		Role:      input.Role,
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
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}
	orgID := chi.URLParam(r, "orgID")
	userIDToDelete := chi.URLParam(r, "userID")

	err = h.organizationUserService.DeleteUserFromOrganization(context.Background(), services.DeleteOrganizationUserParams{
		ActingUserID:   identity.UserID,
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
