package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/specvital/worker/internal/adapter/repository/postgres"
	"github.com/specvital/worker/internal/infra/db"
	retentionuc "github.com/specvital/worker/internal/usecase/retention"
)

// DefaultRetentionTimeout allows sufficient time for batch processing
// large datasets while preventing runaway cron jobs.
const DefaultRetentionTimeout = 30 * time.Minute

// RetentionConfig holds configuration for the retention cleanup service.
type RetentionConfig struct {
	BatchSize   int
	BatchSleep  time.Duration
	DatabaseURL string
	ServiceName string
	Timeout     time.Duration
}

// Validate checks that required retention configuration fields are set.
func (c *RetentionConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}
	return nil
}

// applyDefaults sets default values for optional retention configuration.
func (c *RetentionConfig) applyDefaults() {
	if c.Timeout <= 0 {
		c.Timeout = DefaultRetentionTimeout
	}
}

// RunRetentionCleanup executes the retention cleanup process.
// This is designed to run as a Railway Cron job.
// Returns the cleanup result on success.
func RunRetentionCleanup(cfg RetentionConfig) (*retentionuc.CleanupResult, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	cfg.applyDefaults()

	slog.Info("starting service", "name", cfg.ServiceName)
	slog.Info("config loaded", "database_url", maskURL(cfg.DatabaseURL))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("database connection: %w", err)
	}
	defer pool.Close()

	slog.Info("postgres connected")

	repo := postgres.NewRetentionRepository(pool)

	var opts []retentionuc.Option
	if cfg.BatchSize > 0 {
		opts = append(opts, retentionuc.WithBatchSize(cfg.BatchSize))
	}
	if cfg.BatchSleep > 0 {
		opts = append(opts, retentionuc.WithBatchSleep(cfg.BatchSleep))
	}

	usecase := retentionuc.NewCleanupUseCase(repo, opts...)

	result, err := usecase.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("execute cleanup: %w", err)
	}

	slog.Info("service completed",
		"name", cfg.ServiceName,
		"total_deleted", result.TotalDeleted(),
		"duration", result.Duration(),
	)

	return &result, nil
}
