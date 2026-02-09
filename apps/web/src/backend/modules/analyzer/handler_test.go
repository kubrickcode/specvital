package analyzer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

func TestAnalyzeRepository(t *testing.T) {
	t.Run("returns 400 when owner is missing", func(t *testing.T) {
		_, r := setupTestHandler()

		req := httptest.NewRequest(http.MethodGet, "/api/analyze//repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		// When owner is empty, the router pattern doesn't match properly
		// and the generated handler returns 400 from oapi-codegen's error handler
		if rec.Code != http.StatusBadRequest && rec.Code != http.StatusNotFound {
			t.Errorf("expected status 400 or 404, got %d", rec.Code)
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
		gitClient := &mockGitClient{}
		tokenProvider := &mockTokenProvider{}
		_, r := setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
		}

		if !queue.enqueueCalled {
			t.Error("expected queue.Enqueue to be called")
		}

		var resp api.QueuedResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != "queued" {
			t.Errorf("expected status %q, got %q", "queued", resp.Status)
		}
	})

	t.Run("returns 200 with completed analysis when exists", func(t *testing.T) {
		queue := &mockQueueService{}
		parserVersion := "v1.0.0" // Must match mockSystemConfigReader.parserVersion
		repo := &mockRepository{
			completedAnalysis: &port.CompletedAnalysis{
				ID:            "550e8400-e29b-41d4-a716-446655440000",
				Owner:         "owner",
				Repo:          "repo",
				CommitSHA:     "abc123",
				ParserVersion: &parserVersion,
				CompletedAt:   time.Now(),
				TotalSuites:   5,
				TotalTests:    10,
			},
		}
		gitClient := &mockGitClient{commitSHA: "abc123"}
		tokenProvider := &mockTokenProvider{}
		_, r := setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp api.CompletedResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != "completed" {
			t.Errorf("expected status %q, got %q", "completed", resp.Status)
		}

		if resp.Data.CommitSHA != "abc123" {
			t.Errorf("expected commitSha %q, got %q", "abc123", resp.Data.CommitSHA)
		}

		if !repo.lastViewedCalled {
			t.Error("expected UpdateLastViewed to be called")
		}
		if repo.lastViewedOwner != "owner" || repo.lastViewedRepo != "repo" {
			t.Errorf("expected UpdateLastViewed(owner, repo), got (%s, %s)", repo.lastViewedOwner, repo.lastViewedRepo)
		}
	})
}

func TestGetAnalysisStatus(t *testing.T) {
	t.Run("returns 404 when no record exists", func(t *testing.T) {
		queue := &mockQueueService{}
		repo := &mockRepository{}
		gitClient := &mockGitClient{}
		tokenProvider := &mockTokenProvider{}
		_, r := setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("returns 200 with analyzing status when in progress", func(t *testing.T) {
		queue := &mockQueueService{
			findTaskInfo: &port.TaskInfo{
				CommitSHA: "test-commit-sha",
				State:     "running",
			},
		}
		repo := &mockRepository{}
		gitClient := &mockGitClient{}
		tokenProvider := &mockTokenProvider{}
		_, r := setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/repo/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		// Note: GetAnalysisStatus returns 200 for all statuses (including in-progress)
		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resp api.AnalyzingResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Status != "analyzing" {
			t.Errorf("expected status %q, got %q", "analyzing", resp.Status)
		}

		// UpdateLastViewed should NOT be called for in-progress analysis
		if repo.lastViewedCalled {
			t.Error("UpdateLastViewed should not be called for in-progress analysis")
		}
	})

	t.Run("updates last_viewed_at when returning completed analysis", func(t *testing.T) {
		queue := &mockQueueService{}
		repo := &mockRepository{
			completedAnalysis: &port.CompletedAnalysis{
				ID:          "550e8400-e29b-41d4-a716-446655440001",
				Owner:       "testowner",
				Repo:        "testrepo",
				CommitSHA:   "def456",
				CompletedAt: time.Now(),
				TotalSuites: 3,
				TotalTests:  15,
			},
		}
		gitClient := &mockGitClient{}
		tokenProvider := &mockTokenProvider{}
		_, r := setupTestHandlerWithMocks(repo, queue, gitClient, tokenProvider)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/testowner/testrepo/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		if !repo.lastViewedCalled {
			t.Error("expected UpdateLastViewed to be called")
		}
		if repo.lastViewedOwner != "testowner" || repo.lastViewedRepo != "testrepo" {
			t.Errorf("expected UpdateLastViewed(testowner, testrepo), got (%s, %s)", repo.lastViewedOwner, repo.lastViewedRepo)
		}
	})
}
