package adapter

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

const (
	rlDefaultLimit           = 60              // Gemini Flash-Lite: 30 RPM free tier, 60 gives headroom for paid tier
	rlDefaultWindow          = time.Minute     // Standard rate limit window (requests per minute)
	rlDefaultCleanupInterval = 5 * time.Minute // Cleanup stale buckets; 5min balances memory vs CPU
)

type RateLimiterConfig struct {
	Limit    int
	Window   time.Duration
	BurstMax int
}

type MemoryRateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*tokenBucket
	limit    int
	window   time.Duration
	burstMax int
	stopCh   chan struct{}
	stopped  sync.Once
}

type tokenBucket struct {
	tokens    int
	lastReset time.Time
}

var (
	_ port.RateLimiter = (*MemoryRateLimiter)(nil)
	_ io.Closer        = (*MemoryRateLimiter)(nil) // Caller MUST call Close() to prevent goroutine leak
)

func NewMemoryRateLimiter(cfg RateLimiterConfig) *MemoryRateLimiter {
	limit := cfg.Limit
	if limit <= 0 {
		limit = rlDefaultLimit
	}

	window := cfg.Window
	if window <= 0 {
		window = rlDefaultWindow
	}

	burstMax := cfg.BurstMax
	if burstMax <= 0 {
		burstMax = limit
	}

	rl := &MemoryRateLimiter{
		buckets:  make(map[string]*tokenBucket),
		limit:    limit,
		window:   window,
		burstMax: burstMax,
		stopCh:   make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

func (rl *MemoryRateLimiter) Allow(_ context.Context, key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		rl.buckets[key] = &tokenBucket{
			tokens:    rl.burstMax - 1,
			lastReset: time.Now(),
		}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(bucket.lastReset)
	if elapsed >= rl.window {
		bucket.tokens = rl.burstMax - 1
		bucket.lastReset = now
		return true
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

func (rl *MemoryRateLimiter) Remaining(_ context.Context, key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	bucket, exists := rl.buckets[key]
	if !exists {
		return rl.burstMax
	}

	now := time.Now()
	if now.Sub(bucket.lastReset) >= rl.window {
		return rl.burstMax
	}

	return bucket.tokens
}

func (rl *MemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(rlDefaultCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, bucket := range rl.buckets {
				if now.Sub(bucket.lastReset) > 2*rl.window {
					delete(rl.buckets, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *MemoryRateLimiter) Close() error {
	rl.stopped.Do(func() {
		close(rl.stopCh)
	})
	return nil
}
