package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/port"
)

type GetUpdateStatusInput struct {
	Owner  string
	Repo   string
	UserID string
}

type GetUpdateStatusUseCase struct {
	gitClient     port.GitClient
	repository    port.Repository
	systemConfig  port.SystemConfigReader
	tokenProvider port.TokenProvider
}

func NewGetUpdateStatusUseCase(
	gitClient port.GitClient,
	repository port.Repository,
	systemConfig port.SystemConfigReader,
	tokenProvider port.TokenProvider,
) *GetUpdateStatusUseCase {
	return &GetUpdateStatusUseCase{
		gitClient:     gitClient,
		repository:    repository,
		systemConfig:  systemConfig,
		tokenProvider: tokenProvider,
	}
}

func (uc *GetUpdateStatusUseCase) Execute(ctx context.Context, input GetUpdateStatusInput) (*entity.UpdateStatusResult, error) {
	if input.Owner == "" || input.Repo == "" {
		return nil, errors.New("owner and repo are required")
	}

	// Get latest completed analysis for parser version check and display
	completed, err := uc.repository.GetLatestCompletedAnalysis(ctx, input.Owner, input.Repo)
	if err != nil {
		return nil, fmt.Errorf("get latest analysis: %w", err)
	}

	parserOutdated := uc.isParserOutdated(ctx, completed.ParserVersion)

	latestSHA, err := getLatestCommitWithAuth(ctx, uc.gitClient, uc.tokenProvider, input.Owner, input.Repo, input.UserID)
	if err != nil {
		return &entity.UpdateStatusResult{
			AnalyzedCommitSHA: completed.CommitSHA,
			ParserOutdated:    parserOutdated,
			Status:            entity.UpdateStatusUnknown,
		}, nil
	}

	// Check if analysis exists for the GitHub latest commit (not just comparing with DB latest)
	// This handles git reset scenarios where an old commit becomes the new HEAD
	status := entity.UpdateStatusUpToDate
	analysisExists, err := uc.repository.CheckAnalysisExistsByCommitSHA(ctx, input.Owner, input.Repo, latestSHA)
	if err != nil {
		// On error, fall back to comparing with completed analysis
		if latestSHA != completed.CommitSHA {
			status = entity.UpdateStatusNewCommits
		}
	} else if !analysisExists {
		status = entity.UpdateStatusNewCommits
	}

	return &entity.UpdateStatusResult{
		AnalyzedCommitSHA: completed.CommitSHA,
		LatestCommitSHA:   latestSHA,
		ParserOutdated:    parserOutdated,
		Status:            status,
	}, nil
}

// isParserOutdated checks if the cached analysis was created with an older parser version.
func (uc *GetUpdateStatusUseCase) isParserOutdated(ctx context.Context, cachedVersion *string) bool {
	if cachedVersion == nil {
		// Legacy data without parser_version is considered outdated
		return true
	}

	currentVersion, err := uc.systemConfig.GetParserVersion(ctx)
	if err != nil {
		// System config unavailable → assume not outdated to avoid false positives
		return false
	}

	return *cachedVersion != currentVersion
}
