package analyzer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/common/dto"
)

const (
	dbTimeout = 5 * time.Second
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type Handler struct {
	queue QueueService
	repo  Repository
}

func NewHandler(repo Repository, queue QueueService) *Handler {
	return &Handler{
		queue: queue,
		repo:  repo,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/analyze", func(r chi.Router) {
		r.Get("/{owner}/{repo}", h.handleAnalyze)
		r.Get("/{owner}/{repo}/status", h.handleStatus)
	})
}

func (h *Handler) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if err := validateOwnerRepo(owner, repo); err != nil {
		dto.SendProblemDetail(w, r, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), dbTimeout)
	defer cancel()

	analysis, err := h.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err == nil {
		sendAnalysisResponse(w, http.StatusOK, &AnalysisResponse{
			Status: StatusCompleted,
			Data:   toAnalysisResult(analysis),
		})
		return
	}

	if !errors.Is(err, ErrNotFound) {
		slog.Error("failed to get analysis", "owner", owner, "repo", repo, "error", err)
		dto.SendProblemDetail(w, r, http.StatusInternalServerError, "Internal Server Error", "failed to get analysis")
		return
	}

	status, err := h.repo.GetAnalysisStatus(ctx, owner, repo)
	if err == nil {
		sendAnalysisResponse(w, http.StatusAccepted, &AnalysisResponse{
			Status: mapDBStatus(status.Status),
			Error:  ptrToString(status.ErrorMessage),
		})
		return
	}

	if !errors.Is(err, ErrNotFound) {
		slog.Error("failed to get status", "owner", owner, "repo", repo, "error", err)
		dto.SendProblemDetail(w, r, http.StatusInternalServerError, "Internal Server Error", "failed to get status")
		return
	}

	analysisID, err := h.repo.CreatePendingAnalysis(ctx, owner, repo)
	if err != nil {
		slog.Error("failed to create analysis", "owner", owner, "repo", repo, "error", err)
		dto.SendProblemDetail(w, r, http.StatusInternalServerError, "Internal Server Error", "failed to create analysis")
		return
	}

	if err := h.queue.Enqueue(ctx, analysisID, owner, repo); err != nil {
		slog.Error("failed to enqueue", "owner", owner, "repo", repo, "analysisId", analysisID, "error", err)
		if cleanupErr := h.repo.MarkAnalysisFailed(ctx, analysisID, "queue registration failed"); cleanupErr != nil {
			slog.Error("failed to cleanup after enqueue error", "analysisId", analysisID, "error", cleanupErr)
		}
		dto.SendProblemDetail(w, r, http.StatusInternalServerError, "Internal Server Error", "failed to queue analysis")
		return
	}

	sendAnalysisResponse(w, http.StatusAccepted, &AnalysisResponse{
		Status: StatusQueued,
	})
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if err := validateOwnerRepo(owner, repo); err != nil {
		dto.SendProblemDetail(w, r, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), dbTimeout)
	defer cancel()

	analysis, err := h.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err == nil {
		sendAnalysisResponse(w, http.StatusOK, &AnalysisResponse{
			Status: StatusCompleted,
			Data:   toAnalysisResult(analysis),
		})
		return
	}

	if !errors.Is(err, ErrNotFound) {
		slog.Error("failed to get analysis", "owner", owner, "repo", repo, "error", err)
		dto.SendProblemDetail(w, r, http.StatusInternalServerError, "Internal Server Error", "failed to get analysis")
		return
	}

	status, err := h.repo.GetAnalysisStatus(ctx, owner, repo)
	if err == nil {
		httpStatus := http.StatusOK
		if status.Status != "completed" {
			httpStatus = http.StatusAccepted
		}
		sendAnalysisResponse(w, httpStatus, &AnalysisResponse{
			Status: mapDBStatus(status.Status),
			Error:  ptrToString(status.ErrorMessage),
		})
		return
	}

	if errors.Is(err, ErrNotFound) {
		dto.SendProblemDetail(w, r, http.StatusNotFound, "Not Found", "analysis not found")
		return
	}

	slog.Error("failed to get status", "owner", owner, "repo", repo, "error", err)
	dto.SendProblemDetail(w, r, http.StatusInternalServerError, "Internal Server Error", "failed to get status")
}

func validateOwnerRepo(owner, repo string) error {
	if owner == "" || repo == "" {
		return errors.New("owner and repo are required")
	}
	if !validNamePattern.MatchString(owner) {
		return errors.New("invalid owner format")
	}
	if !validNamePattern.MatchString(repo) {
		return errors.New("invalid repo format")
	}
	return nil
}

func sendAnalysisResponse(w http.ResponseWriter, statusCode int, resp *AnalysisResponse) {
	w.Header().Set(dto.ContentTypeHeader, dto.JSONContentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func toAnalysisResult(a *CompletedAnalysis) *AnalysisResult {
	return &AnalysisResult{
		AnalyzedAt: a.CompletedAt.Format("2006-01-02T15:04:05Z"),
		CommitSHA:  a.CommitSHA,
		Owner:      a.Owner,
		Repo:       a.Repo,
		Suites:     []TestSuite{},
		Summary: Summary{
			Frameworks: []FrameworkSummary{},
			Total:      a.TotalTests,
		},
	}
}

func mapDBStatus(status string) AnalysisStatusType {
	switch status {
	case "completed":
		return StatusCompleted
	case "running":
		return StatusAnalyzing
	case "pending":
		return StatusQueued
	case "failed":
		return StatusFailed
	default:
		slog.Warn("unknown analysis status, treating as failed", "status", status)
		return StatusFailed
	}
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
