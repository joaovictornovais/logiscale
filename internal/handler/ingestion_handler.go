package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joaovictornovais/logiscale/internal/domain"
	"github.com/joaovictornovais/logiscale/internal/service"
)

type IngestionHandler struct {
	service *service.IngestionService
}

func NewIngestionHandler(service *service.IngestionService) *IngestionHandler {
	return &IngestionHandler{service: service}
}

func (h *IngestionHandler) HandleIngest(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "id")
	var payload domain.LocationPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	payload.DriverID = driverID

	if payload.SentAt.IsZero() {
		payload.SentAt = time.Now()
	}

	h.service.Ingest(payload)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "queued"}`))
}
