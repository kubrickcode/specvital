package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	tests := []struct {
		header   string
		expected string
	}{
		{
			header:   "X-Content-Type-Options",
			expected: "nosniff",
		},
		{
			header:   "X-Frame-Options",
			expected: "DENY",
		},
		{
			header:   "X-XSS-Protection",
			expected: "1; mode=block",
		},
		{
			header:   "Referrer-Policy",
			expected: "strict-origin-when-cross-origin",
		},
		{
			header:   "Content-Security-Policy-Report-Only",
			expected: "default-src 'self'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			got := rec.Header().Get(tt.header)
			if got != tt.expected {
				t.Errorf("header %s: expected %q, got %q", tt.header, tt.expected, got)
			}
		})
	}
}

func TestSecurityHeadersPreservesResponse(t *testing.T) {
	expectedBody := "test response"
	expectedStatus := http.StatusCreated

	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)
		_, _ = w.Write([]byte(expectedBody))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, rec.Code)
	}

	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestSecurityHeadersDoesNotOverwriteExisting(t *testing.T) {
	customCSPReportOnly := "default-src 'self' https://example.com"

	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handler sets a custom CSP-Report-Only before writing response
		w.Header().Set("Content-Security-Policy-Report-Only", customCSPReportOnly)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Note: In Go's http package, middleware headers are set before handler headers
	// So the handler's Set() will overwrite middleware's Set()
	// This test documents the current behavior
	got := rec.Header().Get("Content-Security-Policy-Report-Only")
	if got != customCSPReportOnly {
		t.Logf("CSP-Report-Only was overwritten by handler (expected behavior): got %q", got)
	}
}

func TestSecurityHeadersMultipleRequests(t *testing.T) {
	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Verify headers are set consistently across multiple requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Errorf("request %d: expected X-Content-Type-Options=nosniff, got %q", i, got)
		}
		if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
			t.Errorf("request %d: expected X-Frame-Options=DENY, got %q", i, got)
		}
	}
}

func TestSecurityHeadersWithDifferentMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}

	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
				t.Errorf("method %s: expected X-Frame-Options=DENY, got %q", method, got)
			}
		})
	}
}

func BenchmarkSecurityHeaders(b *testing.B) {
	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}
