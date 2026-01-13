package app

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/specvital/core/pkg/crypto"
	"github.com/specvital/worker/internal/adapter/ai/gemini"
	"github.com/specvital/worker/internal/adapter/parser"
	"github.com/specvital/worker/internal/adapter/queue/analyze"
	specviewqueue "github.com/specvital/worker/internal/adapter/queue/specview"
	"github.com/specvital/worker/internal/adapter/repository/postgres"
	"github.com/specvital/worker/internal/adapter/vcs"
	infraqueue "github.com/specvital/worker/internal/infra/queue"
	analysisuc "github.com/specvital/worker/internal/usecase/analysis"
	specviewuc "github.com/specvital/worker/internal/usecase/specview"
)

// AnalyzerContainer holds dependencies for the analyzer worker service.
type AnalyzerContainer struct {
	AnalyzeWorker   *analyze.AnalyzeWorker
	GeminiProvider  *gemini.Provider
	QueueClient     *infraqueue.Client
	SpecViewWorker  *specviewqueue.Worker
	Workers         *river.Workers
}

// NewAnalyzerContainer creates and initializes a new analyzer container with all required dependencies.
func NewAnalyzerContainer(ctx context.Context, cfg ContainerConfig) (*AnalyzerContainer, error) {
	if err := cfg.ValidateAnalyzer(); err != nil {
		return nil, fmt.Errorf("invalid container config: %w", err)
	}

	encryptor, err := crypto.NewEncryptorFromBase64(cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("create encryptor: %w", err)
	}

	analysisRepo := postgres.NewAnalysisRepository(cfg.Pool)
	codebaseRepo := postgres.NewCodebaseRepository(cfg.Pool)
	userRepo := postgres.NewUserRepository(cfg.Pool, encryptor)
	gitVCS := vcs.NewGitVCS()
	githubAPIClient := vcs.NewGitHubAPIClient(nil)
	coreParser := parser.NewCoreParser()
	analyzeUC := analysisuc.NewAnalyzeUseCase(
		analysisRepo, codebaseRepo, gitVCS, githubAPIClient, coreParser, userRepo,
		analysisuc.WithParserVersion(cfg.ParserVersion),
	)
	analyzeWorker := analyze.NewAnalyzeWorker(analyzeUC)

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
	river.AddWorker(workers, analyzeWorker)
	river.AddWorker(workers, specViewWorker)

	queueClient, err := infraqueue.NewClient(ctx, cfg.Pool)
	if err != nil {
		return nil, fmt.Errorf("create queue client: %w", err)
	}

	return &AnalyzerContainer{
		AnalyzeWorker:  analyzeWorker,
		GeminiProvider: geminiProvider,
		QueueClient:    queueClient,
		SpecViewWorker: specViewWorker,
		Workers:        workers,
	}, nil
}

// Close releases container resources.
// Uses error accumulation pattern to ensure all resources are cleaned up.
// Resources are closed in reverse initialization order.
func (c *AnalyzerContainer) Close() error {
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
		return fmt.Errorf("close analyzer container: %v", errs)
	}
	return nil
}
