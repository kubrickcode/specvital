package port

import "context"

type Subscriber interface {
	AssignDefaultPlan(ctx context.Context, userID string) error
}
