package domain

import "errors"

var (
	ErrAnalysisNotFound      = errors.New("analysis not found")
	ErrAIProviderUnavailable = errors.New("AI provider unavailable")
	ErrConversionFailed      = errors.New("conversion failed")
	ErrInvalidLanguage       = errors.New("invalid conversion language")
	ErrRateLimited           = errors.New("rate limit exceeded")
)
