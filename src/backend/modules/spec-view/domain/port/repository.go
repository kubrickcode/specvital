package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type SpecViewRepository interface {
	CheckAnalysisExists(ctx context.Context, analysisID string) (bool, error)
	CheckSpecDocumentExistsByLanguage(ctx context.Context, analysisID string, language string) (bool, error)
	// GetAvailableLanguages returns all available languages for an analysis with their latest version info.
	GetAvailableLanguages(ctx context.Context, analysisID string) ([]entity.AvailableLanguageInfo, error)
	// GetGenerationStatus returns the latest generation status for an analysis (any language).
	GetGenerationStatus(ctx context.Context, analysisID string) (*entity.SpecGenerationStatus, error)
	// GetGenerationStatusByLanguage returns status for a specific analysis + language combination.
	GetGenerationStatusByLanguage(ctx context.Context, analysisID string, language string) (*entity.SpecGenerationStatus, error)
	// GetSpecDocumentByLanguage retrieves the latest version spec document for a given analysis and language.
	// If language is empty, returns the most recent document regardless of language.
	GetSpecDocumentByLanguage(ctx context.Context, analysisID string, language string) (*entity.SpecDocument, error)
	// GetSpecDocumentByVersion retrieves a specific version of spec document for a given analysis and language.
	GetSpecDocumentByVersion(ctx context.Context, analysisID string, language string, version int) (*entity.SpecDocument, error)
	// GetVersionsByLanguage returns all versions for a specific analysis and language.
	GetVersionsByLanguage(ctx context.Context, analysisID string, language string) ([]entity.VersionInfo, error)
}

type TierLookup interface {
	GetUserTier(ctx context.Context, userID string) (string, error)
}
