package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/spec-view/domain/entity"
)

type mockCachePredictionRepository struct {
	analysisExists      bool
	analysisExistsErr   error
	currentTestCount    int
	currentTestCountErr error
	predictionData      *entity.CachePredictionData
	predictionDataErr   error
}

func (m *mockCachePredictionRepository) CheckAnalysisExists(_ context.Context, _ string) (bool, error) {
	return m.analysisExists, m.analysisExistsErr
}

func (m *mockCachePredictionRepository) GetAnalysisTestCount(_ context.Context, _ string) (int, error) {
	return 100, nil
}

func (m *mockCachePredictionRepository) CheckSpecDocumentExistsByLanguage(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}

func (m *mockCachePredictionRepository) GetAvailableLanguages(_ context.Context, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetAvailableLanguagesByUser(_ context.Context, _, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByLanguage(_ context.Context, _, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByUser(_ context.Context, _, _, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetGenerationStatus(_ context.Context, _, _ string) (*entity.SpecGenerationStatus, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetGenerationStatusByLanguage(_ context.Context, _, _, _ string) (*entity.SpecGenerationStatus, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByVersion(_ context.Context, _, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByUserAndVersion(_ context.Context, _, _, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetVersionsByLanguage(_ context.Context, _, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetVersionsByUser(_ context.Context, _, _, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) HasPreviousSpecByLanguage(_ context.Context, _, _, _ string) (bool, error) {
	return false, nil
}

func (m *mockCachePredictionRepository) GetLanguagesWithPreviousSpec(_ context.Context, _, _ string) ([]string, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) CheckCodebaseExists(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByRepository(_ context.Context, _, _, _, _ string) (*entity.RepoSpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByRepositoryAndVersion(_ context.Context, _, _, _, _ string, _ int) (*entity.RepoSpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetSpecDocumentByRepositoryAndDocumentId(_ context.Context, _, _, _, _ string) (*entity.RepoSpecDocument, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetVersionHistoryByRepository(_ context.Context, _, _, _, _ string) ([]entity.RepoVersionInfo, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetAvailableLanguagesByRepository(_ context.Context, _, _, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockCachePredictionRepository) GetCachePredictionData(_ context.Context, _, _, _ string) (*entity.CachePredictionData, error) {
	return m.predictionData, m.predictionDataErr
}

func (m *mockCachePredictionRepository) GetCurrentAnalysisTestCount(_ context.Context, _ string) (int, error) {
	return m.currentTestCount, m.currentTestCountErr
}

func TestGetCachePredictionUseCase_Execute(t *testing.T) {
	t.Run("returns ErrUnauthorized when userID is empty", func(t *testing.T) {
		uc := NewGetCachePredictionUseCase(&mockCachePredictionRepository{})
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrUnauthorized)
		}
	})

	t.Run("returns ErrInvalidAnalysisID when analysisID is empty", func(t *testing.T) {
		uc := NewGetCachePredictionUseCase(&mockCachePredictionRepository{})
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			UserID:   "user-1",
			Language: "Korean",
		})
		if !errors.Is(err, domain.ErrInvalidAnalysisID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidAnalysisID)
		}
	})

	t.Run("returns ErrInvalidLanguage when language is empty", func(t *testing.T) {
		uc := NewGetCachePredictionUseCase(&mockCachePredictionRepository{})
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns ErrInvalidLanguage when language is invalid", func(t *testing.T) {
		uc := NewGetCachePredictionUseCase(&mockCachePredictionRepository{})
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "InvalidLang",
			UserID:     "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns ErrAnalysisNotFound when analysis does not exist", func(t *testing.T) {
		mock := &mockCachePredictionRepository{analysisExists: false}
		uc := NewGetCachePredictionUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-1",
		})
		if !errors.Is(err, domain.ErrAnalysisNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrAnalysisNotFound)
		}
	})

	t.Run("returns prediction data with cache hit", func(t *testing.T) {
		mock := &mockCachePredictionRepository{
			analysisExists:   true,
			currentTestCount: 100,
			predictionData: &entity.CachePredictionData{
				TotalBehaviors:     80,
				CacheableBehaviors: 70,
			},
		}
		uc := NewGetCachePredictionUseCase(mock)
		result, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.TotalBehaviors != 80 {
			t.Errorf("TotalBehaviors = %d, want 80", result.TotalBehaviors)
		}
		if result.CacheableBehaviors != 70 {
			t.Errorf("CacheableBehaviors = %d, want 70", result.CacheableBehaviors)
		}
		if result.NewBehaviors != 30 {
			t.Errorf("NewBehaviors = %d, want 30", result.NewBehaviors)
		}
		if result.EstimatedCost != 30 {
			t.Errorf("EstimatedCost = %d, want 30", result.EstimatedCost)
		}
	})

	t.Run("returns all new behaviors when no previous spec exists", func(t *testing.T) {
		mock := &mockCachePredictionRepository{
			analysisExists:   true,
			currentTestCount: 50,
			predictionData: &entity.CachePredictionData{
				TotalBehaviors:     0,
				CacheableBehaviors: 0,
			},
		}
		uc := NewGetCachePredictionUseCase(mock)
		result, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "English",
			UserID:     "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.TotalBehaviors != 0 {
			t.Errorf("TotalBehaviors = %d, want 0", result.TotalBehaviors)
		}
		if result.CacheableBehaviors != 0 {
			t.Errorf("CacheableBehaviors = %d, want 0", result.CacheableBehaviors)
		}
		if result.NewBehaviors != 50 {
			t.Errorf("NewBehaviors = %d, want 50", result.NewBehaviors)
		}
	})

	t.Run("handles repository error for analysis check", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mock := &mockCachePredictionRepository{
			analysisExistsErr: expectedErr,
		}
		uc := NewGetCachePredictionUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-1",
		})
		if !errors.Is(err, expectedErr) {
			t.Errorf("Execute() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("handles repository error for test count", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mock := &mockCachePredictionRepository{
			analysisExists:      true,
			currentTestCountErr: expectedErr,
		}
		uc := NewGetCachePredictionUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-1",
		})
		if !errors.Is(err, expectedErr) {
			t.Errorf("Execute() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("handles repository error for prediction data", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mock := &mockCachePredictionRepository{
			analysisExists:    true,
			currentTestCount:  100,
			predictionDataErr: expectedErr,
		}
		uc := NewGetCachePredictionUseCase(mock)
		_, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-1",
		})
		if !errors.Is(err, expectedErr) {
			t.Errorf("Execute() error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("clamps newBehaviors to zero when cacheable exceeds current tests", func(t *testing.T) {
		mock := &mockCachePredictionRepository{
			analysisExists:   true,
			currentTestCount: 50,
			predictionData: &entity.CachePredictionData{
				TotalBehaviors:     100,
				CacheableBehaviors: 80,
			},
		}
		uc := NewGetCachePredictionUseCase(mock)
		result, err := uc.Execute(context.Background(), GetCachePredictionInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.NewBehaviors != 0 {
			t.Errorf("NewBehaviors = %d, want 0 (clamped)", result.NewBehaviors)
		}
		if result.EstimatedCost != 0 {
			t.Errorf("EstimatedCost = %d, want 0", result.EstimatedCost)
		}
	})
}
