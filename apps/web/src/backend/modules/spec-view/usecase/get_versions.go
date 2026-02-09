package usecase

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type GetVersionsInput struct {
	AnalysisID string
	Language   string
	// UserID is required. Empty userID returns ErrUnauthorized.
	UserID string
}

type GetVersionsOutput struct {
	Language      string
	LatestVersion int
	Versions      []entity.VersionInfo
}

type GetVersionsUseCase struct {
	repo port.SpecViewRepository
}

func NewGetVersionsUseCase(repo port.SpecViewRepository) *GetVersionsUseCase {
	return &GetVersionsUseCase{repo: repo}
}

func (uc *GetVersionsUseCase) Execute(ctx context.Context, input GetVersionsInput) (*GetVersionsOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.AnalysisID == "" {
		return nil, domain.ErrInvalidAnalysisID
	}
	if input.Language == "" {
		return nil, domain.ErrInvalidLanguage
	}
	if !entity.IsValidLanguage(input.Language) {
		return nil, domain.ErrInvalidLanguage
	}

	versions, err := uc.repo.GetVersionsByUser(ctx, input.UserID, input.AnalysisID, input.Language)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, domain.ErrDocumentNotFound
	}

	return &GetVersionsOutput{
		Language:      input.Language,
		LatestVersion: versions[0].Version,
		Versions:      versions,
	}, nil
}
