package analyzer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAnalyze(t *testing.T) {
	_, r := setupTestHandler()

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
