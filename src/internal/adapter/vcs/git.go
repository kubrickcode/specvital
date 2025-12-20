package vcs

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/specvital/collector/internal/domain/analysis"
	"github.com/specvital/core/pkg/source"
)

// GitVCS implements analysis.VCS using specvital/core's GitSource.
// It is a thin, stateless adapter that delegates to the underlying source package.
// Concurrency control (semaphore) is managed by the use case layer, not here.
type GitVCS struct{}

// NewGitVCS creates a new GitVCS.
func NewGitVCS() *GitVCS {
	return &GitVCS{}
}

// Clone implements analysis.VCS by cloning a Git repository.
func (v *GitVCS) Clone(ctx context.Context, url string, token *string) (analysis.Source, error) {
	if url == "" {
		return nil, fmt.Errorf("clone repository: URL is required")
	}

	var opts *source.GitOptions
	if token != nil {
		opts = &source.GitOptions{
			Credentials: &source.GitCredentials{
				Username: "x-access-token",
				Password: *token,
			},
		}
	}

	gitSrc, err := source.NewGitSource(ctx, url, opts)
	if err != nil {
		return nil, fmt.Errorf("clone repository %q: %w", url, err)
	}

	return &gitSourceAdapter{gitSrc: gitSrc}, nil
}

// GetHeadCommit returns the HEAD commit SHA of the default branch using git ls-remote.
func (v *GitVCS) GetHeadCommit(ctx context.Context, url string, token *string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("get head commit: URL is required")
	}

	args := []string{"ls-remote", url, "HEAD"}

	cmd := exec.CommandContext(ctx, "git", args...)
	if token != nil {
		authURL := strings.Replace(url, "https://", fmt.Sprintf("https://x-access-token:%s@", *token), 1)
		cmd.Args = []string{"git", "ls-remote", authURL, "HEAD"}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git ls-remote %q: %s: %w", url, stderr.String(), err)
	}

	// Output format: "<sha>\tHEAD\n"
	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return "", fmt.Errorf("git ls-remote %q: empty response", url)
	}

	parts := strings.Fields(output)
	if len(parts) < 1 {
		return "", fmt.Errorf("git ls-remote %q: unexpected output format: %s", url, output)
	}

	return parts[0], nil
}

// gitSourceAdapter adapts source.GitSource to implement analysis.Source.
// It also provides access to the underlying source.Source for parser integration.
type gitSourceAdapter struct {
	gitSrc *source.GitSource
}

func (a *gitSourceAdapter) Branch() string {
	return a.gitSrc.Branch()
}

func (a *gitSourceAdapter) CommitSHA() string {
	return a.gitSrc.CommitSHA()
}

func (a *gitSourceAdapter) Close(_ context.Context) error {
	return a.gitSrc.Close()
}

// CoreSource returns the underlying source.Source for use by the parser adapter.
// This allows the parser to access the core source interface without exposing
// implementation details in the domain layer.
func (a *gitSourceAdapter) CoreSource() source.Source {
	return a.gitSrc
}
