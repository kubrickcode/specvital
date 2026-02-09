package port

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
)

var (
	ErrGitHubInsufficientScope = errors.New("github token lacks required permissions")
	ErrGitHubNotFound          = errors.New("github resource not found")
	ErrGitHubUnauthorized      = errors.New("github token expired or invalid")
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
