package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/api"
)

const (
	defaultRateLimit       = 5
	defaultBurst           = 1
	defaultCleanupTTL      = 10 * time.Minute
	defaultCleanupInterval = 5 * time.Minute
)

type RateLimiter interface {
	Allow(key string) bool
	Stop()
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	limiters map[string]*limiterEntry
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
	ttl      time.Duration
	stopCh   chan struct{}
}

func NewIPRateLimiter(requestsPerMinute int) RateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = defaultRateLimit
	}

	ratePerSecond := rate.Limit(float64(requestsPerMinute) / 60.0)

	rl := &ipRateLimiter{
		limiters: make(map[string]*limiterEntry),
		rate:     ratePerSecond,
		burst:    defaultBurst,
		ttl:      defaultCleanupTTL,
		stopCh:   make(chan struct{}),
	}

	go rl.cleanupLoop()

	return rl
}

func (rl *ipRateLimiter) Allow(ip string) bool {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.limiters[ip]
	if !exists {
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(rl.rate, rl.burst),
			lastSeen: now,
		}
		rl.limiters[ip] = entry
	} else {
		entry.lastSeen = now
	}

	return entry.limiter.Allow()
}

func (rl *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(defaultCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *ipRateLimiter) Stop() {
	close(rl.stopCh)
}

func (rl *ipRateLimiter) cleanup() {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, entry := range rl.limiters {
		if now.Sub(entry.lastSeen) > rl.ttl {
			delete(rl.limiters, ip)
		}
	}
}

func RateLimit(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.Allow(ip) {
				writeTooManyRequests(w, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func writeTooManyRequests(w http.ResponseWriter, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.Header().Set("Retry-After", "60")
	w.WriteHeader(http.StatusTooManyRequests)
	_ = json.NewEncoder(w).Encode(api.NewTooManyRequests(detail))
}
