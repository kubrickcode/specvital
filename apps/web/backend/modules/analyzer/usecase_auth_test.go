package analyzer

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/kubrickcode/specvital/apps/web/backend/internal/client"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/port"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/usecase"
	authdomain "github.com/kubrickcode/specvital/apps/web/backend/modules/auth/domain"
)

func TestAnalyzeRepositoryUseCaseWithAuth(t *testing.T) {
	t.Run("uses token for authenticated user", func(t *testing.T) {
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHAToken: "auth-sha"}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)

		ctx := context.Background()
		result, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner:  "owner",
			Repo:   "repo",
			UserID: "user-123",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Progress == nil {
			t.Fatal("expected progress result")
		}

		if !queue.enqueueCalled {
			t.Error("expected queue to be called")
		}
	})

	t.Run("falls back to public access when no token", func(t *testing.T) {
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHA: "public-sha"}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}
		tokenProvider := &mockTokenProvider{err: authdomain.ErrNoGitHubToken}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)

		ctx := context.Background()
		result, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner: "owner",
			Repo:  "repo",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Progress == nil {
			t.Fatal("expected progress result")
		}

		if !queue.enqueueCalled {
			t.Error("expected queue to be called")
		}
	})

	t.Run("falls back to public access when authenticated API fails with non-auth error", func(t *testing.T) {
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{
			commitSHA: "public-sha",
			errToken:  errors.New("temporary API error"),
		}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)

		ctx := context.Background()
		result, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner:  "owner",
			Repo:   "repo",
			UserID: "user-123",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Progress == nil {
			t.Fatal("expected progress result")
		}
	})

	t.Run("returns error when authenticated API returns forbidden", func(t *testing.T) {
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{
			errToken: client.ErrForbidden,
			err:      client.ErrForbidden,
		}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)

		ctx := context.Background()
		_, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner:  "owner",
			Repo:   "repo",
			UserID: "user-123",
		})

		if !errors.Is(err, client.ErrForbidden) {
			t.Errorf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("returns error when authenticated API returns not found", func(t *testing.T) {
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{
			errToken: client.ErrRepoNotFound,
			err:      client.ErrRepoNotFound,
		}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)

		ctx := context.Background()
		_, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner:  "owner",
			Repo:   "repo",
			UserID: "user-123",
		})

		if !errors.Is(err, client.ErrRepoNotFound) {
			t.Errorf("expected ErrRepoNotFound, got %v", err)
		}
	})

	t.Run("returns cached result with authenticated SHA", func(t *testing.T) {
		parserVersion := "v1.0.0"
		completedAnalysis := &port.CompletedAnalysis{
			ID:            "analysis-123",
			Owner:         "owner",
			Repo:          "repo",
			CommitSHA:     "auth-sha",
			ParserVersion: &parserVersion,
			TotalSuites:   5,
			TotalTests:    10,
		}
		repo := &mockRepository{completedAnalysis: completedAnalysis}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHAToken: "auth-sha"}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, tokenProvider, nil, nil)

		ctx := context.Background()
		result, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner:  "owner",
			Repo:   "repo",
			UserID: "user-123",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Analysis == nil {
			t.Fatal("expected analysis result")
		}

		if result.Analysis.CommitSHA != "auth-sha" {
			t.Errorf("expected commitSHA auth-sha, got %s", result.Analysis.CommitSHA)
		}

		if queue.enqueueCalled {
			t.Error("expected queue not to be called for cached result")
		}
	})

	t.Run("handles nil token provider gracefully", func(t *testing.T) {
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHA: "public-sha"}
		systemConfig := &mockSystemConfigReader{parserVersion: "v1.0.0"}

		uc := usecase.NewAnalyzeRepositoryUseCase(gitClient, queue, repo, systemConfig, nil, nil, nil)

		ctx := context.Background()
		result, err := uc.Execute(ctx, usecase.AnalyzeRepositoryInput{
			Owner: "owner",
			Repo:  "repo",
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Progress == nil {
			t.Fatal("expected progress result")
		}
	})
}
