package auth

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/auth/domain"
)

type BookmarkService interface {
	AddBookmark(ctx context.Context, userID, owner, repo string) (bool, error)
	GetUserBookmarkedRepos(ctx context.Context, userID string) ([]*domain.BookmarkedRepository, error)
	RemoveBookmark(ctx context.Context, userID, owner, repo string) (bool, error)
}

type bookmarkServiceImpl struct {
	repo BookmarkRepository
}

func NewBookmarkService(repo BookmarkRepository) BookmarkService {
	return &bookmarkServiceImpl{repo: repo}
}

func (s *bookmarkServiceImpl) AddBookmark(ctx context.Context, userID, owner, repo string) (bool, error) {
	codebaseID, err := s.repo.GetCodebaseIDByOwnerRepo(ctx, owner, repo)
	if err != nil {
		if errors.Is(err, domain.ErrCodebaseNotFound) {
			return false, domain.ErrCodebaseNotFound
		}
		return false, fmt.Errorf("get codebase: %w", err)
	}

	if err := s.repo.AddBookmark(ctx, userID, codebaseID); err != nil {
		return false, fmt.Errorf("add bookmark: %w", err)
	}

	return true, nil
}

func (s *bookmarkServiceImpl) GetUserBookmarkedRepos(ctx context.Context, userID string) ([]*domain.BookmarkedRepository, error) {
	repos, err := s.repo.GetUserBookmarks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get bookmarks: %w", err)
	}

	return repos, nil
}

func (s *bookmarkServiceImpl) RemoveBookmark(ctx context.Context, userID, owner, repo string) (bool, error) {
	codebaseID, err := s.repo.GetCodebaseIDByOwnerRepo(ctx, owner, repo)
	if err != nil {
		if errors.Is(err, domain.ErrCodebaseNotFound) {
			return false, domain.ErrCodebaseNotFound
		}
		return false, fmt.Errorf("get codebase: %w", err)
	}

	if err := s.repo.RemoveBookmark(ctx, userID, codebaseID); err != nil {
		return false, fmt.Errorf("remove bookmark: %w", err)
	}

	return false, nil
}
