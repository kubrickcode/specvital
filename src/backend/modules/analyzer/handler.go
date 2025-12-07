package analyzer

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/common/clients/github"
	"github.com/specvital/web/src/backend/common/dto"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/analyze", func(r chi.Router) {
		r.Get("/{owner}/{repo}", h.handleAnalyze)
	})
}

func (h *Handler) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if owner == "" || repo == "" {
		dto.SendProblemDetail(w, r, http.StatusBadRequest, "Bad Request", "owner and repo are required")
		return
	}

	result, err := h.service.Analyze(r.Context(), owner, repo)
	if err != nil {
		h.handleAnalyzeError(w, r, err)
		return
	}

	sendJSON(w, result)
}

func (h *Handler) handleAnalyzeError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("analysis failed", "error", err)

	rateLimit := h.service.GetRateLimit()
	var rateLimitExt *dto.RateLimitExtension
	if rateLimit.Limit > 0 {
		rateLimitExt = &dto.RateLimitExtension{
			Limit:     rateLimit.Limit,
			Remaining: rateLimit.Remaining,
			ResetAt:   rateLimit.ResetAt,
		}
	}

	switch {
	case errors.Is(err, github.ErrNotFound):
		dto.SendProblemDetailWithRateLimit(w, r, http.StatusNotFound, "Not Found", "Repository not found", rateLimitExt)
	case errors.Is(err, github.ErrForbidden):
		dto.SendProblemDetailWithRateLimit(w, r, http.StatusForbidden, "Forbidden", "This repository is private or you don't have access", rateLimitExt)
	case errors.Is(err, github.ErrRateLimited):
		dto.SendProblemDetailWithRateLimit(w, r, http.StatusTooManyRequests, "Too Many Requests", "GitHub API rate limit exceeded. Please try again later", rateLimitExt)
	case errors.Is(err, github.ErrTreeTruncated):
		dto.SendProblemDetailWithRateLimit(w, r, http.StatusUnprocessableEntity, "Repository Too Large", "This repository is too large to analyze", rateLimitExt)
	case errors.Is(err, ErrRateLimitTooLow):
		dto.SendProblemDetailWithRateLimit(w, r, http.StatusTooManyRequests, "Too Many Requests", "GitHub API rate limit too low. Please try again later", rateLimitExt)
	default:
		dto.SendProblemDetailWithRateLimit(w, r, http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred during analysis", rateLimitExt)
	}
}

func sendJSON(w http.ResponseWriter, data any) {
	w.Header().Set(dto.ContentTypeHeader, dto.JSONContentType)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}
