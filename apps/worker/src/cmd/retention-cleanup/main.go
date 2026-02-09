package main

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/specvital/worker/internal/app/bootstrap"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := bootstrap.RetentionConfig{
		BatchSize:   getEnvInt("RETENTION_BATCH_SIZE", 0),
		BatchSleep:  getEnvDuration("RETENTION_BATCH_SLEEP", 0),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServiceName: "retention-cleanup",
		Timeout:     getEnvDuration("RETENTION_TIMEOUT", 0),
	}

	if _, err := bootstrap.RunRetentionCleanup(cfg); err != nil {
		slog.Error("retention cleanup failed", "error", err)
		os.Exit(1)
	}
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		slog.Warn("invalid integer env var, using default",
			"key", key,
			"value", val,
			"default", defaultValue,
			"error", err,
		)
		return defaultValue
	}
	return parsed
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		slog.Warn("invalid duration env var, using default",
			"key", key,
			"value", val,
			"default", defaultValue,
			"error", err,
		)
		return defaultValue
	}
	return parsed
}
