package domain

import "github.com/cockroachdb/errors"

var (
	ErrAnalysisNotFound = errors.New("analysis not found")
	ErrCodebaseNotFound = errors.New("codebase not found")
	ErrInvalidCursor    = errors.New("invalid cursor format")
)
