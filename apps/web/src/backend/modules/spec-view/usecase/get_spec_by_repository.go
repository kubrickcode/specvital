package usecase

import (
	"context"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/port"
)

type GetSpecByRepositoryInput struct {
	// DocumentID is optional. If set, returns the specific document (takes precedence over Version).
	DocumentID string
	// Language is optional. If empty, defaults to "English".
	Language string
	// Name is the repository name. Required.
	Name string
	// Owner is the repository owner. Required.
	Owner string
	// UserID is required. Empty userID returns ErrUnauthorized.
	UserID string
	// Version is optional. If 0, returns the latest version.
	Version int
}

type GetSpecByRepositoryOutput struct {
	Document *entity.RepoSpecDocument
}

type GetSpecByRepositoryUseCase struct {
	repo port.SpecViewRepository
}

func NewGetSpecByRepositoryUseCase(repo port.SpecViewRepository) *GetSpecByRepositoryUseCase {
	return &GetSpecByRepositoryUseCase{repo: repo}
}

func (uc *GetSpecByRepositoryUseCase) Execute(ctx context.Context, input GetSpecByRepositoryInput) (*GetSpecByRepositoryOutput, error) {
	if input.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if input.Owner == "" || input.Name == "" {
		return nil, domain.ErrInvalidRepository
	}

	if !entity.IsValidRepositoryName(input.Owner) || !entity.IsValidRepositoryName(input.Name) {
		return nil, domain.ErrInvalidRepository
	}

	exists, err := uc.repo.CheckCodebaseExists(ctx, input.Owner, input.Name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrCodebaseNotFound
	}

	// Determine language: use provided, or fallback to any available language
	language := input.Language
	if language == "" {
		// Get available languages first to find any existing spec
		availableLanguages, err := uc.repo.GetAvailableLanguagesByRepository(ctx, input.UserID, input.Owner, input.Name)
		if err != nil {
			return nil, err
		}
		if len(availableLanguages) == 0 {
			return nil, domain.ErrDocumentNotFound
		}
		// Use the first available language (alphabetical order from DB query)
		language = availableLanguages[0].Language
	}

	var doc *entity.RepoSpecDocument
	if input.DocumentID != "" {
		// Validate document ID format before querying
		if !entity.IsValidDocumentID(input.DocumentID) {
			return nil, domain.ErrInvalidDocumentID
		}
		// DocumentID takes precedence - fetch specific document directly
		doc, err = uc.repo.GetSpecDocumentByRepositoryAndDocumentId(ctx, input.UserID, input.Owner, input.Name, input.DocumentID)
		if err != nil {
			return nil, err
		}
	} else {
		if !entity.IsValidLanguage(language) {
			return nil, domain.ErrInvalidLanguage
		}

		if input.Version > 0 {
			doc, err = uc.repo.GetSpecDocumentByRepositoryAndVersion(ctx, input.UserID, input.Owner, input.Name, language, input.Version)
		} else {
			doc, err = uc.repo.GetSpecDocumentByRepository(ctx, input.UserID, input.Owner, input.Name, language)
		}
		if err != nil {
			return nil, err
		}
	}

	if doc == nil {
		return nil, domain.ErrDocumentNotFound
	}

	availableLanguages, err := uc.repo.GetAvailableLanguagesByRepository(ctx, input.UserID, input.Owner, input.Name)
	if err != nil {
		return nil, err
	}
	doc.AvailableLanguages = availableLanguages

	return &GetSpecByRepositoryOutput{Document: doc}, nil
}
