package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

const (
	TypeAnalyze = "analysis:analyze"

	queueName      = "default"
	maxRetries     = 3
	enqueueTimeout = 5 * time.Second
)

type AnalyzeArgs struct {
	AnalysisID string  `json:"analysis_id"`
	Owner      string  `json:"owner"`
	Repo       string  `json:"repo"`
	CommitSHA  string  `json:"commit_sha"`
	UserID     *string `json:"user_id,omitempty"`
}

func (AnalyzeArgs) Kind() string { return TypeAnalyze }

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

type riverQueueService struct {
	client *river.Client[pgx.Tx]
	repo   Repository
}

func NewQueueService(client *river.Client[pgx.Tx], repo Repository) QueueService {
	return &riverQueueService{client: client, repo: repo}
}

func (s *riverQueueService) Enqueue(ctx context.Context, analysisID, owner, repo, commitSHA string, userID *string) error {
	ctx, cancel := context.WithTimeout(ctx, enqueueTimeout)
	defer cancel()

	args := AnalyzeArgs{
		AnalysisID: analysisID,
		Owner:      owner,
		Repo:       repo,
		CommitSHA:  commitSHA,
		UserID:     userID,
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
		return fmt.Errorf("enqueue task for %s/%s: %w", owner, repo, err)
	}

	return nil
}

func (s *riverQueueService) GetTaskInfo(ctx context.Context, analysisID string) (*TaskInfo, error) {
	info, err := s.repo.GetRiverJobByAnalysisID(ctx, TypeAnalyze, analysisID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return &TaskInfo{
		AnalysisID: info.AnalysisID,
		State:      info.State,
	}, nil
}

func (s *riverQueueService) FindTaskByRepo(ctx context.Context, owner, repo string) (*TaskInfo, error) {
	info, err := s.repo.FindActiveRiverJobByRepo(ctx, TypeAnalyze, owner, repo)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return &TaskInfo{
		AnalysisID: info.AnalysisID,
		State:      info.State,
	}, nil
}

func (s *riverQueueService) Close() error {
	return nil
}
