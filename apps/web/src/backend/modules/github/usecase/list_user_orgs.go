package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/singleflight"

	"github.com/specvital/web/src/backend/modules/github/domain/entity"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

type ListUserOrgsInput struct {
	Refresh bool
	UserID  string
}

type ListUserOrgsUseCase struct {
	clientFactory      port.GitHubClientFactory
	installationLookup port.InstallationLookup
	repository         port.Repository
	sfGroup            singleflight.Group
	tokenProvider      port.TokenProvider
}

func NewListUserOrgsUseCase(
	clientFactory port.GitHubClientFactory,
	repository port.Repository,
	tokenProvider port.TokenProvider,
	installationLookup port.InstallationLookup,
) *ListUserOrgsUseCase {
	return &ListUserOrgsUseCase{
		clientFactory:      clientFactory,
		installationLookup: installationLookup,
		repository:         repository,
		tokenProvider:      tokenProvider,
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

	orgs = uc.mergeWithInstallationOrgs(ctx, input.UserID, orgs)
	return uc.enrichWithAccessStatus(ctx, orgs), nil
}

func (uc *ListUserOrgsUseCase) mergeWithInstallationOrgs(ctx context.Context, userID string, oauthOrgs []entity.Organization) []entity.Organization {
	installationOrgs, err := uc.installationLookup.ListOrganizationsByUserID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "failed to list installation orgs, using OAuth orgs only",
			"error", err,
			"userID", userID)
		return oauthOrgs
	}

	if len(installationOrgs) == 0 {
		return oauthOrgs
	}

	existingIDs := make(map[int64]bool, len(oauthOrgs))
	for _, org := range oauthOrgs {
		existingIDs[org.ID] = true
	}

	for _, instOrg := range installationOrgs {
		if existingIDs[instOrg.AccountID] {
			continue
		}

		avatarURL := ""
		if instOrg.AccountAvatarURL != nil {
			avatarURL = *instOrg.AccountAvatarURL
		}

		oauthOrgs = append(oauthOrgs, entity.Organization{
			AvatarURL: avatarURL,
			ID:        instOrg.AccountID,
			Login:     instOrg.AccountLogin,
		})
	}

	return oauthOrgs
}

func (uc *ListUserOrgsUseCase) enrichWithAccessStatus(ctx context.Context, orgs []entity.Organization) []entity.Organization {
	if len(orgs) == 0 {
		return orgs
	}

	accountIDs := make([]int64, len(orgs))
	for i, org := range orgs {
		accountIDs[i] = org.ID
	}

	installations, err := uc.installationLookup.GetInstallationsByAccountIDs(ctx, accountIDs)
	if err != nil {
		slog.WarnContext(ctx, "failed to batch lookup installations, marking all as restricted",
			"error", err,
			"orgCount", len(orgs))
		for i := range orgs {
			orgs[i].AccessStatus = entity.AccessStatusRestricted
		}
		return orgs
	}

	for i := range orgs {
		installation, exists := installations[orgs[i].ID]
		if !exists || installation == nil {
			orgs[i].AccessStatus = entity.AccessStatusRestricted
			continue
		}

		if installation.IsSuspended {
			orgs[i].AccessStatus = entity.AccessStatusPending
		} else {
			orgs[i].AccessStatus = entity.AccessStatusAccessible
		}
	}

	return orgs
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
