package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/specvital/web/src/backend/common/queue"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

var _ port.QueueService = (*RiverQueueService)(nil)

const (
	TypeSpecGeneration = "specview:generate"

	maxRetries     = 3
	enqueueTimeout = 5 * time.Second
)

// SpecGenerationArgs represents the arguments for spec generation job.
// This must match the worker's Args structure (JSON fields only; river tags are producer-only).
// Unique key: (AnalysisID, Language) - allows different languages for same analysis,
// but prevents duplicate jobs for the same analysis+language combination.
type SpecGenerationArgs struct {
	AnalysisID      string  `json:"analysis_id" river:"unique"`
	Language        string  `json:"language" river:"unique"`
	ModelID         string  `json:"model_id"`
	UserID          *string `json:"user_id,omitempty"`
	ForceRegenerate bool    `json:"force_regenerate,omitempty"`
}

func (SpecGenerationArgs) Kind() string { return TypeSpecGeneration }

type RiverQueueService struct {
	client *river.Client[pgx.Tx]
}

func NewRiverQueueService(client *river.Client[pgx.Tx]) *RiverQueueService {
	return &RiverQueueService{client: client}
}

func (s *RiverQueueService) EnqueueSpecGeneration(ctx context.Context, analysisID string, language string, userID *string, tier subscription.PlanTier, forceRegenerate bool) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := SpecGenerationArgs{
		AnalysisID:      analysisID,
		Language:        language,
		ModelID:         "gemini-2.5-pro",
		UserID:          userID,
		ForceRegenerate: forceRegenerate,
	}

	targetQueue := queue.SelectQueueForSpecView(tier, false)

	_, err := s.client.Insert(ctx, args, &river.InsertOpts{
		MaxAttempts: maxRetries,
		Queue:       targetQueue,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
			ByState: []rivertype.JobState{
				rivertype.JobStateAvailable,
				rivertype.JobStatePending,
				rivertype.JobStateRunning,
				rivertype.JobStateRetryable,
				rivertype.JobStateScheduled,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("enqueue spec generation for analysis %s: %w", analysisID, err)
	}

	return nil
}
