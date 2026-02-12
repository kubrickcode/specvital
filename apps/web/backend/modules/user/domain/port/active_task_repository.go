package port

import (
	"context"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/user/domain/entity"
)

type ActiveTaskRepository interface {
	GetUserActiveTasks(ctx context.Context, userID string) ([]entity.ActiveTask, error)
}
