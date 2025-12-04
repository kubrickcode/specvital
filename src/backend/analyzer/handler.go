package analyzer

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/common/dto"
)

// Mock repo prefixes for testing error states (development only).
// Usage examples:
//   - /api/analyze/owner/private-repo  -> 403 Forbidden
//   - /api/analyze/owner/notfound-repo -> 404 Not Found
//   - /api/analyze/owner/ratelimit-repo -> 429 Too Many Requests
const (
	mockPrivatePrefix   = "private-"
	mockNotFoundPrefix  = "notfound-"
	mockRateLimitPrefix = "ratelimit-"
)

func isMockErrorsEnabled() bool {
	// Never enable mock errors in production environment
	env := os.Getenv("ENV")
	if env == "production" || env == "prod" {
		return false
	}
	return os.Getenv("ENABLE_MOCK_ERRORS") == "true"
}

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

	// Mock error states for testing frontend error handling (development only).
	if isMockErrorsEnabled() {
		if strings.HasPrefix(repo, mockPrivatePrefix) {
			dto.SendProblemDetail(w, r, http.StatusForbidden, "Forbidden", "This repository is private or you don't have access")
			return
		}
		if strings.HasPrefix(repo, mockNotFoundPrefix) {
			dto.SendProblemDetail(w, r, http.StatusNotFound, "Not Found", "Repository not found")
			return
		}
		if strings.HasPrefix(repo, mockRateLimitPrefix) {
			dto.SendProblemDetail(w, r, http.StatusTooManyRequests, "Too Many Requests", "GitHub API rate limit exceeded. Please try again later")
			return
		}
	}

	result := GenerateMockResult(owner, repo)

	w.Header().Set(dto.ContentTypeHeader, dto.JSONContentType)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("failed to encode analysis result", "error", err)
	}
}
