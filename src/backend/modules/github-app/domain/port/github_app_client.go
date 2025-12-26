package port

import (
	"context"
	"time"
)

type GitHubAppClient interface {
	CreateInstallationToken(ctx context.Context, installationID int64) (*InstallationToken, error)
	GetInstallationURL() string
}

type InstallationToken struct {
	ExpiresAt time.Time
	Token     string
}
