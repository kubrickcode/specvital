package activetask

import (
	"context"
	"errors"
	"fmt"

	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	"github.com/specvital/web/src/backend/modules/user/domain/port"
)

var ErrUserIDRequired = errors.New("userID is required")

type GetUserActiveTasksInput struct {
	UserID string
}

type GetUserActiveTasksOutput struct {
	Tasks []entity.ActiveTask
}

type GetUserActiveTasksUseCase struct {
	repository port.ActiveTaskRepository
}

func NewGetUserActiveTasksUseCase(repository port.ActiveTaskRepository) *GetUserActiveTasksUseCase {
	return &GetUserActiveTasksUseCase{repository: repository}
}

func (uc *GetUserActiveTasksUseCase) Execute(ctx context.Context, input GetUserActiveTasksInput) (*GetUserActiveTasksOutput, error) {
	if input.UserID == "" {
		return nil, ErrUserIDRequired
	}

	tasks, err := uc.repository.GetUserActiveTasks(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user active tasks: %w", err)
	}

	return &GetUserActiveTasksOutput{Tasks: tasks}, nil
}
