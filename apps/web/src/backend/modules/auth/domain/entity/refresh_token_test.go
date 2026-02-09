package entity

import (
	"testing"
	"time"
)

func TestRefreshToken_IsExpiredAt(t *testing.T) {
	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "expired token",
			expiresAt: now.Add(-1 * time.Hour),
			want:      true,
		},
		{
			name:      "not expired token",
			expiresAt: now.Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "expires exactly now",
			expiresAt: now,
			want:      false,
		},
		{
			name:      "expires 1 nanosecond ago",
			expiresAt: now.Add(-1 * time.Nanosecond),
			want:      true,
		},
		{
			name:      "zero time",
			expiresAt: time.Time{},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RefreshToken{ExpiresAt: tt.expiresAt}
			if got := rt.IsExpiredAt(now); got != tt.want {
				t.Errorf("RefreshToken.IsExpiredAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken_IsRevoked(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		revokedAt *time.Time
		want      bool
	}{
		{
			name:      "revoked token",
			revokedAt: &now,
			want:      true,
		},
		{
			name:      "not revoked token",
			revokedAt: nil,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RefreshToken{RevokedAt: tt.revokedAt}
			if got := rt.IsRevoked(); got != tt.want {
				t.Errorf("RefreshToken.IsRevoked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken_IsValid(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name      string
		expiresAt time.Time
		revokedAt *time.Time
		want      bool
	}{
		{
			name:      "valid token",
			expiresAt: future,
			revokedAt: nil,
			want:      true,
		},
		{
			name:      "expired token",
			expiresAt: past,
			revokedAt: nil,
			want:      false,
		},
		{
			name:      "revoked token",
			expiresAt: future,
			revokedAt: &now,
			want:      false,
		},
		{
			name:      "expired and revoked token",
			expiresAt: past,
			revokedAt: &now,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &RefreshToken{
				ExpiresAt: tt.expiresAt,
				RevokedAt: tt.revokedAt,
			}
			if got := rt.IsValid(); got != tt.want {
				t.Errorf("RefreshToken.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
