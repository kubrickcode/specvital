package app

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/specvital/worker/internal/adapter/ai/gemini"
	specviewqueue "github.com/specvital/worker/internal/adapter/queue/specview"
	"github.com/specvital/worker/internal/adapter/repository/postgres"
	infraqueue "github.com/specvital/worker/internal/infra/queue"
	specviewuc "github.com/specvital/worker/internal/usecase/specview"
)

// SpecGeneratorContainer holds dependencies for the spec-generator worker service.
type SpecGeneratorContainer struct {
	GeminiProvider *gemini.Provider
	QueueClient    *infraqueue.Client
	SpecViewWorker *specviewqueue.Worker
	Workers        *river.Workers
}

// NewSpecGeneratorContainer creates and initializes a new spec-generator container with all required dependencies.
func NewSpecGeneratorContainer(ctx context.Context, cfg ContainerConfig) (*SpecGeneratorContainer, error) {
	if err := cfg.ValidateSpecGenerator(); err != nil {
		return nil, fmt.Errorf("invalid container config: %w", err)
	}

	geminiProvider, err := gemini.NewProvider(ctx, gemini.Config{
		APIKey:      cfg.GeminiAPIKey,
		Phase1Model: cfg.GeminiPhase1Model,
		Phase2Model: cfg.GeminiPhase2Model,
	})
	if err != nil {
		return nil, fmt.Errorf("create gemini provider: %w", err)
	}

	specDocRepo := postgres.NewSpecDocumentRepository(cfg.Pool)
	defaultModelID := cfg.GeminiPhase1Model
	if defaultModelID == "" {
		defaultModelID = "gemini-2.5-flash"
	}
	specViewUC := specviewuc.NewGenerateSpecViewUseCase(
		specDocRepo,
		geminiProvider,
		defaultModelID,
	)
	specViewWorker := specviewqueue.NewWorker(specViewUC)

	workers := river.NewWorkers()
	river.AddWorker(workers, specViewWorker)

	queueClient, err := infraqueue.NewClient(ctx, cfg.Pool)
	if err != nil {
		return nil, fmt.Errorf("create queue client: %w", err)
	}

	return &SpecGeneratorContainer{
		GeminiProvider: geminiProvider,
		QueueClient:    queueClient,
		SpecViewWorker: specViewWorker,
		Workers:        workers,
	}, nil
}

// Close releases container resources.
func (c *SpecGeneratorContainer) Close() error {
	var errs []error

	if c.QueueClient != nil {
		if err := c.QueueClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close queue client: %w", err))
		}
	}

	if c.GeminiProvider != nil {
		if err := c.GeminiProvider.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close gemini provider: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close spec-generator container: %v", errs)
	}
	return nil
}
