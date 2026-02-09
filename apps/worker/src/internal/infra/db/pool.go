package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// Pool size: (workers * 2) + buffer
	// (5 * 2) + 10 = 20
	defaultMaxConns = 20
	defaultMinConns = 5

	// Connection lifecycle settings for long-running analysis jobs
	defaultConnectTimeout   = 10 * time.Second
	defaultHealthCheckPeriod = 30 * time.Second
	defaultMaxConnIdleTime  = 5 * time.Minute
	defaultMaxConnLifetime  = 30 * time.Minute
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	config.MaxConns = defaultMaxConns
	config.MinConns = defaultMinConns
	config.ConnConfig.ConnectTimeout = defaultConnectTimeout
	config.HealthCheckPeriod = defaultHealthCheckPeriod
	config.MaxConnIdleTime = defaultMaxConnIdleTime
	config.MaxConnLifetime = defaultMaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
