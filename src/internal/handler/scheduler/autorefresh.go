package scheduler

import (
	"context"
	"log/slog"
	"time"

	infrascheduler "github.com/specvital/collector/internal/infra/scheduler"
	"github.com/specvital/collector/internal/usecase/autorefresh"
)

const defaultJobTimeout = 5 * time.Minute
const lockHeartbeatInterval = 3 * time.Minute // < 10 min lock TTL

type AutoRefreshHandler struct {
	lock    *infrascheduler.DistributedLock
	useCase *autorefresh.AutoRefreshUseCase
}

// Pass nil for lock to disable distributed locking (single-instance only).
func NewAutoRefreshHandler(
	useCase *autorefresh.AutoRefreshUseCase,
	lock *infrascheduler.DistributedLock,
) *AutoRefreshHandler {
	return &AutoRefreshHandler{
		lock:    lock,
		useCase: useCase,
	}
}

func (h *AutoRefreshHandler) Run() {
	h.RunWithContext(context.Background())
}

func (h *AutoRefreshHandler) RunWithContext(parentCtx context.Context) {
	ctx, cancel := context.WithTimeout(parentCtx, defaultJobTimeout)
	defer cancel()

	start := time.Now()

	if h.lock != nil {
		acquired, err := h.lock.TryAcquire(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "auto-refresh lock acquisition failed",
				"error", err,
			)
			return
		}
		if !acquired {
			slog.DebugContext(ctx, "auto-refresh skipped: another instance is running")
			return
		}

		heartbeatDone := make(chan struct{})
		go h.heartbeat(ctx, heartbeatDone)

		defer func() {
			close(heartbeatDone)
			if err := h.lock.Release(ctx); err != nil {
				slog.WarnContext(ctx, "auto-refresh lock release failed", "error", err)
			}
		}()
	}

	slog.InfoContext(ctx, "auto-refresh job started")

	if err := h.useCase.Execute(ctx); err != nil {
		slog.ErrorContext(ctx, "auto-refresh job failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return
	}

	slog.InfoContext(ctx, "auto-refresh job completed",
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func (h *AutoRefreshHandler) heartbeat(ctx context.Context, done <-chan struct{}) {
	ticker := time.NewTicker(lockHeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := h.lock.Extend(ctx); err != nil {
				slog.WarnContext(ctx, "auto-refresh lock extend failed", "error", err)
			} else {
				slog.DebugContext(ctx, "auto-refresh lock extended")
			}
		}
	}
}
