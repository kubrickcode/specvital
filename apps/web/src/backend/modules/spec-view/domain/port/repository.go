package port

import (
	"context"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/entity"
)

type SpecViewRepository interface {
	CheckAnalysisExists(ctx context.Context, analysisID string) (bool, error)
	// GetAnalysisTestCount returns the total test count for an analysis.
	// Used for quota calculation before spec generation.
	GetAnalysisTestCount(ctx context.Context, analysisID string) (int, error)
	CheckSpecDocumentExistsByLanguage(ctx context.Context, analysisID string, language string) (bool, error)
	// HasPreviousSpecByLanguage checks if user has a spec document for the same codebase
	// but different analysis. Used to determine if behavior cache might exist.
	// Deprecated: Use GetLanguagesWithPreviousSpec for batch queries to avoid N+1.
	HasPreviousSpecByLanguage(ctx context.Context, userID, analysisID, language string) (bool, error)
	// GetLanguagesWithPreviousSpec returns all languages where user has a spec document
	// for the same codebase but different analysis. Batch version to avoid N+1 queries.
	GetLanguagesWithPreviousSpec(ctx context.Context, userID, analysisID string) ([]string, error)
	// GetAvailableLanguages returns all available languages for an analysis with their latest version info.
	// Deprecated: Use GetAvailableLanguagesByUser for access-controlled access.
	GetAvailableLanguages(ctx context.Context, analysisID string) ([]entity.AvailableLanguageInfo, error)
	// GetAvailableLanguagesByUser returns available languages for documents owned by the user for the given analysis.
	GetAvailableLanguagesByUser(ctx context.Context, userID string, analysisID string) ([]entity.AvailableLanguageInfo, error)
	// GetGenerationStatus returns the latest generation status for a user and analysis (any language).
	GetGenerationStatus(ctx context.Context, userID string, analysisID string) (*entity.SpecGenerationStatus, error)
	// GetGenerationStatusByLanguage returns status for a specific user, analysis, and language combination.
	GetGenerationStatusByLanguage(ctx context.Context, userID string, analysisID string, language string) (*entity.SpecGenerationStatus, error)
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

	// Repository-based queries (cross-analysis access)

	// CheckCodebaseExists checks if a codebase exists by owner/name.
	CheckCodebaseExists(ctx context.Context, owner, name string) (bool, error)
	// GetSpecDocumentByRepository retrieves the latest spec document for a repository.
	// Provides cross-analysis access - returns spec from any analysis of this repository.
	GetSpecDocumentByRepository(ctx context.Context, userID, owner, name, language string) (*entity.RepoSpecDocument, error)
	// GetSpecDocumentByRepositoryAndVersion retrieves a specific version of spec document for a repository.
	GetSpecDocumentByRepositoryAndVersion(ctx context.Context, userID, owner, name, language string, version int) (*entity.RepoSpecDocument, error)
	// GetSpecDocumentByRepositoryAndDocumentId retrieves a specific spec document by its ID for a repository.
	GetSpecDocumentByRepositoryAndDocumentId(ctx context.Context, userID, owner, name, documentID string) (*entity.RepoSpecDocument, error)
	// GetVersionHistoryByRepository returns all spec versions for a repository across all analyses.
	// Each version includes the commit SHA at generation time.
	GetVersionHistoryByRepository(ctx context.Context, userID, owner, name, language string) ([]entity.RepoVersionInfo, error)
	// GetAvailableLanguagesByRepository returns all available languages for a repository.
	GetAvailableLanguagesByRepository(ctx context.Context, userID, owner, name string) ([]entity.AvailableLanguageInfo, error)

	// GetCachePredictionData returns cache prediction statistics by comparing
	// current analysis test cases with previous spec behaviors.
	// Matches by test_case.name + test_file.file_path to predict cache hits.
	GetCachePredictionData(ctx context.Context, userID, analysisID, language string) (*entity.CachePredictionData, error)
	// GetCurrentAnalysisTestCount returns the total test count for an analysis.
	GetCurrentAnalysisTestCount(ctx context.Context, analysisID string) (int, error)
}

type TierLookup interface {
	GetUserTier(ctx context.Context, userID string) (string, error)
}
