package infra

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/specvital/core/pkg/crypto"
	"github.com/specvital/web/src/backend/internal/client"
	authadapter "github.com/specvital/web/src/backend/modules/auth/adapter"
	authport "github.com/specvital/web/src/backend/modules/auth/domain/port"
	ghappport "github.com/specvital/web/src/backend/modules/github-app/domain/port"
)

type Container struct {
	River                  *RiverClient
	CookieDomain           string
	DB                     *pgxpool.Pool
	Encryptor              crypto.Encryptor
	Environment            string
	FrontendURL            string
	GitClient              client.GitClient
	GitHubAppClient        ghappport.GitHubAppClient
	GitHubAppWebhookSecret string
	GitHubOAuth            authport.OAuthClient
	JWTManager             authport.TokenManager
	SecureCookie           bool
}

type Config struct {
	CookieDomain            string
	DatabaseURL             string
	EncryptionKey           string
	Environment             string
	FrontendURL             string
	GitHubAppID             int64
	GitHubAppPrivateKey     []byte
	GitHubAppSlug           string
	GitHubAppWebhookSecret  string
	GitHubOAuthClientID     string
	GitHubOAuthClientSecret string
	GitHubOAuthRedirectURL  string
	JWTSecret               string
	SecureCookie            bool
}

func ConfigFromEnv() Config {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	var ghAppID int64
	if id := os.Getenv("GITHUB_APP_ID"); id != "" {
		var err error
		ghAppID, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Warn("invalid GITHUB_APP_ID, GitHub App will be disabled", "value", id, "error", err)
			ghAppID = 0
		}
	}

	var ghAppPrivateKey []byte
	if key := os.Getenv("GITHUB_APP_PRIVATE_KEY"); key != "" {
		ghAppPrivateKey = []byte(strings.ReplaceAll(key, "\\n", "\n"))
	}

	return Config{
		CookieDomain:            os.Getenv("COOKIE_DOMAIN"),
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		EncryptionKey:           os.Getenv("ENCRYPTION_KEY"),
		Environment:             os.Getenv("ENV"),
		FrontendURL:             frontendURL,
		GitHubAppID:             ghAppID,
		GitHubAppPrivateKey:     ghAppPrivateKey,
		GitHubAppSlug:           os.Getenv("GITHUB_APP_SLUG"),
		GitHubAppWebhookSecret:  os.Getenv("GITHUB_APP_WEBHOOK_SECRET"),
		GitHubOAuthClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		GitHubOAuthClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		GitHubOAuthRedirectURL:  os.Getenv("GITHUB_OAUTH_REDIRECT_URL"),
		JWTSecret:               os.Getenv("JWT_SECRET"),
		SecureCookie:            os.Getenv("SECURE_COOKIE") == "true",
	}
}

func NewContainer(ctx context.Context, cfg Config) (*Container, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	var cleanups []func()
	cleanup := func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			cleanups[i]()
		}
	}

	pool, err := NewPostgresPool(ctx, PostgresConfig{
		URL: cfg.DatabaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	cleanups = append(cleanups, pool.Close)

	jwtManager, err := authadapter.NewJWTTokenManager(cfg.JWTSecret)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("jwt: %w", err)
	}

	riverClient, err := NewRiverClient(pool)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("river: %w", err)
	}

	encryptor, err := crypto.NewEncryptorFromBase64(cfg.EncryptionKey)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("encryptor: %w", err)
	}

	githubClient, err := authadapter.NewGitHubOAuthClient(&authadapter.GitHubOAuthConfig{
		ClientID:     cfg.GitHubOAuthClientID,
		ClientSecret: cfg.GitHubOAuthClientSecret,
		RedirectURL:  cfg.GitHubOAuthRedirectURL,
		Scopes:       []string{"user:email", "repo", "read:org"},
	})
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("github oauth: %w", err)
	}

	gitClient := client.NewGitClient()

	ghAppClient, err := client.NewGitHubAppClient(client.GitHubAppConfig{
		AppID:      cfg.GitHubAppID,
		AppSlug:    cfg.GitHubAppSlug,
		PrivateKey: cfg.GitHubAppPrivateKey,
	})
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("github app: %w", err)
	}

	return &Container{
		River:                  riverClient,
		CookieDomain:           cfg.CookieDomain,
		DB:                     pool,
		Encryptor:              encryptor,
		Environment:            cfg.Environment,
		FrontendURL:            cfg.FrontendURL,
		GitClient:              gitClient,
		GitHubAppClient:        ghAppClient,
		GitHubAppWebhookSecret: cfg.GitHubAppWebhookSecret,
		GitHubOAuth:            githubClient,
		JWTManager:             jwtManager,
		SecureCookie:           cfg.SecureCookie,
	}, nil
}

func validateConfig(cfg Config) error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	if cfg.GitHubAppID == 0 {
		return fmt.Errorf("GITHUB_APP_ID is required")
	}
	if len(cfg.GitHubAppPrivateKey) == 0 {
		return fmt.Errorf("GITHUB_APP_PRIVATE_KEY is required")
	}
	if cfg.GitHubAppSlug == "" {
		return fmt.Errorf("GITHUB_APP_SLUG is required")
	}
	if cfg.GitHubAppWebhookSecret == "" {
		return fmt.Errorf("GITHUB_APP_WEBHOOK_SECRET is required")
	}
	if cfg.GitHubOAuthClientID == "" {
		return fmt.Errorf("GITHUB_OAUTH_CLIENT_ID is required")
	}
	if cfg.GitHubOAuthClientSecret == "" {
		return fmt.Errorf("GITHUB_OAUTH_CLIENT_SECRET is required")
	}
	if cfg.GitHubOAuthRedirectURL == "" {
		return fmt.Errorf("GITHUB_OAUTH_REDIRECT_URL is required")
	}
	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.Environment == "production" && !cfg.SecureCookie {
		return fmt.Errorf("SECURE_COOKIE=true is required when ENV=production")
	}
	return nil
}

func (c *Container) Close() error {
	var errs []error

	if c.DB != nil {
		c.DB.Close()
	}

	if len(errs) > 0 {
		return fmt.Errorf("close container: %v", errs)
	}
	return nil
}
