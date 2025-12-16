package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
	"github.com/specvital/collector/internal/handler/queue"
)

// 2 hours > 1h cron interval, prevents duplicate tasks from cron jitter.
const deduplicationWindow = 2 * time.Hour

type Client struct {
	client *asynq.Client
}

func NewClient(redisURL string) (*Client, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URI: %w", err)
	}

	return &Client{
		client: asynq.NewClient(opt),
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// Uses deduplication to prevent duplicate tasks within 2 hour window.
func (c *Client) EnqueueAnalysis(ctx context.Context, owner, repo string) error {
	payload := queue.AnalyzePayload{
		Owner: owner,
		Repo:  repo,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	task := asynq.NewTask(queue.TypeAnalyze, data)

	_, err = c.client.EnqueueContext(ctx, task,
		asynq.Unique(deduplicationWindow),
	)
	if err != nil {
		if errors.Is(err, asynq.ErrDuplicateTask) {
			slog.DebugContext(ctx, "duplicate task ignored by deduplication",
				"owner", owner,
				"repo", repo,
				"window", deduplicationWindow,
			)
			return nil
		}
		return fmt.Errorf("enqueue task: %w", err)
	}

	return nil
}
