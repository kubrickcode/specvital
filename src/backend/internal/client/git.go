package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	defaultGitHubAPIURL = "https://api.github.com"
	maxResponseBodySize = 1 << 20
)

var defaultHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GitClient interface {
	GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error)
	GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error)
}

var (
	ErrRepoNotFound    = errors.New("repository not found")
	ErrForbidden       = errors.New("access forbidden")
	ErrInvalidResponse = errors.New("invalid response from git")
)

type gitClient struct {
	httpClient HTTPClient
	baseURL    string
}

func NewGitClient() GitClient {
	return &gitClient{
		httpClient: defaultHTTPClient,
		baseURL:    defaultGitHubAPIURL,
	}
}

func NewGitClientWithOptions(httpClient HTTPClient, baseURL string) GitClient {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	if baseURL == "" {
		baseURL = defaultGitHubAPIURL
	}
	return &gitClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *gitClient) GetLatestCommitSHA(ctx context.Context, owner, repo string) (string, error) {
	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)

	cmd := exec.CommandContext(ctx, "git", "ls-remote", repoURL, "HEAD")
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

func (c *gitClient) GetLatestCommitSHAWithToken(ctx context.Context, owner, repo, token string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/HEAD", c.baseURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "execute request")
	}
	defer resp.Body.Close()

	limitedBody := io.LimitReader(resp.Body, maxResponseBodySize)

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.Wrapf(ErrRepoNotFound, "%s/%s", owner, repo)
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return "", errors.Wrapf(ErrForbidden, "%s/%s", owner, repo)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(limitedBody)
		return "", errors.Wrapf(ErrInvalidResponse, "unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var commit struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(limitedBody).Decode(&commit); err != nil {
		return "", errors.Wrap(err, "decode response")
	}

	if commit.SHA == "" {
		return "", errors.Wrap(ErrInvalidResponse, "empty SHA in response")
	}

	return commit.SHA, nil
}
