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
