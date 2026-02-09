package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/user/domain/entity"
)

type BookmarkRepository interface {
	AddBookmark(ctx context.Context, userID, codebaseID string) error
	GetCodebaseIDByOwnerRepo(ctx context.Context, owner, repo string) (string, error)
	GetUserBookmarks(ctx context.Context, userID string) ([]*entity.BookmarkedRepository, error)
	RemoveBookmark(ctx context.Context, userID, codebaseID string) error
}
