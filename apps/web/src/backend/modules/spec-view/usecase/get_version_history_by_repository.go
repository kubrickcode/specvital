package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type GetVersionHistoryByRepositoryInput struct {
	// Language is required. Returns versions for this language.
	Language string
	// Name is the repository name. Required.
	Name string
	// Owner is the repository owner. Required.
	Owner string
	// UserID is required. Empty userID returns ErrUnauthorized.
	UserID string
}

type GetVersionHistoryByRepositoryOutput struct {
	Language string
	Versions []entity.RepoVersionInfo
}

type GetVersionHistoryByRepositoryUseCase struct {
	repo port.SpecViewRepository
}

func NewGetVersionHistoryByRepositoryUseCase(repo port.SpecViewRepository) *GetVersionHistoryByRepositoryUseCase {
	return &GetVersionHistoryByRepositoryUseCase{repo: repo}
}

func (uc *GetVersionHistoryByRepositoryUseCase) Execute(ctx context.Context, input GetVersionHistoryByRepositoryInput) (*GetVersionHistoryByRepositoryOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.Owner == "" || input.Name == "" {
		return nil, domain.ErrInvalidRepository
	}

	if !entity.IsValidRepositoryName(input.Owner) || !entity.IsValidRepositoryName(input.Name) {
		return nil, domain.ErrInvalidRepository
	}

	if input.Language == "" || !entity.IsValidLanguage(input.Language) {
		return nil, domain.ErrInvalidLanguage
	}

	exists, err := uc.repo.CheckCodebaseExists(ctx, input.Owner, input.Name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrCodebaseNotFound
	}

	versions, err := uc.repo.GetVersionHistoryByRepository(ctx, input.UserID, input.Owner, input.Name, input.Language)
	if err != nil {
		return nil, err
	}

	return &GetVersionHistoryByRepositoryOutput{
		Language: input.Language,
		Versions: versions,
	}, nil
}
