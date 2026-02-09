package port

import "context"

type StateStore interface {
	Create(ctx context.Context) (string, error)
	Validate(ctx context.Context, state string) error
}
