package health

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const statusOK = "ok"

type Response struct {
	Status string `json:"status"`
}

func RegisterRoutes(r chi.Router) {
	r.Get("/health", handleHealth)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Status: statusOK}); err != nil {
		slog.Error("failed to encode health response", "error", err)
	}
}
