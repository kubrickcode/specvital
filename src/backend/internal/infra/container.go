package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	DB    *pgxpool.Pool
	Queue *asynq.Client
}

type Config struct {
	DatabaseURL string
	RedisURL    string
}

func ConfigFromEnv() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}
}

func NewContainer(ctx context.Context, cfg Config) (*Container, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
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

	queueClient, err := NewAsynqClient(cfg.RedisURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("asynq: %w", err)
	}

	return &Container{
		DB:    pool,
		Queue: queueClient,
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
