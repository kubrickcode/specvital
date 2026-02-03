package domain

import (
	"errors"
	"fmt"

	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
)

var (
	ErrInvalidCursor              = entity.ErrInvalidCursor
	ErrInvalidInput               = errors.New("invalid input")
	ErrNotFound                   = errors.New("analysis not found")
	ErrParserVersionNotConfigured = errors.New("parser_version not configured in system_config")
)

func WrapNotFound(owner, repo string) error {
	return fmt.Errorf("%s/%s: %w", owner, repo, ErrNotFound)
}
