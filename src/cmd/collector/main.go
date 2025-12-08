package main

import (
	"log/slog"
	"os"

	"github.com/specvital/collector/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting collector")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("config loaded",
		"database_url", maskURL(cfg.DatabaseURL),
		"redis_url", maskURL(cfg.RedisURL),
	)

	slog.Info("collector initialized")
}

func maskURL(url string) string {
	if len(url) > 20 {
		return url[:20] + "..."
	}
	return url
}
