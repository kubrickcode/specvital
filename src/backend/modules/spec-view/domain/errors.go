package domain

import "errors"

var (
	ErrAlreadyExists     = errors.New("spec document already exists")
	ErrAnalysisNotFound  = errors.New("analysis not found")
	ErrDocumentNotFound  = errors.New("spec document not found")
	ErrGenerationPending = errors.New("generation already pending")
	ErrGenerationRunning = errors.New("generation already running")
	ErrInvalidAnalysisID = errors.New("invalid analysis ID")
	ErrInvalidLanguage   = errors.New("invalid language")
	ErrNotAuthorized     = errors.New("not authorized to access this resource")
	ErrQuotaExceeded     = errors.New("quota exceeded")
)
