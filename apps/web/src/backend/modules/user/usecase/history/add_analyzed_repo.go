package history

import (
	"context"
	"fmt"
	"time"

	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type AddAnalyzedRepoInput struct {
	Owner  string
	Repo   string
	UserID string
}

type AddAnalyzedRepoOutput struct {
	AnalysisID string
	UpdatedAt  time.Time
}

type AddAnalyzedRepoUseCase struct {
	repository port.AnalysisHistoryRepository
}

func NewAddAnalyzedRepoUseCase(repository port.AnalysisHistoryRepository) *AddAnalyzedRepoUseCase {
	return &AddAnalyzedRepoUseCase{repository: repository}
}

func (uc *AddAnalyzedRepoUseCase) Execute(ctx context.Context, input AddAnalyzedRepoInput) (*AddAnalyzedRepoOutput, error) {
	result, err := uc.repository.AddUserAnalyzedRepository(ctx, input.UserID, input.Owner, input.Repo)
	if err != nil {
		return nil, fmt.Errorf("add analyzed repository: %w", err)
	}

	return &AddAnalyzedRepoOutput{
		AnalysisID: result.AnalysisID,
		UpdatedAt:  result.UpdatedAt,
	}, nil
}
