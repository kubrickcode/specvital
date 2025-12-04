package analyzer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/specvital/web/src/backend/common/dto"
)

func TestHandleAnalyze(t *testing.T) {
	r := chi.NewRouter()
	RegisterRoutes(r)

	t.Run("returns mock analysis result", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/analyze/facebook/react", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}

		var result AnalysisResult
		if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if result.Owner != "facebook" {
			t.Errorf("expected owner facebook, got %s", result.Owner)
		}

		if result.Repo != "react" {
			t.Errorf("expected repo react, got %s", result.Repo)
		}

		if result.Summary.Total == 0 {
			t.Error("expected non-zero total tests")
		}

		if len(result.Suites) == 0 {
			t.Error("expected non-empty test suites")
		}
	})

	t.Run("includes framework summary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/analyze/vercel/next.js", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		var result AnalysisResult
		if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(result.Summary.Frameworks) == 0 {
			t.Error("expected non-empty framework summaries")
		}

		for _, fw := range result.Summary.Frameworks {
			if fw.Total != fw.Active+fw.Skipped+fw.Todo {
				t.Errorf("framework %s: total mismatch, expected %d, got %d",
					fw.Framework, fw.Active+fw.Skipped+fw.Todo, fw.Total)
			}
		}
	})

	t.Run("returns 400 when owner is missing", func(t *testing.T) {
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

		var problem dto.ProblemDetail
		if err := json.NewDecoder(rec.Body).Decode(&problem); err != nil {
			t.Fatalf("failed to decode problem response: %v", err)
		}

		if problem.Status != http.StatusBadRequest {
			t.Errorf("expected problem status %d, got %d", http.StatusBadRequest, problem.Status)
		}
	})

	t.Run("returns 404 when repo path is incomplete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		// Chi router returns 404 when the route pattern doesn't match
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

}

func TestHandleAnalyzeMockErrors(t *testing.T) {
	// Enable mock errors for this test
	t.Setenv("ENABLE_MOCK_ERRORS", "true")

	r := chi.NewRouter()
	RegisterRoutes(r)

	t.Run("returns 403 for private repository mock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/private-repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != dto.ProblemJSONType {
			t.Errorf("expected Content-Type %s, got %s", dto.ProblemJSONType, contentType)
		}

		var problem dto.ProblemDetail
		if err := json.NewDecoder(rec.Body).Decode(&problem); err != nil {
			t.Fatalf("failed to decode problem response: %v", err)
		}

		if problem.Status != http.StatusForbidden {
			t.Errorf("expected problem status %d, got %d", http.StatusForbidden, problem.Status)
		}
	})

	t.Run("returns 404 for not found repository mock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/notfound-repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != dto.ProblemJSONType {
			t.Errorf("expected Content-Type %s, got %s", dto.ProblemJSONType, contentType)
		}

		var problem dto.ProblemDetail
		if err := json.NewDecoder(rec.Body).Decode(&problem); err != nil {
			t.Fatalf("failed to decode problem response: %v", err)
		}

		if problem.Status != http.StatusNotFound {
			t.Errorf("expected problem status %d, got %d", http.StatusNotFound, problem.Status)
		}
	})

	t.Run("returns 429 for rate limited repository mock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/ratelimit-repo", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if contentType != dto.ProblemJSONType {
			t.Errorf("expected Content-Type %s, got %s", dto.ProblemJSONType, contentType)
		}

		var problem dto.ProblemDetail
		if err := json.NewDecoder(rec.Body).Decode(&problem); err != nil {
			t.Fatalf("failed to decode problem response: %v", err)
		}

		if problem.Status != http.StatusTooManyRequests {
			t.Errorf("expected problem status %d, got %d", http.StatusTooManyRequests, problem.Status)
		}
	})

	t.Run("mock errors disabled returns normal response", func(t *testing.T) {
		t.Setenv("ENABLE_MOCK_ERRORS", "false")

		testRouter := chi.NewRouter()
		RegisterRoutes(testRouter)

		req := httptest.NewRequest(http.MethodGet, "/api/analyze/owner/private-repo", nil)
		rec := httptest.NewRecorder()

		testRouter.ServeHTTP(rec, req)

		// Without mock errors enabled, should return 200 with mock data
		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
