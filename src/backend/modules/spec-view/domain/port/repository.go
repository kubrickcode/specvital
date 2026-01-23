package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type SpecViewRepository interface {
	CheckAnalysisExists(ctx context.Context, analysisID string) (bool, error)
	CheckSpecDocumentExistsByLanguage(ctx context.Context, analysisID string, language string) (bool, error)
	// CheckSpecDocumentOwnership returns the first document owner for the given analysis.
	// Returns document ID and user ID of the original creator, nil if no document exists.
	CheckSpecDocumentOwnership(ctx context.Context, analysisID string) (*entity.DocumentOwnership, error)
	// GetAvailableLanguages returns all available languages for an analysis with their latest version info.
	// Deprecated: Use GetAvailableLanguagesByUser for access-controlled access.
	GetAvailableLanguages(ctx context.Context, analysisID string) ([]entity.AvailableLanguageInfo, error)
	// GetAvailableLanguagesByUser returns available languages for documents owned by the user for the given analysis.
	GetAvailableLanguagesByUser(ctx context.Context, userID string, analysisID string) ([]entity.AvailableLanguageInfo, error)
	// GetGenerationStatus returns the latest generation status for an analysis (any language).
	GetGenerationStatus(ctx context.Context, analysisID string) (*entity.SpecGenerationStatus, error)
	// GetGenerationStatusByLanguage returns status for a specific analysis + language combination.
	GetGenerationStatusByLanguage(ctx context.Context, analysisID string, language string) (*entity.SpecGenerationStatus, error)
	// GetSpecDocumentByLanguage retrieves the latest version spec document for a given analysis and language.
	// If language is empty, returns the most recent document regardless of language.
	// Deprecated: Use GetSpecDocumentByUser for access-controlled access.
	GetSpecDocumentByLanguage(ctx context.Context, analysisID string, language string) (*entity.SpecDocument, error)
	// GetSpecDocumentByUser retrieves document owned by the user for the given analysis and language.
	// If language is empty, returns the most recent document regardless of language.
	GetSpecDocumentByUser(ctx context.Context, userID string, analysisID string, language string) (*entity.SpecDocument, error)
	// GetSpecDocumentByVersion retrieves a specific version of spec document for a given analysis and language.
	// Deprecated: Use GetSpecDocumentByUserAndVersion for access-controlled access.
	GetSpecDocumentByVersion(ctx context.Context, analysisID string, language string, version int) (*entity.SpecDocument, error)
	// GetSpecDocumentByUserAndVersion retrieves a specific version of document owned by the user.
	GetSpecDocumentByUserAndVersion(ctx context.Context, userID string, analysisID string, language string, version int) (*entity.SpecDocument, error)
	// GetVersionsByLanguage returns all versions for a specific analysis and language.
	// Deprecated: Use GetVersionsByUser for access-controlled access.
	GetVersionsByLanguage(ctx context.Context, analysisID string, language string) ([]entity.VersionInfo, error)
	// GetVersionsByUser returns all versions for documents owned by the user for the given analysis and language.
	GetVersionsByUser(ctx context.Context, userID string, analysisID string, language string) ([]entity.VersionInfo, error)
}

type TierLookup interface {
	GetUserTier(ctx context.Context, userID string) (string, error)
}
