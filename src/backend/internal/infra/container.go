package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/specvital/web/src/backend/common/crypto"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/auth/github"
	"github.com/specvital/web/src/backend/modules/auth/jwt"
)

type Container struct {
	Asynq        *AsynqComponents
	DB           *pgxpool.Pool
	Encryptor    crypto.Encryptor
	FrontendURL  string
	GitClient    client.GitClient
	GitHubOAuth  github.Client
	JWTManager   *jwt.Manager
	SecureCookie bool
}

type Config struct {
	DatabaseURL        string
	EncryptionKey      string
	Environment        string
	FrontendURL        string
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
	JWTSecret          string
	RedisURL           string
	SecureCookie       bool
}

func ConfigFromEnv() Config {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	return Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		EncryptionKey:      os.Getenv("ENCRYPTION_KEY"),
		Environment:        os.Getenv("ENV"),
		FrontendURL:        frontendURL,
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		RedisURL:           os.Getenv("REDIS_URL"),
		SecureCookie:       os.Getenv("SECURE_COOKIE") == "true",
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

	asynqComponents, err := NewAsynqComponents(cfg.RedisURL)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("asynq: %w", err)
	}
	cleanups = append(cleanups, func() { _ = asynqComponents.Close() })

	encryptor, err := crypto.NewEncryptorFromBase64(cfg.EncryptionKey)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("encryptor: %w", err)
	}

	githubClient, err := github.NewClient(&github.Config{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		RedirectURL:  cfg.GitHubRedirectURL,
	})
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("github oauth: %w", err)
	}

	gitClient := client.NewGitClient()

	return &Container{
		Asynq:        asynqComponents,
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
	if cfg.GitHubClientID == "" {
		return fmt.Errorf("GITHUB_CLIENT_ID is required")
	}
	if cfg.GitHubClientSecret == "" {
		return fmt.Errorf("GITHUB_CLIENT_SECRET is required")
	}
	if cfg.GitHubRedirectURL == "" {
		return fmt.Errorf("GITHUB_REDIRECT_URL is required")
	}
	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if cfg.Environment == "production" && !cfg.SecureCookie {
		return fmt.Errorf("SECURE_COOKIE=true is required when ENV=production")
	}
	return nil
}

func (c *Container) Close() error {
	var errs []error

	if c.Asynq != nil {
		if err := c.Asynq.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close asynq: %w", err))
		}
	}
	if c.DB != nil {
		c.DB.Close()
	}

	if len(errs) > 0 {
		return fmt.Errorf("close container: %v", errs)
	}
	return nil
}
