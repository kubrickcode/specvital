package fairness

import (
	"fmt"
	"sync"
)

type userSlot struct {
	count int
	jobs  map[int64]struct{} // Track job IDs (prevent double release)
}

// PerUserLimiter enforces per-user concurrent job limits based on subscription tier.
// Safe for concurrent use.
type PerUserLimiter struct {
	mu     sync.Mutex
	slots  map[string]*userSlot
	limits map[PlanTier]int
}

// NewPerUserLimiter creates a new limiter with tier-based limits from config.
// Returns error if any limit value is invalid.
func NewPerUserLimiter(cfg *Config) (*PerUserLimiter, error) {
	if cfg.FreeConcurrentLimit <= 0 {
		return nil, fmt.Errorf("FreeConcurrentLimit must be positive, got %d", cfg.FreeConcurrentLimit)
	}
	if cfg.ProConcurrentLimit <= 0 {
		return nil, fmt.Errorf("ProConcurrentLimit must be positive, got %d", cfg.ProConcurrentLimit)
	}
	if cfg.EnterpriseConcurrentLimit <= 0 {
		return nil, fmt.Errorf("EnterpriseConcurrentLimit must be positive, got %d", cfg.EnterpriseConcurrentLimit)
	}

	return &PerUserLimiter{
		slots: make(map[string]*userSlot),
		limits: map[PlanTier]int{
			TierFree:       cfg.FreeConcurrentLimit,
			TierPro:        cfg.ProConcurrentLimit,
			TierProPlus:    cfg.ProConcurrentLimit,
			TierEnterprise: cfg.EnterpriseConcurrentLimit,
		},
	}, nil
}

// TryAcquire attempts to acquire a concurrent execution slot for the user.
// Returns true if acquired, false if limit exceeded.
// Idempotent: calling with the same jobID multiple times returns true without incrementing count.
func (l *PerUserLimiter) TryAcquire(userID string, tier PlanTier, jobID int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	limit := l.limits[tier]
	if limit <= 0 {
		limit = l.limits[TierFree] // fallback to free tier
	}

	slot := l.slots[userID]
	if slot == nil {
		slot = &userSlot{
			count: 0,
			jobs:  make(map[int64]struct{}),
		}
		l.slots[userID] = slot
	}

	// Check if already acquired (idempotent)
	if _, exists := slot.jobs[jobID]; exists {
		return true
	}

	// Check limit
	if slot.count >= limit {
		return false
	}

	// Acquire slot
	slot.jobs[jobID] = struct{}{}
	slot.count++
	return true
}

// Release releases a concurrent execution slot for the user.
// Idempotent: calling with the same jobID multiple times is safely ignored.
// Automatically cleans up user entry when count reaches zero.
func (l *PerUserLimiter) Release(userID string, jobID int64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	slot := l.slots[userID]
	if slot == nil {
		return
	}

	// Check if job exists (ignore double release)
	if _, exists := slot.jobs[jobID]; !exists {
		return
	}

	delete(slot.jobs, jobID)
	slot.count--

	// Cleanup when count reaches 0
	if slot.count <= 0 {
		delete(l.slots, userID)
	}
}

// ActiveCount returns the current number of active jobs for the user.
// Returns 0 if the user has no active jobs.
func (l *PerUserLimiter) ActiveCount(userID string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	slot := l.slots[userID]
	if slot == nil {
		return 0
	}
	return slot.count
}
