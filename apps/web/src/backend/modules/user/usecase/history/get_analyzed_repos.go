package history

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type GetAnalyzedReposInput struct {
	Cursor    *string
	Limit     int
	Ownership entity.OwnershipFilter
	UserID    string
}

type GetAnalyzedReposUseCase struct {
	repository port.AnalysisHistoryRepository
}

func NewGetAnalyzedReposUseCase(repository port.AnalysisHistoryRepository) *GetAnalyzedReposUseCase {
	return &GetAnalyzedReposUseCase{repository: repository}
}

func (uc *GetAnalyzedReposUseCase) Execute(ctx context.Context, input GetAnalyzedReposInput) (*entity.AnalyzedReposResult, error) {
	var cursorData *entity.CursorData
	if input.Cursor != nil {
		decoded, err := DecodeCursor(*input.Cursor)
		if err != nil {
			return nil, err
		}
		cursorData = decoded
	}

	fetchLimit := input.Limit + 1
	fetchParams := port.AnalyzedReposParams{
		Cursor:    cursorData,
		Limit:     fetchLimit,
		Ownership: input.Ownership,
	}

	repos, err := uc.repository.GetUserAnalyzedRepositories(ctx, input.UserID, fetchParams)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCursor) {
			return nil, err
		}
		return nil, fmt.Errorf("get user analyzed repositories: %w", err)
	}

	var nextCursor *string
	hasNext := false

	if len(repos) > input.Limit {
		hasNext = true
		repos = repos[:input.Limit]
		lastRepo := repos[len(repos)-1]
		cursorValue, err := EncodeCursor(lastRepo.UpdatedAt, lastRepo.HistoryID)
		if err != nil {
			return nil, fmt.Errorf("encode cursor: %w", err)
		}
		nextCursor = &cursorValue
	}

	return &entity.AnalyzedReposResult{
		Data:       repos,
		HasNext:    hasNext,
		NextCursor: nextCursor,
	}, nil
}
