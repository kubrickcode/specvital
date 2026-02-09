package port

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/user/domain/entity"
)

type AnalyzedReposParams struct {
	Cursor    *entity.CursorData
	Limit     int
	Ownership entity.OwnershipFilter
}

type AddHistoryResult struct {
	AnalysisID string
	UpdatedAt  time.Time
}

type AnalysisHistoryRepository interface {
	AddUserAnalyzedRepository(ctx context.Context, userID, owner, repo string) (*AddHistoryResult, error)
	CheckUserHistoryExists(ctx context.Context, userID, owner, repo string) (bool, error)
	GetUserAnalyzedRepositories(ctx context.Context, userID string, params AnalyzedReposParams) ([]*entity.AnalyzedRepository, error)
}
