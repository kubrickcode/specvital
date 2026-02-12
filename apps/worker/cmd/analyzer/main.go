package main

import (
	"log/slog"
	"os"

	"github.com/kubrickcode/specvital/apps/worker/internal/app/bootstrap"
	"github.com/kubrickcode/specvital/apps/worker/internal/infra/config"

	_ "github.com/kubrickcode/specvital/lib/parser/strategies/all"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := bootstrap.StartAnalyzer(bootstrap.AnalyzerConfig{
		ServiceName:   "analyzer",
		DatabaseURL:   cfg.DatabaseURL,
		EncryptionKey: cfg.EncryptionKey,
		Fairness:      cfg.Fairness,
		QueueWorkers:  cfg.Queue.Analyzer,
		Streaming:     cfg.Streaming,
	}); err != nil {
		slog.Error("analyzer failed", "error", err)
		os.Exit(1)
	}
}
