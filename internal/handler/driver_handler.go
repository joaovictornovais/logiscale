package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joaovictornovais/logiscale/internal/service"
)

type DriverHandler struct {
	service *service.DriverService
}

func NewDriverHandler(service *service.DriverService) *DriverHandler {
	return &DriverHandler{service: service}
}

type CreateDriverRequest struct {
	Name    string `json:"name"`
	License string `json:"license"`
}

type DriverResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	License string `json:"license"`
}

func (h *DriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	var req CreateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	driver, err := h.service.CreateDriver(r.Context(), req.Name, req.License)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(DriverResponse{
		ID:      driver.ID,
		Name:    driver.Name,
		License: driver.License,
	})
}

func (h *DriverHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	driver, err := h.service.GetDriverByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DriverResponse{
		ID:      driver.ID,
		Name:    driver.Name,
		License: driver.License,
	})
}
