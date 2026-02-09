package port

import (
	"context"
	"time"
)

type GitHubClient interface {
	GetOrganization(ctx context.Context, org string) (*GitHubOrganization, error)
	ListOrgRepositories(ctx context.Context, org string, maxResults int) ([]GitHubRepository, error)
	ListUserOrganizations(ctx context.Context) ([]GitHubOrganization, error)
	ListUserRepositories(ctx context.Context, maxResults int) ([]GitHubRepository, error)
}

type GitHubClientFactory func(token string) GitHubClient

type GitHubRepository struct {
	Archived      bool
	DefaultBranch string
	Description   string
	Disabled      bool
	Fork          bool
	FullName      string
	HTMLURL       string
	ID            int64
	Language      string
	Name          string
	Owner         string
	Private       bool
	PushedAt      *time.Time
	StarCount     int
	Visibility    string
}

type GitHubOrganization struct {
	AvatarURL   string
	Description string
	HTMLURL     string
	ID          int64
	Login       string
}
