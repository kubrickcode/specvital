package github

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrForbidden     = errors.New("access forbidden")
	ErrRateLimited   = errors.New("rate limit exceeded")
	ErrTreeTruncated = errors.New("repository tree truncated due to size")
	ErrInvalidInput  = errors.New("invalid input parameters")
)
