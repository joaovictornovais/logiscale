package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joaovictornovais/logiscale/internal/service"
)

type DispatchHandler struct {
	service *service.DispatchService
}

func NewDispatchHandler(service *service.DispatchService) *DispatchHandler {
	return &DispatchHandler{service: service}
}

func (h *DispatchHandler) FindNearest(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	radius, _ := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)

	if lat == 0 || lng == 0 || radius == 0 {
		http.Error(w, "missing lat, lng or radius", http.StatusBadRequest)
	}

	drivers, err := h.service.FindNearestDrivers(r.Context(), lat, lng, radius)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)

}
