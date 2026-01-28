package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/worker/internal/adapter/repository/postgres"
	"github.com/specvital/worker/internal/infra/buildinfo"
	infraqueue "github.com/specvital/worker/internal/infra/queue"
)

// maskURL returns a sanitized URL for logging (hides credentials).
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

// registerParserVersion writes the given core module version to system_config.
// This allows tracking which parser version produced each analysis.
func registerParserVersion(ctx context.Context, pool *pgxpool.Pool, version string) error {
	if version == "unknown" {
		slog.Warn("parser version unknown, skipping registration")
		return nil
	}

	repo := postgres.NewSystemConfigRepository(pool)
	if err := repo.Upsert(ctx, postgres.ConfigKeyParserVersion, version); err != nil {
		return fmt.Errorf("upsert parser version: %w", err)
	}

	slog.Info("parser version registered",
		"version", version,
		"display", buildinfo.FormatVersionDisplay(version))
	return nil
}

// logQueueSubscription logs the queues that a service is subscribing to.
func logQueueSubscription(service string, queues []infraqueue.QueueAllocation) {
	var queueInfo []string
	totalWorkers := 0
	for _, q := range queues {
		queueInfo = append(queueInfo, fmt.Sprintf("%s(%d)", q.Name, q.MaxWorkers))
		totalWorkers += q.MaxWorkers
	}

	slog.Info("subscribing to queues",
		"service", service,
		"queues", strings.Join(queueInfo, ", "),
		"total_workers", totalWorkers,
	)
}
