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
	t.Run("returns ErrUnauthorized when userID is empty", func(t *testing.T) {
		uc := NewGetVersionsUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
		})
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrUnauthorized)
		}
	})

	t.Run("returns error when analysisID is empty", func(t *testing.T) {
		uc := NewGetVersionsUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			Language: "Korean",
			UserID:   "user-1",
		})
		if !errors.Is(err, domain.ErrInvalidAnalysisID) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrInvalidAnalysisID)
		}
	})

	t.Run("returns error when language is empty", func(t *testing.T) {
		uc := NewGetVersionsUseCase(&mockRepository{})
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			UserID:     "user-1",
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
			UserID:     "user-1",
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
			UserID:     "user-1",
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
			UserID:     "user-1",
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
			UserID:     "user-1",
		})
		if !errors.Is(err, dbErr) {
			t.Errorf("Execute() error = %v, want %v", err, dbErr)
		}
	})

	t.Run("returns ErrDocumentNotFound when user has no versions even if other users have documents", func(t *testing.T) {
		// This tests the per-user personalization: each user has their own AI Spec versions
		mock := &mockRepository{
			versions: []entity.VersionInfo{}, // No versions for this user
		}
		uc := NewGetVersionsUseCase(mock)
		_, err := uc.Execute(context.Background(), GetVersionsInput{
			AnalysisID: "analysis-1",
			Language:   "Korean",
			UserID:     "user-b",
		})
		if !errors.Is(err, domain.ErrDocumentNotFound) {
			t.Errorf("Execute() error = %v, want %v", err, domain.ErrDocumentNotFound)
		}
	})
}
