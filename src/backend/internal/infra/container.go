package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/specvital/core/pkg/crypto"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/auth/github"
	"github.com/specvital/web/src/backend/modules/auth/jwt"
)

type Container struct {
	River        *RiverClient
	CookieDomain string
	DB           *pgxpool.Pool
	Encryptor    crypto.Encryptor
	FrontendURL  string
	GitClient    client.GitClient
	GitHubOAuth  github.Client
	JWTManager   *jwt.Manager
	SecureCookie bool
}

type Config struct {
	CookieDomain            string
	DatabaseURL             string
	EncryptionKey           string
	Environment             string
	FrontendURL             string
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
	return Config{
		CookieDomain:            os.Getenv("COOKIE_DOMAIN"),
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		EncryptionKey:           os.Getenv("ENCRYPTION_KEY"),
		Environment:             os.Getenv("ENV"),
		FrontendURL:             frontendURL,
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

	jwtManager, err := jwt.NewManager(cfg.JWTSecret)
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

	githubClient, err := github.NewClient(&github.Config{
		ClientID:     cfg.GitHubOAuthClientID,
		ClientSecret: cfg.GitHubOAuthClientSecret,
		RedirectURL:  cfg.GitHubOAuthRedirectURL,
		Scopes:       []string{"user:email", "repo"},
	})
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("github oauth: %w", err)
	}

	gitClient := client.NewGitClient()

	return &Container{
		River:        riverClient,
		CookieDomain: cfg.CookieDomain,
		DB:           pool,
		Encryptor:    encryptor,
		FrontendURL:  cfg.FrontendURL,
		GitClient:    gitClient,
		GitHubOAuth:  githubClient,
		JWTManager:   jwtManager,
		SecureCookie: cfg.SecureCookie,
	}, nil
}

func validateConfig(cfg Config) error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
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
