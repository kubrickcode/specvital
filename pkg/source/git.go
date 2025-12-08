package source

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// GitCredentials contains authentication information for Git operations.
type GitCredentials struct {
	Username string
	Password string
}

// GitOptions configures how a Git repository is cloned.
type GitOptions struct {
	Branch      string
	Depth       int
	Credentials *GitCredentials
}

// Default clone depth for shallow clones.
const defaultCloneDepth = 1

// GitSource implements Source for Git repository access.
// It clones the repository to a temporary directory and provides
// filesystem-like access to its contents.
type GitSource struct {
	closeErr  error
	closeOnce sync.Once
	local     *LocalSource
	tempDir   string
}

// NewGitSource clones a Git repository and returns a Source for accessing its files.
// The repository is cloned to a temporary directory with shallow clone (depth 1) by default.
// The caller must call Close() to clean up the cloned repository.
func NewGitSource(ctx context.Context, repoURL string, opts *GitOptions) (*GitSource, error) {
	if err := VerifyGitInstalled(); err != nil {
		return nil, err
	}

	if opts == nil {
		opts = &GitOptions{}
	}
	if opts.Depth == 0 {
		opts.Depth = defaultCloneDepth
	}

	cloneURL, err := injectCredentials(repoURL, opts.Credentials)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid repository URL: %v", ErrInvalidPath, err)
	}

	tempDir, err := os.MkdirTemp("", "gitsource-*")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create temp directory: %v", ErrGitCloneFailed, err)
	}

	if err := os.Chmod(tempDir, 0700); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("%w: failed to secure temp directory: %v", ErrGitCloneFailed, err)
	}

	if err := cloneRepository(ctx, cloneURL, tempDir, opts); err != nil {
		os.RemoveAll(tempDir)
		return nil, sanitizeError(err, repoURL, opts.Credentials)
	}

	local, err := NewLocalSource(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, sanitizeError(
			fmt.Errorf("%w: failed to create local source: %v", ErrGitCloneFailed, err),
			repoURL, opts.Credentials,
		)
	}

	return &GitSource{
		local:   local,
		tempDir: tempDir,
	}, nil
}

// Root returns the path to the cloned repository.
func (s *GitSource) Root() string {
	return s.local.Root()
}

// Open opens the file at the given path for reading.
func (s *GitSource) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.local.Open(ctx, path)
}

// Stat returns file info for the given path.
func (s *GitSource) Stat(ctx context.Context, path string) (fs.FileInfo, error) {
	return s.local.Stat(ctx, path)
}

// Close removes the cloned repository and releases resources.
// Close is idempotent; calling it multiple times has no additional effect.
func (s *GitSource) Close() error {
	s.closeOnce.Do(func() {
		s.closeErr = os.RemoveAll(s.tempDir)
	})
	return s.closeErr
}

// VerifyGitInstalled checks if git is available in the system PATH.
func VerifyGitInstalled() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("%w: git not found in PATH", ErrGitCloneFailed)
	}
	return nil
}

// VerifyRepository checks if a repository exists and is accessible.
// It uses git ls-remote which doesn't count against API rate limits.
func VerifyRepository(ctx context.Context, repoURL string, creds *GitCredentials) error {
	if err := VerifyGitInstalled(); err != nil {
		return err
	}

	checkURL, err := injectCredentials(repoURL, creds)
	if err != nil {
		return fmt.Errorf("%w: invalid repository URL: %v", ErrInvalidPath, err)
	}

	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--exit-code", checkURL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("%w: %v", ErrRepositoryNotFound, ctx.Err())
		}
		return fmt.Errorf("%w: %s", ErrRepositoryNotFound, sanitizeOutput(stderr.String(), creds))
	}

	return nil
}

// validateBranchName checks if a branch name is safe to use in git commands.
// It rejects branch names that could be interpreted as command-line options.
func validateBranchName(branch string) error {
	if strings.HasPrefix(branch, "-") {
		return fmt.Errorf("%w: branch name cannot start with '-'", ErrInvalidPath)
	}
	if strings.ContainsAny(branch, "\x00") {
		return fmt.Errorf("%w: branch name contains invalid characters", ErrInvalidPath)
	}
	return nil
}

// cloneRepository executes git clone with the given options.
func cloneRepository(ctx context.Context, cloneURL, destDir string, opts *GitOptions) error {
	args := []string{"clone", "--single-branch"}

	if opts.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", opts.Depth))
	}

	if opts.Branch != "" {
		if err := validateBranchName(opts.Branch); err != nil {
			return err
		}
		args = append(args, "--branch", opts.Branch)
	}

	args = append(args, cloneURL, destDir)

	cmd := exec.CommandContext(ctx, "git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("%w: %v", ErrGitCloneFailed, ctx.Err())
		}
		return fmt.Errorf("%w: %s", ErrGitCloneFailed, stderr.String())
	}

	return nil
}

// injectCredentials adds credentials to the repository URL if provided.
func injectCredentials(repoURL string, creds *GitCredentials) (string, error) {
	if creds == nil || (creds.Username == "" && creds.Password == "") {
		return repoURL, nil
	}

	parsed, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}

	if creds.Password != "" {
		if creds.Username != "" {
			parsed.User = url.UserPassword(creds.Username, creds.Password)
		} else {
			parsed.User = url.UserPassword("oauth2", creds.Password)
		}
	} else if creds.Username != "" {
		parsed.User = url.User(creds.Username)
	}

	return parsed.String(), nil
}

// sanitizeError removes credentials from error messages while preserving error chain.
func sanitizeError(err error, originalURL string, creds *GitCredentials) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errMsg = sanitizeOutput(errMsg, creds)

	parsed, parseErr := url.Parse(originalURL)
	if parseErr == nil && parsed.User != nil {
		cleanURL := *parsed
		cleanURL.User = nil
		errMsg = strings.ReplaceAll(errMsg, originalURL, cleanURL.String())
	}

	if errors.Is(err, ErrInvalidPath) {
		return fmt.Errorf("%w: %s", ErrInvalidPath, errMsg)
	}
	if errors.Is(err, ErrGitCloneFailed) {
		return fmt.Errorf("%w: %s", ErrGitCloneFailed, errMsg)
	}
	if errors.Is(err, ErrRepositoryNotFound) {
		return fmt.Errorf("%w: %s", ErrRepositoryNotFound, errMsg)
	}

	return fmt.Errorf("%s", errMsg)
}

// sanitizeOutput removes credential information from command output.
func sanitizeOutput(output string, creds *GitCredentials) string {
	if creds == nil {
		return output
	}

	if len(creds.Username) >= 3 {
		output = strings.ReplaceAll(output, creds.Username, "[REDACTED]")
	}
	if len(creds.Password) >= 3 {
		output = strings.ReplaceAll(output, creds.Password, "[REDACTED]")
		if encoded := url.QueryEscape(creds.Password); encoded != creds.Password {
			output = strings.ReplaceAll(output, encoded, "[REDACTED]")
		}
	}

	return output
}

// RepoPathFromURL extracts a filesystem-safe path from a repository URL.
// Example: "https://github.com/owner/repo.git" -> "github.com/owner/repo"
// Credentials in the URL are stripped and not included in error messages.
func RepoPathFromURL(repoURL string) (string, error) {
	parsed, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	parsed.User = nil

	host := parsed.Host
	path := strings.TrimSuffix(parsed.Path, ".git")
	path = strings.TrimPrefix(path, "/")

	if host == "" || path == "" {
		return "", fmt.Errorf("URL must have host and path: %s", parsed.String())
	}

	return filepath.Join(host, path), nil
}
