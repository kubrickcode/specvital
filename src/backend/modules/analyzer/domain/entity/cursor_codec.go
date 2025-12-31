package entity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidCursor = errors.New("invalid cursor format")

type cursorPayload struct {
	AnalyzedAt time.Time `json:"at,omitempty"`
	ID         string    `json:"id"`
	Name       string    `json:"n,omitempty"`
	SortBy     string    `json:"sb"`
	TestCount  int       `json:"tc,omitempty"`
}

func EncodeCursor(c RepositoryCursor) string {
	payload := cursorPayload{
		AnalyzedAt: c.AnalyzedAt,
		ID:         c.ID,
		Name:       c.Name,
		SortBy:     string(c.SortBy),
		TestCount:  c.TestCount,
	}
	b, _ := json.Marshal(payload)
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCursor(encoded string, expectedSortBy SortBy) (*RepositoryCursor, error) {
	if encoded == "" {
		return nil, nil
	}

	b, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	var payload cursorPayload
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, ErrInvalidCursor
	}

	if SortBy(payload.SortBy) != expectedSortBy {
		return nil, ErrInvalidCursor
	}

	return &RepositoryCursor{
		AnalyzedAt: payload.AnalyzedAt,
		ID:         payload.ID,
		Name:       payload.Name,
		SortBy:     SortBy(payload.SortBy),
		TestCount:  payload.TestCount,
	}, nil
}
