package mapper

import (
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/github/domain/entity"
)

func ToGitHubOrganization(org entity.Organization) api.GitHubOrganization {
	o := api.GitHubOrganization{
		AccessStatus: toAPIAccessStatus(org.AccessStatus),
		ID:           org.ID,
		Login:        org.Login,
	}

	if org.AvatarURL != "" {
		o.AvatarURL = &org.AvatarURL
	}
	if org.Description != "" {
		o.Description = &org.Description
	}

	return o
}

func toAPIAccessStatus(status entity.AccessStatus) api.OrganizationAccessStatus {
	switch status {
	case entity.AccessStatusAccessible:
		return api.OrganizationAccessStatusAccessible
	case entity.AccessStatusPending:
		return api.OrganizationAccessStatusPending
	default:
		return api.OrganizationAccessStatusRestricted
	}
}

func ToGitHubOrganizations(orgs []entity.Organization) []api.GitHubOrganization {
	result := make([]api.GitHubOrganization, len(orgs))
	for i, org := range orgs {
		result[i] = ToGitHubOrganization(org)
	}
	return result
}

func ToGitHubRepositories(repos []entity.Repository) []api.GitHubRepository {
	result := make([]api.GitHubRepository, len(repos))
	for i, repo := range repos {
		result[i] = ToGitHubRepository(repo)
	}
	return result
}

func ToGitHubRepository(repo entity.Repository) api.GitHubRepository {
	r := api.GitHubRepository{
		DefaultBranch: repo.DefaultBranch,
		FullName:      repo.FullName,
		ID:            repo.ID,
		IsPrivate:     repo.Private,
		Name:          repo.Name,
		Owner:         repo.Owner,
	}

	if repo.Description != "" {
		r.Description = &repo.Description
	}
	if repo.PushedAt != nil {
		r.PushedAt = repo.PushedAt
	}

	return r
}
