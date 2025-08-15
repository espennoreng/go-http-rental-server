package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/espennoreng/go-http-rental-server/internal/auth"
	"github.com/espennoreng/go-http-rental-server/internal/services"
	"github.com/go-chi/chi/v5"
)

type organizationHandler struct {
	organizationService services.OrganizationService
}

func NewOrganizationHandler(organizationService services.OrganizationService) *organizationHandler {
	return &organizationHandler{
		organizationService: organizationService,
	}
}

func (h *organizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	identity, err := auth.FromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user ID not found in context")
		return
	}

	var input CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := input.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	org, err := h.organizationService.CreateOrganization(r.Context(), services.CreateOrganizationParams{
		Name:      input.Name,
		CreatedBy: identity.UserID,
	})
	
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

	response := toOrganizationResponse(org)

	respondJSON(w, http.StatusCreated, response)
}

func (h *organizationHandler) GetOrganizationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	org, err := h.organizationService.GetOrganizationByID(r.Context(), services.GetOrganizationByIDParams{ID: id})
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
