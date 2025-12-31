package entity_test

import (
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
)

func TestEncodeCursor(t *testing.T) {
	t.Parallel()

	cursor := entity.RepositoryCursor{
		AnalyzedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		Name:       "test-repo",
		SortBy:     entity.SortByRecent,
		TestCount:  150,
	}

	encoded := entity.EncodeCursor(cursor)
	if encoded == "" {
		t.Error("expected non-empty encoded string")
	}
}

func TestDecodeCursor_ValidCursor(t *testing.T) {
	t.Parallel()

	original := entity.RepositoryCursor{
		AnalyzedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		Name:       "test-repo",
		SortBy:     entity.SortByRecent,
		TestCount:  150,
	}

	encoded := entity.EncodeCursor(original)
	decoded, err := entity.DecodeCursor(encoded, entity.SortByRecent)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if decoded == nil {
		t.Fatal("expected decoded cursor")
	}
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, original.ID)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, original.Name)
	}
	if decoded.SortBy != original.SortBy {
		t.Errorf("SortBy mismatch: got %s, want %s", decoded.SortBy, original.SortBy)
	}
	if decoded.TestCount != original.TestCount {
		t.Errorf("TestCount mismatch: got %d, want %d", decoded.TestCount, original.TestCount)
	}
	if !decoded.AnalyzedAt.Equal(original.AnalyzedAt) {
		t.Errorf("AnalyzedAt mismatch: got %v, want %v", decoded.AnalyzedAt, original.AnalyzedAt)
	}
}

func TestDecodeCursor_EmptyString(t *testing.T) {
	t.Parallel()

	decoded, err := entity.DecodeCursor("", entity.SortByRecent)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if decoded != nil {
		t.Error("expected nil cursor for empty string")
	}
}

func TestDecodeCursor_SortByMismatch(t *testing.T) {
	t.Parallel()

	cursor := entity.RepositoryCursor{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
		SortBy: entity.SortByRecent,
	}
	encoded := entity.EncodeCursor(cursor)

	_, err := entity.DecodeCursor(encoded, entity.SortByName)

	if err == nil {
		t.Error("expected error for sortBy mismatch")
	}
	if err != entity.ErrInvalidCursor {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestDecodeCursor_InvalidBase64(t *testing.T) {
	t.Parallel()

	_, err := entity.DecodeCursor("not-valid-base64!!!", entity.SortByRecent)

	if err == nil {
		t.Error("expected error for invalid base64")
	}
	if err != entity.ErrInvalidCursor {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestDecodeCursor_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := entity.DecodeCursor("bm90LWpzb24", entity.SortByRecent) // "not-json" in base64

	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if err != entity.ErrInvalidCursor {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestEncodeDecode_AllSortBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		sortBy entity.SortBy
	}{
		{"recent", entity.SortByRecent},
		{"name", entity.SortByName},
		{"tests", entity.SortByTests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			original := entity.RepositoryCursor{
				AnalyzedAt: time.Now().UTC().Truncate(time.Second),
				ID:         "test-id",
				Name:       "test-repo",
				SortBy:     tt.sortBy,
				TestCount:  100,
			}

			encoded := entity.EncodeCursor(original)
			decoded, err := entity.DecodeCursor(encoded, tt.sortBy)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if decoded.SortBy != tt.sortBy {
				t.Errorf("SortBy mismatch: got %s, want %s", decoded.SortBy, tt.sortBy)
			}
		})
	}
}
