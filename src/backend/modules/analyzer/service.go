package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	authdomain "github.com/specvital/web/src/backend/modules/auth/domain"
)

const dbTimeout = 5 * time.Second

type TokenProvider interface {
	GetUserGitHubToken(ctx context.Context, userID string) (string, error)
}

type AnalyzeResult struct {
	Analysis *domain.Analysis
	Progress *domain.AnalysisProgress
}

type AnalyzerService interface {
	AnalyzeRepository(ctx context.Context, owner, repo string) (*AnalyzeResult, error)
	GetAnalysisStatus(ctx context.Context, owner, repo string) (*AnalyzeResult, error)
}

type analyzerService struct {
	gitClient     client.GitClient
	logger        *logger.Logger
	queue         QueueService
	repo          Repository
	tokenProvider TokenProvider
}

func NewAnalyzerService(logger *logger.Logger, repo Repository, queue QueueService, gitClient client.GitClient, tokenProvider TokenProvider) AnalyzerService {
	return &analyzerService{
		gitClient:     gitClient,
		logger:        logger,
		queue:         queue,
		repo:          repo,
		tokenProvider: tokenProvider,
	}
}

func (s *analyzerService) AnalyzeRepository(ctx context.Context, owner, repo string) (*AnalyzeResult, error) {
	log := s.logger.With("owner", owner, "repo", repo)

	latestSHA, err := s.getLatestCommitWithAuth(ctx, owner, repo)
	if err != nil {
		log.Error(ctx, "failed to get latest commit", "error", err)
		return nil, fmt.Errorf("get latest commit: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	completed, err := s.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err == nil {
		if completed.CommitSHA == latestSHA {
			analysis, buildErr := s.buildAnalysisFromCompleted(ctx, completed)
			if buildErr != nil {
				log.Error(ctx, "failed to build analysis", "error", buildErr)
				return nil, fmt.Errorf("build analysis: %w", buildErr)
			}
			log.Info(ctx, "cache hit", "commitSHA", latestSHA)
			return &AnalyzeResult{Analysis: analysis}, nil
		}
		log.Info(ctx, "cache miss - commit changed",
			"cachedSHA", completed.CommitSHA,
			"latestSHA", latestSHA)
	}

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Error(ctx, "failed to get analysis", "error", err)
		return nil, fmt.Errorf("get analysis: %w", err)
	}

	status, err := s.repo.GetAnalysisStatus(ctx, owner, repo)
	if err == nil {
		progress := s.buildProgressFromStatus(ctx, log, status)
		return &AnalyzeResult{Progress: progress}, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		log.Error(ctx, "failed to get status", "error", err)
		return nil, fmt.Errorf("get status: %w", err)
	}

	analysisID, err := s.repo.CreatePendingAnalysis(ctx, owner, repo, latestSHA)
	if err != nil {
		log.Error(ctx, "failed to create analysis", "error", err)
		return nil, fmt.Errorf("create analysis: %w", err)
	}

	if err := s.queue.Enqueue(ctx, analysisID, owner, repo); err != nil {
		log.Error(ctx, "failed to enqueue", "analysisId", analysisID, "error", err)
		if cleanupErr := s.repo.MarkAnalysisFailed(ctx, analysisID, "queue registration failed"); cleanupErr != nil {
			log.Error(ctx, "failed to cleanup after enqueue error", "analysisId", analysisID, "error", cleanupErr)
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
	log := s.logger.With("owner", owner, "repo", repo)

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	completed, err := s.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err == nil {
		analysis, buildErr := s.buildAnalysisFromCompleted(ctx, completed)
		if buildErr != nil {
			log.Error(ctx, "failed to build analysis", "error", buildErr)
			return nil, fmt.Errorf("build analysis: %w", buildErr)
		}
		return &AnalyzeResult{Analysis: analysis}, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		log.Error(ctx, "failed to get analysis", "error", err)
		return nil, fmt.Errorf("get analysis: %w", err)
	}

	status, err := s.repo.GetAnalysisStatus(ctx, owner, repo)
	if err == nil {
		progress := s.buildProgressFromStatus(ctx, log, status)
		return &AnalyzeResult{Progress: progress}, nil
	}

	if errors.Is(err, domain.ErrNotFound) {
		return nil, domain.ErrNotFound
	}

	log.Error(ctx, "failed to get status", "error", err)
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

func (s *analyzerService) buildProgressFromStatus(ctx context.Context, log *logger.Logger, status *AnalysisStatus) *domain.AnalysisProgress {
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
		log.Warn(ctx, "unknown analysis status, treating as failed", "status", status.Status)
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

func (s *analyzerService) getLatestCommitWithAuth(ctx context.Context, owner, repo string) (string, error) {
	token, err := s.getUserToken(ctx)
	if err != nil && !errors.Is(err, authdomain.ErrNoGitHubToken) {
		return "", fmt.Errorf("get user token: %w", err)
	}

	if token != "" {
		sha, err := s.gitClient.GetLatestCommitSHAWithToken(ctx, owner, repo, token)
		if err == nil {
			return sha, nil
		}

		if errors.Is(err, client.ErrForbidden) || errors.Is(err, client.ErrRepoNotFound) {
			return "", err
		}

		s.logger.Warn(ctx, "authenticated GitHub API call failed, falling back to public access",
			"error", err, "owner", owner, "repo", repo)
	}

	return s.gitClient.GetLatestCommitSHA(ctx, owner, repo)
}

func (s *analyzerService) getUserToken(ctx context.Context) (string, error) {
	if s.tokenProvider == nil {
		return "", authdomain.ErrNoGitHubToken
	}

	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return "", authdomain.ErrNoGitHubToken
	}

	return s.tokenProvider.GetUserGitHubToken(ctx, userID)
}
