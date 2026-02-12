package port

import "context"

type GitClient interface {
	GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error)
	GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error)
}
