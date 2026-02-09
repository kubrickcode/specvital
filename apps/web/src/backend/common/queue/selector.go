package queue

import (
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

// Queue Suffix constants
// Note: River queue names only allow letters, numbers, underscores, and hyphens (no colons)
const (
	SuffixPriority  = "_priority"
	SuffixDefault   = "_default"
	SuffixScheduled = "_scheduled"
)

// Base queue name constants
const (
	BaseQueueAnalysis = "analysis"
	BaseQueueSpecView = "specview"
)

// Full queue name constants for reference
const (
	QueueAnalysisPriority  = BaseQueueAnalysis + SuffixPriority
	QueueAnalysisDefault   = BaseQueueAnalysis + SuffixDefault
	QueueAnalysisScheduled = BaseQueueAnalysis + SuffixScheduled

	QueueSpecViewPriority  = BaseQueueSpecView + SuffixPriority
	QueueSpecViewDefault   = BaseQueueSpecView + SuffixDefault
	QueueSpecViewScheduled = BaseQueueSpecView + SuffixScheduled
)

// SelectQueue determines the target queue based on plan tier and scheduling status.
// Priority queue is for paying users (pro, pro_plus, enterprise).
// Default queue is for free tier users.
// Scheduled queue is for background/cron jobs.
func SelectQueue(baseQueue string, tier subscription.PlanTier, isScheduled bool) string {
	if isScheduled {
		return baseQueue + SuffixScheduled
	}

	switch tier {
	case subscription.PlanTierPro, subscription.PlanTierProPlus, subscription.PlanTierEnterprise:
		return baseQueue + SuffixPriority
	default:
		// Fallback for free tier, empty tier (anonymous users), or unrecognized values
		return baseQueue + SuffixDefault
	}
}

// SelectQueueForAnalysis is a convenience function for analysis queue selection.
func SelectQueueForAnalysis(tier subscription.PlanTier, isScheduled bool) string {
	return SelectQueue(BaseQueueAnalysis, tier, isScheduled)
}

// SelectQueueForSpecView is a convenience function for spec-view queue selection.
func SelectQueueForSpecView(tier subscription.PlanTier, isScheduled bool) string {
	return SelectQueue(BaseQueueSpecView, tier, isScheduled)
}
