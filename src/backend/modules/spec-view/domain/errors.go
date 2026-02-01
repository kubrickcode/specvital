package domain

import "errors"

var (
	ErrAlreadyExists        = errors.New("spec document already exists")
	ErrAnalysisNotCompleted = errors.New("analysis not completed")
	ErrAnalysisNotFound     = errors.New("analysis not found")
	ErrCodebaseNotFound     = errors.New("codebase not found")
	ErrDocumentNotFound     = errors.New("spec document not found")
	ErrForbidden            = errors.New("access denied to this resource")
	ErrGenerationPending    = errors.New("generation already pending")
	ErrGenerationRunning    = errors.New("generation already running")
	ErrInvalidAnalysisID    = errors.New("invalid analysis ID")
	ErrInvalidLanguage      = errors.New("invalid language")
	ErrInvalidRepository    = errors.New("invalid repository (owner or name empty)")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrUnauthorized         = errors.New("authentication required")
)
