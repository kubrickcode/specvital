package adapter

import (
	"context"
	"testing"
	"time"
)

func TestMemoryRateLimiter_AllowWithinLimit(t *testing.T) {
	rl := NewMemoryRateLimiter(RateLimiterConfig{
		Limit:    3,
		Window:   time.Second,
		BurstMax: 3,
	})
	defer rl.Close()

	ctx := context.Background()
	key := "test-user"

	for i := range 3 {
		if !rl.Allow(ctx, key) {
			t.Errorf("expected Allow() to return true for request %d", i+1)
		}
	}

	if rl.Allow(ctx, key) {
		t.Error("expected Allow() to return false after limit exceeded")
	}
}

func TestMemoryRateLimiter_ResetAfterWindow(t *testing.T) {
	rl := NewMemoryRateLimiter(RateLimiterConfig{
		Limit:    2,
		Window:   50 * time.Millisecond,
		BurstMax: 2,
	})
	defer rl.Close()

	ctx := context.Background()
	key := "test-user"

	rl.Allow(ctx, key)
	rl.Allow(ctx, key)

	if rl.Allow(ctx, key) {
		t.Error("expected Allow() to return false after limit exceeded")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow(ctx, key) {
		t.Error("expected Allow() to return true after window reset")
	}
}

func TestMemoryRateLimiter_SeparateKeys(t *testing.T) {
	rl := NewMemoryRateLimiter(RateLimiterConfig{
		Limit:    1,
		Window:   time.Second,
		BurstMax: 1,
	})
	defer rl.Close()

	ctx := context.Background()

	if !rl.Allow(ctx, "user1") {
		t.Error("expected Allow() for user1 to return true")
	}
	if !rl.Allow(ctx, "user2") {
		t.Error("expected Allow() for user2 to return true")
	}

	if rl.Allow(ctx, "user1") {
		t.Error("expected Allow() for user1 to return false after limit")
	}
	if rl.Allow(ctx, "user2") {
		t.Error("expected Allow() for user2 to return false after limit")
	}
}

func TestMemoryRateLimiter_Remaining(t *testing.T) {
	rl := NewMemoryRateLimiter(RateLimiterConfig{
		Limit:    5,
		Window:   time.Second,
		BurstMax: 5,
	})
	defer rl.Close()

	ctx := context.Background()
	key := "test-user"

	if remaining := rl.Remaining(ctx, key); remaining != 5 {
		t.Errorf("expected Remaining() to return 5 for new key, got %d", remaining)
	}

	rl.Allow(ctx, key)
	rl.Allow(ctx, key)

	if remaining := rl.Remaining(ctx, key); remaining != 3 {
		t.Errorf("expected Remaining() to return 3 after 2 allows, got %d", remaining)
	}
}

func TestMemoryRateLimiter_RemainingAfterWindowReset(t *testing.T) {
	rl := NewMemoryRateLimiter(RateLimiterConfig{
		Limit:    3,
		Window:   50 * time.Millisecond,
		BurstMax: 3,
	})
	defer rl.Close()

	ctx := context.Background()
	key := "test-user"

	rl.Allow(ctx, key)
	rl.Allow(ctx, key)

	time.Sleep(60 * time.Millisecond)

	if remaining := rl.Remaining(ctx, key); remaining != 3 {
		t.Errorf("expected Remaining() to return 3 after window reset, got %d", remaining)
	}
}

func TestMemoryRateLimiter_DefaultValues(t *testing.T) {
	rl := NewMemoryRateLimiter(RateLimiterConfig{})
	defer rl.Close()

	ctx := context.Background()
	key := "test-user"

	for range 60 {
		rl.Allow(ctx, key)
	}

	if rl.Allow(ctx, key) {
		t.Error("expected default limit to be 60")
	}
}
