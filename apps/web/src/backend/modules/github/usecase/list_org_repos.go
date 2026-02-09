package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/singleflight"

	"github.com/specvital/web/src/backend/modules/github/domain/entity"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

type ListOrgReposInput struct {
	OrgLogin string
	Refresh  bool
	UserID   string
}

type ListOrgReposUseCase struct {
	clientFactory             port.GitHubClientFactory
	installationLookup        port.InstallationLookup
	installationTokenProvider port.InstallationTokenProvider
	repository                port.Repository
	sfGroup                   singleflight.Group
	tokenProvider             port.TokenProvider
}

func NewListOrgReposUseCase(
	clientFactory port.GitHubClientFactory,
	repository port.Repository,
	tokenProvider port.TokenProvider,
	installationLookup port.InstallationLookup,
	installationTokenProvider port.InstallationTokenProvider,
) *ListOrgReposUseCase {
	return &ListOrgReposUseCase{
		clientFactory:             clientFactory,
		installationLookup:        installationLookup,
		installationTokenProvider: installationTokenProvider,
		repository:                repository,
		tokenProvider:             tokenProvider,
	}
}

func (uc *ListOrgReposUseCase) Execute(ctx context.Context, input ListOrgReposInput) ([]entity.Repository, error) {
	key := fmt.Sprintf("org-repos:%s:%s:refresh=%t", input.UserID, input.OrgLogin, input.Refresh)

	result, err, _ := uc.sfGroup.Do(key, func() (any, error) {
		return uc.executeWithCache(ctx, input)
	})

	if err != nil {
		return nil, err
	}

	repos, ok := result.([]entity.Repository)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
	return repos, nil
}

func (uc *ListOrgReposUseCase) executeWithCache(ctx context.Context, input ListOrgReposInput) ([]entity.Repository, error) {
	orgID, err := uc.repository.GetOrgIDByLogin(ctx, input.OrgLogin)
	if err != nil {
		return uc.fetchAndCacheOrgRepos(ctx, input.UserID, input.OrgLogin)
	}

	if input.Refresh {
		if err := uc.repository.DeleteOrgRepositories(ctx, input.UserID, orgID); err != nil {
			return nil, fmt.Errorf("delete cached org repositories: %w", err)
		}
	}

	hasData, err := uc.repository.HasOrgRepositories(ctx, input.UserID, orgID)
	if err != nil {
		return nil, fmt.Errorf("check cached org repositories: %w", err)
	}

	if hasData {
		return uc.getFromCache(ctx, input.UserID, orgID)
	}

	return uc.fetchAndCacheOrgRepos(ctx, input.UserID, input.OrgLogin)
}

func (uc *ListOrgReposUseCase) getFromCache(ctx context.Context, userID, orgID string) ([]entity.Repository, error) {
	records, err := uc.repository.GetOrgRepositories(ctx, userID, orgID)
	if err != nil {
		return nil, fmt.Errorf("get cached org repositories: %w", err)
	}
	return mapRepositoryRecordsToEntities(records), nil
}

func (uc *ListOrgReposUseCase) fetchAndCacheOrgRepos(ctx context.Context, userID, orgLogin string) ([]entity.Repository, error) {
	ghClient, err := getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	if err != nil {
		return nil, err
	}

	ghOrg, err := ghClient.GetOrganization(ctx, orgLogin)
	if err != nil {
		return nil, mapClientError(err)
	}

	org := mapGitHubOrganizationToEntity(*ghOrg)
	orgRecord := mapEntityToOrganizationRecord(org)
	if err := uc.repository.UpsertUserOrganizations(ctx, userID, []port.OrganizationRecord{orgRecord}); err != nil {
		return nil, fmt.Errorf("save organization: %w", err)
	}

	orgID, err := uc.repository.GetOrgIDByLogin(ctx, orgLogin)
	if err != nil {
		return nil, fmt.Errorf("get org id: %w", err)
	}

	repoClient, err := uc.getOrgRepoClient(ctx, userID, ghOrg.ID)
	if err != nil {
		return nil, err
	}

	ghRepos, err := repoClient.ListOrgRepositories(ctx, orgLogin, maxReposPerFetch)
	if err != nil {
		return nil, mapClientError(err)
	}

	repos := mapGitHubRepositoriesToEntities(ghRepos)

	repoRecords := mapEntitiesToRepositoryRecords(repos)
	if err := uc.repository.UpsertOrgRepositories(ctx, userID, orgID, repoRecords); err != nil {
		return nil, fmt.Errorf("save org repositories: %w", err)
	}

	return repos, nil
}

func (uc *ListOrgReposUseCase) getOrgRepoClient(ctx context.Context, userID string, orgAccountID int64) (port.GitHubClient, error) {
	installation, err := uc.installationLookup.GetInstallationByAccountID(ctx, orgAccountID)
	if err != nil {
		slog.WarnContext(ctx, "failed to lookup installation, falling back to OAuth token",
			"orgAccountID", orgAccountID,
			"error", err)
		return getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	}

	if installation == nil {
		return getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	}

	if installation.IsSuspended {
		slog.DebugContext(ctx, "installation is suspended, falling back to OAuth token",
			"orgAccountID", orgAccountID,
			"installationID", installation.InstallationID)
		return getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	}

	token, err := uc.installationTokenProvider.GetInstallationToken(ctx, installation.InstallationID)
	if err != nil {
		slog.WarnContext(ctx, "failed to get installation token, falling back to OAuth token",
			"orgAccountID", orgAccountID,
			"installationID", installation.InstallationID,
			"error", err)
		return getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	}

	return uc.clientFactory(token), nil
}
