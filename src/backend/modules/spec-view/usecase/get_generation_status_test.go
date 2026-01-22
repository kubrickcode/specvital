package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type mockStatusRepository struct {
	status              *entity.SpecGenerationStatus
	statusErr           error
	statusByLanguage    *entity.SpecGenerationStatus
	statusByLanguageErr error

	calledAnalysisID string
	calledLanguage   string
}

func (m *mockStatusRepository) CheckAnalysisExists(_ context.Context, _ string) (bool, error) {
	return true, nil
}

func (m *mockStatusRepository) CheckSpecDocumentExistsByLanguage(_ context.Context, _ string, _ string) (bool, error) {
	return false, nil
}

func (m *mockStatusRepository) DeleteSpecDocumentByLanguage(_ context.Context, _ string, _ string) error {
	return nil
}

func (m *mockStatusRepository) GetAvailableLanguages(_ context.Context, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}

func (m *mockStatusRepository) GetSpecDocumentByLanguage(_ context.Context, _ string, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockStatusRepository) GetGenerationStatus(_ context.Context, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

func (m *mockStatusRepository) GetGenerationStatusByLanguage(_ context.Context, analysisID string, language string) (*entity.SpecGenerationStatus, error) {
	m.calledAnalysisID = analysisID
	m.calledLanguage = language
	return m.statusByLanguage, m.statusByLanguageErr
}

func (m *mockStatusRepository) GetSpecDocumentByVersion(_ context.Context, _ string, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}

func (m *mockStatusRepository) GetVersionsByLanguage(_ context.Context, _ string, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func TestGetGenerationStatusUseCase_Execute(t *testing.T) {
	t.Run("returns error when analysisID is empty", func(t *testing.T) {
		uc := NewGetGenerationStatusUseCase(&mockStatusRepository{})
		_, err := uc.Execute(context.Background(), GetGenerationStatusInput{})
		if !errors.Is(err, domain.ErrInvalidAnalysisID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidAnalysisID)
		}
	})

	t.Run("returns status when found without language", func(t *testing.T) {
		status := &entity.SpecGenerationStatus{
			AnalysisID: "analysis-1",
			Status:     entity.StatusRunning,
		}
		mock := &mockStatusRepository{status: status}
		uc := NewGetGenerationStatusUseCase(mock)
		result, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Status == nil {
			t.Fatal("Execute() Status is nil")
		}
		if result.Status.Status != entity.StatusRunning {
			t.Errorf("Execute() Status.Status = %v, want %v", result.Status.Status, entity.StatusRunning)
		}
	})

	t.Run("returns status when found with language", func(t *testing.T) {
		status := &entity.SpecGenerationStatus{
			AnalysisID: "analysis-1",
			Status:     entity.StatusPending,
		}
		mock := &mockStatusRepository{statusByLanguage: status}
		uc := NewGetGenerationStatusUseCase(mock)
		result, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if mock.calledAnalysisID != "analysis-1" {
			t.Errorf("Repository called with analysisID = %q, want %q", mock.calledAnalysisID, "analysis-1")
		}
		if mock.calledLanguage != "Korean" {
			t.Errorf("Repository called with language = %q, want %q", mock.calledLanguage, "Korean")
		}
		if result.Status == nil {
			t.Fatal("Execute() Status is nil")
		}
		if result.Status.Status != entity.StatusPending {
			t.Errorf("Execute() Status.Status = %v, want %v", result.Status.Status, entity.StatusPending)
		}
	})

	t.Run("returns not_found status when no status exists without language", func(t *testing.T) {
		mock := &mockStatusRepository{status: nil}
		uc := NewGetGenerationStatusUseCase(mock)
		result, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Status == nil {
			t.Fatal("Execute() Status is nil")
		}
		if result.Status.Status != entity.StatusNotFound {
			t.Errorf("Execute() Status.Status = %v, want %v", result.Status.Status, entity.StatusNotFound)
		}
	})

	t.Run("returns not_found status when no status exists with language", func(t *testing.T) {
		mock := &mockStatusRepository{statusByLanguage: nil}
		uc := NewGetGenerationStatusUseCase(mock)
		result, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
			Language:   "English",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if mock.calledAnalysisID != "analysis-1" {
			t.Errorf("Repository called with analysisID = %q, want %q", mock.calledAnalysisID, "analysis-1")
		}
		if mock.calledLanguage != "English" {
			t.Errorf("Repository called with language = %q, want %q", mock.calledLanguage, "English")
		}
		if result.Status == nil {
			t.Fatal("Execute() Status is nil")
		}
		if result.Status.Status != entity.StatusNotFound {
			t.Errorf("Execute() Status.Status = %v, want %v", result.Status.Status, entity.StatusNotFound)
		}
	})

	t.Run("propagates repository error without language", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &mockStatusRepository{statusErr: dbErr}
		uc := NewGetGenerationStatusUseCase(mock)
		_, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("propagates repository error with language", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &mockStatusRepository{statusByLanguageErr: dbErr}
		uc := NewGetGenerationStatusUseCase(mock)
		_, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("returns error when language is invalid", func(t *testing.T) {
		mock := &mockStatusRepository{}
		uc := NewGetGenerationStatusUseCase(mock)
		_, err := uc.Execute(context.Background(), GetGenerationStatusInput{
			AnalysisID: "analysis-1",
			Language:   "InvalidLanguage",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})
}
