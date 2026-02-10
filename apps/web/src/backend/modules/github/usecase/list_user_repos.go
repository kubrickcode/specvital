package usecase

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/domain/port"
)

type ListUserReposInput struct {
	Refresh bool
	UserID  string
}

type ListUserReposUseCase struct {
	clientFactory port.GitHubClientFactory
	repository    port.Repository
	sfGroup       singleflight.Group
	tokenProvider port.TokenProvider
}

func NewListUserReposUseCase(
	clientFactory port.GitHubClientFactory,
	repository port.Repository,
	tokenProvider port.TokenProvider,
) *ListUserReposUseCase {
	return &ListUserReposUseCase{
		clientFactory: clientFactory,
		repository:    repository,
		tokenProvider: tokenProvider,
	}
}

func (uc *ListUserReposUseCase) Execute(ctx context.Context, input ListUserReposInput) ([]entity.Repository, error) {
	key := fmt.Sprintf("user-repos:%s:refresh=%t", input.UserID, input.Refresh)

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

func (uc *ListUserReposUseCase) executeWithCache(ctx context.Context, input ListUserReposInput) ([]entity.Repository, error) {
	if input.Refresh {
		if err := uc.repository.DeleteUserRepositories(ctx, input.UserID); err != nil {
			return nil, fmt.Errorf("delete cached repositories: %w", err)
		}
	}

	hasData, err := uc.repository.HasUserRepositories(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("check cached repositories: %w", err)
	}

	if hasData {
		return uc.getFromCache(ctx, input.UserID)
	}

	return uc.fetchAndCache(ctx, input.UserID)
}

func (uc *ListUserReposUseCase) getFromCache(ctx context.Context, userID string) ([]entity.Repository, error) {
	records, err := uc.repository.GetUserRepositories(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cached repositories: %w", err)
	}
	return mapRepositoryRecordsToEntities(records), nil
}

func (uc *ListUserReposUseCase) fetchAndCache(ctx context.Context, userID string) ([]entity.Repository, error) {
	ghClient, err := getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	if err != nil {
		return nil, err
	}

	ghRepos, err := ghClient.ListUserRepositories(ctx, maxReposPerFetch)
	if err != nil {
		return nil, mapClientError(err)
	}

	repos := mapGitHubRepositoriesToEntities(ghRepos)

	records := mapEntitiesToRepositoryRecords(repos)
	if err := uc.repository.UpsertUserRepositories(ctx, userID, records); err != nil {
		return nil, fmt.Errorf("save repositories: %w", err)
	}

	return repos, nil
}
