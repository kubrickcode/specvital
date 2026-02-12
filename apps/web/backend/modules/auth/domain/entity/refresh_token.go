package entity

import "time"

type RefreshToken struct {
	CreatedAt time.Time
	ExpiresAt time.Time
	FamilyID  string
	ID        string
	Replaces  *string
	RevokedAt *time.Time
	TokenHash string
	UserID    string
}

func (rt *RefreshToken) IsExpired() bool {
	return rt.IsExpiredAt(time.Now())
}

func (rt *RefreshToken) IsExpiredAt(now time.Time) bool {
	return now.After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}
