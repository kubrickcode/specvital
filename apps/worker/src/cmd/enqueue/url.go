package main

import (
	"fmt"
	"strings"
)

// ParseGitHubURL parses various GitHub URL formats and extracts owner and repo.
// Supported formats:
//   - github.com/owner/repo
//   - https://github.com/owner/repo
//   - http://github.com/owner/repo
//   - github.com/owner/repo.git
//   - https://github.com/owner/repo.git
func ParseGitHubURL(url string) (owner, repo string, err error) {
	if url == "" {
		return "", "", fmt.Errorf("URL cannot be empty")
	}

	normalized := url
	normalized = strings.TrimPrefix(normalized, "https://")
	normalized = strings.TrimPrefix(normalized, "http://")
	normalized = strings.TrimPrefix(normalized, "github.com/")
	normalized = strings.TrimSuffix(normalized, ".git")
	normalized = strings.TrimSuffix(normalized, "/")

	parts := strings.Split(normalized, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format: expected owner/repo, got %q", url)
	}

	owner = parts[0]
	repo = parts[1]

	if owner == "" {
		return "", "", fmt.Errorf("owner cannot be empty in URL: %q", url)
	}
	if repo == "" {
		return "", "", fmt.Errorf("repo cannot be empty in URL: %q", url)
	}

	return owner, repo, nil
}
