package fairness

import "time"

// Config defines per-tier concurrent job limits and snooze parameters for fairness control.
type Config struct {
	FreeConcurrentLimit       int           // Free tier concurrent limit
	ProConcurrentLimit        int           // Pro/ProPlus tier concurrent limit
	EnterpriseConcurrentLimit int           // Enterprise tier concurrent limit
	SnoozeDuration            time.Duration // Base delay for River JobSnooze when limit exceeded
	SnoozeJitter              time.Duration // Random jitter added to SnoozeDuration
}

// DefaultConfig returns a Config with recommended production values.
func DefaultConfig() *Config {
	return &Config{
		FreeConcurrentLimit:       1,
		ProConcurrentLimit:        3,
		EnterpriseConcurrentLimit: 5,
		SnoozeDuration:            30 * time.Second,
		SnoozeJitter:              10 * time.Second,
	}
}

// PlanTier represents a user's subscription tier for concurrent limit enforcement.
type PlanTier string

const (
	TierFree       PlanTier = "free"       // Free tier users
	TierPro        PlanTier = "pro"        // Pro tier users
	TierProPlus    PlanTier = "pro_plus"   // Pro Plus tier users
	TierEnterprise PlanTier = "enterprise" // Enterprise tier users
)
