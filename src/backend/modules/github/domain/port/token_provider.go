package port

import "context"

type TokenProvider interface {
	GetUserGitHubToken(ctx context.Context, userID string) (string, error)
}
