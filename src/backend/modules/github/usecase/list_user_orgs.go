package usecase

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"

	"github.com/specvital/web/src/backend/modules/github/domain/entity"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

type ListUserOrgsInput struct {
	Refresh bool
	UserID  string
}

type ListUserOrgsUseCase struct {
	clientFactory port.GitHubClientFactory
	repository    port.Repository
	sfGroup       singleflight.Group
	tokenProvider port.TokenProvider
}

func NewListUserOrgsUseCase(
	clientFactory port.GitHubClientFactory,
	repository port.Repository,
	tokenProvider port.TokenProvider,
) *ListUserOrgsUseCase {
	return &ListUserOrgsUseCase{
		clientFactory: clientFactory,
		repository:    repository,
		tokenProvider: tokenProvider,
	}
}

func (uc *ListUserOrgsUseCase) Execute(ctx context.Context, input ListUserOrgsInput) ([]entity.Organization, error) {
	key := fmt.Sprintf("user-orgs:%s:refresh=%t", input.UserID, input.Refresh)

	result, err, _ := uc.sfGroup.Do(key, func() (any, error) {
		return uc.executeWithCache(ctx, input)
	})

	if err != nil {
		return nil, err
	}

	orgs, ok := result.([]entity.Organization)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
	return orgs, nil
}

func (uc *ListUserOrgsUseCase) executeWithCache(ctx context.Context, input ListUserOrgsInput) ([]entity.Organization, error) {
	if input.Refresh {
		if err := uc.repository.DeleteUserOrganizations(ctx, input.UserID); err != nil {
			return nil, fmt.Errorf("delete cached organizations: %w", err)
		}
	}

	hasData, err := uc.repository.HasUserOrganizations(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("check cached organizations: %w", err)
	}

	if hasData {
		return uc.getFromCache(ctx, input.UserID)
	}

	return uc.fetchAndCache(ctx, input.UserID)
}

func (uc *ListUserOrgsUseCase) getFromCache(ctx context.Context, userID string) ([]entity.Organization, error) {
	records, err := uc.repository.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cached organizations: %w", err)
	}
	return mapOrganizationRecordsToEntities(records), nil
}

func (uc *ListUserOrgsUseCase) fetchAndCache(ctx context.Context, userID string) ([]entity.Organization, error) {
	ghClient, err := getGitHubClient(ctx, uc.clientFactory, uc.tokenProvider, userID)
	if err != nil {
		return nil, err
	}

	ghOrgs, err := ghClient.ListUserOrganizations(ctx)
	if err != nil {
		return nil, mapClientError(err)
	}

	orgs := mapGitHubOrganizationsToEntities(ghOrgs)

	records := mapEntitiesToOrganizationRecords(orgs)
	if err := uc.repository.UpsertUserOrganizations(ctx, userID, records); err != nil {
		return nil, fmt.Errorf("save organizations: %w", err)
	}

	return orgs, nil
}
