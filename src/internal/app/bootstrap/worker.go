// Package bootstrap provides application startup utilities for worker services.
package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/specvital/collector/internal/app"
	"github.com/specvital/collector/internal/handler/queue"
	"github.com/specvital/collector/internal/infra/db"
	infraqueue "github.com/specvital/collector/internal/infra/queue"
)

const defaultConcurrency = 5

type WorkerConfig struct {
	ServiceName     string
	Concurrency     int
	ShutdownTimeout time.Duration
	DatabaseURL     string
	RedisURL        string
}

func (c *WorkerConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.RedisURL == "" {
		return fmt.Errorf("redis URL is required")
	}
	return nil
}

func (c *WorkerConfig) applyDefaults() {
	if c.Concurrency <= 0 {
		c.Concurrency = defaultConcurrency
	}
	if c.ShutdownTimeout <= 0 {
		c.ShutdownTimeout = infraqueue.DefaultShutdownTimeout
	}
}

const autoRefreshSchedule = "@every 1h"
const schedulerShutdownTimeout = 30 * time.Second

// StartWorker starts the worker service with queue processing and scheduled jobs.
//
// Scheduler jobs use Redis-based distributed locking to prevent duplicate execution
// across multiple worker instances. Horizontal scaling is safe for queue processing,
// but each instance will attempt to acquire the scheduler lock (only one succeeds).
func StartWorker(cfg WorkerConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	cfg.applyDefaults()

	slog.Info("starting service", "name", cfg.ServiceName)
	slog.Info("config loaded",
		"database_url", maskURL(cfg.DatabaseURL),
		"redis_url", maskURL(cfg.RedisURL),
	)

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database connection: %w", err)
	}

	slog.Info("postgres connected")

	srv, err := infraqueue.NewServer(infraqueue.ServerConfig{
		RedisURL:        cfg.RedisURL,
		Concurrency:     cfg.Concurrency,
		ShutdownTimeout: cfg.ShutdownTimeout,
	})
	if err != nil {
		return fmt.Errorf("queue server: %w", err)
	}

	container, err := app.NewContainer(app.ContainerConfig{
		Pool:     pool,
		RedisURL: cfg.RedisURL,
	})
	if err != nil {
		return fmt.Errorf("container: %w", err)
	}

	mux := infraqueue.NewServeMux()
	mux.HandleFunc(queue.TypeAnalyze, container.AnalyzeHandler.ProcessTask)

	if err := container.Scheduler.AddFunc(autoRefreshSchedule, container.AutoRefreshHandler.Run); err != nil {
		container.Close()
		pool.Close()
		return fmt.Errorf("add auto-refresh schedule: %w", err)
	}
	container.Scheduler.Start()
	slog.Info("scheduler started", "schedule", autoRefreshSchedule)

	slog.Info("worker starting", "concurrency", cfg.Concurrency)
	if err := srv.Start(mux); err != nil {
		srv.Shutdown()
		_ = container.Scheduler.StopWithTimeout(schedulerShutdownTimeout)
		container.Close()
		pool.Close()
		return fmt.Errorf("start server: %w", err)
	}
	slog.Info("worker ready", "concurrency", cfg.Concurrency)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	sig := <-shutdown
	slog.Info("shutdown signal received", "signal", sig.String())

	if err := container.Scheduler.StopWithTimeout(schedulerShutdownTimeout); err != nil {
		slog.Warn("scheduler shutdown timeout", "error", err)
	}
	slog.Info("scheduler stopped")

	srv.Shutdown()
	slog.Info("queue server stopped")

	if err := container.Close(); err != nil {
		slog.Error("failed to close container", "error", err)
	}

	pool.Close()
	slog.Info("database pool closed")

	slog.Info("service shutdown complete", "name", cfg.ServiceName)
	return nil
}

func maskURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "[invalid-url]"
	}

	host := parsed.Host
	if len(host) > 30 {
		host = host[:30] + "..."
	}

	userPart := ""
	if parsed.User != nil {
		userPart = parsed.User.Username() + ":****@"
	}

	return fmt.Sprintf("%s://%s%s/...", parsed.Scheme, userPart, host)
}
