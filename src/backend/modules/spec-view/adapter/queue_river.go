package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

var _ port.QueueService = (*RiverQueueService)(nil)

const (
	TypeSpecGeneration = "specview:generate"

	queueName      = "specview"
	maxRetries     = 3
	enqueueTimeout = 5 * time.Second
)

// SpecGenerationArgs represents the arguments for spec generation job.
// This must match the worker's Args structure.
type SpecGenerationArgs struct {
	AnalysisID string `json:"analysis_id" river:"unique"`
	Language   string `json:"language"`
	ModelID    string `json:"model_id"`
}

func (SpecGenerationArgs) Kind() string { return TypeSpecGeneration }

type RiverQueueService struct {
	client *river.Client[pgx.Tx]
}

func NewRiverQueueService(client *river.Client[pgx.Tx]) *RiverQueueService {
	return &RiverQueueService{client: client}
}

func (s *RiverQueueService) EnqueueSpecGeneration(ctx context.Context, analysisID string, language string) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := SpecGenerationArgs{
		AnalysisID: analysisID,
		Language:   language,
		ModelID:    "gemini-2.5-pro",
	}

	_, err := s.client.Insert(ctx, args, &river.InsertOpts{
		MaxAttempts: maxRetries,
		Queue:       queueName,
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
