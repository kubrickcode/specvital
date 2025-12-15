package analyzer

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/client"
	authdomain "github.com/specvital/web/src/backend/modules/auth/domain"
)

func TestAnalyzeRepositoryWithAuth(t *testing.T) {
	t.Run("uses token for authenticated user", func(t *testing.T) {
		// Given: authenticated user with GitHub token
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHAToken: "auth-sha"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: analyzing repository
		result, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: uses authenticated method and returns result
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
		// Given: unauthenticated user
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHA: "public-sha"}
		tokenProvider := &mockTokenProvider{err: authdomain.ErrNoGitHubToken}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)

		ctx := context.Background()

		// When: analyzing repository
		result, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: falls back to public method
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
		// Given: authenticated user but API fails temporarily
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{
			commitSHA: "public-sha",
			errToken:  errors.New("temporary API error"),
		}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: analyzing repository
		result, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: falls back to public method
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Progress == nil {
			t.Fatal("expected progress result")
		}
	})

	t.Run("returns error when authenticated API returns forbidden", func(t *testing.T) {
		// Given: authenticated user accessing private repo
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{
			errToken: client.ErrForbidden,
		}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: analyzing repository
		_, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: returns forbidden error without fallback
		if !errors.Is(err, client.ErrForbidden) {
			t.Errorf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("returns error when authenticated API returns not found", func(t *testing.T) {
		// Given: authenticated user accessing non-existent repo
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{
			errToken: client.ErrRepoNotFound,
		}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: analyzing repository
		_, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: returns not found error without fallback
		if !errors.Is(err, client.ErrRepoNotFound) {
			t.Errorf("expected ErrRepoNotFound, got %v", err)
		}
	})

	t.Run("returns cached result with authenticated SHA", func(t *testing.T) {
		// Given: authenticated user with cached analysis
		completedAnalysis := &CompletedAnalysis{
			ID:          "analysis-123",
			Owner:       "owner",
			Repo:        "repo",
			CommitSHA:   "auth-sha",
			TotalSuites: 5,
			TotalTests:  10,
		}
		repo := &mockRepository{completedAnalysis: completedAnalysis}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHAToken: "auth-sha"}
		tokenProvider := &mockTokenProvider{token: "github-token"}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, tokenProvider)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: analyzing repository
		result, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: returns cached result
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
		// Given: service without token provider
		repo := &mockRepository{}
		queue := &mockQueueService{}
		gitClient := &mockGitClient{commitSHA: "public-sha"}

		log := logger.New()
		service := NewAnalyzerService(log, repo, queue, gitClient, nil)

		ctx := context.Background()

		// When: analyzing repository
		result, err := service.AnalyzeRepository(ctx, "owner", "repo")

		// Then: falls back to public method
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result == nil || result.Progress == nil {
			t.Fatal("expected progress result")
		}
	})
}

func TestGetUserToken(t *testing.T) {
	t.Run("returns token for authenticated user", func(t *testing.T) {
		// Given: authenticated user with token
		tokenProvider := &mockTokenProvider{token: "github-token"}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, &mockGitClient{}, tokenProvider).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting user token
		token, err := service.getUserToken(ctx)

		// Then: returns token
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if token != "github-token" {
			t.Errorf("expected token github-token, got %s", token)
		}
	})

	t.Run("returns error when user has no token", func(t *testing.T) {
		// Given: authenticated user without GitHub token
		tokenProvider := &mockTokenProvider{err: authdomain.ErrNoGitHubToken}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, &mockGitClient{}, tokenProvider).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting user token
		_, err := service.getUserToken(ctx)

		// Then: returns error
		if !errors.Is(err, authdomain.ErrNoGitHubToken) {
			t.Errorf("expected ErrNoGitHubToken, got %v", err)
		}
	})

	t.Run("returns error when no user in context", func(t *testing.T) {
		// Given: unauthenticated request
		tokenProvider := &mockTokenProvider{token: "github-token"}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, &mockGitClient{}, tokenProvider).(*analyzerService)

		ctx := context.Background()

		// When: getting user token
		_, err := service.getUserToken(ctx)

		// Then: returns error
		if !errors.Is(err, authdomain.ErrNoGitHubToken) {
			t.Errorf("expected ErrNoGitHubToken, got %v", err)
		}
	})

	t.Run("returns error when token provider is nil", func(t *testing.T) {
		// Given: service without token provider
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, &mockGitClient{}, nil).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting user token
		_, err := service.getUserToken(ctx)

		// Then: returns error
		if !errors.Is(err, authdomain.ErrNoGitHubToken) {
			t.Errorf("expected ErrNoGitHubToken, got %v", err)
		}
	})
}

func TestGetLatestCommitWithAuth(t *testing.T) {
	t.Run("uses authenticated method when token available", func(t *testing.T) {
		// Given: authenticated user
		tokenProvider := &mockTokenProvider{token: "github-token"}
		gitClient := &mockGitClient{commitSHAToken: "auth-sha"}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, gitClient, tokenProvider).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting latest commit
		sha, err := service.getLatestCommitWithAuth(ctx, "owner", "repo")

		// Then: uses authenticated method
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if sha != "auth-sha" {
			t.Errorf("expected SHA auth-sha, got %s", sha)
		}
	})

	t.Run("falls back to public method when token not available", func(t *testing.T) {
		// Given: unauthenticated user
		tokenProvider := &mockTokenProvider{err: authdomain.ErrNoGitHubToken}
		gitClient := &mockGitClient{commitSHA: "public-sha"}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, gitClient, tokenProvider).(*analyzerService)

		ctx := context.Background()

		// When: getting latest commit
		sha, err := service.getLatestCommitWithAuth(ctx, "owner", "repo")

		// Then: uses public method
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if sha != "public-sha" {
			t.Errorf("expected SHA public-sha, got %s", sha)
		}
	})

	t.Run("falls back to public method on auth API error", func(t *testing.T) {
		// Given: authenticated user but API error
		tokenProvider := &mockTokenProvider{token: "github-token"}
		gitClient := &mockGitClient{
			commitSHA: "public-sha",
			errToken:  errors.New("API error"),
		}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, gitClient, tokenProvider).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting latest commit
		sha, err := service.getLatestCommitWithAuth(ctx, "owner", "repo")

		// Then: falls back to public method
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if sha != "public-sha" {
			t.Errorf("expected SHA public-sha, got %s", sha)
		}
	})

	t.Run("does not fall back on forbidden error", func(t *testing.T) {
		// Given: authenticated user accessing forbidden repo
		tokenProvider := &mockTokenProvider{token: "github-token"}
		gitClient := &mockGitClient{errToken: client.ErrForbidden}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, gitClient, tokenProvider).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting latest commit
		_, err := service.getLatestCommitWithAuth(ctx, "owner", "repo")

		// Then: returns forbidden error
		if !errors.Is(err, client.ErrForbidden) {
			t.Errorf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("does not fall back on not found error", func(t *testing.T) {
		// Given: authenticated user accessing non-existent repo
		tokenProvider := &mockTokenProvider{token: "github-token"}
		gitClient := &mockGitClient{errToken: client.ErrRepoNotFound}
		log := logger.New()
		service := NewAnalyzerService(log, &mockRepository{}, &mockQueueService{}, gitClient, tokenProvider).(*analyzerService)

		ctx := context.Background()
		claims := &authdomain.Claims{}
		claims.Subject = "user-123"
		ctx = middleware.WithClaims(ctx, claims)

		// When: getting latest commit
		_, err := service.getLatestCommitWithAuth(ctx, "owner", "repo")

		// Then: returns not found error
		if !errors.Is(err, client.ErrRepoNotFound) {
			t.Errorf("expected ErrRepoNotFound, got %v", err)
		}
	})
}
