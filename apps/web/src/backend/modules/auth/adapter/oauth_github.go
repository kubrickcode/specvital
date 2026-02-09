package adapter

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cockroachdb/errors"
	"golang.org/x/oauth2"
	oauthGithub "golang.org/x/oauth2/github"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

const (
	githubUserURL = "https://api.github.com/user"
	DefaultScope  = "user:email"
)

type GitHubOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type GitHubOAuthClient struct {
	config      *oauth2.Config
	oauthConfig *GitHubOAuthConfig
}

var _ port.OAuthClient = (*GitHubOAuthClient)(nil)

var ErrEmptyState = errors.New("state parameter is required for CSRF protection")

func NewGitHubOAuthClient(config *GitHubOAuthConfig) (*GitHubOAuthClient, error) {
	if config.ClientID == "" {
		return nil, errors.New("github client ID is required")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("github client secret is required")
	}
	if config.RedirectURL == "" {
		return nil, errors.New("github redirect URL is required")
	}

	scopes := config.Scopes
	if len(scopes) == 0 {
		scopes = []string{DefaultScope}
	}

	return &GitHubOAuthClient{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes:       scopes,
			Endpoint:     oauthGithub.Endpoint,
		},
		oauthConfig: config,
	}, nil
}

func (c *GitHubOAuthClient) GenerateAuthURL(state string) (string, error) {
	if state == "" {
		return "", ErrEmptyState
	}
	return c.config.AuthCodeURL(state), nil
}

func (c *GitHubOAuthClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", errors.Wrap(domain.ErrInvalidCode, "code is empty")
	}

	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		var retrieveErr *oauth2.RetrieveError
		if errors.As(err, &retrieveErr) && retrieveErr.Response != nil {
			switch retrieveErr.Response.StatusCode {
			case http.StatusUnauthorized, http.StatusBadRequest:
				return "", errors.Wrap(domain.ErrInvalidCode, "bad verification code")
			case http.StatusForbidden:
				return "", errors.Wrap(domain.ErrAccessDenied, "access denied")
			}
		}
		return "", errors.Wrap(domain.ErrNetworkFailure, err.Error())
	}

	if token.AccessToken == "" {
		return "", errors.Wrap(domain.ErrInvalidCode, "empty access token")
	}

	return token.AccessToken, nil
}

func (c *GitHubOAuthClient) GetUserInfo(ctx context.Context, accessToken string) (*entity.OAuthUserInfo, error) {
	if accessToken == "" {
		return nil, errors.Wrap(domain.ErrInvalidGitHubToken, "access token is empty")
	}

	token := &oauth2.Token{AccessToken: accessToken}
	httpClient := c.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
	if err != nil {
		return nil, errors.Wrap(domain.ErrNetworkFailure, "failed to create request")
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(domain.ErrNetworkFailure, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.Wrap(domain.ErrInvalidGitHubToken, "unauthorized")
	}
	if resp.StatusCode == http.StatusForbidden {
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return nil, errors.Wrap(domain.ErrRateLimited, "rate limit exceeded")
		}
		return nil, errors.Wrap(domain.ErrInvalidGitHubToken, "forbidden")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(domain.ErrNetworkFailure, "unexpected status: %d", resp.StatusCode)
	}

	var userResp githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, errors.Wrap(domain.ErrNetworkFailure, "failed to parse response")
	}

	return &entity.OAuthUserInfo{
		AvatarURL:  userResp.AvatarURL,
		Email:      userResp.Email,
		ExternalID: strconv.FormatInt(userResp.ID, 10),
		Username:   userResp.Login,
	}, nil
}

type githubUserResponse struct {
	AvatarURL string  `json:"avatar_url"`
	Email     *string `json:"email"`
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
}
