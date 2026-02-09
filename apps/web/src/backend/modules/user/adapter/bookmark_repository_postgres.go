package adapter

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

const defaultHost = "github.com"

type BookmarkRepositoryPostgres struct {
	queries *db.Queries
}

var _ port.BookmarkRepository = (*BookmarkRepositoryPostgres)(nil)

func NewBookmarkRepository(queries *db.Queries) *BookmarkRepositoryPostgres {
	return &BookmarkRepositoryPostgres{queries: queries}
}

func (r *BookmarkRepositoryPostgres) AddBookmark(ctx context.Context, userID, codebaseID string) error {
	userUUID, err := stringToUUID(userID)
	if err != nil {
		return fmt.Errorf("parse user ID: %w", err)
	}

	codebaseUUID, err := stringToUUID(codebaseID)
	if err != nil {
		return fmt.Errorf("parse codebase ID: %w", err)
	}

	params := db.AddBookmarkParams{
		UserID:     userUUID,
		CodebaseID: codebaseUUID,
	}

	if err := r.queries.AddBookmark(ctx, params); err != nil {
		return fmt.Errorf("add bookmark: %w", err)
	}

	return nil
}

func (r *BookmarkRepositoryPostgres) GetCodebaseIDByOwnerRepo(ctx context.Context, owner, repo string) (string, error) {
	params := db.GetCodebaseByOwnerRepoParams{
		Host:  defaultHost,
		Owner: owner,
		Name:  repo,
	}

	id, err := r.queries.GetCodebaseByOwnerRepo(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrCodebaseNotFound
		}
		return "", fmt.Errorf("get codebase: %w", err)
	}

	return uuidToString(id), nil
}

func (r *BookmarkRepositoryPostgres) GetUserBookmarks(ctx context.Context, userID string) ([]*entity.BookmarkedRepository, error) {
	userUUID, err := stringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user ID: %w", err)
	}

	rows, err := r.queries.GetUserBookmarks(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("get user bookmarks: %w", err)
	}

	return mapBookmarkedRepositoriesFromDB(rows), nil
}

func (r *BookmarkRepositoryPostgres) RemoveBookmark(ctx context.Context, userID, codebaseID string) error {
	userUUID, err := stringToUUID(userID)
	if err != nil {
		return fmt.Errorf("parse user ID: %w", err)
	}

	codebaseUUID, err := stringToUUID(codebaseID)
	if err != nil {
		return fmt.Errorf("parse codebase ID: %w", err)
	}

	params := db.RemoveBookmarkParams{
		UserID:     userUUID,
		CodebaseID: codebaseUUID,
	}

	if err := r.queries.RemoveBookmark(ctx, params); err != nil {
		return fmt.Errorf("remove bookmark: %w", err)
	}

	return nil
}

func mapBookmarkedRepositoriesFromDB(rows []db.GetUserBookmarksRow) []*entity.BookmarkedRepository {
	result := make([]*entity.BookmarkedRepository, 0, len(rows))

	for _, row := range rows {
		repo := &entity.BookmarkedRepository{
			BookmarkedAt: row.BookmarkedAt.Time,
			CodebaseID:   uuidToString(row.CodebaseID),
			CommitSHA:    row.CommitSha,
			Name:         row.Name,
			Owner:        row.Owner,
			TotalTests:   int(row.TotalTests),
		}

		if row.AnalyzedAt.Valid {
			t := row.AnalyzedAt.Time
			repo.AnalyzedAt = &t
		}

		result = append(result, repo)
	}

	return result
}
