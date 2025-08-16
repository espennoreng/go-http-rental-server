package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type organizationUserHandler struct {
	organizationUserService services.OrganizationUserService
	log *slog.Logger
}

func NewOrganizationUserHandler(organizationUserService services.OrganizationUserService, log *slog.Logger) *organizationUserHandler {
	return &organizationUserHandler{
		organizationUserService: organizationUserService,
		log:                      log.With(slog.String("component", "organization_user_handler")),
	}
}

func (h *organizationUserHandler) AddUserToOrganization(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		h.log.Error("Failed to retrieve user ID from context", slog.Any("error", err))
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}
	orgID := chi.URLParam(r, "orgID")

	if orgID == "" {
		h.log.Warn("Organization ID is required for adding user", slog.String("orgID", orgID))
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	var input AddUserToOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("Failed to decode request body", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.log.Warn("Validation failed for adding user to organization", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	log := h.log.With(slog.String("user_id", identity.UserID), slog.String("org_id", orgID))
	log.Info("Adding user to organization", slog.String("user_id", input.UserID), slog.String("role", string(input.Role)))

	newOrgUser, err := h.organizationUserService.CreateOrganizationUser(context.Background(), services.CreateOrganizationUserParams{
		ActingUserID: identity.UserID,
		OrgID:        orgID,
		UserID:       input.UserID,
		Role:         input.Role,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			log.Warn("Invalid input for adding user to organization", slog.Any("error", err))
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Error("Failed to add user to organization", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	log.Info("User added to organization successfully", slog.String("user_id", input.UserID), slog.String("org_id", orgID))

	response := NewOrganizationUserResponse(newOrgUser)

	respondJSON(w, http.StatusCreated, response)
}

func (h *organizationUserHandler) GetUsersByOrganizationID(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())

	if err != nil {
		h.log.Error("Failed to retrieve user ID from context", slog.Any("error", err))
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")

	if orgID == "" {
		h.log.Warn("Organization ID is required for fetching users", slog.String("orgID", orgID))
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	log := h.log.With(slog.String("user_id", identity.UserID), slog.String("org_id", orgID))
	log.Info("Fetching users for organization", slog.String("org_id", orgID))

	users, err := h.organizationUserService.GetUsersByOrganizationID(context.Background(), services.GetUsersByOrganizationIDParams{
		OrgID:        orgID,
		ActingUserID: identity.UserID,
	})

	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			log.Warn("Unauthorized access attempt", slog.Any("error", err))
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		log.Error("Failed to fetch users for organization", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := NewOrganizationMembersResponse(users)

	respondJSON(w, http.StatusOK, response)
}

func (h *organizationUserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		h.log.Error("Failed to retrieve user ID from context", slog.Any("error", err))
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")

	if orgID == "" {
		h.log.Warn("Organization ID is required for updating user role", slog.String("orgID", orgID))
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	var input UpdateUserRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("Failed to decode request body", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		h.log.Warn("Validation failed for updating user role", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	log := h.log.With(slog.String("user_id", identity.UserID), slog.String("org_id", orgID))
	log.Info("Updating user role in organization", slog.String("user_id", input.UserID), slog.String("role", string(input.Role)))

	err = h.organizationUserService.UpdateUserRole(context.Background(), services.UpdateUserRoleParams{
		OrgID:        orgID,
		ActingUserID: identity.UserID,
		Role:         input.Role,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			log.Warn("Invalid input for updating user role", slog.Any("error", err))
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			log.Warn("Unauthorized access attempt", slog.Any("error", err))
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		log.Error("Failed to update user role in organization", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *organizationUserHandler) DeleteUserFromOrganization(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		h.log.Error("Failed to retrieve user ID from context", slog.Any("error", err))
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")

	if orgID == "" {
		h.log.Warn("Organization ID is required for deleting user", slog.String("orgID", orgID))
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	userIDToDelete := chi.URLParam(r, "userID")
	if userIDToDelete == "" {
		h.log.Warn("User ID is required for deleting user from organization", slog.String("userID", userIDToDelete))
		respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	log := h.log.With(slog.String("user_id", identity.UserID), slog.String("org_id", orgID), slog.String("user_id_to_delete", userIDToDelete))
	log.Info("Deleting user from organization", slog.String("user_id_to_delete", userIDToDelete))

	err = h.organizationUserService.DeleteUserFromOrganization(context.Background(), services.DeleteOrganizationUserParams{
		ActingUserID:   identity.UserID,
		OrgID:          orgID,
		UserIDToDelete: userIDToDelete,
	})

	if err != nil {
		if errors.Is(err, services.ErrUnauthorized) {
			log.Warn("Unauthorized access attempt", slog.Any("error", err))
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		log.Error("Internal error while deleting user from organization", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	log.Info("User deleted from organization successfully", slog.String("user_id_to_delete", userIDToDelete))

	w.WriteHeader(http.StatusNoContent)
}
