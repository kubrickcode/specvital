package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

	formattedVersion := FormatParserVersion(completed.ParserVersion)
	var parserVersion *string
	if formattedVersion != "" {
		parserVersion = &formattedVersion
	}

	return &entity.Analysis{
		BranchName:    completed.BranchName,
		CommitSHA:     completed.CommitSHA,
		CommittedAt:   completed.CommittedAt,
		CompletedAt:   completed.CompletedAt,
		ID:            completed.ID,
		Owner:         completed.Owner,
		ParserVersion: parserVersion,
		Repo:          completed.Repo,
		TestSuites:    suites,
		TotalSuites:   completed.TotalSuites,
		TotalTests:    completed.TotalTests,
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

// FormatParserVersion transforms Go module version to display format.
// Input: "v1.5.1-0.20260112121406-deacdda09e17" -> Output: "v1.5.1 (deacdda)"
// Input: "v1.5.1" -> Output: "v1.5.1"
// Input: nil/empty -> Output: ""
func FormatParserVersion(version *string) string {
	if version == nil || *version == "" {
		return ""
	}

	v := *version

	// Go pseudo-version formats:
	// 1. vX.Y.Z-0.YYYYMMDDHHMMSS-commitHash (tagged release with subsequent commits)
	// 2. vX.0.0-YYYYMMDDHHMMSS-commitHash (no prior tag, timestamp without "0." prefix)
	parts := strings.Split(v, "-")

	// Standard version or not enough parts for pseudo-version
	if len(parts) < 3 {
		return v
	}

	// Check if this is a pseudo-version format
	// Pattern 1: second part starts with "0." (e.g., "0.20260112121406")
	// Pattern 2: second part is a timestamp (14 digits)
	isPseudoVersion := strings.HasPrefix(parts[1], "0.") || isTimestamp(parts[1])
	if !isPseudoVersion {
		return v
	}

	// Extract base version (first part)
	baseVersion := parts[0]

	// Extract commit hash (last part) and take first 7 chars
	commitHash := parts[len(parts)-1]
	if len(commitHash) > 7 {
		commitHash = commitHash[:7]
	}

	return fmt.Sprintf("%s (%s)", baseVersion, commitHash)
}

// isTimestamp checks if a string is a 14-digit timestamp (YYYYMMDDHHMMSS)
func isTimestamp(s string) bool {
	if len(s) != 14 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
