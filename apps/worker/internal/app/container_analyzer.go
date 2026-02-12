package app

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/kubrickcode/specvital/lib/crypto"
	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/parser"
	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/queue/analyze"
	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/queue/fairness"
	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/repository/postgres"
	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/vcs"
	"github.com/kubrickcode/specvital/apps/worker/internal/infra/db"
	infraqueue "github.com/kubrickcode/specvital/apps/worker/internal/infra/queue"
	analysisuc "github.com/kubrickcode/specvital/apps/worker/internal/usecase/analysis"
)

// AnalyzerContainer holds dependencies for the analyzer worker service.
type AnalyzerContainer struct {
	AnalyzeWorker *analyze.AnalyzeWorker
	Middleware    []rivertype.WorkerMiddleware
	QueueClient   *infraqueue.Client
	Workers       *river.Workers
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
	quotaRepo := postgres.NewQuotaReservationRepository(cfg.Pool)
	userRepo := postgres.NewUserRepository(cfg.Pool, encryptor)
	gitVCS := vcs.NewGitVCS()
	githubAPIClient := vcs.NewGitHubAPIClient(nil)
	coreParser := parser.NewCoreParser()
	analyzeUC := analysisuc.NewAnalyzeUseCase(
		analysisRepo, codebaseRepo, gitVCS, githubAPIClient, coreParser, userRepo,
		analysisuc.WithParserVersion(cfg.ParserVersion),
		analysisuc.WithBatchSize(cfg.Streaming.BatchSize),
	)
	analyzeWorker := analyze.NewAnalyzeWorker(analyzeUC, quotaRepo)

	workers := river.NewWorkers()
	river.AddWorker(workers, analyzeWorker)

	queueClient, err := infraqueue.NewClient(ctx, cfg.Pool)
	if err != nil {
		return nil, fmt.Errorf("create queue client: %w", err)
	}

	var middleware []rivertype.WorkerMiddleware
	queries := db.New(cfg.Pool)
	tierResolver := fairness.NewDBTierResolver(queries)
	fm, err := NewFairnessMiddleware(cfg.Fairness, tierResolver)
	if err != nil {
		return nil, fmt.Errorf("create fairness middleware: %w", err)
	}
	if fm != nil {
		middleware = append(middleware, fm)
	}

	return &AnalyzerContainer{
		AnalyzeWorker: analyzeWorker,
		Middleware:    middleware,
		QueueClient:   queueClient,
		Workers:       workers,
	}, nil
}

// Close releases container resources.
func (c *AnalyzerContainer) Close() error {
	if c.QueueClient != nil {
		if err := c.QueueClient.Close(); err != nil {
			return fmt.Errorf("close queue client: %w", err)
		}
	}
	return nil
}
