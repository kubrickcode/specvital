package reliability

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter_AllowsFirstRequest(t *testing.T) {
	config := RateLimiterConfig{
		RPM:         60, // 1 per second
		BurstFactor: 1,
	}
	rl := NewRateLimiter(config)

	ctx := context.Background()
	err := rl.Wait(ctx)
	if err != nil {
		t.Errorf("expected first request to be allowed, got error: %v", err)
	}
}

func TestRateLimiter_RespectsRateLimit(t *testing.T) {
	config := RateLimiterConfig{
		RPM:         600, // 10 per second
		BurstFactor: 1,
	}
	rl := NewRateLimiter(config)

	ctx := context.Background()

	// First request should be immediate
	start := time.Now()
	err := rl.Wait(ctx)
	if err != nil {
		t.Errorf("expected first request to succeed, got error: %v", err)
	}
	firstDuration := time.Since(start)
	if firstDuration > 10*time.Millisecond {
		t.Errorf("first request took too long: %v", firstDuration)
	}

	// Second request should wait
	start = time.Now()
	err = rl.Wait(ctx)
	if err != nil {
		t.Errorf("expected second request to succeed, got error: %v", err)
	}
	secondDuration := time.Since(start)
	// Should wait approximately 100ms (1/10 second)
	if secondDuration < 50*time.Millisecond {
		t.Errorf("second request should have waited, but only took: %v", secondDuration)
	}
}

func TestRateLimiter_RespectsContextCancellation(t *testing.T) {
	config := RateLimiterConfig{
		RPM:         1, // Very slow
		BurstFactor: 1,
	}
	rl := NewRateLimiter(config)

	ctx := context.Background()
	// Consume the initial token
	_ = rl.Wait(ctx)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Second request should be cancelled
	err := rl.Wait(ctx)
	if err == nil {
		t.Error("expected context deadline exceeded error")
	}
}

func TestRateLimiter_Available(t *testing.T) {
	config := RateLimiterConfig{
		RPM:         60,
		BurstFactor: 2,
	}
	rl := NewRateLimiter(config)

	available := rl.Available()
	if available != 2.0 {
		t.Errorf("expected 2.0 tokens available, got %f", available)
	}

	ctx := context.Background()
	_ = rl.Wait(ctx)

	available = rl.Available()
	if available >= 2.0 {
		t.Errorf("expected less than 2.0 tokens after wait, got %f", available)
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	config := RateLimiterConfig{
		RPM:         600, // 10 per second
		BurstFactor: 1,
	}
	rl := NewRateLimiter(config)

	ctx := context.Background()
	// Consume initial token
	_ = rl.Wait(ctx)

	available := rl.Available()
	if available >= 1.0 {
		t.Errorf("expected less than 1.0 tokens after consumption, got %f", available)
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	available = rl.Available()
	if available < 1.0 {
		t.Errorf("expected at least 1.0 tokens after refill, got %f", available)
	}
}

func TestRateLimiter_DefaultConfig(t *testing.T) {
	config := DefaultRateLimiterConfig()

	if config.RPM != 1800 {
		t.Errorf("expected default RPM to be 1800, got %d", config.RPM)
	}
	if config.BurstFactor != 1 {
		t.Errorf("expected default BurstFactor to be 1, got %d", config.BurstFactor)
	}
}

func TestGetGlobalRateLimiter_ReturnsSingleton(t *testing.T) {
	t.Cleanup(ResetGlobalRateLimiter)

	rl1 := GetGlobalRateLimiter()
	rl2 := GetGlobalRateLimiter()

	if rl1 != rl2 {
		t.Error("expected GetGlobalRateLimiter to return the same instance")
	}
}
