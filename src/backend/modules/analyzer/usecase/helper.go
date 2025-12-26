package usecase

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

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
