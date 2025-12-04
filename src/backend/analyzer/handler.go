package analyzer

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/common/dto"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/api/analyze", func(r chi.Router) {
		r.Get("/{owner}/{repo}", handleAnalyze)
	})
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if owner == "" || repo == "" {
		dto.SendProblemDetail(w, r, http.StatusBadRequest, "Bad Request", "owner and repo are required")
		return
	}

	result := GenerateMockResult(owner, repo)

	w.Header().Set(dto.ContentTypeHeader, dto.JSONContentType)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("failed to encode analysis result", "error", err)
	}
}
