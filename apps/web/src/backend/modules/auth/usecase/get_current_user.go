package usecase

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type GetCurrentUserInput struct {
	UserID string
}

type GetCurrentUserOutput struct {
	User *entity.User
}

type GetCurrentUserUseCase struct {
	repository port.Repository
}

func NewGetCurrentUserUseCase(repository port.Repository) *GetCurrentUserUseCase {
	if repository == nil {
		panic("repository is required")
	}
	return &GetCurrentUserUseCase{
		repository: repository,
	}
}

func (uc *GetCurrentUserUseCase) Execute(ctx context.Context, input GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	user, err := uc.repository.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &GetCurrentUserOutput{User: user}, nil
}
