package ratelimit

import (
	"sync"
	"time"
)

// IPRateLimiter provides in-memory rate limiting per IP address.
// Uses a simple fixed window counter with automatic cleanup of stale entries.
type IPRateLimiter struct {
	counters    map[string]*counter
	limit       int
	window      time.Duration
	mu          sync.RWMutex
	cleanupTick time.Duration
	stopCleanup chan struct{}
}

type counter struct {
	count     int
	expiresAt time.Time
}

// NewIPRateLimiter creates a new rate limiter.
// limit: maximum requests per window
// window: time window duration (e.g., 1 minute, 1 hour)
func NewIPRateLimiter(limit int, window time.Duration) *IPRateLimiter {
	rl := &IPRateLimiter{
		counters:    make(map[string]*counter),
		limit:       limit,
		window:      window,
		cleanupTick: window * 2,
		stopCleanup: make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

// Allow checks if the request from the given IP should be allowed.
// Returns true if allowed, false if rate limited.
func (rl *IPRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	c, exists := rl.counters[ip]
	if !exists || now.After(c.expiresAt) {
		rl.counters[ip] = &counter{
			count:     1,
			expiresAt: now.Add(rl.window),
		}
		return true
	}

	if c.count >= rl.limit {
		return false
	}

	c.count++
	return true
}

// cleanup periodically removes expired entries to prevent memory leaks.
func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupTick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.removeExpired()
		case <-rl.stopCleanup:
			return
		}
	}
}

func (rl *IPRateLimiter) removeExpired() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, c := range rl.counters {
		if now.After(c.expiresAt) {
			delete(rl.counters, ip)
		}
	}
}

// Close stops the cleanup goroutine. Implements io.Closer.
func (rl *IPRateLimiter) Close() error {
	close(rl.stopCleanup)
	return nil
}
