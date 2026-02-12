package usecase

import (
	"context"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/domain/port"
)

type ListInstallationsInput struct {
	UserID string
}

type ListInstallationsUseCase struct {
	repository port.InstallationRepository
}

func NewListInstallationsUseCase(repository port.InstallationRepository) *ListInstallationsUseCase {
	return &ListInstallationsUseCase{repository: repository}
}

func (uc *ListInstallationsUseCase) Execute(ctx context.Context, input ListInstallationsInput) ([]entity.Installation, error) {
	return uc.repository.ListByUserID(ctx, input.UserID)
}
