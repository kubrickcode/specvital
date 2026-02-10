package mapper

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/entity"
)

type CompletedResponseOptions struct {
	IsInMyHistory *bool
}

func ToCompletedResponse(analysis *entity.Analysis, opts ...CompletedResponseOptions) (api.AnalysisResponse, error) {
	if analysis == nil {
		return api.AnalysisResponse{}, fmt.Errorf("analysis is nil")
	}

	var options CompletedResponseOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	suites := make([]api.TestSuite, len(analysis.TestSuites))
	frameworkStats := make(map[string]*api.FrameworkSummary)

	for i, suite := range analysis.TestSuites {
		tests := make([]api.TestCase, len(suite.TestCases))
		for j, testCase := range suite.TestCases {
			tests[j] = api.TestCase{
				FilePath:  suite.FilePath,
				Framework: suite.Framework,
				Line:      testCase.Line,
				Name:      testCase.Name,
				Status:    toAPITestStatus(testCase.Status),
			}

			if _, exists := frameworkStats[suite.Framework]; !exists {
				frameworkStats[suite.Framework] = &api.FrameworkSummary{
					Framework: suite.Framework,
				}
			}
			fs := frameworkStats[suite.Framework]
			fs.Total++
			switch testCase.Status {
			case entity.TestStatusActive:
				fs.Active++
			case entity.TestStatusFocused:
				fs.Focused++
			case entity.TestStatusSkipped:
				fs.Skipped++
			case entity.TestStatusTodo:
				fs.Todo++
			case entity.TestStatusXfail:
				fs.Xfail++
			}
		}

		suites[i] = api.TestSuite{
			FilePath:  suite.FilePath,
			Framework: suite.Framework,
			SuiteName: suite.Name,
			Tests:     tests,
		}
	}

	frameworks := make([]api.FrameworkSummary, 0, len(frameworkStats))
	var totalActive, totalFocused, totalSkipped, totalTodo, totalXfail int
	for _, fs := range frameworkStats {
		frameworks = append(frameworks, *fs)
		totalActive += fs.Active
		totalFocused += fs.Focused
		totalSkipped += fs.Skipped
		totalTodo += fs.Todo
		totalXfail += fs.Xfail
	}

	sort.Slice(frameworks, func(i, j int) bool {
		return frameworks[i].Framework < frameworks[j].Framework
	})

	result := api.AnalysisResult{
		AnalyzedAt:    analysis.CompletedAt,
		BranchName:    analysis.BranchName,
		CommitSHA:     analysis.CommitSHA,
		CommittedAt:   analysis.CommittedAt,
		ID:            uuid.MustParse(analysis.ID),
		IsInMyHistory: options.IsInMyHistory,
		Owner:         analysis.Owner,
		ParserVersion: analysis.ParserVersion,
		Repo:          analysis.Repo,
		Suites:        suites,
		Summary: api.Summary{
			Active:     totalActive,
			Focused:    totalFocused,
			Frameworks: frameworks,
			Skipped:    totalSkipped,
			Todo:       totalTodo,
			Total:      analysis.TotalTests,
			Xfail:      totalXfail,
		},
	}

	var response api.AnalysisResponse
	if err := response.FromCompletedResponse(api.CompletedResponse{Data: result}); err != nil {
		return api.AnalysisResponse{}, fmt.Errorf("marshal completed response: %w", err)
	}

	return response, nil
}

func ToStatusResponse(progress *entity.AnalysisProgress) (api.AnalysisResponse, error) {
	if progress == nil {
		return api.AnalysisResponse{}, fmt.Errorf("progress is nil")
	}

	var response api.AnalysisResponse
	var err error

	switch progress.Status {
	case entity.AnalysisStatusCompleted:
		err = response.FromCompletedResponse(api.CompletedResponse{})
	case entity.AnalysisStatusRunning:
		err = response.FromAnalyzingResponse(api.AnalyzingResponse{
			StartedAt: progress.StartedAt,
		})
	case entity.AnalysisStatusPending:
		err = response.FromQueuedResponse(api.QueuedResponse{})
	case entity.AnalysisStatusFailed:
		errorMsg := "analysis failed"
		if progress.ErrorMessage != nil {
			errorMsg = *progress.ErrorMessage
		}
		err = response.FromFailedResponse(api.FailedResponse{
			Error: errorMsg,
		})
	default:
		err = response.FromFailedResponse(api.FailedResponse{
			Error: "unknown status",
		})
	}

	if err != nil {
		return api.AnalysisResponse{}, fmt.Errorf("marshal status response: %w", err)
	}

	return response, nil
}

func ToAnalysisHistoryResponse(items []AnalysisHistoryInput, headCommitSHA string) (api.AnalysisHistoryResponse, error) {
	data := make([]api.AnalysisHistoryItem, len(items))
	for i, item := range items {
		id, err := uuid.Parse(item.ID)
		if err != nil {
			return api.AnalysisHistoryResponse{}, fmt.Errorf("invalid analysis ID %s: %w", item.ID, err)
		}
		isHead := item.CommitSHA == headCommitSHA
		data[i] = api.AnalysisHistoryItem{
			BranchName:  item.BranchName,
			CommitSHA:   item.CommitSHA,
			CommittedAt: item.CommittedAt,
			CompletedAt: item.CompletedAt,
			ID:          id,
			IsHead:      &isHead,
			TotalTests:  item.TotalTests,
		}
	}
	return api.AnalysisHistoryResponse{Data: data}, nil
}

type AnalysisHistoryInput struct {
	BranchName  *string
	CommitSHA   string
	CommittedAt *time.Time
	CompletedAt time.Time
	ID          string
	TotalTests  int
}

func toAPITestStatus(status entity.TestStatus) api.TestStatus {
	switch status {
	case entity.TestStatusActive:
		return api.Active
	case entity.TestStatusFocused:
		return api.Focused
	case entity.TestStatusSkipped:
		return api.Skipped
	case entity.TestStatusTodo:
		return api.Todo
	case entity.TestStatusXfail:
		return api.Xfail
	default:
		return api.Active
	}
}
