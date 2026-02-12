package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
)

type GitClient interface {
	GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error)
	GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error)
}

var (
	ErrRepoNotFound    = errors.New("repository not found")
	ErrForbidden       = errors.New("access forbidden")
	ErrInvalidResponse = errors.New("invalid response from git")
)

type gitClient struct{}

func NewGitClient() GitClient {
	return &gitClient{}
}

func (c *gitClient) GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error) {
	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)
	return c.runLsRemote(ctx, repoURL, owner, repo)
}

func (c *gitClient) GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error) {
	repoURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", token, owner, repo)
	return c.runLsRemote(ctx, repoURL, owner, repo)
}

func (c *gitClient) runLsRemote(ctx context.Context, repoURL, owner, repo string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "ls-remote", repoURL, "HEAD")
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GIT_ASKPASS=",
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
	)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "not found") || strings.Contains(stderr, "Repository not found") {
				return "", errors.Wrap(ErrRepoNotFound, fmt.Sprintf("%s/%s", owner, repo))
			}
			if strings.Contains(stderr, "could not read Username") || strings.Contains(stderr, "Authentication failed") {
				return "", errors.Wrap(ErrForbidden, fmt.Sprintf("%s/%s", owner, repo))
			}
		}
		return "", errors.Wrapf(ErrInvalidResponse, "git ls-remote failed for %s/%s: %v", owner, repo, err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			return parts[0], nil
		}
	}

	return "", errors.Wrap(ErrInvalidResponse, "no commit SHA in output")
}
