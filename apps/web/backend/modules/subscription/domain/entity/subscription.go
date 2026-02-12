package entity

import "time"

type Plan struct {
	ID                   string
	Tier                 PlanTier
	SpecviewMonthlyLimit *int32
	AnalysisMonthlyLimit *int32
	RetentionDays        *int32
}

type PricingPlan struct {
	Tier                 PlanTier
	MonthlyPrice         *int32
	SpecviewMonthlyLimit *int32
	AnalysisMonthlyLimit *int32
	RetentionDays        *int32
}

func (p *Plan) IsUnlimited() bool {
	return p.Tier.IsUnlimited()
}

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
)

type Subscription struct {
	ID                 string
	UserID             string
	PlanID             string
	Status             SubscriptionStatus
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	CanceledAt         *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive
}

type SubscriptionWithPlan struct {
	Subscription
	Plan Plan
}
