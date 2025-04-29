package handlers

import (
	"encoding/json"
	"net/http"

	"targeting-engine/internal/models"
	"targeting-engine/internal/service"
)

// This handler deals with getting the right ads to show
type DeliveryHandler struct {
	service service.Service
}

// Create a new handler for delivering ads
func NewDeliveryHandler(service service.Service) http.Handler {
	return &DeliveryHandler{
		service: service,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *DeliveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get the info we need from the URL
	query := r.URL.Query()
	req := models.DeliveryRequest{
		App:     query.Get("app"),
		OS:      query.Get("os"),
		Country: query.Get("country"),
	}

	// Make sure we have all the info we need
	if req.App == "" {
		h.respondWithError(w, http.StatusBadRequest, "missing app param")
		return
	}
	if req.OS == "" {
		h.respondWithError(w, http.StatusBadRequest, "missing os param")
		return
	}
	if req.Country == "" {
		h.respondWithError(w, http.StatusBadRequest, "missing country param")
		return
	}

	// Find ads that match what we're looking for
	campaigns, err := h.service.GetMatchingCampaigns(r.Context(), req)
	if err != nil {
		if err == service.ErrInvalidRequest {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Oops, something went wrong
		h.respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// If we didn't find any matching ads
	if len(campaigns) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Send back the matching ads
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaigns)
}

// Send back an error message
func (h *DeliveryHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
