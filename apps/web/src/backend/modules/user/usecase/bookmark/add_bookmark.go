package bookmark

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type AddBookmarkInput struct {
	Owner  string
	Repo   string
	UserID string
}

type AddBookmarkOutput struct {
	IsBookmarked bool
}

type AddBookmarkUseCase struct {
	repository port.BookmarkRepository
}

func NewAddBookmarkUseCase(repository port.BookmarkRepository) *AddBookmarkUseCase {
	return &AddBookmarkUseCase{repository: repository}
}

func (uc *AddBookmarkUseCase) Execute(ctx context.Context, input AddBookmarkInput) (*AddBookmarkOutput, error) {
	codebaseID, err := uc.repository.GetCodebaseIDByOwnerRepo(ctx, input.Owner, input.Repo)
	if err != nil {
		if errors.Is(err, domain.ErrCodebaseNotFound) {
			return nil, domain.ErrCodebaseNotFound
		}
		return nil, fmt.Errorf("get codebase: %w", err)
	}

	if err := uc.repository.AddBookmark(ctx, input.UserID, codebaseID); err != nil {
		return nil, fmt.Errorf("add bookmark: %w", err)
	}

	return &AddBookmarkOutput{IsBookmarked: true}, nil
}
