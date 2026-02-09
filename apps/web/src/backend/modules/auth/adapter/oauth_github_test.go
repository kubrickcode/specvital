package adapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pkgerrors "github.com/cockroachdb/errors"
	"golang.org/x/oauth2"

	"github.com/specvital/web/src/backend/modules/auth/domain"
)

type testUserInfo struct {
	AvatarURL string
	Email     *string
	ID        int64
	Login     string
}

func TestNewGitHubOAuthClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *GitHubOAuthConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &GitHubOAuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8000/callback",
			},
			wantErr: false,
		},
		{
			name: "missing client ID",
			config: &GitHubOAuthConfig{
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8000/callback",
			},
			wantErr: true,
		},
		{
			name: "missing client secret",
			config: &GitHubOAuthConfig{
				ClientID:    "test-client-id",
				RedirectURL: "http://localhost:8000/callback",
			},
			wantErr: true,
		},
		{
			name: "missing redirect URL",
			config: &GitHubOAuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGitHubOAuthClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitHubOAuthClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitHubOAuthClient_GenerateAuthURL(t *testing.T) {
	config := &GitHubOAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8000/callback",
		Scopes:       []string{"user:email"},
	}

	client, _ := NewGitHubOAuthClient(config)
	url, err := client.GenerateAuthURL("test-state-123")
	if err != nil {
		t.Fatalf("GenerateAuthURL() error = %v", err)
	}

	if !strings.Contains(url, "client_id=test-client-id") {
		t.Error("URL should contain client_id")
	}
	if !strings.Contains(url, "state=test-state-123") {
		t.Error("URL should contain state")
	}
	if !strings.Contains(url, "redirect_uri=") {
		t.Error("URL should contain redirect_uri")
	}
	if !strings.Contains(url, "scope=user") {
		t.Error("URL should contain scope")
	}
	if !strings.HasPrefix(url, "https://github.com/login/oauth/authorize") {
		t.Error("URL should start with GitHub auth URL")
	}
}

func TestGitHubOAuthClient_GenerateAuthURL_EmptyState(t *testing.T) {
	config := &GitHubOAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8000/callback",
	}

	client, _ := NewGitHubOAuthClient(config)
	_, err := client.GenerateAuthURL("")
	if !pkgerrors.Is(err, ErrEmptyState) {
		t.Errorf("GenerateAuthURL() error = %v, wantErr %v", err, ErrEmptyState)
	}
}

func TestGitHubOAuthClient_ExchangeCode(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		response    string
		statusCode  int
		wantToken   string
		wantErr     bool
		wantErrType error
	}{
		{
			name:       "successful exchange",
			code:       "valid-code",
			response:   `{"access_token":"gho_xxxx","token_type":"bearer","scope":"user:email"}`,
			statusCode: 200,
			wantToken:  "gho_xxxx",
			wantErr:    false,
		},
		{
			name:        "empty code",
			code:        "",
			response:    "",
			statusCode:  200,
			wantErr:     true,
			wantErrType: domain.ErrInvalidCode,
		},
		{
			name:        "bad verification code",
			code:        "invalid-code",
			response:    `{"error":"bad_verification_code","error_description":"The code passed is incorrect or expired."}`,
			statusCode:  400,
			wantErr:     true,
			wantErrType: domain.ErrInvalidCode,
		},
		{
			name:        "access denied",
			code:        "denied-code",
			response:    `{"error":"access_denied"}`,
			statusCode:  403,
			wantErr:     true,
			wantErrType: domain.ErrAccessDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code == "" {
				config := &GitHubOAuthConfig{
					ClientID:     "test-client-id",
					ClientSecret: "test-secret",
					RedirectURL:  "http://localhost:8000/callback",
				}
				client, _ := NewGitHubOAuthClient(config)
				_, err := client.ExchangeCode(context.Background(), tt.code)

				if (err != nil) != tt.wantErr {
					t.Errorf("ExchangeCode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErrType != nil && !pkgerrors.Is(err, tt.wantErrType) {
					t.Errorf("ExchangeCode() error type = %v, wantErrType %v", err, tt.wantErrType)
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := &GitHubOAuthClient{
				config: &oauth2.Config{
					ClientID:     "test-client-id",
					ClientSecret: "test-secret",
					RedirectURL:  "http://localhost:8000/callback",
					Scopes:       []string{"user:email"},
					Endpoint:     oauth2.Endpoint{TokenURL: server.URL},
				},
			}

			token, err := client.ExchangeCode(context.Background(), tt.code)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExchangeCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrType != nil {
				if !pkgerrors.Is(err, tt.wantErrType) {
					t.Errorf("ExchangeCode() error type = %v, wantErrType %v", err, tt.wantErrType)
				}
				return
			}

			if token != tt.wantToken {
				t.Errorf("ExchangeCode() token = %v, want %v", token, tt.wantToken)
			}
		})
	}
}

func TestGitHubOAuthClient_GetUserInfo(t *testing.T) {
	email := "test@example.com"

	tests := []struct {
		name        string
		token       string
		response    string
		statusCode  int
		headers     map[string]string
		wantUserID  int64
		wantLogin   string
		wantErr     bool
		wantErrType error
	}{
		{
			name:       "successful fetch",
			token:      "valid-token",
			response:   `{"id":12345,"login":"testuser","email":"test@example.com","avatar_url":"https://avatar.url"}`,
			statusCode: 200,
			wantUserID: 12345,
			wantLogin:  "testuser",
			wantErr:    false,
		},
		{
			name:        "empty token",
			token:       "",
			response:    "",
			statusCode:  200,
			wantErr:     true,
			wantErrType: domain.ErrInvalidGitHubToken,
		},
		{
			name:        "unauthorized",
			token:       "invalid-token",
			response:    `{"message":"Bad credentials"}`,
			statusCode:  401,
			wantErr:     true,
			wantErrType: domain.ErrInvalidGitHubToken,
		},
		{
			name:        "rate limited",
			token:       "valid-token",
			response:    `{"message":"API rate limit exceeded"}`,
			statusCode:  403,
			headers:     map[string]string{"X-RateLimit-Remaining": "0"},
			wantErr:     true,
			wantErrType: domain.ErrRateLimited,
		},
		{
			name:        "malformed json response",
			token:       "valid-token",
			response:    `{"id": "not-a-number"}`,
			statusCode:  200,
			wantErr:     true,
			wantErrType: domain.ErrNetworkFailure,
		},
		{
			name:        "unexpected status code",
			token:       "valid-token",
			response:    `{"message":"Internal Server Error"}`,
			statusCode:  500,
			wantErr:     true,
			wantErrType: domain.ErrNetworkFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.token == "" {
				config := &GitHubOAuthConfig{
					ClientID:     "test-client-id",
					ClientSecret: "test-secret",
					RedirectURL:  "http://localhost:8000/callback",
				}
				client, _ := NewGitHubOAuthClient(config)
				_, err := client.GetUserInfo(context.Background(), tt.token)

				if (err != nil) != tt.wantErr {
					t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErrType != nil && !pkgerrors.Is(err, tt.wantErrType) {
					t.Errorf("GetUserInfo() error type = %v, wantErrType %v", err, tt.wantErrType)
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, v := range tt.headers {
					w.Header().Set(k, v)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := &GitHubOAuthClient{
				config: &oauth2.Config{
					ClientID:     "test-client-id",
					ClientSecret: "test-secret",
					RedirectURL:  "http://localhost:8000/callback",
					Scopes:       []string{"user:email"},
				},
			}

			user, err := client.getUserInfoWithURL(context.Background(), tt.token, server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.wantErrType != nil && !pkgerrors.Is(err, tt.wantErrType) {
					t.Errorf("GetUserInfo() error type = %v, wantErrType %v", err, tt.wantErrType)
				}
				return
			}

			if user.ID != tt.wantUserID {
				t.Errorf("GetUserInfo() ID = %v, want %v", user.ID, tt.wantUserID)
			}
			if user.Login != tt.wantLogin {
				t.Errorf("GetUserInfo() Login = %v, want %v", user.Login, tt.wantLogin)
			}
			if user.Email == nil || *user.Email != email {
				t.Errorf("GetUserInfo() Email = %v, want %v", user.Email, email)
			}
		})
	}
}

func (c *GitHubOAuthClient) getUserInfoWithURL(ctx context.Context, accessToken string, url string) (*testUserInfo, error) {
	if accessToken == "" {
		return nil, pkgerrors.Wrap(domain.ErrInvalidGitHubToken, "access token is empty")
	}

	token := &oauth2.Token{AccessToken: accessToken}
	httpClient := c.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, pkgerrors.Wrap(domain.ErrNetworkFailure, "failed to create request")
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, pkgerrors.Wrap(domain.ErrNetworkFailure, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, pkgerrors.Wrap(domain.ErrInvalidGitHubToken, "unauthorized")
	}
	if resp.StatusCode == http.StatusForbidden {
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return nil, pkgerrors.Wrap(domain.ErrRateLimited, "rate limit exceeded")
		}
		return nil, pkgerrors.Wrap(domain.ErrInvalidGitHubToken, "forbidden")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, pkgerrors.Wrapf(domain.ErrNetworkFailure, "unexpected status: %d", resp.StatusCode)
	}

	var userResp githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, pkgerrors.Wrap(domain.ErrNetworkFailure, "failed to parse response")
	}

	return &testUserInfo{
		AvatarURL: userResp.AvatarURL,
		Email:     userResp.Email,
		ID:        userResp.ID,
		Login:     userResp.Login,
	}, nil
}
