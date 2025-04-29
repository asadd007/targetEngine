package handlers

import (
	"encoding/json"
	"net/http"

	"targeting-engine/internal/models"
	"targeting-engine/internal/service"
)

type DeliveryHandler struct {
	service service.Service
}

func NewDeliveryHandler(service service.Service) http.Handler {
	return &DeliveryHandler{
		service: service,
	}
}

func (h *DeliveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	req := models.DeliveryRequest{
		App:     query.Get("app"),
		OS:      query.Get("os"),
		Country: query.Get("country"),
	}

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
	campaigns, err := h.service.GetMatchingCampaigns(r.Context(), req)
	if err != nil {
		if err == service.ErrInvalidRequest {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if len(campaigns) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaigns)
}

func (h *DeliveryHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
