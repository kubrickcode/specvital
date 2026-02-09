package bookmark

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

type GetBookmarksInput struct {
	UserID string
}

type GetBookmarksUseCase struct {
	repository port.BookmarkRepository
}

func NewGetBookmarksUseCase(repository port.BookmarkRepository) *GetBookmarksUseCase {
	return &GetBookmarksUseCase{repository: repository}
}

func (uc *GetBookmarksUseCase) Execute(ctx context.Context, input GetBookmarksInput) ([]*entity.BookmarkedRepository, error) {
	repos, err := uc.repository.GetUserBookmarks(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get bookmarks: %w", err)
	}
	return repos, nil
}
