package port

import (
	"context"
	"time"

	"github.com/specvital/web/src/backend/modules/github-app/domain/entity"
)

type InstallationRepository interface {
	Delete(ctx context.Context, installationID int64) error
	GetByAccountID(ctx context.Context, accountID int64) (*entity.Installation, error)
	GetByInstallationID(ctx context.Context, installationID int64) (*entity.Installation, error)
	ListByAccountIDs(ctx context.Context, accountIDs []int64) ([]entity.Installation, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.Installation, error)
	ListOrganizations(ctx context.Context) ([]entity.Installation, error)
	UpdateSuspended(ctx context.Context, installationID int64, suspendedAt *time.Time) error
	Upsert(ctx context.Context, installation *entity.Installation) error
}
