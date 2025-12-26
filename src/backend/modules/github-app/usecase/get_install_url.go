package usecase

import (
	"github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

type GetInstallURLUseCase struct {
	ghAppClient port.GitHubAppClient
}

func NewGetInstallURLUseCase(ghAppClient port.GitHubAppClient) *GetInstallURLUseCase {
	return &GetInstallURLUseCase{ghAppClient: ghAppClient}
}

func (uc *GetInstallURLUseCase) Execute() string {
	return uc.ghAppClient.GetInstallationURL()
}
