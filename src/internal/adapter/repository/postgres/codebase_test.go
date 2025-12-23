package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/specvital/collector/internal/domain/analysis"
)

func TestCodebaseRepository_FindByExternalID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool, cleanup := setupTestDB(t)
	defer cleanup()

	analysisRepo := NewAnalysisRepository(pool)
	codebaseRepo := NewCodebaseRepository(pool)
	ctx := context.Background()

	t.Run("should find codebase by external ID", func(t *testing.T) {
		_, err := analysisRepo.CreateAnalysisRecord(ctx, analysis.CreateAnalysisRecordParams{
			Owner:          "find-ext-owner",
			Repo:           "find-ext-repo",
			CommitSHA:      "find-ext-sha",
			Branch:         "main",
			ExternalRepoID: "ext-id-12345",
		})
		if err != nil {
			t.Fatalf("CreateAnalysisRecord failed: %v", err)
		}

		codebase, err := codebaseRepo.FindByExternalID(ctx, "github.com", "ext-id-12345")
		if err != nil {
			t.Fatalf("FindByExternalID failed: %v", err)
		}

		if codebase.Owner != "find-ext-owner" {
			t.Errorf("expected owner 'find-ext-owner', got '%s'", codebase.Owner)
		}
		if codebase.Name != "find-ext-repo" {
			t.Errorf("expected name 'find-ext-repo', got '%s'", codebase.Name)
		}
		if codebase.ExternalRepoID != "ext-id-12345" {
			t.Errorf("expected external repo ID 'ext-id-12345', got '%s'", codebase.ExternalRepoID)
		}
	})

	t.Run("should return ErrCodebaseNotFound when not exists", func(t *testing.T) {
		_, err := codebaseRepo.FindByExternalID(ctx, "github.com", "non-existent-id")
		if !errors.Is(err, analysis.ErrCodebaseNotFound) {
			t.Errorf("expected ErrCodebaseNotFound, got %v", err)
		}
	})
}

func TestCodebaseRepository_FindByOwnerName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool, cleanup := setupTestDB(t)
	defer cleanup()

	analysisRepo := NewAnalysisRepository(pool)
	codebaseRepo := NewCodebaseRepository(pool)
	ctx := context.Background()

	t.Run("should find codebase by owner and name", func(t *testing.T) {
		_, err := analysisRepo.CreateAnalysisRecord(ctx, analysis.CreateAnalysisRecordParams{
			Owner:          "find-owner",
			Repo:           "find-repo",
			CommitSHA:      "find-sha",
			Branch:         "main",
			ExternalRepoID: "owner-name-id-1",
		})
		if err != nil {
			t.Fatalf("CreateAnalysisRecord failed: %v", err)
		}

		codebase, err := codebaseRepo.FindByOwnerName(ctx, "github.com", "find-owner", "find-repo")
		if err != nil {
			t.Fatalf("FindByOwnerName failed: %v", err)
		}

		if codebase.Owner != "find-owner" {
			t.Errorf("expected owner 'find-owner', got '%s'", codebase.Owner)
		}
		if codebase.Name != "find-repo" {
			t.Errorf("expected name 'find-repo', got '%s'", codebase.Name)
		}
		if codebase.ExternalRepoID != "owner-name-id-1" {
			t.Errorf("expected external repo ID 'owner-name-id-1', got '%s'", codebase.ExternalRepoID)
		}
	})

	t.Run("should return ErrCodebaseNotFound when not exists", func(t *testing.T) {
		_, err := codebaseRepo.FindByOwnerName(ctx, "github.com", "non-existent", "repo")
		if !errors.Is(err, analysis.ErrCodebaseNotFound) {
			t.Errorf("expected ErrCodebaseNotFound, got %v", err)
		}
	})

	t.Run("should not find stale codebase by owner name", func(t *testing.T) {
		_, err := analysisRepo.CreateAnalysisRecord(ctx, analysis.CreateAnalysisRecordParams{
			Owner:          "stale-owner",
			Repo:           "stale-repo",
			CommitSHA:      "stale-sha",
			Branch:         "main",
			ExternalRepoID: "stale-ext-id",
		})
		if err != nil {
			t.Fatalf("CreateAnalysisRecord failed: %v", err)
		}

		_, err = pool.Exec(ctx, "UPDATE codebases SET is_stale = true WHERE owner = 'stale-owner' AND name = 'stale-repo'")
		if err != nil {
			t.Fatalf("failed to mark codebase stale: %v", err)
		}

		_, err = codebaseRepo.FindByOwnerName(ctx, "github.com", "stale-owner", "stale-repo")
		if !errors.Is(err, analysis.ErrCodebaseNotFound) {
			t.Errorf("expected ErrCodebaseNotFound for stale codebase, got %v", err)
		}
	})
}
