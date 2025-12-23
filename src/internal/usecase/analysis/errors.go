package analysis

import "errors"

var (
	ErrCloneFailed              = errors.New("clone failed")
	ErrCodebaseResolutionFailed = errors.New("codebase resolution failed")
	ErrRaceConditionDetected    = errors.New("race condition detected: repository state changed during analysis")
	ErrSaveFailed               = errors.New("save failed")
	ErrScanFailed               = errors.New("scan failed")
	ErrTokenLookupFailed        = errors.New("token lookup failed")
)
