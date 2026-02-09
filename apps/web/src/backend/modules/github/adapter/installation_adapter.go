package adapter

import (
	"context"
	"errors"

	ghappdomain "github.com/specvital/web/src/backend/modules/github-app/domain"
	ghappport "github.com/specvital/web/src/backend/modules/github-app/domain/port"
	ghappusecase "github.com/specvital/web/src/backend/modules/github-app/usecase"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

var (
	_ port.InstallationLookup        = (*InstallationLookupAdapter)(nil)
	_ port.InstallationTokenProvider = (*InstallationTokenProviderAdapter)(nil)
)

type InstallationLookupAdapter struct {
	repository ghappport.InstallationRepository
}

func NewInstallationLookupAdapter(repository ghappport.InstallationRepository) *InstallationLookupAdapter {
	return &InstallationLookupAdapter{repository: repository}
}

func (a *InstallationLookupAdapter) GetInstallationByAccountID(ctx context.Context, accountID int64) (*port.InstallationInfo, error) {
	installation, err := a.repository.GetByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, ghappdomain.ErrInstallationNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &port.InstallationInfo{
		AccountID:      installation.AccountID,
		InstallationID: installation.InstallationID,
		IsSuspended:    installation.IsSuspended(),
	}, nil
}

func (a *InstallationLookupAdapter) GetInstallationsByAccountIDs(ctx context.Context, accountIDs []int64) (map[int64]*port.InstallationInfo, error) {
	if len(accountIDs) == 0 {
		return make(map[int64]*port.InstallationInfo), nil
	}

	installations, err := a.repository.ListByAccountIDs(ctx, accountIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]*port.InstallationInfo, len(installations))
	for _, inst := range installations {
		result[inst.AccountID] = &port.InstallationInfo{
			AccountID:      inst.AccountID,
			InstallationID: inst.InstallationID,
			IsSuspended:    inst.IsSuspended(),
		}
	}
	return result, nil
}

func (a *InstallationLookupAdapter) ListOrganizationsByUserID(ctx context.Context, _ string) ([]port.OrganizationInstallationInfo, error) {
	installations, err := a.repository.ListOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	orgs := make([]port.OrganizationInstallationInfo, 0, len(installations))
	for _, inst := range installations {
		orgs = append(orgs, port.OrganizationInstallationInfo{
			AccountAvatarURL: inst.AccountAvatarURL,
			AccountID:        inst.AccountID,
			AccountLogin:     inst.AccountLogin,
			IsSuspended:      inst.IsSuspended(),
		})
	}
	return orgs, nil
}

type InstallationTokenProviderAdapter struct {
	usecase *ghappusecase.GetInstallationTokenUseCase
}

func NewInstallationTokenProviderAdapter(usecase *ghappusecase.GetInstallationTokenUseCase) *InstallationTokenProviderAdapter {
	return &InstallationTokenProviderAdapter{usecase: usecase}
}

func (a *InstallationTokenProviderAdapter) GetInstallationToken(ctx context.Context, installationID int64) (string, error) {
	output, err := a.usecase.Execute(ctx, ghappusecase.GetInstallationTokenInput{
		InstallationID: installationID,
	})
	if err != nil {
		return "", err
	}
	return output.Token, nil
}
