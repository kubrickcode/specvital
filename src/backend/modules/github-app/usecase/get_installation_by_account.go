package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/github-app/domain/entity"
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

type GetInstallationByAccountInput struct {
	AccountID int64
}

type GetInstallationByAccountUseCase struct {
	repository port.InstallationRepository
}

func NewGetInstallationByAccountUseCase(repository port.InstallationRepository) *GetInstallationByAccountUseCase {
	return &GetInstallationByAccountUseCase{repository: repository}
}

func (uc *GetInstallationByAccountUseCase) Execute(ctx context.Context, input GetInstallationByAccountInput) (*entity.Installation, error) {
	return uc.repository.GetByAccountID(ctx, input.AccountID)
}
