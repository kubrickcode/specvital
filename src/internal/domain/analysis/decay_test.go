package analysis

import (
	"testing"
	"time"
)

func TestCalculateRefreshIntervalAt(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		lastViewedAt time.Time
		want         time.Duration
	}{
		{
			name:         "0 days idle - 6 hours",
			lastViewedAt: now,
			want:         6 * time.Hour,
		},
		{
			name:         "3 days idle - 6 hours",
			lastViewedAt: now.AddDate(0, 0, -3),
			want:         6 * time.Hour,
		},
		{
			name:         "7 days idle - 6 hours",
			lastViewedAt: now.AddDate(0, 0, -7),
			want:         6 * time.Hour,
		},
		{
			name:         "8 days idle - 12 hours",
			lastViewedAt: now.AddDate(0, 0, -8),
			want:         12 * time.Hour,
		},
		{
			name:         "14 days idle - 12 hours",
			lastViewedAt: now.AddDate(0, 0, -14),
			want:         12 * time.Hour,
		},
		{
			name:         "15 days idle - 24 hours",
			lastViewedAt: now.AddDate(0, 0, -15),
			want:         24 * time.Hour,
		},
		{
			name:         "30 days idle - 24 hours",
			lastViewedAt: now.AddDate(0, 0, -30),
			want:         24 * time.Hour,
		},
		{
			name:         "31 days idle - 3 days",
			lastViewedAt: now.AddDate(0, 0, -31),
			want:         3 * 24 * time.Hour,
		},
		{
			name:         "60 days idle - 3 days",
			lastViewedAt: now.AddDate(0, 0, -60),
			want:         3 * 24 * time.Hour,
		},
		{
			name:         "61 days idle - 7 days",
			lastViewedAt: now.AddDate(0, 0, -61),
			want:         7 * 24 * time.Hour,
		},
		{
			name:         "90 days idle - 7 days",
			lastViewedAt: now.AddDate(0, 0, -90),
			want:         7 * 24 * time.Hour,
		},
		{
			name:         "91 days idle - stop (0)",
			lastViewedAt: now.AddDate(0, 0, -91),
			want:         0,
		},
		{
			name:         "120 days idle - stop (0)",
			lastViewedAt: now.AddDate(0, 0, -120),
			want:         0,
		},
		// Edge cases: future timestamp
		{
			name:         "future timestamp - stop (0)",
			lastViewedAt: now.Add(1 * time.Hour),
			want:         0,
		},
		{
			name:         "future timestamp 1 day - stop (0)",
			lastViewedAt: now.AddDate(0, 0, 1),
			want:         0,
		},
		// Edge cases: boundary precision (floor division behavior)
		{
			name:         "7 days 23 hours - still 6 hours (floor)",
			lastViewedAt: now.Add(-7*24*time.Hour - 23*time.Hour),
			want:         6 * time.Hour,
		},
		{
			name:         "8 days 1 second - 12 hours",
			lastViewedAt: now.Add(-8*24*time.Hour - time.Second),
			want:         12 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateRefreshIntervalAt(tt.lastViewedAt, now)
			if got != tt.want {
				t.Errorf("CalculateRefreshIntervalAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldRefreshAt(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name                string
		lastViewedAt        time.Time
		lastCompletedAt     *time.Time
		consecutiveFailures int
		want                bool
	}{
		{
			name:                "never completed - should refresh",
			lastViewedAt:        now,
			lastCompletedAt:     nil,
			consecutiveFailures: 0,
			want:                true,
		},
		{
			name:                "recently completed within interval - no refresh",
			lastViewedAt:        now,
			lastCompletedAt:     timePtr(now.Add(-3 * time.Hour)),
			consecutiveFailures: 0,
			want:                false,
		},
		{
			name:                "completed beyond interval - should refresh",
			lastViewedAt:        now,
			lastCompletedAt:     timePtr(now.Add(-7 * time.Hour)),
			consecutiveFailures: 0,
			want:                true,
		},
		{
			name:                "exactly at interval boundary - should refresh",
			lastViewedAt:        now,
			lastCompletedAt:     timePtr(now.Add(-6 * time.Hour)),
			consecutiveFailures: 0,
			want:                true,
		},
		{
			name:                "5 consecutive failures - skip refresh",
			lastViewedAt:        now,
			lastCompletedAt:     nil,
			consecutiveFailures: 5,
			want:                false,
		},
		{
			name:                "4 consecutive failures - still refresh",
			lastViewedAt:        now,
			lastCompletedAt:     nil,
			consecutiveFailures: 4,
			want:                true,
		},
		{
			name:                "91 days idle - stop auto-refresh",
			lastViewedAt:        now.AddDate(0, 0, -91),
			lastCompletedAt:     nil,
			consecutiveFailures: 0,
			want:                false,
		},
		{
			name:                "30 days idle with old completion - should refresh",
			lastViewedAt:        now.AddDate(0, 0, -30),
			lastCompletedAt:     timePtr(now.AddDate(0, 0, -2)),
			consecutiveFailures: 0,
			want:                true,
		},
		{
			name:                "30 days idle with recent completion - no refresh",
			lastViewedAt:        now.AddDate(0, 0, -30),
			lastCompletedAt:     timePtr(now.Add(-12 * time.Hour)),
			consecutiveFailures: 0,
			want:                false,
		},
		{
			name:                "60 days idle - 3 day interval check",
			lastViewedAt:        now.AddDate(0, 0, -60),
			lastCompletedAt:     timePtr(now.AddDate(0, 0, -4)),
			consecutiveFailures: 0,
			want:                true,
		},
		{
			name:                "high failures with 91+ days idle - skip for both reasons",
			lastViewedAt:        now.AddDate(0, 0, -100),
			lastCompletedAt:     nil,
			consecutiveFailures: 10,
			want:                false,
		},
		// Edge cases: negative consecutive failures (treated as 0)
		{
			name:                "negative consecutive failures - should refresh",
			lastViewedAt:        now,
			lastCompletedAt:     nil,
			consecutiveFailures: -1,
			want:                true,
		},
		{
			name:                "negative consecutive failures with completed - check interval",
			lastViewedAt:        now,
			lastCompletedAt:     timePtr(now.Add(-7 * time.Hour)),
			consecutiveFailures: -5,
			want:                true,
		},
		// Edge case: future lastViewedAt
		{
			name:                "future lastViewedAt - no refresh",
			lastViewedAt:        now.Add(1 * time.Hour),
			lastCompletedAt:     nil,
			consecutiveFailures: 0,
			want:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldRefreshAt(tt.lastViewedAt, tt.lastCompletedAt, tt.consecutiveFailures, now)
			if got != tt.want {
				t.Errorf("ShouldRefreshAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
