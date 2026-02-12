package domain

import "errors"

var (
	ErrNoActiveSubscription = errors.New("no active subscription found")
	ErrPlanNotFound         = errors.New("plan not found")
	ErrInvalidPlanTier      = errors.New("invalid plan tier")
)
