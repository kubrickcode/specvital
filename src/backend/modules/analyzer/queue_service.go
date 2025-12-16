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
	enqueueTimeout = 5 * time.Second
)

type AnalyzePayload struct {
	AnalysisID string  `json:"analysisId"`
	Owner      string  `json:"owner"`
	Repo       string  `json:"repo"`
	CommitSHA  string  `json:"commitSha"`
	UserID     *string `json:"user_id"`
}

type QueueService interface {
	Enqueue(ctx context.Context, analysisID, owner, repo, commitSHA string, userID *string) error
	GetTaskInfo(ctx context.Context, analysisID string) (*TaskInfo, error)
	FindTaskByRepo(ctx context.Context, owner, repo string) (*TaskInfo, error)
	Close() error
}

type TaskInfo struct {
	AnalysisID string
	State      string
}

type queueService struct {
	client    *asynq.Client
	inspector *asynq.Inspector
}

func NewQueueService(client *asynq.Client, inspector *asynq.Inspector) QueueService {
	return &queueService{client: client, inspector: inspector}
}

func (s *queueService) Enqueue(ctx context.Context, analysisID, owner, repo, commitSHA string, userID *string) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	payload, err := json.Marshal(AnalyzePayload{
		AnalysisID: analysisID,
		Owner:      owner,
		Repo:       repo,
		CommitSHA:  commitSHA,
		UserID:     userID,
	})
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	taskID := fmt.Sprintf("analyze:%s", analysisID)
	task := asynq.NewTask(TypeAnalyze, payload)

	_, err = s.client.EnqueueContext(ctx, task,
		asynq.TaskID(taskID),
		asynq.MaxRetry(maxRetries),
		asynq.Timeout(taskTimeout),
		asynq.Queue(queueName),
	)
	if err != nil {
		return fmt.Errorf("enqueue task for %s/%s: %w", owner, repo, err)
	}

	return nil
}

func (s *queueService) GetTaskInfo(ctx context.Context, analysisID string) (*TaskInfo, error) {
	taskID := fmt.Sprintf("analyze:%s", analysisID)

	info, err := s.inspector.GetTaskInfo(queueName, taskID)
	if err != nil {
		return nil, err
	}

	return &TaskInfo{
		AnalysisID: analysisID,
		State:      info.State.String(),
	}, nil
}

func (s *queueService) FindTaskByRepo(ctx context.Context, owner, repo string) (*TaskInfo, error) {
	info, err := s.findInTasks(s.inspector.ListPendingTasks, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("search pending tasks: %w", err)
	}
	if info != nil {
		return info, nil
	}

	info, err = s.findInTasks(s.inspector.ListActiveTasks, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("search active tasks: %w", err)
	}
	if info != nil {
		return info, nil
	}

	info, err = s.findInTasks(s.inspector.ListRetryTasks, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("search retry tasks: %w", err)
	}
	if info != nil {
		return info, nil
	}

	return nil, nil
}

type taskLister func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)

func (s *queueService) findInTasks(lister taskLister, owner, repo string) (*TaskInfo, error) {
	tasks, err := lister(queueName)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		var payload AnalyzePayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			continue
		}
		if payload.Owner == owner && payload.Repo == repo {
			return &TaskInfo{
				AnalysisID: payload.AnalysisID,
				State:      task.State.String(),
			}, nil
		}
	}

	return nil, nil
}

func (s *queueService) Close() error {
	var errs []error
	if err := s.client.Close(); err != nil {
		errs = append(errs, fmt.Errorf("client: %w", err))
	}
	if err := s.inspector.Close(); err != nil {
		errs = append(errs, fmt.Errorf("inspector: %w", err))
	}
	if len(errs) > 0 {
		return fmt.Errorf("close queue service: %v", errs)
	}
	return nil
}
