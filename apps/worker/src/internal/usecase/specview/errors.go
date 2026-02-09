package specview

import "errors"

var (
	ErrAIProcessingFailed    = errors.New("AI processing failed")
	ErrLoadInventoryFailed   = errors.New("failed to load test inventory")
	ErrPartialFeatureFailure = errors.New("partial feature conversion failure exceeds threshold")
	ErrSaveFailed            = errors.New("failed to save document")
)
