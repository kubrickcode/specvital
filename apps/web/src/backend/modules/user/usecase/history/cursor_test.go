package history

import (
	"testing"
	"time"

	"github.com/specvital/web/src/backend/modules/user/domain"
)

func TestEncodeCursor(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	historyID := "123e4567-e89b-12d3-a456-426614174000"

	cursor, err := EncodeCursor(testTime, historyID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cursor == "" {
		t.Error("expected non-empty cursor")
	}
}

func TestDecodeCursor_Valid(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	historyID := "123e4567-e89b-12d3-a456-426614174000"

	encoded, err := EncodeCursor(testTime, historyID)
	if err != nil {
		t.Fatalf("unexpected encode error: %v", err)
	}

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if !decoded.Time.Equal(testTime) {
		t.Errorf("expected time %v, got %v", testTime, decoded.Time)
	}
	if decoded.ID != historyID {
		t.Errorf("expected ID %s, got %s", historyID, decoded.ID)
	}
}

func TestDecodeCursor_InvalidBase64(t *testing.T) {
	_, err := DecodeCursor("not-valid-base64!!!")
	if err != domain.ErrInvalidCursor {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestDecodeCursor_InvalidJSON(t *testing.T) {
	_, err := DecodeCursor("aW52YWxpZC1qc29u")
	if err != domain.ErrInvalidCursor {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestDecodeCursor_EmptyString(t *testing.T) {
	_, err := DecodeCursor("")
	if err != domain.ErrInvalidCursor {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestCursorRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		time      time.Time
		historyID string
	}{
		{
			name:      "standard case",
			time:      time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC),
			historyID: "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:      "zero time",
			time:      time.Time{},
			historyID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:      "with nanoseconds",
			time:      time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
			historyID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := EncodeCursor(tt.time, tt.historyID)
			if err != nil {
				t.Fatalf("encode error: %v", err)
			}

			decoded, err := DecodeCursor(encoded)
			if err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if !decoded.Time.Equal(tt.time) {
				t.Errorf("time mismatch: expected %v, got %v", tt.time, decoded.Time)
			}
			if decoded.ID != tt.historyID {
				t.Errorf("ID mismatch: expected %s, got %s", tt.historyID, decoded.ID)
			}
		})
	}
}
