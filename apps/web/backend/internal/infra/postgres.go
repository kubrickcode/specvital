package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns = 10
	defaultMinConns = 2
)

type PostgresConfig struct {
	URL      string
	MaxConns int32
	MinConns int32
}

func NewPostgresPool(ctx context.Context, cfg PostgresConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	maxConns := cfg.MaxConns
	if maxConns == 0 {
		maxConns = defaultMaxConns
	}

	minConns := cfg.MinConns
	if minConns == 0 {
		minConns = defaultMinConns
	}

	poolConfig.MaxConns = maxConns
	poolConfig.MinConns = minConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}
