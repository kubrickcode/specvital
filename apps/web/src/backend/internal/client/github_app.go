package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v75/github"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github-app/domain/port"
)

var _ port.GitHubAppClient = (*GitHubAppClient)(nil)

type GitHubAppClient struct {
	appID        int64
	appSlug      string
	appTransport *ghinstallation.AppsTransport
}

type GitHubAppConfig struct {
	AppID      int64
	AppSlug    string
	PrivateKey []byte
}

func NewGitHubAppClient(cfg GitHubAppConfig) (*GitHubAppClient, error) {
	if cfg.AppID == 0 {
		return nil, domain.ErrMissingAppID
	}
	if len(cfg.PrivateKey) == 0 {
		return nil, domain.ErrInvalidPrivateKey
	}

	appTransport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, cfg.AppID, cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidPrivateKey, err)
	}

	return &GitHubAppClient{
		appID:        cfg.AppID,
		appSlug:      cfg.AppSlug,
		appTransport: appTransport,
	}, nil
}

func (c *GitHubAppClient) CreateInstallationToken(ctx context.Context, installationID int64) (*port.InstallationToken, error) {
	itr := ghinstallation.NewFromAppsTransport(c.appTransport, installationID)

	token, err := itr.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	expiresAt, _, err := itr.Expiry()
	if err != nil {
		expiresAt = time.Now().Add(55 * time.Minute)
	}

	return &port.InstallationToken{
		ExpiresAt: expiresAt,
		Token:     token,
	}, nil
}

func (c *GitHubAppClient) GetInstallationURL() string {
	if c.appSlug != "" {
		return fmt.Sprintf("https://github.com/apps/%s/installations/new", c.appSlug)
	}
	return fmt.Sprintf("https://github.com/settings/apps/%d/installations", c.appID)
}

func (c *GitHubAppClient) NewInstallationClient(installationID int64) *gh.Client {
	itr := ghinstallation.NewFromAppsTransport(c.appTransport, installationID)
	return gh.NewClient(&http.Client{Transport: itr})
}
