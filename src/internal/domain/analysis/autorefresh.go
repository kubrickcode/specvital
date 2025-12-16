package analysis

import (
	"context"
	"time"
)

type AutoRefreshRepository interface {
	GetCodebasesForAutoRefresh(ctx context.Context) ([]CodebaseRefreshInfo, error)
}

type CodebaseRefreshInfo struct {
	ConsecutiveFailures int
	Host                string
	ID                  UUID
	LastCompletedAt     *time.Time
	LastViewedAt        time.Time
	Name                string
	Owner               string
}

type TaskQueue interface {
	EnqueueAnalysis(ctx context.Context, owner, repo string) error
}
