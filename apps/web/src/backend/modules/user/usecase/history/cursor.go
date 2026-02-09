package history

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
)

type cursorJSON struct {
	ID   string    `json:"i"`
	Time time.Time `json:"t"`
}

func EncodeCursor(updatedAt time.Time, historyID string) (string, error) {
	data := cursorJSON{Time: updatedAt, ID: historyID}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeCursor(s string) (*entity.CursorData, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, domain.ErrInvalidCursor
	}
	var data cursorJSON
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, domain.ErrInvalidCursor
	}
	return &entity.CursorData{ID: data.ID, Time: data.Time}, nil
}
