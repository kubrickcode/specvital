package bookmark

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type RemoveBookmarkInput struct {
	Owner  string
	Repo   string
	UserID string
}

type RemoveBookmarkOutput struct {
	IsBookmarked bool
}

type RemoveBookmarkUseCase struct {
	repository port.BookmarkRepository
}

func NewRemoveBookmarkUseCase(repository port.BookmarkRepository) *RemoveBookmarkUseCase {
	return &RemoveBookmarkUseCase{repository: repository}
}

func (uc *RemoveBookmarkUseCase) Execute(ctx context.Context, input RemoveBookmarkInput) (*RemoveBookmarkOutput, error) {
	codebaseID, err := uc.repository.GetCodebaseIDByOwnerRepo(ctx, input.Owner, input.Repo)
	if err != nil {
		if errors.Is(err, domain.ErrCodebaseNotFound) {
			return nil, domain.ErrCodebaseNotFound
		}
		return nil, fmt.Errorf("get codebase: %w", err)
	}

	if err := uc.repository.RemoveBookmark(ctx, input.UserID, codebaseID); err != nil {
		return nil, fmt.Errorf("remove bookmark: %w", err)
	}

	return &RemoveBookmarkOutput{IsBookmarked: false}, nil
}
