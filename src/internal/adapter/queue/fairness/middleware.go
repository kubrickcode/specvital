package fairness

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// FairnessMiddleware enforces per-user concurrent job limits using River middleware.
// Jobs exceeding their tier's limit are snoozed with jitter to prevent thundering herd.
type FairnessMiddleware struct {
	river.MiddlewareDefaults
	limiter   *PerUserLimiter
	extractor UserJobExtractor
	config    *Config
}

// NewFairnessMiddleware creates a new fairness middleware with the given limiter and extractor.
func NewFairnessMiddleware(limiter *PerUserLimiter, extractor UserJobExtractor, config *Config) *FairnessMiddleware {
	return &FairnessMiddleware{
		limiter:   limiter,
		extractor: extractor,
		config:    config,
	}
}

// Work implements river.WorkerMiddleware by enforcing per-user concurrent limits.
// System jobs (empty userID) bypass limits. Jobs exceeding limits are snoozed with jitter.
func (m *FairnessMiddleware) Work(
	ctx context.Context,
	job *rivertype.JobRow,
	doInner func(ctx context.Context) error,
) error {
	userID := m.extractor.ExtractUserID(job.EncodedArgs)
	if userID == "" {
		return doInner(ctx)
	}

	tier := m.extractor.ExtractTier(job.EncodedArgs)

	if !m.limiter.TryAcquire(userID, tier, job.ID) {
		jitterNanos := rand.Int64N(int64(m.config.SnoozeJitter))
		snoozeDuration := m.config.SnoozeDuration + time.Duration(jitterNanos)

		slog.InfoContext(ctx, "user at concurrency limit, snoozing",
			"user_id", userID,
			"job_id", job.ID,
			"tier", tier,
			"active_jobs", m.limiter.ActiveCount(userID),
			"snooze_duration", snoozeDuration,
		)
		return river.JobSnooze(snoozeDuration)
	}

	defer m.limiter.Release(userID, job.ID)
	return doInner(ctx)
}
