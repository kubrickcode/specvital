package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

const (
	jobKindAnalysis = "analysis:analyze"
)

var _ port.ActiveTaskRepository = (*ActiveTaskRepository)(nil)

type ActiveTaskRepository struct {
	queries *db.Queries
}

func NewActiveTaskRepository(queries *db.Queries) *ActiveTaskRepository {
	return &ActiveTaskRepository{queries: queries}
}

type analyzeArgs struct {
	CommitSHA string  `json:"commit_sha"`
	Owner     string  `json:"owner"`
	Repo      string  `json:"repo"`
	UserID    *string `json:"user_id,omitempty"`
}

func (r *ActiveTaskRepository) GetUserActiveTasks(ctx context.Context, userID string) ([]entity.ActiveTask, error) {
	rows, err := r.queries.GetUserActiveJobs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user active jobs: %w", err)
	}

	tasks := make([]entity.ActiveTask, 0, len(rows))
	for _, row := range rows {
		task, err := r.mapRowToTask(row)
		if err != nil {
			slog.Warn("skipping unmappable job", "jobID", row.ID, "kind", row.Kind, "error", err)
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *ActiveTaskRepository) mapRowToTask(row db.GetUserActiveJobsRow) (entity.ActiveTask, error) {
	taskType, err := mapKindToTaskType(row.Kind)
	if err != nil {
		return entity.ActiveTask{}, err
	}

	var args analyzeArgs
	if err := json.Unmarshal(row.Args, &args); err != nil {
		return entity.ActiveTask{}, fmt.Errorf("unmarshal args: %w", err)
	}

	var startedAt *time.Time
	if row.AttemptedAt.Valid {
		startedAt = &row.AttemptedAt.Time
	}

	return entity.ActiveTask{
		CreatedAt: row.CreatedAt.Time,
		ID:        strconv.FormatInt(row.ID, 10),
		Owner:     args.Owner,
		Repo:      args.Repo,
		StartedAt: startedAt,
		Status:    entity.MapRiverStateToTaskStatus(row.State),
		Type:      taskType,
	}, nil
}

func mapKindToTaskType(kind string) (entity.TaskType, error) {
	switch kind {
	case jobKindAnalysis:
		return entity.TaskTypeAnalysis, nil
	default:
		return "", fmt.Errorf("unknown job kind: %s", kind)
	}
}
