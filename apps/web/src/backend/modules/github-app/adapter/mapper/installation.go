package mapper

import (
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/github-app/domain/entity"
)

func ToAPIInstallation(e entity.Installation) api.GitHubAppInstallation {
	return api.GitHubAppInstallation{
		AccountAvatarURL: e.AccountAvatarURL,
		AccountID:        e.AccountID,
		AccountLogin:     e.AccountLogin,
		AccountType:      toAPIAccountType(e.AccountType),
		CreatedAt:        e.CreatedAt,
		ID:               e.ID,
		InstallationID:   e.InstallationID,
		IsSuspended:      e.IsSuspended(),
	}
}

func ToAPIInstallations(entities []entity.Installation) []api.GitHubAppInstallation {
	result := make([]api.GitHubAppInstallation, len(entities))
	for i, e := range entities {
		result[i] = ToAPIInstallation(e)
	}
	return result
}

func toAPIAccountType(t entity.AccountType) api.GitHubAppInstallationAccountType {
	switch t {
	case entity.AccountTypeOrganization:
		return api.GitHubAppInstallationAccountTypeOrganization
	default:
		return api.GitHubAppInstallationAccountTypeUser
	}
}
