package analyzer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
)

const dbTimeout = 5 * time.Second

type AnalyzeResult struct {
	Analysis *domain.Analysis
	Progress *domain.AnalysisProgress
}

type AnalyzerService interface {
	AnalyzeRepository(ctx context.Context, owner, repo string) (*AnalyzeResult, error)
	GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalyzeResult, error)
}

type analyzerService struct {
	queue     QueueService
	repo      Repository
	gitClient client.GitClient
}

func NewAnalyzerService(repo Repository, queue QueueService, gitClient client.GitClient) AnalyzerService {
	return &analyzerService{
		queue:     queue,
		repo:      repo,
		gitClient: gitClient,
	}
}

func (s *analyzerService) AnalyzeRepository(ctx context.Context, owner, repo string) (*AnalyzeResult, error) {
	latestSHA, err := s.gitClient.GetLatestCommitSHA(ctx, owner, repo)
	if err != nil {
		slog.Error("failed to get latest commit", "owner", owner, "repo", repo, "error", err)
		return nil, fmt.Errorf("get latest commit: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	completed, err := s.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err == nil {
		if completed.CommitSHA == latestSHA {
			analysis, buildErr := s.buildAnalysisFromCompleted(ctx, completed)
			if buildErr != nil {
				slog.Error("failed to build analysis", "owner", owner, "repo", repo, "error", buildErr)
				return nil, fmt.Errorf("build analysis: %w", buildErr)
			}
			slog.Info("cache hit", "owner", owner, "repo", repo, "commitSHA", latestSHA)
			return &AnalyzeResult{Analysis: analysis}, nil
		}
		slog.Info("cache miss - commit changed",
			"owner", owner, "repo", repo,
			"cachedSHA", completed.CommitSHA,
			"latestSHA", latestSHA)
	}

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		slog.Error("failed to get analysis", "owner", owner, "repo", repo, "error", err)
		return nil, fmt.Errorf("get analysis: %w", err)
	}

	status, err := s.repo.GetAnalysisStatus(ctx, owner, repo)
	if err == nil {
		progress := s.buildProgressFromStatus(status)
		return &AnalyzeResult{Progress: progress}, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		slog.Error("failed to get status", "owner", owner, "repo", repo, "error", err)
		return nil, fmt.Errorf("get status: %w", err)
	}

	analysisID, err := s.repo.CreatePendingAnalysis(ctx, owner, repo, latestSHA)
	if err != nil {
		slog.Error("failed to create analysis", "owner", owner, "repo", repo, "error", err)
		return nil, fmt.Errorf("create analysis: %w", err)
	}

	if err := s.queue.Enqueue(ctx, analysisID, owner, repo); err != nil {
		slog.Error("failed to enqueue", "owner", owner, "repo", repo, "analysisId", analysisID, "error", err)
		if cleanupErr := s.repo.MarkAnalysisFailed(ctx, analysisID, "queue registration failed"); cleanupErr != nil {
			slog.Error("failed to cleanup after enqueue error", "analysisId", analysisID, "error", cleanupErr)
		}
		return nil, fmt.Errorf("queue analysis: %w", err)
	}

	now := time.Now()
	progress := &domain.AnalysisProgress{
		CreatedAt: now,
		ID:        analysisID,
		Status:    domain.StatusPending,
	}
	return &AnalyzeResult{Progress: progress}, nil
}

func (s *analyzerService) GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalyzeResult, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	completed, err := s.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err == nil {
		analysis, buildErr := s.buildAnalysisFromCompleted(ctx, completed)
		if buildErr != nil {
			slog.Error("failed to build analysis", "owner", owner, "repo", repo, "error", buildErr)
			return nil, fmt.Errorf("build analysis: %w", buildErr)
		}
		return &AnalyzeResult{Analysis: analysis}, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		slog.Error("failed to get analysis", "owner", owner, "repo", repo, "error", err)
		return nil, fmt.Errorf("get analysis: %w", err)
	}

	status, err := s.repo.GetAnalysisStatus(ctx, owner, repo)
	if err == nil {
		progress := s.buildProgressFromStatus(status)
		return &AnalyzeResult{Progress: progress}, nil
	}

	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrNotFound
	}

	slog.Error("failed to get status", "owner", owner, "repo", repo, "error", err)
	return nil, fmt.Errorf("get status: %w", err)
}

func (s *analyzerService) buildAnalysisFromCompleted(ctx context.Context, completed *CompletedAnalysis) (*domain.Analysis, error) {
	suitesWithCases, err := s.repo.GetTestSuitesWithCases(ctx, completed.ID)
	if err != nil {
		return nil, fmt.Errorf("get test suites: %w", err)
	}

	suites := make([]domain.TestSuite, len(suitesWithCases))
	for i, suite := range suitesWithCases {
		testCases := make([]domain.TestCase, len(suite.Tests))
		for j, t := range suite.Tests {
			testCases[j] = domain.TestCase{
				Line:   t.Line,
				Name:   t.Name,
				Status: mapToDomainTestStatus(t.Status),
			}
		}

		suites[i] = domain.TestSuite{
			FilePath:  suite.FilePath,
			Framework: suite.Framework,
			ID:        suite.ID,
			TestCases: testCases,
		}
	}

	return &domain.Analysis{
		CommitSHA:   completed.CommitSHA,
		CompletedAt: completed.CompletedAt,
		ID:          completed.ID,
		Owner:       completed.Owner,
		Repo:        completed.Repo,
		TestSuites:  suites,
		TotalSuites: completed.TotalSuites,
		TotalTests:  completed.TotalTests,
	}, nil
}

func (s *analyzerService) buildProgressFromStatus(status *AnalysisStatus) *domain.AnalysisProgress {
	var domainStatus domain.Status
	switch status.Status {
	case "completed":
		domainStatus = domain.StatusCompleted
	case "running":
		domainStatus = domain.StatusRunning
	case "pending":
		domainStatus = domain.StatusPending
	case "failed":
		domainStatus = domain.StatusFailed
	default:
		slog.Warn("unknown analysis status, treating as failed", "status", status.Status)
		domainStatus = domain.StatusFailed
	}

	return &domain.AnalysisProgress{
		CompletedAt:  status.CompletedAt,
		CreatedAt:    status.CreatedAt,
		ErrorMessage: status.ErrorMessage,
		ID:           status.ID,
		Status:       domainStatus,
	}
}

func mapToDomainTestStatus(status string) domain.TestStatus {
	switch status {
	case "active":
		return domain.TestStatusActive
	case "focused":
		return domain.TestStatusFocused
	case "skipped":
		return domain.TestStatusSkipped
	case "todo":
		return domain.TestStatusTodo
	case "xfail":
		return domain.TestStatusXfail
	default:
		return domain.TestStatusActive
	}
}
