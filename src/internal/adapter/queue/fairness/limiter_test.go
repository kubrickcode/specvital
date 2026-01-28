package fairness

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerUserLimiter_BasicAcquireRelease(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "user123"
	jobID := int64(100)

	// should acquire when under limit
	acquired := limiter.TryAcquire(userID, TierFree, jobID)
	assert.True(t, acquired, "should acquire first job")
	assert.Equal(t, 1, limiter.ActiveCount(userID))

	// should release successfully
	limiter.Release(userID, jobID)
	assert.Equal(t, 0, limiter.ActiveCount(userID))
}

func TestPerUserLimiter_PerTierLimits(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	tests := []struct {
		name           string
		tier           PlanTier
		expectedLimit  int
		jobsToAcquire  int
		expectedResult []bool
	}{
		{
			name:          "Free tier limit 1",
			tier:          TierFree,
			expectedLimit: 1,
			jobsToAcquire: 3,
			// [0]=true (under limit), [1]=false, [2]=false (over limit)
			expectedResult: []bool{true, false, false},
		},
		{
			name:          "Pro tier limit 3",
			tier:          TierPro,
			expectedLimit: 3,
			jobsToAcquire: 5,
			// [0-2]=true (under limit), [3-4]=false (over limit)
			expectedResult: []bool{true, true, true, false, false},
		},
		{
			name:          "Enterprise tier limit 5",
			tier:          TierEnterprise,
			expectedLimit: 5,
			jobsToAcquire: 7,
			// [0-4]=true (under limit), [5-6]=false (over limit)
			expectedResult: []bool{true, true, true, true, true, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := "user_" + tt.name
			results := make([]bool, tt.jobsToAcquire)

			for i := range tt.jobsToAcquire {
				jobID := int64(i + 1)
				results[i] = limiter.TryAcquire(userID, tt.tier, jobID)
			}

			assert.Equal(t, tt.expectedResult, results)
			assert.Equal(t, tt.expectedLimit, limiter.ActiveCount(userID))

			// cleanup
			for i := range tt.jobsToAcquire {
				limiter.Release(userID, int64(i+1))
			}
		})
	}
}

func TestPerUserLimiter_Concurrency(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       10,
		ProConcurrentLimit:        10,
		EnterpriseConcurrentLimit: 10,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "concurrent_user"
	goroutines := 100
	jobsPerGoroutine := 5

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(offset int) {
			defer wg.Done()
			for j := range jobsPerGoroutine {
				jobID := int64(offset*jobsPerGoroutine + j)
				limiter.TryAcquire(userID, TierFree, jobID)
				time.Sleep(time.Microsecond)
				limiter.Release(userID, jobID)
			}
		}(i)
	}

	wg.Wait()

	// should cleanup all slots
	assert.Equal(t, 0, limiter.ActiveCount(userID))
}

func TestPerUserLimiter_DoubleReleaseIgnored(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       3,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 3,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "user_double_release"
	jobID := int64(1)

	// acquire
	acquired := limiter.TryAcquire(userID, TierFree, jobID)
	assert.True(t, acquired)
	assert.Equal(t, 1, limiter.ActiveCount(userID))

	// first release
	limiter.Release(userID, jobID)
	assert.Equal(t, 0, limiter.ActiveCount(userID))

	// second release should be ignored (not panic, not negative count)
	limiter.Release(userID, jobID)
	assert.Equal(t, 0, limiter.ActiveCount(userID))
}

func TestPerUserLimiter_CleanupWhenCountZero(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       2,
		ProConcurrentLimit:        2,
		EnterpriseConcurrentLimit: 2,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "user_cleanup"
	job1 := int64(1)
	job2 := int64(2)

	// acquire 2 jobs
	limiter.TryAcquire(userID, TierFree, job1)
	limiter.TryAcquire(userID, TierFree, job2)
	assert.Equal(t, 2, limiter.ActiveCount(userID))

	// release first job
	limiter.Release(userID, job1)
	assert.Equal(t, 1, limiter.ActiveCount(userID))

	// release second job - should cleanup userID from slots map
	limiter.Release(userID, job2)
	assert.Equal(t, 0, limiter.ActiveCount(userID))

	// verify slot is removed from internal map
	limiter.mu.Lock()
	_, exists := limiter.slots[userID]
	limiter.mu.Unlock()
	assert.False(t, exists, "slot should be removed when count reaches 0")
}

func TestPerUserLimiter_IdempotentAcquire(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       2,
		ProConcurrentLimit:        2,
		EnterpriseConcurrentLimit: 2,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "user_idempotent"
	jobID := int64(1)

	// first acquire
	acquired1 := limiter.TryAcquire(userID, TierFree, jobID)
	assert.True(t, acquired1)
	assert.Equal(t, 1, limiter.ActiveCount(userID))

	// second acquire with same jobID should be idempotent
	acquired2 := limiter.TryAcquire(userID, TierFree, jobID)
	assert.True(t, acquired2)
	assert.Equal(t, 1, limiter.ActiveCount(userID), "count should not increase")

	// release once
	limiter.Release(userID, jobID)
	assert.Equal(t, 0, limiter.ActiveCount(userID))
}

func TestPerUserLimiter_FallbackToFreeTier(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       2,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "user_unknown_tier"

	// unknown tier should fallback to free tier limit
	unknownTier := PlanTier("unknown")

	job1 := limiter.TryAcquire(userID, unknownTier, 1)
	job2 := limiter.TryAcquire(userID, unknownTier, 2)
	job3 := limiter.TryAcquire(userID, unknownTier, 3)

	assert.True(t, job1)
	assert.True(t, job2)
	assert.False(t, job3, "should use free tier limit (2)")
	assert.Equal(t, 2, limiter.ActiveCount(userID))
}

func TestPerUserLimiter_NegativeJobID(t *testing.T) {
	cfg := &Config{
		FreeConcurrentLimit:       2,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
	}
	limiter, err := NewPerUserLimiter(cfg)
	require.NoError(t, err)

	userID := "user_negative_jobid"

	// negative jobID should work (map keys allow any int64)
	acquired := limiter.TryAcquire(userID, TierFree, -1)
	assert.True(t, acquired, "should handle negative jobID")
	assert.Equal(t, 1, limiter.ActiveCount(userID))

	limiter.Release(userID, -1)
	assert.Equal(t, 0, limiter.ActiveCount(userID), "should properly release negative jobID")
}

func TestNewPerUserLimiter_ValidationErrors(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *Config
		errMsg string
	}{
		{
			name: "negative free limit",
			cfg: &Config{
				FreeConcurrentLimit:       -1,
				ProConcurrentLimit:        3,
				EnterpriseConcurrentLimit: 5,
			},
			errMsg: "FreeConcurrentLimit must be positive",
		},
		{
			name: "zero free limit",
			cfg: &Config{
				FreeConcurrentLimit:       0,
				ProConcurrentLimit:        3,
				EnterpriseConcurrentLimit: 5,
			},
			errMsg: "FreeConcurrentLimit must be positive",
		},
		{
			name: "negative pro limit",
			cfg: &Config{
				FreeConcurrentLimit:       1,
				ProConcurrentLimit:        -3,
				EnterpriseConcurrentLimit: 5,
			},
			errMsg: "ProConcurrentLimit must be positive",
		},
		{
			name: "negative enterprise limit",
			cfg: &Config{
				FreeConcurrentLimit:       1,
				ProConcurrentLimit:        3,
				EnterpriseConcurrentLimit: -5,
			},
			errMsg: "EnterpriseConcurrentLimit must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter, err := NewPerUserLimiter(tt.cfg)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
			assert.Nil(t, limiter)
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, 1, cfg.FreeConcurrentLimit)
	assert.Equal(t, 3, cfg.ProConcurrentLimit)
	assert.Equal(t, 5, cfg.EnterpriseConcurrentLimit)
	assert.Equal(t, 30*time.Second, cfg.SnoozeDuration)
	assert.Equal(t, 10*time.Second, cfg.SnoozeJitter)
}
