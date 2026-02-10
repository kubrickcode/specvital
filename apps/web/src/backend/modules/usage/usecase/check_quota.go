package usecase

import (
	"context"
	"time"

	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/port"
	usageentity "github.com/kubrickcode/specvital/apps/web/src/backend/modules/usage/domain/entity"
	usageport "github.com/kubrickcode/specvital/apps/web/src/backend/modules/usage/domain/port"
)

type CheckQuotaInput struct {
	UserID    string
	EventType usageentity.EventType
	Amount    int
}

type CheckQuotaOutput struct {
	IsAllowed bool
	Used      int64
	Reserved  int64
	Limit     *int32
	ResetAt   time.Time
}

type CheckQuotaUseCase struct {
	subscriptionRepo port.SubscriptionRepository
	usageRepo        usageport.UsageRepository
	reservationRepo  usageport.QuotaReservationRepository
}

func NewCheckQuotaUseCase(
	subscriptionRepo port.SubscriptionRepository,
	usageRepo usageport.UsageRepository,
	reservationRepo usageport.QuotaReservationRepository,
) *CheckQuotaUseCase {
	return &CheckQuotaUseCase{
		subscriptionRepo: subscriptionRepo,
		usageRepo:        usageRepo,
		reservationRepo:  reservationRepo,
	}
}

func (uc *CheckQuotaUseCase) Execute(ctx context.Context, input CheckQuotaInput) (*CheckQuotaOutput, error) {
	subscription, err := uc.subscriptionRepo.GetActiveSubscriptionWithPlan(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	periodStart := subscription.CurrentPeriodStart
	periodEnd := subscription.CurrentPeriodEnd

	used, err := uc.usageRepo.GetMonthlyUsage(ctx, input.UserID, input.EventType, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	reserved, err := uc.reservationRepo.GetTotalReservedAmount(ctx, input.UserID, input.EventType)
	if err != nil {
		return nil, err
	}

	var limit *int32
	isAllowed := true

	if subscription.Plan.IsUnlimited() {
		limit = nil
	} else {
		switch input.EventType {
		case usageentity.EventTypeSpecview:
			limit = subscription.Plan.SpecviewMonthlyLimit
		case usageentity.EventTypeAnalysis:
			limit = subscription.Plan.AnalysisMonthlyLimit
		}

		if limit != nil {
			effective := used + reserved
			if effective+int64(input.Amount) > int64(*limit) {
				isAllowed = false
			}
		}
	}

	return &CheckQuotaOutput{
		IsAllowed: isAllowed,
		Used:      used,
		Reserved:  reserved,
		Limit:     limit,
		ResetAt:   periodEnd,
	}, nil
}
