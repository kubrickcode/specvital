package analyzer

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/common/clients/github"
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

var (
	analyzerService *Service
	serviceOnce     sync.Once
)

func getService() *Service {
	serviceOnce.Do(func() {
		token := os.Getenv("GITHUB_TOKEN")
		slog.Info("initializing analyzer service", "hasToken", token != "")
		githubClient := github.NewClient(token)
		analyzerService = NewService(githubClient)
	})
	return analyzerService
}

func isMockMode() bool {
	return os.Getenv("MOCK_MODE") == "true"
}

func isMockErrorsEnabled() bool {
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

	if isMockMode() {
		result := GenerateMockResult(owner, repo)
		sendJSON(w, result)
		return
	}

	result, err := getService().Analyze(r.Context(), owner, repo)
	if err != nil {
		handleAnalyzeError(w, r, err)
		return
	}

	sendJSON(w, result)
}

func handleAnalyzeError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("analysis failed", "error", err)

	rateLimit := getService().GetRateLimit()
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
