package port

import (
	"context"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/user/domain/entity"
)

type ActiveTaskRepository interface {
	GetUserActiveTasks(ctx context.Context, userID string) ([]entity.ActiveTask, error)
}
