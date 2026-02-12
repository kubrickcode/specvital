package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInsufficientScope    = errors.New("github token lacks required permissions")
	ErrNoGitHubToken        = errors.New("no github token available")
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrUnauthorized         = errors.New("github token expired or invalid")
)

type RateLimitError struct {
	Limit     int
	Remaining int
	ResetAt   time.Time
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("github api rate limited (limit=%d, remaining=%d, reset=%s)",
		e.Limit, e.Remaining, e.ResetAt.Format(time.RFC3339))
}

func IsRateLimitError(err error) bool {
	var rateLimitErr *RateLimitError
	return errors.As(err, &rateLimitErr)
}
