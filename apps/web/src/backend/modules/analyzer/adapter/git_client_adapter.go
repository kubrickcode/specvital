package adapter

import (
	"context"

	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

var _ port.GitClient = (*GitClientAdapter)(nil)

type GitClientAdapter struct {
	client client.GitClient
}

func NewGitClientAdapter(c client.GitClient) *GitClientAdapter {
	return &GitClientAdapter{client: c}
}

func (a *GitClientAdapter) GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error) {
	return a.client.GetLatestCommitSHA(ctx, owner, repo)
}

func (a *GitClientAdapter) GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error) {
	return a.client.GetLatestCommitSHAWithToken(ctx, owner, repo, token)
}
