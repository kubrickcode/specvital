package domain

import (
	"errors"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
)

var (
	ErrInvalidCursor = entity.ErrInvalidCursor
	ErrNotFound      = errors.New("analysis not found")
)

func WrapNotFound(owner, repo string) error {
	return fmt.Errorf("%s/%s: %w", owner, repo, ErrNotFound)
}
