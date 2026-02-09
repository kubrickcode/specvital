package client

import (
	"context"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	gh "github.com/google/go-github/v75/github"
	"golang.org/x/oauth2"
)

var (
	ErrGitHubUnauthorized      = errors.New("github token expired or invalid")
	ErrGitHubInsufficientScope = errors.New("github token lacks required permissions")
	ErrGitHubNotFound          = errors.New("github resource not found")
)

type RateLimitError struct {
	Limit     int
	Remaining int
	ResetAt   time.Time
}

func (e *RateLimitError) Error() string {
	return "github api rate limited"
}

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

type GitHubClient interface {
	GetOrganization(ctx context.Context, org string) (*GitHubOrganization, error)
	ListOrgRepositories(ctx context.Context, org string, maxResults int) ([]GitHubRepository, error)
	ListUserOrganizations(ctx context.Context) ([]GitHubOrganization, error)
	ListUserRepositories(ctx context.Context, maxResults int) ([]GitHubRepository, error)
}

type GitHubClientFactory func(token string) GitHubClient

func NewGitHubClientFactory() GitHubClientFactory {
	return func(token string) GitHubClient {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(context.Background(), ts)
		return &gitHubClient{client: gh.NewClient(tc)}
	}
}

type gitHubClient struct {
	client *gh.Client
}

const githubPageSize = 100

func (c *gitHubClient) ListUserRepositories(ctx context.Context, maxResults int) ([]GitHubRepository, error) {
	var allRepos []GitHubRepository
	opts := &gh.RepositoryListByAuthenticatedUserOptions{
		Sort:        "pushed",
		Direction:   "desc",
		Visibility:  "all",
		ListOptions: gh.ListOptions{PerPage: githubPageSize},
	}

	for len(allRepos) < maxResults {
		repos, resp, err := c.client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			return nil, handleGitHubError(err)
		}

		for _, repo := range repos {
			allRepos = append(allRepos, mapRepository(repo))
		}

		if resp.NextPage == 0 || len(allRepos) >= maxResults {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func (c *gitHubClient) ListUserOrganizations(ctx context.Context) ([]GitHubOrganization, error) {
	var allOrgs []GitHubOrganization
	opts := &gh.ListOptions{PerPage: githubPageSize}

	for {
		orgs, resp, err := c.client.Organizations.List(ctx, "", opts)
		if err != nil {
			return nil, handleGitHubError(err)
		}

		for _, org := range orgs {
			allOrgs = append(allOrgs, mapOrganization(org))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allOrgs, nil
}

func (c *gitHubClient) ListOrgRepositories(ctx context.Context, org string, maxResults int) ([]GitHubRepository, error) {
	var allRepos []GitHubRepository
	opts := &gh.RepositoryListByOrgOptions{
		Sort:        "pushed",
		Direction:   "desc",
		ListOptions: gh.ListOptions{PerPage: githubPageSize},
	}

	for len(allRepos) < maxResults {
		repos, resp, err := c.client.Repositories.ListByOrg(ctx, org, opts)
		if err != nil {
			return nil, handleGitHubError(err)
		}

		for _, repo := range repos {
			allRepos = append(allRepos, mapRepository(repo))
		}

		if resp.NextPage == 0 || len(allRepos) >= maxResults {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

func (c *gitHubClient) GetOrganization(ctx context.Context, org string) (*GitHubOrganization, error) {
	ghOrg, _, err := c.client.Organizations.Get(ctx, org)
	if err != nil {
		return nil, handleGitHubError(err)
	}

	result := mapOrganization(ghOrg)
	return &result, nil
}

func handleGitHubError(err error) error {
	var rateLimitErr *gh.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return &RateLimitError{
			Limit:     rateLimitErr.Rate.Limit,
			Remaining: rateLimitErr.Rate.Remaining,
			ResetAt:   rateLimitErr.Rate.Reset.Time,
		}
	}

	var abuseErr *gh.AbuseRateLimitError
	if errors.As(err, &abuseErr) {
		return &RateLimitError{}
	}

	var ghErr *gh.ErrorResponse
	if errors.As(err, &ghErr) {
		switch ghErr.Response.StatusCode {
		case http.StatusUnauthorized:
			return ErrGitHubUnauthorized
		case http.StatusForbidden:
			return ErrGitHubInsufficientScope
		case http.StatusNotFound:
			return ErrGitHubNotFound
		}
	}

	return errors.Wrap(err, "github api")
}

func mapRepository(repo *gh.Repository) GitHubRepository {
	var owner string
	if repo.Owner != nil {
		owner = repo.Owner.GetLogin()
	}

	r := GitHubRepository{
		Archived:  repo.GetArchived(),
		Disabled:  repo.GetDisabled(),
		Fork:      repo.GetFork(),
		FullName:  repo.GetFullName(),
		HTMLURL:   repo.GetHTMLURL(),
		ID:        repo.GetID(),
		Name:      repo.GetName(),
		Owner:     owner,
		Private:   repo.GetPrivate(),
		StarCount: repo.GetStargazersCount(),
	}

	if repo.Description != nil {
		r.Description = *repo.Description
	}
	if repo.DefaultBranch != nil {
		r.DefaultBranch = *repo.DefaultBranch
	}
	if repo.Language != nil {
		r.Language = *repo.Language
	}
	if repo.Visibility != nil {
		r.Visibility = *repo.Visibility
	}
	if repo.PushedAt != nil {
		t := repo.PushedAt.Time
		r.PushedAt = &t
	}

	return r
}

func mapOrganization(org *gh.Organization) GitHubOrganization {
	o := GitHubOrganization{
		ID:    org.GetID(),
		Login: org.GetLogin(),
	}

	if org.AvatarURL != nil {
		o.AvatarURL = *org.AvatarURL
	}
	if org.HTMLURL != nil {
		o.HTMLURL = *org.HTMLURL
	}
	if org.Description != nil {
		o.Description = *org.Description
	}

	return o
}
