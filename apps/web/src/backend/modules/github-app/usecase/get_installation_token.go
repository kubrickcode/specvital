package usecase

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/github-app/domain"
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

type GetInstallationTokenInput struct {
	InstallationID int64
}

type GetInstallationTokenOutput struct {
	ExpiresAt time.Time
	Token     string
}

type GetInstallationTokenUseCase struct {
	ghAppClient port.GitHubAppClient
	repository  port.InstallationRepository
}

func NewGetInstallationTokenUseCase(
	ghAppClient port.GitHubAppClient,
	repository port.InstallationRepository,
) *GetInstallationTokenUseCase {
	return &GetInstallationTokenUseCase{
		ghAppClient: ghAppClient,
		repository:  repository,
	}
}

func (uc *GetInstallationTokenUseCase) Execute(ctx context.Context, input GetInstallationTokenInput) (*GetInstallationTokenOutput, error) {
	installation, err := uc.repository.GetByInstallationID(ctx, input.InstallationID)
	if err != nil {
		return nil, err
	}

	if installation.IsSuspended() {
		return nil, domain.ErrInstallationSuspended
	}

	token, err := uc.ghAppClient.CreateInstallationToken(ctx, input.InstallationID)
	if err != nil {
		return nil, err
	}

	return &GetInstallationTokenOutput{
		ExpiresAt: token.ExpiresAt,
		Token:     token.Token,
	}, nil
}
