package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/api"
)

// mockRateLimiter is a test double for RateLimiter interface.
type mockRateLimiter struct {
	allowFunc func(key string) bool
}

func (m *mockRateLimiter) Allow(key string) bool {
	return m.allowFunc(key)
}

func (m *mockRateLimiter) Stop() {
	// No-op for mock
}

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name           string
		allowFunc      func(string) bool
		expectedStatus int
		expectError    bool
	}{
		{
			name: "allow request when under rate limit",
			allowFunc: func(string) bool {
				return true
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "reject request when over rate limit",
			allowFunc: func(string) bool {
				return false
			},
			expectedStatus: http.StatusTooManyRequests,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := &mockRateLimiter{allowFunc: tt.allowFunc}

			handler := RateLimit(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectError {
				contentType := rec.Header().Get("Content-Type")
				if contentType != "application/problem+json" {
					t.Errorf("expected content-type application/problem+json, got %s", contentType)
				}

				var problem api.ProblemDetail
				if err := json.NewDecoder(rec.Body).Decode(&problem); err != nil {
					t.Fatalf("failed to decode problem response: %v", err)
				}

				if problem.Status != http.StatusTooManyRequests {
					t.Errorf("expected problem status %d, got %d", http.StatusTooManyRequests, problem.Status)
				}
				if problem.Title != "Too Many Requests" {
					t.Errorf("expected title 'Too Many Requests', got %q", problem.Title)
				}
				if problem.Detail == "" {
					t.Error("expected non-empty detail")
				}
			}
		})
	}
}

func TestRateLimitWithRealIP(t *testing.T) {
	limiter := &mockRateLimiter{
		allowFunc: func(ip string) bool {
			return ip == "192.168.1.1"
		},
	}

	handler := RateLimit(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name           string
		realIP         string
		expectedStatus int
	}{
		{
			name:           "allowed IP",
			realIP:         "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "blocked IP",
			realIP:         "10.0.0.1",
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("X-Real-IP", tt.realIP)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestIPRateLimiter(t *testing.T) {
	t.Run("allow requests within rate limit", func(t *testing.T) {
		limiter := NewIPRateLimiter(60) // 60 requests per minute = 1 per second
		defer limiter.Stop()

		ip := "192.168.1.1"
		if !limiter.Allow(ip) {
			t.Error("expected first request to be allowed")
		}
	})

	t.Run("block requests exceeding burst", func(t *testing.T) {
		limiter := NewIPRateLimiter(60) // 60 requests per minute
		defer limiter.Stop()

		ip := "192.168.1.2"

		// First request should be allowed (within burst)
		if !limiter.Allow(ip) {
			t.Error("expected first request to be allowed")
		}

		// Second immediate request should be blocked (burst = 1)
		if limiter.Allow(ip) {
			t.Error("expected second immediate request to be blocked")
		}
	})

	t.Run("separate rate limits per IP", func(t *testing.T) {
		limiter := NewIPRateLimiter(60)
		defer limiter.Stop()

		ip1 := "192.168.1.1"
		ip2 := "192.168.1.2"

		// Both IPs should have their own limiters
		if !limiter.Allow(ip1) {
			t.Error("expected first request from IP1 to be allowed")
		}
		if !limiter.Allow(ip2) {
			t.Error("expected first request from IP2 to be allowed")
		}

		// Second requests should be blocked for both
		if limiter.Allow(ip1) {
			t.Error("expected second immediate request from IP1 to be blocked")
		}
		if limiter.Allow(ip2) {
			t.Error("expected second immediate request from IP2 to be blocked")
		}
	})

	t.Run("use default rate when invalid", func(t *testing.T) {
		limiter := NewIPRateLimiter(0) // Invalid rate
		defer limiter.Stop()

		// Should use default rate (5 per minute)
		ipLimiter, ok := limiter.(*ipRateLimiter)
		if !ok {
			t.Fatal("expected ipRateLimiter type")
		}

		expectedRate := rate.Limit(float64(defaultRateLimit) / 60.0)
		if ipLimiter.rate != expectedRate {
			t.Errorf("expected rate %v, got %v", expectedRate, ipLimiter.rate)
		}
	})
}

func TestIPRateLimiterConcurrency(t *testing.T) {
	limiter := NewIPRateLimiter(60)
	defer limiter.Stop()
	ip := "192.168.1.1"

	done := make(chan bool)
	for range 10 {
		go func() {
			limiter.Allow(ip)
			done <- true
		}()
	}

	for range 10 {
		<-done
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		realIP     string
		remoteAddr string
		expected   string
	}{
		{
			name:       "prefer X-Real-IP header",
			realIP:     "192.168.1.1",
			remoteAddr: "10.0.0.1:1234",
			expected:   "192.168.1.1",
		},
		{
			name:       "fallback to RemoteAddr when no header",
			realIP:     "",
			remoteAddr: "10.0.0.1:1234",
			expected:   "10.0.0.1:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}
			req.RemoteAddr = tt.remoteAddr

			got := getClientIP(req)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func BenchmarkRateLimit(b *testing.B) {
	limiter := NewIPRateLimiter(6000)
	defer limiter.Stop()

	handler := RateLimit(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for b.Loop() {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Real-IP", "192.168.1.1")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}
}

func TestIPRateLimiterIntegration(t *testing.T) {
	// Integration test: verify actual rate limiting behavior over time
	limiter := NewIPRateLimiter(60) // 60 requests per minute = 1 per second
	defer limiter.Stop()
	ip := "192.168.1.1"

	// First request allowed
	if !limiter.Allow(ip) {
		t.Fatal("expected first request to be allowed")
	}

	// Immediate second request blocked (burst exhausted)
	if limiter.Allow(ip) {
		t.Fatal("expected immediate second request to be blocked")
	}

	// Wait for token to refill (rate is 1 per second)
	time.Sleep(1100 * time.Millisecond)

	// Request should be allowed after waiting
	if !limiter.Allow(ip) {
		t.Error("expected request to be allowed after token refill")
	}
}

func TestIPRateLimiterStop(t *testing.T) {
	limiter := NewIPRateLimiter(60)

	// Should be able to stop without panic
	limiter.Stop()

	// After stop, Allow should still work (graceful degradation)
	// Note: cleanup goroutine is stopped, but Allow still functions
	if !limiter.Allow("192.168.1.1") {
		t.Error("expected request to be allowed after stop")
	}
}
