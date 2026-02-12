package retention

import (
	"testing"
	"time"
)

func TestNewPolicy(t *testing.T) {
	tests := []struct {
		name     string
		days     int
		wantDays int
	}{
		{"positive days", 90, 90},
		{"zero uses default", 0, DefaultRetentionDays},
		{"negative uses default", -1, DefaultRetentionDays},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPolicy(tt.days)
			if got := p.RetentionDays(); got != tt.wantDays {
				t.Errorf("RetentionDays() = %d, want %d", got, tt.wantDays)
			}
		})
	}
}

func TestDefaultPolicy(t *testing.T) {
	p := DefaultPolicy()
	if got := p.RetentionDays(); got != DefaultRetentionDays {
		t.Errorf("RetentionDays() = %d, want %d", got, DefaultRetentionDays)
	}
}

func TestPolicy_ExpirationDate(t *testing.T) {
	p := NewPolicy(30)
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	want := time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC)
	got := p.ExpirationDate(createdAt)

	if !got.Equal(want) {
		t.Errorf("ExpirationDate() = %v, want %v", got, want)
	}
}

func TestPolicy_IsExpired(t *testing.T) {
	p := NewPolicy(30)
	createdAt := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{
			name: "before expiration",
			now:  time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "exactly at expiration",
			now:  time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "after expiration",
			now:  time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.IsExpired(createdAt, tt.now); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicy_CutoffTime(t *testing.T) {
	p := NewPolicy(30)
	now := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)

	want := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	got := p.CutoffTime(now)

	if !got.Equal(want) {
		t.Errorf("CutoffTime() = %v, want %v", got, want)
	}
}
