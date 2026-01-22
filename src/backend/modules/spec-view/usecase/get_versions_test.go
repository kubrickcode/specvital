package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

func TestGetVersionsUseCase_Execute(t *testing.T) {
	t.Run("returns error when analysisID is empty", func(t *testing.T) {
		uc := NewGetVersionsUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			Language: "Korean",
		})
		if !errors.Is(err, domain.ErrInvalidAnalysisID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidAnalysisID)
		}
	})

	t.Run("returns error when language is empty", func(t *testing.T) {
		uc := NewGetVersionsUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns error when language is invalid", func(t *testing.T) {
		uc := NewGetVersionsUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			Language:   "InvalidLang",
		})
		if !errors.Is(err, domain.ErrInvalidLanguage) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidLanguage)
		}
	})

	t.Run("returns versions when found", func(t *testing.T) {
		versions := []entity.VersionInfo{
			{Version: 3, CreatedAt: time.Now(), ModelID: "model-3"},
			{Version: 2, CreatedAt: time.Now().Add(-time.Hour), ModelID: "model-2"},
			{Version: 1, CreatedAt: time.Now().Add(-2 * time.Hour), ModelID: "model-1"},
		}
		mock := &mockRepository{versions: versions}
		uc := NewGetVersionsUseCase(mock)
		result, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if len(result.Versions) != 3 {
			t.Errorf("Execute() Versions length = %d, want 3", len(result.Versions))
		}
		if result.LatestVersion != 3 {
			t.Errorf("Execute() LatestVersion = %d, want 3", result.LatestVersion)
		}
		if result.Language != "Korean" {
			t.Errorf("Execute() Language = %q, want %q", result.Language, "Korean")
		}
		if mock.calledLanguage != "Korean" {
			t.Errorf("Repository called with language = %q, want %q", mock.calledLanguage, "Korean")
		}
	})

	t.Run("returns ErrDocumentNotFound when no versions found", func(t *testing.T) {
		mock := &mockRepository{versions: []entity.VersionInfo{}}
		uc := NewGetVersionsUseCase(mock)
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			Language:   "English",
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		dbErr := errors.New("database error")
		mock := &mockRepository{versionsErr: dbErr}
		uc := NewGetVersionsUseCase(mock)
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})
}
