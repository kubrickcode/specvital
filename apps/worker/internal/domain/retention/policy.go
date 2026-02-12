package retention

import "time"

const (
	DefaultRetentionDays = 30
	DefaultBatchSize     = 1000
)

// Policy defines retention rules for user data cleanup.
type Policy struct {
	retentionDays int
}

// NewPolicy creates a Policy with the specified retention days.
// If days <= 0, DefaultRetentionDays is used.
func NewPolicy(days int) Policy {
	if days <= 0 {
		days = DefaultRetentionDays
	}
	return Policy{retentionDays: days}
}

// DefaultPolicy returns a Policy with default settings.
func DefaultPolicy() Policy {
	return NewPolicy(DefaultRetentionDays)
}

// RetentionDays returns the configured retention period.
func (p Policy) RetentionDays() int {
	return p.retentionDays
}

// ExpirationDate calculates when data created at createdAt should expire.
func (p Policy) ExpirationDate(createdAt time.Time) time.Time {
	return createdAt.AddDate(0, 0, p.retentionDays)
}

// IsExpired checks if data created at createdAt has expired as of now.
func (p Policy) IsExpired(createdAt time.Time, now time.Time) bool {
	return now.After(p.ExpirationDate(createdAt))
}

// CutoffTime calculates the cutoff time for cleanup.
// Records created before this time are eligible for deletion.
func (p Policy) CutoffTime(now time.Time) time.Time {
	return now.AddDate(0, 0, -p.retentionDays)
}
