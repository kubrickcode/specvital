package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/auth/jwt"
)

type Container struct {
	DB         *pgxpool.Pool
	GitClient  client.GitClient
	JWTManager *jwt.Manager
	Queue      *asynq.Client
}

type Config struct {
	DatabaseURL string
	JWTSecret   string
	RedisURL    string
}

func ConfigFromEnv() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}
}

func NewContainer(ctx context.Context, cfg Config) (*Container, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}

	pool, err := NewPostgresPool(ctx, PostgresConfig{
		URL: cfg.DatabaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	jwtManager, err := jwt.NewManager(cfg.JWTSecret)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("jwt: %w", err)
	}

	queueClient, err := NewAsynqClient(cfg.RedisURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("asynq: %w", err)
	}

	gitClient := client.NewGitClient()

	return &Container{
		DB:         pool,
		GitClient:  gitClient,
		JWTManager: jwtManager,
		Queue:      queueClient,
	}, nil
}

func (c *Container) Close() error {
	var errs []error

	if c.Queue != nil {
		if err := c.Queue.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close queue: %w", err))
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
