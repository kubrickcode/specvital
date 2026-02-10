package postgres

import (
	"context"
	"errors"
	"testing"

	testdb "github.com/kubrickcode/specvital/apps/worker/src/internal/testutil/postgres"
)

func TestSystemConfigRepository_Upsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool, cleanup := testdb.SetupTestDB(t)
	defer cleanup()

	repo := NewSystemConfigRepository(pool)
	ctx := context.Background()

	t.Run("should insert new config", func(t *testing.T) {
		err := repo.Upsert(ctx, "test_key", "test_value")
		if err != nil {
			t.Fatalf("Upsert failed: %v", err)
		}

		value, err := repo.Get(ctx, "test_key")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if value != "test_value" {
			t.Errorf("expected 'test_value', got %q", value)
		}
	})

	t.Run("should update existing config", func(t *testing.T) {
		err := repo.Upsert(ctx, "update_key", "initial_value")
		if err != nil {
			t.Fatalf("initial Upsert failed: %v", err)
		}

		err = repo.Upsert(ctx, "update_key", "updated_value")
		if err != nil {
			t.Fatalf("update Upsert failed: %v", err)
		}

		value, err := repo.Get(ctx, "update_key")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if value != "updated_value" {
			t.Errorf("expected 'updated_value', got %q", value)
		}
	})

	t.Run("should reject empty key", func(t *testing.T) {
		err := repo.Upsert(ctx, "", "some_value")
		if err == nil {
			t.Error("expected error for empty key, got nil")
		}
	})
}

func TestSystemConfigRepository_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool, cleanup := testdb.SetupTestDB(t)
	defer cleanup()

	repo := NewSystemConfigRepository(pool)
	ctx := context.Background()

	t.Run("should return ErrConfigNotFound for non-existent key", func(t *testing.T) {
		_, err := repo.Get(ctx, "non_existent_key")
		if !errors.Is(err, ErrConfigNotFound) {
			t.Errorf("expected ErrConfigNotFound, got %v", err)
		}
	})

	t.Run("should return value for existing key", func(t *testing.T) {
		err := repo.Upsert(ctx, "existing_key", "existing_value")
		if err != nil {
			t.Fatalf("Upsert failed: %v", err)
		}

		value, err := repo.Get(ctx, "existing_key")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if value != "existing_value" {
			t.Errorf("expected 'existing_value', got %q", value)
		}
	})
}

func TestSystemConfigRepository_ParserVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pool, cleanup := testdb.SetupTestDB(t)
	defer cleanup()

	repo := NewSystemConfigRepository(pool)
	ctx := context.Background()

	t.Run("should store and retrieve parser_version", func(t *testing.T) {
		parserVersion := "v1.5.1-0.20260112121406-deacdda09e17"

		err := repo.Upsert(ctx, ConfigKeyParserVersion, parserVersion)
		if err != nil {
			t.Fatalf("Upsert parser_version failed: %v", err)
		}

		value, err := repo.Get(ctx, ConfigKeyParserVersion)
		if err != nil {
			t.Fatalf("Get parser_version failed: %v", err)
		}

		if value != parserVersion {
			t.Errorf("expected %q, got %q", parserVersion, value)
		}
	})

	t.Run("should update parser_version on re-upsert", func(t *testing.T) {
		oldVersion := "v1.5.0-0.20260101000000-abc123"
		newVersion := "v1.5.1-0.20260112121406-deacdda09e17"

		err := repo.Upsert(ctx, "parser_version_update", oldVersion)
		if err != nil {
			t.Fatalf("initial Upsert failed: %v", err)
		}

		err = repo.Upsert(ctx, "parser_version_update", newVersion)
		if err != nil {
			t.Fatalf("update Upsert failed: %v", err)
		}

		value, err := repo.Get(ctx, "parser_version_update")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if value != newVersion {
			t.Errorf("expected %q, got %q", newVersion, value)
		}
	})
}
