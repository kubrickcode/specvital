package main

import (
	"log/slog"
	"os"

	"github.com/specvital/worker/internal/app/bootstrap"
	"github.com/specvital/worker/internal/infra/config"

	_ "github.com/specvital/core/pkg/parser/strategies/all"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := bootstrap.StartWorker(bootstrap.WorkerConfig{
		ServiceName:   "analyzer",
		DatabaseURL:   cfg.DatabaseURL,
		EncryptionKey: cfg.EncryptionKey,
	}); err != nil {
		slog.Error("analyzer failed", "error", err)
		os.Exit(1)
	}
}
