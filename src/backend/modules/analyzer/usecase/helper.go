package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"

	authdomain "github.com/specvital/web/src/backend/modules/auth/domain"
)

func getLatestCommitWithAuth(
	ctx context.Context,
	gitClient port.GitClient,
	tokenProvider port.TokenProvider,
	owner, repo, userID string,
) (string, error) {
	token, err := getUserToken(ctx, tokenProvider, userID)
	if err != nil && !errors.Is(err, authdomain.ErrNoGitHubToken) && !errors.Is(err, authdomain.ErrUserNotFound) {
		return "", fmt.Errorf("get user token: %w", err)
	}

	if token != "" {
		sha, err := gitClient.GetLatestCommitSHAWithToken(ctx, owner, repo, token)
		if err == nil {
			return sha, nil
		}
	}

	return gitClient.GetLatestCommitSHA(ctx, owner, repo)
}

func getUserToken(ctx context.Context, tokenProvider port.TokenProvider, userID string) (string, error) {
	if tokenProvider == nil {
		return "", authdomain.ErrNoGitHubToken
	}

	if userID == "" {
		return "", authdomain.ErrNoGitHubToken
	}

	return tokenProvider.GetUserGitHubToken(ctx, userID)
}

func buildAnalysisFromCompleted(ctx context.Context, repository port.Repository, completed *port.CompletedAnalysis) (*entity.Analysis, error) {
	suitesWithCases, err := repository.GetTestSuitesWithCases(ctx, completed.ID)
	if err != nil {
		return nil, fmt.Errorf("get test suites: %w", err)
	}

	suites := make([]entity.TestSuite, len(suitesWithCases))
	for i, suite := range suitesWithCases {
		testCases := make([]entity.TestCase, len(suite.Tests))
		for j, t := range suite.Tests {
			testCases[j] = entity.TestCase{
				Line:   t.Line,
				Name:   t.Name,
				Status: mapToTestStatus(t.Status),
			}
		}

		suites[i] = entity.TestSuite{
			FilePath:  suite.FilePath,
			Framework: suite.Framework,
			ID:        suite.ID,
			Name:      suite.Name,
			TestCases: testCases,
		}
	}

	return &entity.Analysis{
		BranchName:  completed.BranchName,
		CommitSHA:   completed.CommitSHA,
		CommittedAt: completed.CommittedAt,
		CompletedAt: completed.CompletedAt,
		ID:          completed.ID,
		Owner:       completed.Owner,
		Repo:        completed.Repo,
		TestSuites:  suites,
		TotalSuites: completed.TotalSuites,
		TotalTests:  completed.TotalTests,
	}, nil
}

func mapToTestStatus(status string) entity.TestStatus {
	switch status {
	case "active":
		return entity.TestStatusActive
	case "focused":
		return entity.TestStatusFocused
	case "skipped":
		return entity.TestStatusSkipped
	case "todo":
		return entity.TestStatusTodo
	case "xfail":
		return entity.TestStatusXfail
	default:
		// Unknown statuses default to active (most common state)
		return entity.TestStatusActive
	}
}

func mapQueueStateToAnalysisStatus(state string) entity.AnalysisStatus {
	switch state {
	case "available", "pending", "retryable", "scheduled":
		return entity.AnalysisStatusPending
	case "running":
		return entity.AnalysisStatusRunning
	case "cancelled", "discarded":
		return entity.AnalysisStatusFailed
	case "completed":
		return entity.AnalysisStatusCompleted
	default:
		// Unknown queue states conservatively treated as pending
		return entity.AnalysisStatusPending
	}
}
