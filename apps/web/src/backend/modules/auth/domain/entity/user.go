package entity

import "time"

type User struct {
	AvatarURL   string
	CreatedAt   time.Time
	Email       *string
	ID          string
	LastLoginAt *time.Time
	UpdatedAt   time.Time
	Username    string
}

type OAuthAccount struct {
	AccessToken      *string
	CreatedAt        time.Time
	ID               string
	Provider         string
	ProviderUserID   string
	ProviderUsername *string
	Scope            *string
	UpdatedAt        time.Time
	UserID           string
}

const ProviderGitHub = "github"

type OAuthUserInfo struct {
	AvatarURL  string
	Email      *string
	ExternalID string
	Username   string
}
