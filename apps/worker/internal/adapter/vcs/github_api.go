package vcs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kubrickcode/specvital/apps/worker/internal/domain/analysis"
)

const (
	gitHubHost    = "github.com"
	gitHubAPIBase = "https://api.github.com"
)

type GitHubAPIClient struct {
	apiBase    string
	httpClient *http.Client
}

var _ analysis.VCSAPIClient = (*GitHubAPIClient)(nil)

func NewGitHubAPIClient(httpClient *http.Client) *GitHubAPIClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &GitHubAPIClient{
		apiBase:    gitHubAPIBase,
		httpClient: httpClient,
	}
}

func (c *GitHubAPIClient) GetRepoInfo(ctx context.Context, host, owner, repo string, token *string) (analysis.RepoInfo, error) {
	if host != gitHubHost {
		return analysis.RepoInfo{}, fmt.Errorf("%w: unsupported host %q (only %q is supported)", analysis.ErrInvalidInput, host, gitHubHost)
	}
	if owner == "" {
		return analysis.RepoInfo{}, fmt.Errorf("%w: owner is required", analysis.ErrInvalidInput)
	}
	if repo == "" {
		return analysis.RepoInfo{}, fmt.Errorf("%w: repo is required", analysis.ErrInvalidInput)
	}

	url := fmt.Sprintf("%s/repos/%s/%s", c.apiBase, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return analysis.RepoInfo{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token != nil && *token != "" {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return analysis.RepoInfo{}, fmt.Errorf("get repository %s/%s: %w", owner, repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return analysis.RepoInfo{}, fmt.Errorf("%w: %s/%s", analysis.ErrRepoNotFound, owner, repo)
	}
	if resp.StatusCode != http.StatusOK {
		return analysis.RepoInfo{}, fmt.Errorf("get repository %s/%s: unexpected status %d", owner, repo, resp.StatusCode)
	}

	var result struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return analysis.RepoInfo{}, fmt.Errorf("decode response: %w", err)
	}

	return analysis.RepoInfo{
		ExternalRepoID: strconv.FormatInt(result.ID, 10),
		Name:           result.Name,
		Owner:          result.Owner.Login,
	}, nil
}
