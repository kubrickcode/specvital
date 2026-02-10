package usecase

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	authdomain "github.com/kubrickcode/specvital/apps/web/src/backend/modules/auth/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/domain/port"
)

const maxReposPerFetch = 1000

func getGitHubClient(
	ctx context.Context,
	clientFactory port.GitHubClientFactory,
	tokenProvider port.TokenProvider,
	userID string,
) (port.GitHubClient, error) {
	token, err := tokenProvider.GetUserGitHubToken(ctx, userID)
	if err != nil {
		if errors.Is(err, authdomain.ErrUserNotFound) || errors.Is(err, authdomain.ErrNoGitHubToken) {
			return nil, domain.ErrNoGitHubToken
		}
		return nil, fmt.Errorf("get github token: %w", err)
	}

	return clientFactory(token), nil
}

func mapGitHubRepositoryToEntity(r port.GitHubRepository) entity.Repository {
	return entity.Repository{
		Archived:      r.Archived,
		DefaultBranch: r.DefaultBranch,
		Description:   r.Description,
		Disabled:      r.Disabled,
		Fork:          r.Fork,
		FullName:      r.FullName,
		HTMLURL:       r.HTMLURL,
		ID:            r.ID,
		Language:      r.Language,
		Name:          r.Name,
		Owner:         r.Owner,
		Private:       r.Private,
		PushedAt:      r.PushedAt,
		StarCount:     r.StarCount,
		Visibility:    r.Visibility,
	}
}

func mapGitHubRepositoriesToEntities(repos []port.GitHubRepository) []entity.Repository {
	result := make([]entity.Repository, len(repos))
	for i, r := range repos {
		result[i] = mapGitHubRepositoryToEntity(r)
	}
	return result
}

func mapGitHubOrganizationToEntity(org port.GitHubOrganization) entity.Organization {
	return entity.Organization{
		AvatarURL:   org.AvatarURL,
		Description: org.Description,
		HTMLURL:     org.HTMLURL,
		ID:          org.ID,
		Login:       org.Login,
	}
}

func mapGitHubOrganizationsToEntities(orgs []port.GitHubOrganization) []entity.Organization {
	result := make([]entity.Organization, len(orgs))
	for i, o := range orgs {
		result[i] = mapGitHubOrganizationToEntity(o)
	}
	return result
}

func mapRepositoryRecordToEntity(r port.RepositoryRecord) entity.Repository {
	return entity.Repository{
		Archived:      r.Archived,
		DefaultBranch: r.DefaultBranch,
		Description:   r.Description,
		Disabled:      r.Disabled,
		Fork:          r.Fork,
		FullName:      r.FullName,
		HTMLURL:       r.HTMLURL,
		ID:            r.ID,
		Language:      r.Language,
		Name:          r.Name,
		Owner:         r.Owner,
		Private:       r.Private,
		PushedAt:      r.PushedAt,
		StarCount:     r.StarCount,
		Visibility:    r.Visibility,
	}
}

func mapRepositoryRecordsToEntities(records []port.RepositoryRecord) []entity.Repository {
	result := make([]entity.Repository, len(records))
	for i, r := range records {
		result[i] = mapRepositoryRecordToEntity(r)
	}
	return result
}

func mapOrganizationRecordToEntity(r port.OrganizationRecord) entity.Organization {
	return entity.Organization{
		AvatarURL:   r.AvatarURL,
		Description: r.Description,
		HTMLURL:     r.HTMLURL,
		ID:          r.ID,
		Login:       r.Login,
		Role:        r.Role,
	}
}

func mapOrganizationRecordsToEntities(records []port.OrganizationRecord) []entity.Organization {
	result := make([]entity.Organization, len(records))
	for i, r := range records {
		result[i] = mapOrganizationRecordToEntity(r)
	}
	return result
}

func mapEntityToRepositoryRecord(e entity.Repository) port.RepositoryRecord {
	return port.RepositoryRecord{
		Archived:      e.Archived,
		DefaultBranch: e.DefaultBranch,
		Description:   e.Description,
		Disabled:      e.Disabled,
		Fork:          e.Fork,
		FullName:      e.FullName,
		HTMLURL:       e.HTMLURL,
		ID:            e.ID,
		Language:      e.Language,
		Name:          e.Name,
		Owner:         e.Owner,
		Private:       e.Private,
		PushedAt:      e.PushedAt,
		StarCount:     e.StarCount,
		Visibility:    e.Visibility,
	}
}

func mapEntitiesToRepositoryRecords(entities []entity.Repository) []port.RepositoryRecord {
	result := make([]port.RepositoryRecord, len(entities))
	for i, e := range entities {
		result[i] = mapEntityToRepositoryRecord(e)
	}
	return result
}

func mapEntityToOrganizationRecord(e entity.Organization) port.OrganizationRecord {
	return port.OrganizationRecord{
		AvatarURL:   e.AvatarURL,
		Description: e.Description,
		HTMLURL:     e.HTMLURL,
		ID:          e.ID,
		Login:       e.Login,
		Role:        e.Role,
	}
}

func mapEntitiesToOrganizationRecords(entities []entity.Organization) []port.OrganizationRecord {
	result := make([]port.OrganizationRecord, len(entities))
	for i, e := range entities {
		result[i] = mapEntityToOrganizationRecord(e)
	}
	return result
}

func mapClientError(err error) error {
	switch {
	case errors.Is(err, port.ErrGitHubUnauthorized):
		return domain.ErrUnauthorized
	case errors.Is(err, port.ErrGitHubInsufficientScope):
		return domain.ErrInsufficientScope
	case errors.Is(err, port.ErrGitHubNotFound):
		return domain.ErrOrganizationNotFound
	case port.IsRateLimitError(err):
		var rateLimitErr *port.RateLimitError
		if errors.As(err, &rateLimitErr) {
			return &domain.RateLimitError{
				Limit:     rateLimitErr.Limit,
				Remaining: rateLimitErr.Remaining,
				ResetAt:   rateLimitErr.ResetAt,
			}
		}
		return &domain.RateLimitError{}
	default:
		return err
	}
}
