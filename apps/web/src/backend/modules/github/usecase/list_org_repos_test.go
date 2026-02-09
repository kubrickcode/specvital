package usecase

import (
	"context"
	"testing"

	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

type mockInstallationTokenProvider struct {
	token string
	err   error
}

func (m *mockInstallationTokenProvider) GetInstallationToken(_ context.Context, _ int64) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.token, nil
}

func TestListOrgReposUseCase_UsesInstallationToken(t *testing.T) {
	orgAccountID := int64(123)

	repo := &mockRepository{
		hasOrgRepos: false,
		orgID:       "org-internal-id",
	}
	provider := &mockTokenProvider{token: "oauth-token"}

	oauthClient := &mockGitHubClient{
		org:   &port.GitHubOrganization{ID: orgAccountID, Login: "test-org"},
		repos: []port.GitHubRepository{},
	}
	installationClient := &mockGitHubClient{
		repos: []port.GitHubRepository{
			{ID: 1, Name: "private-repo", FullName: "test-org/private-repo"},
		},
	}

	calledWithToken := ""
	factory := func(token string) port.GitHubClient {
		calledWithToken = token
		if token == "installation-token" {
			return installationClient
		}
		return oauthClient
	}

	lookup := &mockInstallationLookup{
		installations: map[int64]*port.InstallationInfo{
			orgAccountID: {AccountID: orgAccountID, InstallationID: 456, IsSuspended: false},
		},
	}
	tokenProvider := &mockInstallationTokenProvider{token: "installation-token"}

	uc := NewListOrgReposUseCase(factory, repo, provider, lookup, tokenProvider)

	repos, err := uc.Execute(context.Background(), ListOrgReposInput{
		UserID:   "user-123",
		OrgLogin: "test-org",
		Refresh:  false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calledWithToken != "installation-token" {
		t.Errorf("expected installation-token, got %s", calledWithToken)
	}

	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
}

func TestListOrgReposUseCase_FallsBackToOAuthToken(t *testing.T) {
	orgAccountID := int64(123)

	repo := &mockRepository{
		hasOrgRepos: false,
		orgID:       "org-internal-id",
	}
	provider := &mockTokenProvider{token: "oauth-token"}

	oauthClient := &mockGitHubClient{
		org: &port.GitHubOrganization{ID: orgAccountID, Login: "test-org"},
		repos: []port.GitHubRepository{
			{ID: 1, Name: "public-repo", FullName: "test-org/public-repo"},
		},
	}

	lastToken := ""
	factory := func(token string) port.GitHubClient {
		lastToken = token
		return oauthClient
	}

	lookup := &mockInstallationLookup{}
	tokenProvider := &mockInstallationTokenProvider{token: "installation-token"}

	uc := NewListOrgReposUseCase(factory, repo, provider, lookup, tokenProvider)

	repos, err := uc.Execute(context.Background(), ListOrgReposInput{
		UserID:   "user-123",
		OrgLogin: "test-org",
		Refresh:  false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lastToken != "oauth-token" {
		t.Errorf("expected oauth-token for fallback, got %s", lastToken)
	}

	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
}

func TestListOrgReposUseCase_FallsBackWhenSuspended(t *testing.T) {
	orgAccountID := int64(123)

	repo := &mockRepository{
		hasOrgRepos: false,
		orgID:       "org-internal-id",
	}
	provider := &mockTokenProvider{token: "oauth-token"}

	oauthClient := &mockGitHubClient{
		org:   &port.GitHubOrganization{ID: orgAccountID, Login: "test-org"},
		repos: []port.GitHubRepository{},
	}

	lastToken := ""
	factory := func(token string) port.GitHubClient {
		lastToken = token
		return oauthClient
	}

	lookup := &mockInstallationLookup{
		installations: map[int64]*port.InstallationInfo{
			orgAccountID: {AccountID: orgAccountID, InstallationID: 456, IsSuspended: true},
		},
	}
	tokenProvider := &mockInstallationTokenProvider{token: "installation-token"}

	uc := NewListOrgReposUseCase(factory, repo, provider, lookup, tokenProvider)

	_, err := uc.Execute(context.Background(), ListOrgReposInput{
		UserID:   "user-123",
		OrgLogin: "test-org",
		Refresh:  false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lastToken != "oauth-token" {
		t.Errorf("expected oauth-token when suspended, got %s", lastToken)
	}
}
