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

type organizationHandler struct {
	organizationService services.OrganizationService
	log *slog.Logger
}

func NewOrganizationHandler(organizationService services.OrganizationService, log *slog.Logger) *organizationHandler {
	return &organizationHandler{
		organizationService: organizationService,
		log: log.With(slog.String("component", "organization_handler")),
	}
}

func (h *organizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		h.log.Error("User identity not found in request context", slog.Any("error", err))
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	log := h.log.With(slog.String("user_id", identity.UserID))

	log.Info("Attempting to create a new organization", slog.String("user_id", identity.UserID))

	var input CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Warn("Failed to decode request body", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		log.Warn("Validation failed for organization creation", slog.Any("error", err))
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	org, err := h.organizationService.CreateOrganization(r.Context(), services.CreateOrganizationParams{
		Name:      input.Name,
		CreatedBy: identity.UserID,
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			log.Warn("Invalid input for organization creation", slog.Any("error", err))
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, services.ErrOrganizationWithDuplicateDetailsExists) {
			log.Warn("Organization with duplicate details already exists", slog.Any("error", err))
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		log.Error("Failed to create organization due to internal error", slog.Any("error", err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	log.Info("Organization created successfully", slog.String("organization_id", org.ID))

	response := NewOrganizationResponse(org)

	respondJSON(w, http.StatusCreated, response)
}

func (h *organizationHandler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {
	_, err := auth.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	orgID := chi.URLParam(r, "orgID")
	if orgID == "" {
		respondError(w, http.StatusBadRequest, "organization ID is required")
		return
	}

	org, err := h.organizationService.GetOrganizationByID(r.Context(), services.GetOrganizationByIDParams{ID: orgID})
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
