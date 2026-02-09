package analysis

import (
	"context"
	"time"
)

// CommitInfo contains commit SHA and visibility information.
type CommitInfo struct {
	IsPrivate bool
	SHA       string
}

type VCS interface {
	Clone(ctx context.Context, url string, token *string) (Source, error)
	// GetHeadCommit returns the HEAD commit info (SHA and visibility) of the default branch.
	// It determines visibility by trying unauthenticated access first:
	// - Success without token = public repository (IsPrivate=false)
	// - Failure without token, success with token = private repository (IsPrivate=true)
	GetHeadCommit(ctx context.Context, url string, token *string) (CommitInfo, error)
}

type Source interface {
	Branch() string
	CommitSHA() string
	CommittedAt() time.Time
	Close(ctx context.Context) error
	// VerifyCommitExists checks if a commit SHA exists in the remote repository
	// by running "git fetch --depth 1 origin <sha>" on the cloned repository.
	// Returns true if the commit exists, false if not found (e.g., "not our ref").
	// This enables reanalysis verification without API calls.
	VerifyCommitExists(ctx context.Context, sha string) (bool, error)
}

type RepoInfo struct {
	ExternalRepoID string
	Name           string
	Owner          string
}

type VCSAPIClient interface {
	// Returns ErrRepoNotFound if the repository does not exist.
	GetRepoInfo(ctx context.Context, host, owner, repo string, token *string) (RepoInfo, error)
}
