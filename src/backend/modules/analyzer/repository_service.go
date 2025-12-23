package analyzer

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
)

type RepositoryService interface {
	GetRecentRepositories(ctx context.Context, limit int) ([]domain.RepositoryCard, error)
	GetRepositoryStats(ctx context.Context) (*domain.RepositoryStats, error)
	GetUpdateStatus(ctx context.Context, owner, repo string) (*domain.UpdateStatusResult, error)
}

type repositoryService struct {
	gitClient     client.GitClient
	logger        *logger.Logger
	repo          Repository
	tokenProvider TokenProvider
}

func NewRepositoryService(
	logger *logger.Logger,
	repo Repository,
	gitClient client.GitClient,
	tokenProvider TokenProvider,
) RepositoryService {
	return &repositoryService{
		gitClient:     gitClient,
		logger:        logger,
		repo:          repo,
		tokenProvider: tokenProvider,
	}
}

func (s *repositoryService) GetRecentRepositories(ctx context.Context, limit int) ([]domain.RepositoryCard, error) {
	repos, err := s.repo.GetRecentRepositories(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent repositories: %w", err)
	}

	cards := make([]domain.RepositoryCard, len(repos))
	for i, r := range repos {
		var analysis *domain.AnalysisSummary
		if r.AnalysisID != "" {
			change := 0
			prevAnalysis, err := s.repo.GetPreviousAnalysis(ctx, r.CodebaseID, r.AnalysisID)
			if err != nil {
				s.logger.Warn(ctx, "failed to get previous analysis for delta", "owner", r.Owner, "repo", r.Name, "error", err)
			} else if prevAnalysis != nil {
				change = r.TotalTests - prevAnalysis.TotalTests
			}

			analysis = &domain.AnalysisSummary{
				AnalyzedAt: r.AnalyzedAt,
				Change:     change,
				CommitSHA:  r.CommitSHA,
				TestCount:  r.TotalTests,
			}
		}

		cards[i] = domain.RepositoryCard{
			FullName:       fmt.Sprintf("%s/%s", r.Owner, r.Name),
			ID:             r.CodebaseID,
			IsBookmarked:   false,
			LatestAnalysis: analysis,
			Name:           r.Name,
			Owner:          r.Owner,
			UpdateStatus:   domain.UpdateStatusUnknown,
		}
	}

	return cards, nil
}

func (s *repositoryService) GetRepositoryStats(ctx context.Context) (*domain.RepositoryStats, error) {
	return s.repo.GetRepositoryStats(ctx)
}

func (s *repositoryService) GetUpdateStatus(ctx context.Context, owner, repo string) (*domain.UpdateStatusResult, error) {
	completed, err := s.repo.GetLatestCompletedAnalysis(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("get latest analysis: %w", err)
	}

	latestSHA, err := s.getLatestCommitWithAuth(ctx, owner, repo)
	if err != nil {
		s.logger.Warn(ctx, "failed to get latest commit, returning unknown status", "owner", owner, "repo", repo, "error", err)
		return &domain.UpdateStatusResult{
			AnalyzedCommitSHA: completed.CommitSHA,
			Status:            domain.UpdateStatusUnknown,
		}, nil
	}

	status := domain.UpdateStatusUpToDate
	if latestSHA != completed.CommitSHA {
		status = domain.UpdateStatusNewCommits
	}

	return &domain.UpdateStatusResult{
		AnalyzedCommitSHA: completed.CommitSHA,
		LatestCommitSHA:   latestSHA,
		Status:            status,
	}, nil
}

func (s *repositoryService) getLatestCommitWithAuth(ctx context.Context, owner, repo string) (string, error) {
	token, _ := s.getUserToken(ctx)
	if token != "" {
		sha, err := s.gitClient.GetLatestCommitSHAWithToken(ctx, owner, repo, token)
		if err == nil {
			return sha, nil
		}
		s.logger.Warn(ctx, "authenticated GitHub API call failed, falling back to public access", "error", err)
	}
	return s.gitClient.GetLatestCommitSHA(ctx, owner, repo)
}

func (s *repositoryService) getUserToken(ctx context.Context) (string, error) {
	if s.tokenProvider == nil {
		return "", nil
	}

	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return "", nil
	}

	return s.tokenProvider.GetUserGitHubToken(ctx, userID)
}
