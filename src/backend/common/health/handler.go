package health

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
)

const statusOK = "ok"

type Response struct {
	Status string `json:"status"`
}

type Handler struct {
	logger *logger.Logger
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.handleHealth)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Status: statusOK}); err != nil {
		h.logger.Error(r.Context(), "failed to encode health response", "error", err)
	}
}
