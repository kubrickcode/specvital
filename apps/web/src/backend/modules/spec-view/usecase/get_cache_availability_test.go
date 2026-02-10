package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/entity"
)

type mockCacheAvailabilityRepository struct {
	analysisExists        bool
	analysisExistsErr     error
	langsWithPreviousSpec []string
	langsWithPreviousErr  error
}

func (m *mockCacheAvailabilityRepository) CheckAnalysisExists(_ context.Context, _ string) (bool, error) {
	return m.analysisExists, m.analysisExistsErr
}

func (m *mockCacheAvailabilityRepository) GetAnalysisTestCount(_ context.Context, _ string) (int, error) {
	return 100, nil
}

func (m *mockCacheAvailabilityRepository) CheckSpecDocumentExistsByLanguage(_ context.Context, _ string, _ string) (bool, error) {
	return false, nil
}

func (m *mockCacheAvailabilityRepository) GetAvailableLanguages(_ context.Context, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetAvailableLanguagesByUser(_ context.Context, _ string, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetSpecDocumentByLanguage(_ context.Context, _ string, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetSpecDocumentByUser(_ context.Context, _ string, _ string, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetGenerationStatus(_ context.Context, _ string, _ string) (*entity.SpecGenerationStatus, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetGenerationStatusByLanguage(_ context.Context, _ string, _ string, _ string) (*entity.SpecGenerationStatus, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetSpecDocumentByVersion(_ context.Context, _ string, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetSpecDocumentByUserAndVersion(_ context.Context, _ string, _ string, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetVersionsByLanguage(_ context.Context, _ string, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) GetVersionsByUser(_ context.Context, _ string, _ string, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func (m *mockCacheAvailabilityRepository) HasPreviousSpecByLanguage(_ context.Context, _ string, _ string, _ string) (bool, error) {
	return false, nil
}

func (m *mockCacheAvailabilityRepository) GetLanguagesWithPreviousSpec(_ context.Context, _ string, _ string) ([]string, error) {
	return m.langsWithPreviousSpec, m.langsWithPreviousErr
}

func (m *mockCacheAvailabilityRepository) CheckCodebaseExists(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}
func (m *mockCacheAvailabilityRepository) GetSpecDocumentByRepository(_ context.Context, _, _, _, _ string) (*entity.RepoSpecDocument, error) {
	return nil, nil
}
func (m *mockCacheAvailabilityRepository) GetSpecDocumentByRepositoryAndVersion(_ context.Context, _, _, _, _ string, _ int) (*entity.RepoSpecDocument, error) {
	return nil, nil
}
func (m *mockCacheAvailabilityRepository) GetSpecDocumentByRepositoryAndDocumentId(_ context.Context, _, _, _, _ string) (*entity.RepoSpecDocument, error) {
	return nil, nil
}
func (m *mockCacheAvailabilityRepository) GetVersionHistoryByRepository(_ context.Context, _, _, _, _ string) ([]entity.RepoVersionInfo, error) {
	return nil, nil
}
func (m *mockCacheAvailabilityRepository) GetAvailableLanguagesByRepository(_ context.Context, _, _, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}
func (m *mockCacheAvailabilityRepository) GetCachePredictionData(_ context.Context, _, _, _ string) (*entity.CachePredictionData, error) {
	return nil, nil
}
func (m *mockCacheAvailabilityRepository) GetCurrentAnalysisTestCount(_ context.Context, _ string) (int, error) {
	return 0, nil
}

func TestGetCacheAvailabilityUseCase_Execute(t *testing.T) {
	t.Run("returns ErrUnauthorized when userID is empty", func(t *testing.T) {
		uc := NewGetCacheAvailabilityUseCase(&mockCacheAvailabilityRepository{})
		_, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrUnauthorized)
		}
	})

	t.Run("returns ErrInvalidAnalysisID when analysisID is empty", func(t *testing.T) {
		uc := NewGetCacheAvailabilityUseCase(&mockCacheAvailabilityRepository{})
		_, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidAnalysisID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidAnalysisID)
		}
	})

	t.Run("returns ErrAnalysisNotFound when analysis does not exist", func(t *testing.T) {
		mock := &mockCacheAvailabilityRepository{
			analysisExists: false,
		}
		uc := NewGetCacheAvailabilityUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
		})
		if !errors.Is(err, domain.ErrAnalysisNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrAnalysisNotFound)
		}
	})

	t.Run("propagates analysis check error", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &mockCacheAvailabilityRepository{
			analysisExistsErr: dbErr,
		}
		uc := NewGetCacheAvailabilityUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("returns cache availability for all supported languages", func(t *testing.T) {
		mock := &mockCacheAvailabilityRepository{
			analysisExists:        true,
			langsWithPreviousSpec: []string{"Korean", "English"},
		}
		uc := NewGetCacheAvailabilityUseCase(mock)
		result, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result == nil {
			t.Fatal("Execute() result is nil")
		}
		// Should have all supported languages
		if len(result.Languages) != len(entity.SupportedLanguages) {
			t.Errorf("Execute() Languages count = %d, want %d", len(result.Languages), len(entity.SupportedLanguages))
		}
		// Korean and English should be true
		if !result.Languages["Korean"] {
			t.Error("Execute() Languages[Korean] = false, want true")
		}
		if !result.Languages["English"] {
			t.Error("Execute() Languages[English] = false, want true")
		}
		// Japanese should be false (not in langsWithPreviousSpec)
		if result.Languages["Japanese"] {
			t.Error("Execute() Languages[Japanese] = true, want false")
		}
	})

	t.Run("returns all false when no previous specs exist", func(t *testing.T) {
		mock := &mockCacheAvailabilityRepository{
			analysisExists:        true,
			langsWithPreviousSpec: []string{},
		}
		uc := NewGetCacheAvailabilityUseCase(mock)
		result, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		// All languages should be false
		for lang, hasPrev := range result.Languages {
			if hasPrev {
				t.Errorf("Execute() Languages[%s] = true, want false", lang)
			}
		}
	})

	t.Run("propagates GetLanguagesWithPreviousSpec error", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &mockCacheAvailabilityRepository{
			analysisExists:       true,
			langsWithPreviousErr: dbErr,
		}
		uc := NewGetCacheAvailabilityUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCacheAvailabilityInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})
}
