package analyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeAnalyze = "analysis:analyze"

	queueName      = "default"
	maxRetries     = 3
	taskTimeout    = 10 * time.Minute
	uniqueDuration = 1 * time.Hour
	enqueueTimeout = 5 * time.Second
)

type AnalyzePayload struct {
	AnalysisID string  `json:"analysisId"`
	Owner      string  `json:"owner"`
	Repo       string  `json:"repo"`
	UserID     *string `json:"user_id"`
}

type QueueService interface {
	Enqueue(ctx context.Context, analysisID, owner, repo string, userID *string) error
	Close() error
}

type queueService struct {
	client *asynq.Client
}

func NewQueueService(client *asynq.Client) QueueService {
	return &queueService{client: client}
}

func (s *queueService) Enqueue(ctx context.Context, analysisID, owner, repo string, userID *string) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	payload, err := json.Marshal(AnalyzePayload{
		AnalysisID: analysisID,
		Owner:      owner,
		Repo:       repo,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	taskID := fmt.Sprintf("analyze:%s:%s", owner, repo)
	task := asynq.NewTask(TypeAnalyze, payload)

	_, err = s.client.EnqueueContext(ctx, task,
		asynq.TaskID(taskID),
		asynq.Unique(uniqueDuration),
		asynq.MaxRetry(maxRetries),
		asynq.Timeout(taskTimeout),
		asynq.Queue(queueName),
	)
	if err != nil {
		return fmt.Errorf("enqueue task for %s/%s: %w", owner, repo, err)
	}

	return nil
}

func (s *queueService) Close() error {
	return s.client.Close()
}
