package adapter

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

var _ port.GitHubClient = (*GitHubClientAdapter)(nil)

type GitHubClientAdapter struct {
	client client.GitHubClient
}

func NewGitHubClientFactory(factory client.GitHubClientFactory) port.GitHubClientFactory {
	return func(token string) port.GitHubClient {
		return &GitHubClientAdapter{client: factory(token)}
	}
}

func (a *GitHubClientAdapter) GetOrganization(ctx context.Context, org string) (*port.GitHubOrganization, error) {
	ghOrg, err := a.client.GetOrganization(ctx, org)
	if err != nil {
		return nil, mapClientError(err)
	}
	return toPortOrganization(ghOrg), nil
}

func (a *GitHubClientAdapter) ListOrgRepositories(ctx context.Context, org string, maxResults int) ([]port.GitHubRepository, error) {
	repos, err := a.client.ListOrgRepositories(ctx, org, maxResults)
	if err != nil {
		return nil, mapClientError(err)
	}
	return toPortRepositories(repos), nil
}

func (a *GitHubClientAdapter) ListUserOrganizations(ctx context.Context) ([]port.GitHubOrganization, error) {
	orgs, err := a.client.ListUserOrganizations(ctx)
	if err != nil {
		return nil, mapClientError(err)
	}
	return toPortOrganizations(orgs), nil
}

func (a *GitHubClientAdapter) ListUserRepositories(ctx context.Context, maxResults int) ([]port.GitHubRepository, error) {
	repos, err := a.client.ListUserRepositories(ctx, maxResults)
	if err != nil {
		return nil, mapClientError(err)
	}
	return toPortRepositories(repos), nil
}

func mapClientError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, client.ErrGitHubUnauthorized) {
		return port.ErrGitHubUnauthorized
	}
	if errors.Is(err, client.ErrGitHubInsufficientScope) {
		return port.ErrGitHubInsufficientScope
	}
	if errors.Is(err, client.ErrGitHubNotFound) {
		return port.ErrGitHubNotFound
	}

	var clientRateLimitErr *client.RateLimitError
	if errors.As(err, &clientRateLimitErr) {
		return &port.RateLimitError{
			Limit:     clientRateLimitErr.Limit,
			Remaining: clientRateLimitErr.Remaining,
			ResetAt:   clientRateLimitErr.ResetAt,
		}
	}

	return err
}

func toPortRepository(r client.GitHubRepository) port.GitHubRepository {
	return port.GitHubRepository{
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

func toPortRepositories(repos []client.GitHubRepository) []port.GitHubRepository {
	result := make([]port.GitHubRepository, len(repos))
	for i, r := range repos {
		result[i] = toPortRepository(r)
	}
	return result
}

func toPortOrganization(o *client.GitHubOrganization) *port.GitHubOrganization {
	if o == nil {
		return nil
	}
	return &port.GitHubOrganization{
		AvatarURL:   o.AvatarURL,
		Description: o.Description,
		HTMLURL:     o.HTMLURL,
		ID:          o.ID,
		Login:       o.Login,
	}
}

func toPortOrganizations(orgs []client.GitHubOrganization) []port.GitHubOrganization {
	result := make([]port.GitHubOrganization, len(orgs))
	for i, o := range orgs {
		result[i] = port.GitHubOrganization{
			AvatarURL:   o.AvatarURL,
			Description: o.Description,
			HTMLURL:     o.HTMLURL,
			ID:          o.ID,
			Login:       o.Login,
		}
	}
	return result
}
