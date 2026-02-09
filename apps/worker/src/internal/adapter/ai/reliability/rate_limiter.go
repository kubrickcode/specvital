package reliability

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter wraps golang.org/x/time/rate.Limiter with project-specific config.
type RateLimiter struct {
	limiter *rate.Limiter
}

// RateLimiterConfig holds configuration for the rate limiter.
type RateLimiterConfig struct {
	BurstFactor int // Burst capacity multiplier (default: 1)
	RPM         int // Requests per minute
}

// DefaultRateLimiterConfig returns the default global rate limiter config.
// Uses 90% of Gemini's free tier limit (2000 RPM for Flash).
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		BurstFactor: 1,
		RPM:         1800, // 90% of 2000 RPM
	}
}

// NewRateLimiter creates a new rate limiter using golang.org/x/time/rate.
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.BurstFactor < 1 {
		config.BurstFactor = 1
	}

	// Convert RPM to rate.Limit (requests per second)
	limit := rate.Limit(float64(config.RPM) / 60.0)

	return &RateLimiter{
		limiter: rate.NewLimiter(limit, config.BurstFactor),
	}
}

// Wait blocks until a token is available or context is cancelled.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// Allow reports whether an event may happen now.
func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

// Available returns the current number of available tokens (approximate).
func (rl *RateLimiter) Available() float64 {
	return float64(rl.limiter.Tokens())
}

// GlobalRateLimiter is a singleton rate limiter for all AI operations.
var (
	globalLimiter   *RateLimiter
	globalLimiterMu sync.RWMutex
)

// GetGlobalRateLimiter returns the singleton rate limiter.
// Creates a new instance with default config if not set.
func GetGlobalRateLimiter() *RateLimiter {
	globalLimiterMu.RLock()
	if globalLimiter != nil {
		defer globalLimiterMu.RUnlock()
		return globalLimiter
	}
	globalLimiterMu.RUnlock()

	globalLimiterMu.Lock()
	defer globalLimiterMu.Unlock()
	// Double-check after acquiring write lock
	if globalLimiter == nil {
		globalLimiter = NewRateLimiter(DefaultRateLimiterConfig())
	}
	return globalLimiter
}

// SetGlobalRateLimiter sets a custom global rate limiter (for testing).
// Must be called before any GetGlobalRateLimiter calls in production.
func SetGlobalRateLimiter(rl *RateLimiter) {
	globalLimiterMu.Lock()
	defer globalLimiterMu.Unlock()
	globalLimiter = rl
}

// ResetGlobalRateLimiter resets the global rate limiter (for testing).
func ResetGlobalRateLimiter() {
	globalLimiterMu.Lock()
	defer globalLimiterMu.Unlock()
	globalLimiter = nil
}
