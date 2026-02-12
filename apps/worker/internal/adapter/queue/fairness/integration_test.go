package fairness

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFairnessIntegration_ProUserConcurrency verifies Pro user can run 3 jobs concurrently
func TestFairnessIntegration_ProUserConcurrency(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
	limiter, err := NewPerUserLimiter(cfg)
	assert.NoError(t, err)

	userID := "pro-user-123"
	tier := TierPro

	// Acquire 3 jobs (should succeed)
	assert.True(t, limiter.TryAcquire(userID, tier, 1))
	assert.True(t, limiter.TryAcquire(userID, tier, 2))
	assert.True(t, limiter.TryAcquire(userID, tier, 3))

	// 4th job should fail (limit=3)
	assert.False(t, limiter.TryAcquire(userID, tier, 4))

	// Release one job
	limiter.Release(userID, 1)

	// 4th job should now succeed
	assert.True(t, limiter.TryAcquire(userID, tier, 4))
}

// TestFairnessIntegration_FreeUserConcurrency verifies Free user can run 1 job only
func TestFairnessIntegration_FreeUserConcurrency(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
	limiter, err := NewPerUserLimiter(cfg)
	assert.NoError(t, err)

	userID := "free-user-456"
	tier := TierFree

	// 1st job succeeds
	assert.True(t, limiter.TryAcquire(userID, tier, 1))

	// 2nd job fails (limit=1)
	assert.False(t, limiter.TryAcquire(userID, tier, 2))

	// Release 1st job
	limiter.Release(userID, 1)

	// 2nd job can now run
	assert.True(t, limiter.TryAcquire(userID, tier, 2))
}

// TestFairnessIntegration_DifferentUsersIndependent verifies users don't affect each other
func TestFairnessIntegration_DifferentUsersIndependent(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
	limiter, err := NewPerUserLimiter(cfg)
	assert.NoError(t, err)

	userA := "user-a"
	userB := "user-b"

	// User A acquires 3 Pro slots
	assert.True(t, limiter.TryAcquire(userA, TierPro, 1))
	assert.True(t, limiter.TryAcquire(userA, TierPro, 2))
	assert.True(t, limiter.TryAcquire(userA, TierPro, 3))
	assert.False(t, limiter.TryAcquire(userA, TierPro, 4)) // Exceeds limit

	// User B can still acquire independently
	assert.True(t, limiter.TryAcquire(userB, TierPro, 10))
	assert.True(t, limiter.TryAcquire(userB, TierPro, 11))
	assert.True(t, limiter.TryAcquire(userB, TierPro, 12))
	assert.False(t, limiter.TryAcquire(userB, TierPro, 13)) // Exceeds limit

	// User A and B have separate slot tracking
	assert.Equal(t, 3, limiter.ActiveCount(userA))
	assert.Equal(t, 3, limiter.ActiveCount(userB))
}

// TestFairnessIntegration_EnterpriseUserConcurrency verifies Enterprise user can run 5 jobs
func TestFairnessIntegration_EnterpriseUserConcurrency(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
	limiter, err := NewPerUserLimiter(cfg)
	assert.NoError(t, err)

	userID := "enterprise-user-789"
	tier := TierEnterprise

	// Acquire 5 jobs (should succeed)
	for i := int64(1); i <= 5; i++ {
		assert.True(t, limiter.TryAcquire(userID, tier, i))
	}

	// 6th job should fail (limit=5)
	assert.False(t, limiter.TryAcquire(userID, tier, 6))

	// Active count verification
	assert.Equal(t, 5, limiter.ActiveCount(userID))

	// Release all jobs
	for i := int64(1); i <= 5; i++ {
		limiter.Release(userID, i)
	}

	// After cleanup, active count should be 0
	assert.Equal(t, 0, limiter.ActiveCount(userID))
}

// TestFairnessIntegration_RaceCondition verifies concurrent acquire/release is safe
func TestFairnessIntegration_RaceCondition(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
	limiter, err := NewPerUserLimiter(cfg)
	assert.NoError(t, err)

	userID := "concurrent-user"
	tier := TierPro

	const goroutines = 50
	successCh := make(chan bool, goroutines)

	// 50 goroutines try to acquire concurrently
	for i := 0; i < goroutines; i++ {
		go func(jobID int64) {
			acquired := limiter.TryAcquire(userID, tier, jobID)
			successCh <- acquired

			if acquired {
				// Simulate work
				time.Sleep(10 * time.Millisecond)
				limiter.Release(userID, jobID)
			}
		}(int64(i))
	}

	// Count successful acquisitions
	successCount := 0
	for i := 0; i < goroutines; i++ {
		if <-successCh {
			successCount++
		}
	}

	// Due to concurrent execution, we can't predict exact success count
	// but it should be reasonable (at least some succeeded)
	assert.Greater(t, successCount, 0, "At least some jobs should succeed")
	assert.LessOrEqual(t, successCount, goroutines, "Success count should not exceed total")

	// Final cleanup verification
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, limiter.ActiveCount(userID), "All jobs should be released")
}

// TestFairnessIntegration_MixedTierScenario simulates realistic mixed user scenario
func TestFairnessIntegration_MixedTierScenario(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
	limiter, err := NewPerUserLimiter(cfg)
	assert.NoError(t, err)

	// Free user submits 3 jobs
	freeUser := "free-123"
	assert.True(t, limiter.TryAcquire(freeUser, TierFree, 1))
	assert.False(t, limiter.TryAcquire(freeUser, TierFree, 2)) // Blocked
	assert.False(t, limiter.TryAcquire(freeUser, TierFree, 3)) // Blocked

	// Pro user submits 5 jobs
	proUser := "pro-456"
	assert.True(t, limiter.TryAcquire(proUser, TierPro, 10))
	assert.True(t, limiter.TryAcquire(proUser, TierPro, 11))
	assert.True(t, limiter.TryAcquire(proUser, TierPro, 12))
	assert.False(t, limiter.TryAcquire(proUser, TierPro, 13)) // Blocked
	assert.False(t, limiter.TryAcquire(proUser, TierPro, 14)) // Blocked

	// Enterprise user submits 7 jobs
	enterpriseUser := "enterprise-789"
	for i := int64(20); i < 25; i++ {
		assert.True(t, limiter.TryAcquire(enterpriseUser, TierEnterprise, i))
	}
	assert.False(t, limiter.TryAcquire(enterpriseUser, TierEnterprise, 25)) // Blocked
	assert.False(t, limiter.TryAcquire(enterpriseUser, TierEnterprise, 26)) // Blocked

	// Verify active counts
	assert.Equal(t, 1, limiter.ActiveCount(freeUser))
	assert.Equal(t, 3, limiter.ActiveCount(proUser))
	assert.Equal(t, 5, limiter.ActiveCount(enterpriseUser))

	// Free user completes job → blocked job can proceed
	limiter.Release(freeUser, 1)
	assert.True(t, limiter.TryAcquire(freeUser, TierFree, 2))

	// Pro user completes 2 jobs → 2 blocked jobs can proceed
	limiter.Release(proUser, 10)
	limiter.Release(proUser, 11)
	assert.True(t, limiter.TryAcquire(proUser, TierPro, 13))
	assert.True(t, limiter.TryAcquire(proUser, TierPro, 14))
}
