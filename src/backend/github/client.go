package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	defaultBaseURL   = "https://api.github.com"
	defaultTimeout   = 30 * time.Second
	userAgent        = "specvital-web/1.0"
	rateLimitHeader  = "X-RateLimit-Remaining"
	rateLimitLimit   = "X-RateLimit-Limit"
	rateLimitReset   = "X-RateLimit-Reset"
	acceptHeader     = "application/vnd.github+json"
	apiVersionHeader = "X-GitHub-Api-Version"
	apiVersion       = "2022-11-28"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string

	mu        sync.RWMutex
	rateLimit RateLimitInfo
}

func NewClient(token string) *Client {
	return &Client{
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		token:      token,
	}
}

func (c *Client) WithBaseURL(url string) *Client {
	c.baseURL = url
	return c
}

func (c *Client) GetRateLimit() RateLimitInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rateLimit
}

func validateOwnerRepo(owner, repo string) error {
	if owner == "" || repo == "" {
		return ErrInvalidInput
	}
	if !validNamePattern.MatchString(owner) || !validNamePattern.MatchString(repo) {
		return ErrInvalidInput
	}
	return nil
}

func buildURL(baseURL string, pathSegments ...string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	escaped := make([]string, len(pathSegments))
	for i, seg := range pathSegments {
		escaped[i] = url.PathEscape(seg)
	}
	base.Path = base.Path + "/" + joinPath(escaped...)
	return base.String(), nil
}

func joinPath(segments ...string) string {
	result := ""
	for i, s := range segments {
		if i > 0 {
			result += "/"
		}
		result += s
	}
	return result
}

func (c *Client) ListFiles(ctx context.Context, owner, repo string) ([]FileInfo, error) {
	if err := validateOwnerRepo(owner, repo); err != nil {
		return nil, err
	}

	reqURL, err := buildURL(c.baseURL, "repos", owner, repo, "git", "trees", "HEAD")
	if err != nil {
		return nil, err
	}
	reqURL += "?recursive=1"

	body, err := c.doRequestAndRead(ctx, http.MethodGet, reqURL)
	if err != nil {
		return nil, err
	}

	var tree treeResponse
	if err := json.Unmarshal(body, &tree); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tree.Truncated {
		return nil, ErrTreeTruncated
	}

	files := make([]FileInfo, 0, len(tree.Tree))
	for _, entry := range tree.Tree {
		files = append(files, FileInfo{
			Path: entry.Path,
			SHA:  entry.SHA,
			Size: entry.Size,
			Type: toFileType(entry.Type),
		})
	}

	return files, nil
}

func toFileType(githubType string) string {
	if githubType == "blob" {
		return "file"
	}
	return "dir"
}

func (c *Client) GetFileContent(ctx context.Context, owner, repo, path string) (string, error) {
	if err := validateOwnerRepo(owner, repo); err != nil {
		return "", err
	}
	if path == "" {
		return "", ErrInvalidInput
	}

	reqURL, err := buildURL(c.baseURL, "repos", owner, repo, "contents", path)
	if err != nil {
		return "", err
	}

	body, err := c.doRequestAndRead(ctx, http.MethodGet, reqURL)
	if err != nil {
		return "", err
	}

	var content contentResponse
	if err := json.Unmarshal(body, &content); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if content.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding: %s", content.Encoding)
	}

	decoded, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	return string(decoded), nil
}

func (c *Client) GetLatestCommit(ctx context.Context, owner, repo string) (*CommitInfo, error) {
	if err := validateOwnerRepo(owner, repo); err != nil {
		return nil, err
	}

	repoURL, err := buildURL(c.baseURL, "repos", owner, repo)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequestAndRead(ctx, http.MethodGet, repoURL)
	if err != nil {
		return nil, err
	}

	var repoInfo struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.Unmarshal(body, &repoInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	branchURL, err := buildURL(c.baseURL, "repos", owner, repo, "branches", repoInfo.DefaultBranch)
	if err != nil {
		return nil, err
	}

	branchBody, err := c.doRequestAndRead(ctx, http.MethodGet, branchURL)
	if err != nil {
		return nil, err
	}

	var branch branchResponse
	if err := json.Unmarshal(branchBody, &branch); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &CommitInfo{
		SHA:     branch.Commit.SHA,
		Message: branch.Commit.Commit.Message,
	}, nil
}

func (c *Client) doRequestAndRead(ctx context.Context, method, reqURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", acceptHeader)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set(apiVersionHeader, apiVersion)

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	c.updateRateLimit(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if err := c.checkResponse(resp, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) checkResponse(resp *http.Response, body []byte) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden:
		rateLimit := c.GetRateLimit()
		if rateLimit.Remaining == 0 {
			return ErrRateLimited
		}
		return ErrForbidden
	case http.StatusTooManyRequests:
		return ErrRateLimited
	default:
		var errResp errorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *Client) updateRateLimit(resp *http.Response) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if remaining := resp.Header.Get(rateLimitHeader); remaining != "" {
		if val, err := strconv.Atoi(remaining); err == nil {
			c.rateLimit.Remaining = val
		}
	}

	if limit := resp.Header.Get(rateLimitLimit); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			c.rateLimit.Limit = val
		}
	}

	if reset := resp.Header.Get(rateLimitReset); reset != "" {
		if val, err := strconv.ParseInt(reset, 10, 64); err == nil {
			c.rateLimit.ResetAt = val
		}
	}
}
