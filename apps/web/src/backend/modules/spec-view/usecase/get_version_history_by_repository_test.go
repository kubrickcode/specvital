package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

func TestGetVersionHistoryByRepositoryUseCase_Execute(t *testing.T) {
	t.Run("returns ErrUnauthorized when userID is empty", func(t *testing.T) {
		uc := NewGetVersionHistoryByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
		})
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrUnauthorized)
		}
	})

	t.Run("returns ErrInvalidRepository when owner is empty", func(t *testing.T) {
		uc := NewGetVersionHistoryByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidRepository) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidRepository)
		}
	})

	t.Run("returns ErrInvalidRepository when name is empty", func(t *testing.T) {
		uc := NewGetVersionHistoryByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "facebook",
			Language: "Korean",
			UserID:   "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidRepository) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidRepository)
		}
	})

	t.Run("returns ErrInvalidLanguage when language is empty", func(t *testing.T) {
		uc := NewGetVersionHistoryByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:  "facebook",
			Name:   "react",
			UserID: "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns ErrInvalidLanguage when language is invalid", func(t *testing.T) {
		uc := NewGetVersionHistoryByRepositoryUseCase(&repoMockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
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
		uc := NewGetVersionHistoryByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
		})
		if !errors.Is(err, domain.ErrCodebaseNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrCodebaseNotFound)
		}
	})

	t.Run("returns versions when found", func(t *testing.T) {
		now := time.Now()
		versions := []entity.RepoVersionInfo{
			{ID: "doc-3", AnalysisID: "analysis-3", Version: 3, CommitSHA: "sha3", Language: "Korean", CreatedAt: now},
			{ID: "doc-2", AnalysisID: "analysis-2", Version: 2, CommitSHA: "sha2", Language: "Korean", CreatedAt: now.Add(-time.Hour)},
			{ID: "doc-1", AnalysisID: "analysis-1", Version: 1, CommitSHA: "sha1", Language: "Korean", CreatedAt: now.Add(-2 * time.Hour)},
		}
		mock := &repoMockRepository{
			codebaseExists: true,
			repoVersions:   versions,
		}
		uc := NewGetVersionHistoryByRepositoryUseCase(mock)
		result, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if len(result.Versions) != 3 {
			t.Errorf("Execute() Versions length = %d, want 3", len(result.Versions))
		}
		if result.Language != "Korean" {
			t.Errorf("Execute() Language = %q, want %q", result.Language, "Korean")
		}
		if result.Versions[0].CommitSHA != "sha3" {
			t.Errorf("Execute() first version CommitSHA = %q, want %q", result.Versions[0].CommitSHA, "sha3")
		}
		if mock.calledLanguage != "Korean" {
			t.Errorf("Repository called with language = %q, want %q", mock.calledLanguage, "Korean")
		}
	})

	t.Run("returns empty versions when no versions found", func(t *testing.T) {
		mock := &repoMockRepository{
			codebaseExists: true,
			repoVersions:   []entity.RepoVersionInfo{},
		}
		uc := NewGetVersionHistoryByRepositoryUseCase(mock)
		result, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if len(result.Versions) != 0 {
			t.Errorf("Execute() Versions length = %d, want 0", len(result.Versions))
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &repoMockRepository{
			codebaseExists:  true,
			repoVersionsErr: dbErr,
		}
		uc := NewGetVersionHistoryByRepositoryUseCase(mock)
		_, err := uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "facebook",
			Name:     "react",
			Language: "Korean",
			UserID:   "user-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("passes correct parameters to repository", func(t *testing.T) {
		mock := &repoMockRepository{
			codebaseExists: true,
			repoVersions:   []entity.RepoVersionInfo{},
		}
		uc := NewGetVersionHistoryByRepositoryUseCase(mock)
		_, _ = uc.Execute(context.Background(), GetVersionHistoryByRepositoryInput{
			Owner:    "specvital",
			Name:     "core",
			Language: "English",
			UserID:   "user-123",
		})
		if mock.calledOwner != "specvital" {
			t.Errorf("Repository called with owner = %q, want %q", mock.calledOwner, "specvital")
		}
		if mock.calledName != "core" {
			t.Errorf("Repository called with name = %q, want %q", mock.calledName, "core")
		}
		if mock.calledLanguage != "English" {
			t.Errorf("Repository called with language = %q, want %q", mock.calledLanguage, "English")
		}
	})
}
