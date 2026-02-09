package usecase

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	authdomain "github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/github/domain"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

type mockTokenProvider struct {
	token string
	err   error
}

func (m *mockTokenProvider) GetUserGitHubToken(_ context.Context, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.token, nil
}

type mockGitHubClient struct {
	repos    []port.GitHubRepository
	orgs     []port.GitHubOrganization
	org      *port.GitHubOrganization
	reposErr error
	orgsErr  error
	orgErr   error
}

func (m *mockGitHubClient) GetOrganization(_ context.Context, _ string) (*port.GitHubOrganization, error) {
	if m.orgErr != nil {
		return nil, m.orgErr
	}
	return m.org, nil
}

func (m *mockGitHubClient) ListOrgRepositories(_ context.Context, _ string, _ int) ([]port.GitHubRepository, error) {
	if m.reposErr != nil {
		return nil, m.reposErr
	}
	return m.repos, nil
}

func (m *mockGitHubClient) ListUserOrganizations(_ context.Context) ([]port.GitHubOrganization, error) {
	if m.orgsErr != nil {
		return nil, m.orgsErr
	}
	return m.orgs, nil
}

func (m *mockGitHubClient) ListUserRepositories(_ context.Context, _ int) ([]port.GitHubRepository, error) {
	if m.reposErr != nil {
		return nil, m.reposErr
	}
	return m.repos, nil
}

func mockClientFactory(c port.GitHubClient) port.GitHubClientFactory {
	return func(_ string) port.GitHubClient {
		return c
	}
}

type mockRepository struct {
	repos       []port.RepositoryRecord
	orgs        []port.OrganizationRecord
	orgRepos    []port.RepositoryRecord
	hasRepos    bool
	hasOrgs     bool
	hasOrgRepos bool
	orgID       string
	orgIDErr    error
	err         error
}

func (m *mockRepository) DeleteOrgRepositories(_ context.Context, _, _ string) error {
	return m.err
}

func (m *mockRepository) DeleteUserOrganizations(_ context.Context, _ string) error {
	return m.err
}

func (m *mockRepository) DeleteUserRepositories(_ context.Context, _ string) error {
	return m.err
}

func (m *mockRepository) GetOrgIDByLogin(_ context.Context, _ string) (string, error) {
	if m.orgIDErr != nil {
		return "", m.orgIDErr
	}
	return m.orgID, nil
}

func (m *mockRepository) GetOrgRepositories(_ context.Context, _, _ string) ([]port.RepositoryRecord, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.orgRepos, nil
}

func (m *mockRepository) GetUserOrganizations(_ context.Context, _ string) ([]port.OrganizationRecord, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.orgs, nil
}

func (m *mockRepository) GetUserRepositories(_ context.Context, _ string) ([]port.RepositoryRecord, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.repos, nil
}

func (m *mockRepository) HasOrgRepositories(_ context.Context, _, _ string) (bool, error) {
	return m.hasOrgRepos, m.err
}

func (m *mockRepository) HasUserOrganizations(_ context.Context, _ string) (bool, error) {
	return m.hasOrgs, m.err
}

func (m *mockRepository) HasUserRepositories(_ context.Context, _ string) (bool, error) {
	return m.hasRepos, m.err
}

func (m *mockRepository) UpsertOrgRepositories(_ context.Context, _, _ string, _ []port.RepositoryRecord) error {
	return m.err
}

func (m *mockRepository) UpsertUserOrganizations(_ context.Context, _ string, _ []port.OrganizationRecord) error {
	return m.err
}

func (m *mockRepository) UpsertUserRepositories(_ context.Context, _ string, _ []port.RepositoryRecord) error {
	return m.err
}

func TestListUserReposUseCase_FromCache(t *testing.T) {
	repo := &mockRepository{
		hasRepos: true,
		repos:    []port.RepositoryRecord{{ID: 1, Name: "repo1"}},
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{}
	uc := NewListUserReposUseCase(mockClientFactory(ghClient), repo, provider)

	repos, err := uc.Execute(context.Background(), ListUserReposInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
}

func TestListUserReposUseCase_FromGitHub(t *testing.T) {
	repo := &mockRepository{
		hasRepos: false,
	}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{
		repos: []port.GitHubRepository{
			{ID: 1, Name: "repo1", FullName: "user/repo1"},
			{ID: 2, Name: "repo2", FullName: "user/repo2"},
		},
	}
	uc := NewListUserReposUseCase(mockClientFactory(ghClient), repo, provider)

	repos, err := uc.Execute(context.Background(), ListUserReposInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(repos))
	}
}

func TestListUserReposUseCase_NoToken(t *testing.T) {
	repo := &mockRepository{}
	provider := &mockTokenProvider{err: authdomain.ErrNoGitHubToken}
	ghClient := &mockGitHubClient{}
	uc := NewListUserReposUseCase(mockClientFactory(ghClient), repo, provider)

	_, err := uc.Execute(context.Background(), ListUserReposInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if !errors.Is(err, domain.ErrNoGitHubToken) {
		t.Errorf("expected ErrNoGitHubToken, got %v", err)
	}
}

func TestListUserReposUseCase_RateLimited(t *testing.T) {
	repo := &mockRepository{hasRepos: false}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{
		reposErr: &port.RateLimitError{Limit: 5000, Remaining: 0},
	}
	uc := NewListUserReposUseCase(mockClientFactory(ghClient), repo, provider)

	_, err := uc.Execute(context.Background(), ListUserReposInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if !domain.IsRateLimitError(err) {
		t.Errorf("expected RateLimitError, got %v", err)
	}
}

func TestListUserReposUseCase_Unauthorized(t *testing.T) {
	repo := &mockRepository{hasRepos: false}
	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{
		reposErr: port.ErrGitHubUnauthorized,
	}
	uc := NewListUserReposUseCase(mockClientFactory(ghClient), repo, provider)

	_, err := uc.Execute(context.Background(), ListUserReposInput{
		UserID:  "user-123",
		Refresh: false,
	})
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestListUserReposUseCase_Refresh(t *testing.T) {
	deleteCount := 0
	repo := &mockRepository{
		hasRepos: false,
	}
	originalDelete := repo.DeleteUserRepositories
	_ = originalDelete

	provider := &mockTokenProvider{token: "test-token"}
	ghClient := &mockGitHubClient{
		repos: []port.GitHubRepository{
			{ID: 1, Name: "repo1"},
		},
	}

	customRepo := &mockRepositoryWithDeleteCounter{
		mockRepository: repo,
		deleteCount:    &deleteCount,
	}
	uc := NewListUserReposUseCase(mockClientFactory(ghClient), customRepo, provider)

	repos, err := uc.Execute(context.Background(), ListUserReposInput{
		UserID:  "user-123",
		Refresh: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
	if deleteCount != 1 {
		t.Errorf("expected DeleteUserRepositories to be called once, got %d", deleteCount)
	}
}

type mockRepositoryWithDeleteCounter struct {
	*mockRepository
	deleteCount *int
}

func (m *mockRepositoryWithDeleteCounter) DeleteUserRepositories(_ context.Context, _ string) error {
	*m.deleteCount++
	return m.mockRepository.err
}
