package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type mockRepository struct {
	document             *entity.SpecDocument
	documentErr          error
	documentByVersion    *entity.SpecDocument
	documentByVersionErr error
	status               *entity.SpecGenerationStatus
	statusErr            error
	analysisExists       bool
	analysisExistsErr    error
	specDocExists        bool
	specDocExistsErr     error
	availableLanguages   []entity.AvailableLanguageInfo
	availableLangsErr    error
	versions             []entity.VersionInfo
	versionsErr          error
	ownership            *entity.DocumentOwnership
	ownershipErr         error

	// Captured parameters for verification
	calledLanguage string
	calledVersion  int
}

func (m *mockRepository) CheckAnalysisExists(_ context.Context, _ string) (bool, error) {
	return m.analysisExists, m.analysisExistsErr
}

func (m *mockRepository) CheckSpecDocumentExistsByLanguage(_ context.Context, _ string, _ string) (bool, error) {
	return m.specDocExists, m.specDocExistsErr
}

func (m *mockRepository) CheckSpecDocumentOwnership(_ context.Context, _ string) (*entity.DocumentOwnership, error) {
	return m.ownership, m.ownershipErr
}

func (m *mockRepository) GetAvailableLanguages(_ context.Context, _ string) ([]entity.AvailableLanguageInfo, error) {
	return m.availableLanguages, m.availableLangsErr
}

func (m *mockRepository) GetAvailableLanguagesByUser(_ context.Context, _ string, _ string) ([]entity.AvailableLanguageInfo, error) {
	return m.availableLanguages, m.availableLangsErr
}

func (m *mockRepository) GetSpecDocumentByLanguage(_ context.Context, _ string, language string) (*entity.SpecDocument, error) {
	m.calledLanguage = language
	return m.document, m.documentErr
}

func (m *mockRepository) GetSpecDocumentByUser(_ context.Context, _ string, _ string, language string) (*entity.SpecDocument, error) {
	m.calledLanguage = language
	return m.document, m.documentErr
}

func (m *mockRepository) GetGenerationStatus(_ context.Context, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

func (m *mockRepository) GetGenerationStatusByLanguage(_ context.Context, _ string, _ string) (*entity.SpecGenerationStatus, error) {
	return m.status, m.statusErr
}

func (m *mockRepository) GetSpecDocumentByVersion(_ context.Context, _ string, language string, version int) (*entity.SpecDocument, error) {
	m.calledLanguage = language
	m.calledVersion = version
	return m.documentByVersion, m.documentByVersionErr
}

func (m *mockRepository) GetSpecDocumentByUserAndVersion(_ context.Context, _ string, _ string, language string, version int) (*entity.SpecDocument, error) {
	m.calledLanguage = language
	m.calledVersion = version
	return m.documentByVersion, m.documentByVersionErr
}

func (m *mockRepository) GetVersionsByLanguage(_ context.Context, _ string, language string) ([]entity.VersionInfo, error) {
	m.calledLanguage = language
	return m.versions, m.versionsErr
}

func (m *mockRepository) GetVersionsByUser(_ context.Context, _ string, _ string, language string) ([]entity.VersionInfo, error) {
	m.calledLanguage = language
	return m.versions, m.versionsErr
}

func TestGetSpecDocumentUseCase_Execute(t *testing.T) {
	t.Run("returns error when analysisID is empty", func(t *testing.T) {
		uc := NewGetSpecDocumentUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{})
		if !errors.Is(err, domain.ErrInvalidAnalysisID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidAnalysisID)
		}
	})

	t.Run("returns document when found without language filter", func(t *testing.T) {
		doc := &entity.SpecDocument{
			ID:         "doc-1",
			AnalysisID: "analysis-1",
			Language:   "English",
			CreatedAt:  time.Now(),
		}
		mock := &mockRepository{document: doc}
		uc := NewGetSpecDocumentUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
			Language:   "",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Document == nil {
			t.Fatal("Execute() Document is nil")
		}
		if result.Document.ID != "doc-1" {
			t.Errorf("Execute() Document.ID = %q, want %q", result.Document.ID, "doc-1")
		}
		if mock.calledLanguage != "" {
			t.Errorf("Repository called with language = %q, want empty", mock.calledLanguage)
		}
	})

	t.Run("returns document when found with language filter", func(t *testing.T) {
		doc := &entity.SpecDocument{
			ID:         "doc-2",
			AnalysisID: "analysis-1",
			Language:   "Korean",
			CreatedAt:  time.Now(),
		}
		mock := &mockRepository{document: doc}
		uc := NewGetSpecDocumentUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Document == nil {
			t.Fatal("Execute() Document is nil")
		}
		if result.Document.Language != "Korean" {
			t.Errorf("Execute() Document.Language = %q, want %q", result.Document.Language, "Korean")
		}
		if mock.calledLanguage != "Korean" {
			t.Errorf("Repository called with language = %q, want %q", mock.calledLanguage, "Korean")
		}
	})

	t.Run("returns generation status when document not found but generation in progress", func(t *testing.T) {
		status := &entity.SpecGenerationStatus{
			AnalysisID: "analysis-1",
			Status:     entity.StatusRunning,
		}
		uc := NewGetSpecDocumentUseCase(&mockRepository{
			document: nil,
			status:   status,
		})
		result, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.GenerationStatus == nil {
			t.Fatal("Execute() GenerationStatus is nil")
		}
		if result.GenerationStatus.Status != entity.StatusRunning {
			t.Errorf("Execute() GenerationStatus.Status = %v, want %v", result.GenerationStatus.Status, entity.StatusRunning)
		}
	})

	t.Run("returns ErrDocumentNotFound when no document and no generation status", func(t *testing.T) {
		uc := NewGetSpecDocumentUseCase(&mockRepository{
			document: nil,
			status:   nil,
		})
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})

	t.Run("returns ErrDocumentNotFound when status is not_found", func(t *testing.T) {
		status := &entity.SpecGenerationStatus{
			AnalysisID: "analysis-1",
			Status:     entity.StatusNotFound,
		}
		uc := NewGetSpecDocumentUseCase(&mockRepository{
			document: nil,
			status:   status,
		})
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		dbErr := errors.New("database error")
		uc := NewGetSpecDocumentUseCase(&mockRepository{documentErr: dbErr})
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("returns document with availableLanguages when found", func(t *testing.T) {
		availableLanguages := []entity.AvailableLanguageInfo{
			{Language: "English", LatestVersion: 2, CreatedAt: time.Now()},
			{Language: "Korean", LatestVersion: 1, CreatedAt: time.Now()},
		}
		doc := &entity.SpecDocument{
			ID:         "doc-1",
			AnalysisID: "analysis-1",
			Language:   "English",
			Version:    2,
			CreatedAt:  time.Now(),
		}
		mock := &mockRepository{
			document:           doc,
			availableLanguages: availableLanguages,
		}
		uc := NewGetSpecDocumentUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Document == nil {
			t.Fatal("Execute() Document is nil")
		}
		if len(result.Document.AvailableLanguages) != 2 {
			t.Errorf("Execute() Document.AvailableLanguages length = %d, want 2", len(result.Document.AvailableLanguages))
		}
		if result.Document.Version != 2 {
			t.Errorf("Execute() Document.Version = %d, want 2", result.Document.Version)
		}
	})

	t.Run("propagates availableLanguages error", func(t *testing.T) {
		doc := &entity.SpecDocument{
			ID:         "doc-1",
			AnalysisID: "analysis-1",
			Language:   "English",
			CreatedAt:  time.Now(),
		}
		langErr := errors.New("failed to fetch languages")
		mock := &mockRepository{
			document:          doc,
			availableLangsErr: langErr,
		}
		uc := NewGetSpecDocumentUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, langErr) {
			t.Errorf("Execute() error = %v, want %v", err, langErr)
		}
	})

	t.Run("returns error when version is specified without language", func(t *testing.T) {
		uc := NewGetSpecDocumentUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
			Version:    1,
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns error when version is specified with invalid language", func(t *testing.T) {
		uc := NewGetSpecDocumentUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
			Language:   "InvalidLang",
			Version:    1,
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns specific version when version is specified", func(t *testing.T) {
		doc := &entity.SpecDocument{
			ID:         "doc-v1",
			AnalysisID: "analysis-1",
			Language:   "Korean",
			Version:    1,
			CreatedAt:  time.Now(),
		}
		mock := &mockRepository{documentByVersion: doc}
		uc := NewGetSpecDocumentUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			Version:    1,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Document == nil {
			t.Fatal("Execute() Document is nil")
		}
		if result.Document.Version != 1 {
			t.Errorf("Execute() Document.Version = %d, want 1", result.Document.Version)
		}
		if mock.calledVersion != 1 {
			t.Errorf("Repository called with version = %d, want 1", mock.calledVersion)
		}
	})

	t.Run("returns ErrDocumentNotFound when specific version not found", func(t *testing.T) {
		mock := &mockRepository{documentByVersion: nil}
		uc := NewGetSpecDocumentUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecDocumentInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			Version:    999,
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})
}
