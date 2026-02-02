package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type repoMockRepository struct {
	codebaseExists     bool
	codebaseExistsErr  error
	repoDocument       *entity.RepoSpecDocument
	repoDocumentErr    error
	repoVersions       []entity.RepoVersionInfo
	repoVersionsErr    error
	availableLanguages []entity.AvailableLanguageInfo
	availableLangsErr  error

	calledDocumentID string
	calledLanguage   string
	calledVersion    int
	calledOwner      string
	calledName       string
}

func (m *repoMockRepository) CheckCodebaseExists(_ context.Context, owner, name string) (bool, error) {
	m.calledOwner = owner
	m.calledName = name
	return m.codebaseExists, m.codebaseExistsErr
}

func (m *repoMockRepository) GetSpecDocumentByRepository(_ context.Context, _, owner, name, language string) (*entity.RepoSpecDocument, error) {
	m.calledOwner = owner
	m.calledName = name
	m.calledLanguage = language
	return m.repoDocument, m.repoDocumentErr
}

func (m *repoMockRepository) GetSpecDocumentByRepositoryAndVersion(_ context.Context, _, owner, name, language string, version int) (*entity.RepoSpecDocument, error) {
	m.calledOwner = owner
	m.calledName = name
	m.calledLanguage = language
	m.calledVersion = version
	return m.repoDocument, m.repoDocumentErr
}

func (m *repoMockRepository) GetSpecDocumentByRepositoryAndDocumentId(_ context.Context, _, owner, name, documentID string) (*entity.RepoSpecDocument, error) {
	m.calledOwner = owner
	m.calledName = name
	m.calledDocumentID = documentID
	return m.repoDocument, m.repoDocumentErr
}

func (m *repoMockRepository) GetVersionHistoryByRepository(_ context.Context, _, owner, name, language string) ([]entity.RepoVersionInfo, error) {
	m.calledOwner = owner
	m.calledName = name
	m.calledLanguage = language
	return m.repoVersions, m.repoVersionsErr
}

func (m *repoMockRepository) GetAvailableLanguagesByRepository(_ context.Context, _, owner, name string) ([]entity.AvailableLanguageInfo, error) {
	m.calledOwner = owner
	m.calledName = name
	return m.availableLanguages, m.availableLangsErr
}

// Unused methods from port.SpecViewRepository interface
func (m *repoMockRepository) CheckAnalysisExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *repoMockRepository) GetAnalysisTestCount(_ context.Context, _ string) (int, error) {
	return 100, nil
}
func (m *repoMockRepository) CheckSpecDocumentExistsByLanguage(_ context.Context, _, _ string) (bool, error) {
	return false, nil
}
func (m *repoMockRepository) HasPreviousSpecByLanguage(_ context.Context, _, _, _ string) (bool, error) {
	return false, nil
}
func (m *repoMockRepository) GetLanguagesWithPreviousSpec(_ context.Context, _, _ string) ([]string, error) {
	return nil, nil
}
func (m *repoMockRepository) GetAvailableLanguages(_ context.Context, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}
func (m *repoMockRepository) GetAvailableLanguagesByUser(_ context.Context, _, _ string) ([]entity.AvailableLanguageInfo, error) {
	return nil, nil
}
func (m *repoMockRepository) GetGenerationStatus(_ context.Context, _, _ string) (*entity.SpecGenerationStatus, error) {
	return nil, nil
}
func (m *repoMockRepository) GetGenerationStatusByLanguage(_ context.Context, _, _, _ string) (*entity.SpecGenerationStatus, error) {
	return nil, nil
}
func (m *repoMockRepository) GetSpecDocumentByLanguage(_ context.Context, _, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}
func (m *repoMockRepository) GetSpecDocumentByUser(_ context.Context, _, _, _ string) (*entity.SpecDocument, error) {
	return nil, nil
}
func (m *repoMockRepository) GetSpecDocumentByVersion(_ context.Context, _, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}
func (m *repoMockRepository) GetSpecDocumentByUserAndVersion(_ context.Context, _, _, _ string, _ int) (*entity.SpecDocument, error) {
	return nil, nil
}
func (m *repoMockRepository) GetVersionsByLanguage(_ context.Context, _, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}
func (m *repoMockRepository) GetVersionsByUser(_ context.Context, _, _, _ string) ([]entity.VersionInfo, error) {
	return nil, nil
}

func TestGetSpecByRepositoryUseCase_Execute(t *testing.T) {
	t.Run("returns ErrUnauthorized when userID is empty", func(t *testing.T) {
		uc := NewGetSpecByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner: "facebook",
			Name:  "react",
		})
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrUnauthorized)
		}
	})

	t.Run("returns ErrInvalidRepository when owner is empty", func(t *testing.T) {
		uc := NewGetSpecByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Name:   "react",
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidRepository) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidRepository)
		}
	})

	t.Run("returns ErrInvalidRepository when name is empty", func(t *testing.T) {
		uc := NewGetSpecByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "facebook",
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidRepository) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidRepository)
		}
	})

	t.Run("returns ErrInvalidRepository when owner contains invalid characters", func(t *testing.T) {
		uc := NewGetSpecByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "../../../etc",
			Name:   "react",
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidRepository) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidRepository)
		}
	})

	t.Run("returns ErrInvalidRepository when name is too long", func(t *testing.T) {
		longName := ""
		for i := 0; i < 101; i++ {
			longName += "a"
		}
		uc := NewGetSpecByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "facebook",
			Name:   longName,
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidRepository) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidRepository)
		}
	})

	t.Run("defaults language to English when not specified", func(t *testing.T) {
		doc := &entity.RepoSpecDocument{
			ID:        "doc-1",
			Language:  "English",
			CommitSHA: "abc123",
		}
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   doc,
			availableLanguages: []entity.AvailableLanguageInfo{
				{Language: "English", LatestVersion: 1},
			},
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "facebook",
			Name:   "react",
			UserID: "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if mock.calledLanguage != "English" {
			t.Errorf("Execute() called with language = %q, want %q", mock.calledLanguage, "English")
		}
	})

	t.Run("returns ErrInvalidLanguage when language is invalid", func(t *testing.T) {
		uc := NewGetSpecByRepositoryUseCase(&repoMockRepository{codebaseExists: true})
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "InvalidLang",
			UserID:   "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns ErrCodebaseNotFound when codebase does not exist", func(t *testing.T) {
		mock := &repoMockRepository{codebaseExists: false}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "facebook",
			Name:   "react",
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrCodebaseNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrCodebaseNotFound)
		}
	})

	t.Run("returns document when found", func(t *testing.T) {
		doc := &entity.RepoSpecDocument{
			ID:         "doc-1",
			AnalysisID: "analysis-1",
			Language:   "Korean",
			Version:    2,
			CommitSHA:  "abc123",
			CreatedAt:  time.Now(),
		}
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   doc,
			availableLanguages: []entity.AvailableLanguageInfo{
				{Language: "Korean", LatestVersion: 2},
			},
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
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
		if result.Document.CommitSHA != "abc123" {
			t.Errorf("Execute() Document.CommitSHA = %q, want %q", result.Document.CommitSHA, "abc123")
		}
		if len(result.Document.AvailableLanguages) != 1 {
			t.Errorf("Execute() AvailableLanguages length = %d, want 1", len(result.Document.AvailableLanguages))
		}
	})

	t.Run("returns ErrDocumentNotFound when no document found", func(t *testing.T) {
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   nil,
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "facebook",
			Name:   "react",
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})

	t.Run("fetches specific version when version is specified", func(t *testing.T) {
		doc := &entity.RepoSpecDocument{
			ID:        "doc-1",
			Version:   1,
			CommitSHA: "old123",
		}
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   doc,
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
			Version:  1,
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if mock.calledVersion != 1 {
			t.Errorf("Execute() called with version = %d, want 1", mock.calledVersion)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &repoMockRepository{
			codebaseExists:  true,
			repoDocumentErr: dbErr,
			availableLanguages: []entity.AvailableLanguageInfo{
				{Language: "English", LatestVersion: 1},
			},
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:  "facebook",
			Name:   "react",
			UserID: "user-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("fetches document by documentID when specified", func(t *testing.T) {
		doc := &entity.RepoSpecDocument{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			Version:   1,
			Language:  "Korean",
			CommitSHA: "abc123",
		}
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   doc,
			availableLanguages: []entity.AvailableLanguageInfo{
				{Language: "Korean", LatestVersion: 1},
			},
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:      "facebook",
			Name:       "react",
			UserID:     "user-1",
			DocumentID: "550e8400-e29b-41d4-a716-446655440000",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if result.Document == nil {
			t.Fatal("Execute() Document is nil")
		}
		if result.Document.ID != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("Execute() Document.ID = %q, want %q", result.Document.ID, "550e8400-e29b-41d4-a716-446655440000")
		}
		if mock.calledDocumentID != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("Repository called with documentID = %q, want %q", mock.calledDocumentID, "550e8400-e29b-41d4-a716-446655440000")
		}
	})

	t.Run("returns ErrInvalidDocumentID when documentID format is invalid", func(t *testing.T) {
		mock := &repoMockRepository{
			codebaseExists: true,
			availableLanguages: []entity.AvailableLanguageInfo{
				{Language: "English", LatestVersion: 1},
			},
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:      "facebook",
			Name:       "react",
			UserID:     "user-1",
			DocumentID: "invalid-uuid-format",
		})
		if !errors.Is(err, domain.ErrInvalidDocumentID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidDocumentID)
		}
	})

	t.Run("returns ErrDocumentNotFound when documentID not found", func(t *testing.T) {
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   nil,
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:      "facebook",
			Name:       "react",
			UserID:     "user-1",
			DocumentID: "550e8400-e29b-41d4-a716-446655440000",
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})

	t.Run("documentID takes precedence over version parameter", func(t *testing.T) {
		doc := &entity.RepoSpecDocument{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			Version:   3,
			Language:  "English",
			CommitSHA: "abc123",
		}
		mock := &repoMockRepository{
			codebaseExists: true,
			repoDocument:   doc,
			availableLanguages: []entity.AvailableLanguageInfo{
				{Language: "English", LatestVersion: 5},
			},
		}
		uc := NewGetSpecByRepositoryUseCase(mock)
		result, err := uc.Execute(context.Background(), GetSpecByRepositoryInput{
			Owner:      "facebook",
			Name:       "react",
			UserID:     "user-1",
			DocumentID: "550e8400-e29b-41d4-a716-446655440000",
			Version:    1, // Should be ignored when documentID is provided
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		// Verify documentID was used, not version
		if mock.calledDocumentID != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("Repository called with documentID = %q, want %q", mock.calledDocumentID, "550e8400-e29b-41d4-a716-446655440000")
		}
		if mock.calledVersion != 0 {
			t.Errorf("Repository should not call version-based method, but calledVersion = %d", mock.calledVersion)
		}
		if result.Document.Version != 3 {
			t.Errorf("Execute() Document.Version = %d, want 3 (from documentID lookup)", result.Document.Version)
		}
	})
}
