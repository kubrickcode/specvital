package analyzer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleAnalyze(t *testing.T) {
	t.Run("returns 400 when owner is missing", func(t *testing.T) {
		_, r := setupTestHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/analyze//repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/problem+json" {
			t.Errorf("expected Content-Type application/problem+json, got %s", contentType)
		}
	})

	t.Run("returns 404 when repo path is incomplete", func(t *testing.T) {
		_, r := setupTestHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		// Chi router returns 404 when the route pattern doesn't match
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("returns 202 and queues analysis when no record exists", func(t *testing.T) {
		queue := &mockQueueService{}
		repo := &mockRepository{}
		_, r := setupTestHandlerWithMocks(repo, queue)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
		}

		if !queue.enqueueCalled {
			t.Error("expected queue.Enqueue to be called")
		}

		var resp AnalysisResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != StatusQueued {
			t.Errorf("expected status %q, got %q", StatusQueued, resp.Status)
		}
	})

	t.Run("returns 200 with completed analysis when exists", func(t *testing.T) {
		queue := &mockQueueService{}
		repo := &mockRepository{
			completedAnalysis: &CompletedAnalysis{
				ID:          "test-id",
				Owner:       "owner",
				Repo:        "repo",
				CommitSHA:   "abc123",
				CompletedAt: time.Now(),
				TotalSuites: 5,
				TotalTests:  10,
			},
		}
		_, r := setupTestHandlerWithMocks(repo, queue)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp AnalysisResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != StatusCompleted {
			t.Errorf("expected status %q, got %q", StatusCompleted, resp.Status)
		}

		if resp.Data == nil {
			t.Fatal("expected data to be present")
		}

		if resp.Data.CommitSHA != "abc123" {
			t.Errorf("expected commitSha %q, got %q", "abc123", resp.Data.CommitSHA)
		}
	})
}

func TestHandleStatus(t *testing.T) {
	t.Run("returns 404 when no record exists", func(t *testing.T) {
		queue := &mockQueueService{}
		repo := &mockRepository{}
		_, r := setupTestHandlerWithMocks(repo, queue)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("returns analyzing status when in progress", func(t *testing.T) {
		queue := &mockQueueService{}
		repo := &mockRepository{
			analysisStatus: &AnalysisStatus{
				ID:        "test-id",
				Status:    "running",
				CreatedAt: time.Now(),
			},
		}
		_, r := setupTestHandlerWithMocks(repo, queue)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
		}

		var resp AnalysisResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != StatusAnalyzing {
			t.Errorf("expected status %q, got %q", StatusAnalyzing, resp.Status)
		}
	})
}
