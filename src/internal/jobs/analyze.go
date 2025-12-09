package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

const (
	TypeAnalyze = "analysis:analyze"
)

type AnalyzePayload struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func (p AnalyzePayload) Validate() error {
	if p.Owner == "" {
		return errors.New("owner is required")
	}
	if p.Repo == "" {
		return errors.New("repo is required")
	}
	return nil
}

type AnalyzeHandler struct{}

func NewAnalyzeHandler() *AnalyzeHandler {
	return &AnalyzeHandler{}
}

func (h *AnalyzeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload AnalyzePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	slog.InfoContext(ctx, "processing analyze task",
		"owner", payload.Owner,
		"repo", payload.Repo,
	)

	// TODO(commit-4): Implement full pipeline
	// 1. source.NewGitSource()
	// 2. defer src.Cleanup()
	// 3. scanner.Scan()
	// 4. DB transaction save

	slog.InfoContext(ctx, "analyze task completed",
		"owner", payload.Owner,
		"repo", payload.Repo,
	)

	return nil
}
