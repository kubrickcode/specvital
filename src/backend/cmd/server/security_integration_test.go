package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/common/middleware"
)

func TestSecurityHeadersIntegration(t *testing.T) {
	// Create a minimal router with just security headers
	handler := middleware.SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options":              "nosniff",
		"X-Frame-Options":                     "DENY",
		"X-XSS-Protection":                    "1; mode=block",
		"Referrer-Policy":                     "strict-origin-when-cross-origin",
		"Content-Security-Policy-Report-Only": "default-src 'self'",
	}

	for header, expected := range expectedHeaders {
		if got := rec.Header().Get(header); got != expected {
			t.Errorf("header %s: expected %q, got %q", header, expected, got)
		}
	}
}

func TestRateLimitingIntegration(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(60) // 60 requests per minute
	handler := middleware.RateLimit(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	// First request should succeed
	req1 := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	req1.Header.Set("X-Real-IP", "192.168.1.1")
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("first request: expected status %d, got %d", http.StatusOK, rec1.Code)
	}

	// Second immediate request should be rate limited (burst=1)
	req2 := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	req2.Header.Set("X-Real-IP", "192.168.1.1")
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected status %d, got %d", http.StatusTooManyRequests, rec2.Code)
	}

	// Different IP should not be rate limited
	req3 := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	req3.Header.Set("X-Real-IP", "192.168.1.2")
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Errorf("different IP: expected status %d, got %d", http.StatusOK, rec3.Code)
	}
}

func TestRateLimitWithTokenRefill(t *testing.T) {
	// Skip in short mode (test takes >1 second)
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	limiter := middleware.NewIPRateLimiter(60) // 60 requests per minute = 1 per second
	handler := middleware.RateLimit(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ip := "192.168.1.100"

	// First request: allowed
	req1 := httptest.NewRequest(http.MethodGet, "/api/auth/callback", nil)
	req1.Header.Set("X-Real-IP", ip)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("first request failed: got status %d", rec1.Code)
	}

	// Immediate second request: rate limited
	req2 := httptest.NewRequest(http.MethodGet, "/api/auth/callback", nil)
	req2.Header.Set("X-Real-IP", ip)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request should be rate limited: got status %d", rec2.Code)
	}

	// Wait for token refill (1 second + buffer)
	time.Sleep(1100 * time.Millisecond)

	// Third request after waiting: allowed
	req3 := httptest.NewRequest(http.MethodGet, "/api/auth/callback", nil)
	req3.Header.Set("X-Real-IP", ip)
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Errorf("third request after refill: expected status %d, got %d", http.StatusOK, rec3.Code)
	}
}

func TestCombinedMiddlewares(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(60)

	// Stack middlewares in order: security headers -> rate limit -> handler
	handler := middleware.SecurityHeaders()(
		middleware.RateLimit(limiter)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			}),
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("X-Real-IP", "192.168.1.200")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify security headers are set
	if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("expected X-Frame-Options header, got %q", got)
	}

	// Verify rate limiting works
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Second request should be rate limited
	req2 := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req2.Header.Set("X-Real-IP", "192.168.1.200")
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected status %d, got %d", http.StatusTooManyRequests, rec2.Code)
	}

	// Rate limited response should still have security headers
	if got := rec2.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("rate limited response should have security headers, got %q", got)
	}
}
