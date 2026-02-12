package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kubrickcode/specvital/apps/web/backend/common/queue"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain/port"
	subscription "github.com/kubrickcode/specvital/apps/web/backend/modules/subscription/domain/entity"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
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
	AnalysisID     string  `json:"analysis_id" river:"unique"`
	Language       string  `json:"language" river:"unique"`
	ModelID        string  `json:"model_id"`
	UserID         *string `json:"user_id,omitempty"`
	GenerationMode string  `json:"generation_mode,omitempty"`
	// Backward compatibility: derived from GenerationMode for workers not yet updated.
	ForceRegenerate bool `json:"force_regenerate,omitempty"`
}

func (SpecGenerationArgs) Kind() string { return TypeSpecGeneration }

type RiverQueueService struct {
	client *river.Client[pgx.Tx]
}

func NewRiverQueueService(client *river.Client[pgx.Tx]) *RiverQueueService {
	return &RiverQueueService{client: client}
}

func (s *RiverQueueService) EnqueueSpecGeneration(ctx context.Context, analysisID string, language string, userID *string, tier subscription.PlanTier, mode entity.GenerationMode) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := SpecGenerationArgs{
		AnalysisID:      analysisID,
		Language:        language,
		ModelID:         "gemini-2.5-pro",
		UserID:          userID,
		GenerationMode:  string(mode),
		ForceRegenerate: mode.IsRegeneration(),
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

func (s *RiverQueueService) EnqueueSpecGenerationTx(ctx context.Context, tx pgx.Tx, analysisID string, language string, userID *string, tier subscription.PlanTier, mode entity.GenerationMode) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := SpecGenerationArgs{
		AnalysisID:      analysisID,
		Language:        language,
		ModelID:         "gemini-2.5-pro",
		UserID:          userID,
		GenerationMode:  string(mode),
		ForceRegenerate: mode.IsRegeneration(),
	}

	targetQueue := queue.SelectQueueForSpecView(tier, false)

	result, err := s.client.InsertTx(ctx, tx, args, &river.InsertOpts{
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
		return 0, fmt.Errorf("enqueue spec generation for analysis %s: %w", analysisID, err)
	}

	return result.Job.ID, nil
}
