package port

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
	Remaining(ctx context.Context, key string) int
}
