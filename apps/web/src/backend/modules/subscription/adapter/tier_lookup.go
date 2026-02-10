package adapter

import (
	"context"
	"errors"

	analyzerport "github.com/kubrickcode/specvital/apps/web/src/backend/modules/analyzer/domain/port"
	specviewport "github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/port"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/port"
)

var _ analyzerport.TierLookup = (*TierLookupAdapter)(nil)
var _ specviewport.TierLookup = (*TierLookupAdapter)(nil)

type TierLookupAdapter struct {
	repo port.SubscriptionRepository
}

func NewTierLookupAdapter(repo port.SubscriptionRepository) *TierLookupAdapter {
	return &TierLookupAdapter{repo: repo}
}

func (a *TierLookupAdapter) GetUserTier(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", nil
	}

	sub, err := a.repo.GetActiveSubscriptionWithPlan(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNoActiveSubscription) {
			return "", nil
		}
		return "", err
	}

	return string(sub.Plan.Tier), nil
}
