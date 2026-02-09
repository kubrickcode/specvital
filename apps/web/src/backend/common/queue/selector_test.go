package queue

import (
	"testing"

	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

func TestSelectQueue(t *testing.T) {
	tests := []struct {
		name        string
		baseQueue   string
		tier        subscription.PlanTier
		isScheduled bool
		want        string
	}{
		{
			name:        "pro user gets priority queue",
			baseQueue:   "analysis",
			tier:        subscription.PlanTierPro,
			isScheduled: false,
			want:        "analysis_priority",
		},
		{
			name:        "pro_plus user gets priority queue",
			baseQueue:   "analysis",
			tier:        subscription.PlanTierProPlus,
			isScheduled: false,
			want:        "analysis_priority",
		},
		{
			name:        "enterprise user gets priority queue",
			baseQueue:   "analysis",
			tier:        subscription.PlanTierEnterprise,
			isScheduled: false,
			want:        "analysis_priority",
		},
		{
			name:        "free user gets default queue",
			baseQueue:   "analysis",
			tier:        subscription.PlanTierFree,
			isScheduled: false,
			want:        "analysis_default",
		},
		{
			name:        "empty tier gets default queue",
			baseQueue:   "analysis",
			tier:        subscription.PlanTier(""),
			isScheduled: false,
			want:        "analysis_default",
		},
		{
			name:        "scheduled job gets scheduled queue regardless of tier",
			baseQueue:   "analysis",
			tier:        subscription.PlanTierPro,
			isScheduled: true,
			want:        "analysis_scheduled",
		},
		{
			name:        "scheduled free tier gets scheduled queue",
			baseQueue:   "analysis",
			tier:        subscription.PlanTierFree,
			isScheduled: true,
			want:        "analysis_scheduled",
		},
		{
			name:        "specview pro user gets priority queue",
			baseQueue:   "specview",
			tier:        subscription.PlanTierProPlus,
			isScheduled: false,
			want:        "specview_priority",
		},
		{
			name:        "specview free user gets default queue",
			baseQueue:   "specview",
			tier:        subscription.PlanTierFree,
			isScheduled: false,
			want:        "specview_default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SelectQueue(tt.baseQueue, tt.tier, tt.isScheduled)
			if got != tt.want {
				t.Errorf("SelectQueue(%q, %q, %v) = %q, want %q",
					tt.baseQueue, tt.tier, tt.isScheduled, got, tt.want)
			}
		})
	}
}

func TestSelectQueueForAnalysis(t *testing.T) {
	tests := []struct {
		name        string
		tier        subscription.PlanTier
		isScheduled bool
		want        string
	}{
		{
			name:        "pro user",
			tier:        subscription.PlanTierPro,
			isScheduled: false,
			want:        QueueAnalysisPriority,
		},
		{
			name:        "free user",
			tier:        subscription.PlanTierFree,
			isScheduled: false,
			want:        QueueAnalysisDefault,
		},
		{
			name:        "scheduled",
			tier:        subscription.PlanTierFree,
			isScheduled: true,
			want:        QueueAnalysisScheduled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SelectQueueForAnalysis(tt.tier, tt.isScheduled)
			if got != tt.want {
				t.Errorf("SelectQueueForAnalysis(%q, %v) = %q, want %q",
					tt.tier, tt.isScheduled, got, tt.want)
			}
		})
	}
}

func TestSelectQueueForSpecView(t *testing.T) {
	tests := []struct {
		name        string
		tier        subscription.PlanTier
		isScheduled bool
		want        string
	}{
		{
			name:        "pro_plus user",
			tier:        subscription.PlanTierProPlus,
			isScheduled: false,
			want:        QueueSpecViewPriority,
		},
		{
			name:        "free user",
			tier:        subscription.PlanTierFree,
			isScheduled: false,
			want:        QueueSpecViewDefault,
		},
		{
			name:        "scheduled",
			tier:        subscription.PlanTierPro,
			isScheduled: true,
			want:        QueueSpecViewScheduled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SelectQueueForSpecView(tt.tier, tt.isScheduled)
			if got != tt.want {
				t.Errorf("SelectQueueForSpecView(%q, %v) = %q, want %q",
					tt.tier, tt.isScheduled, got, tt.want)
			}
		})
	}
}

func TestQueueConstants(t *testing.T) {
	// Verify queue constant values match expected patterns
	tests := []struct {
		constant string
		want     string
	}{
		{QueueAnalysisPriority, "analysis_priority"},
		{QueueAnalysisDefault, "analysis_default"},
		{QueueAnalysisScheduled, "analysis_scheduled"},
		{QueueSpecViewPriority, "specview_priority"},
		{QueueSpecViewDefault, "specview_default"},
		{QueueSpecViewScheduled, "specview_scheduled"},
	}

	for _, tt := range tests {
		if tt.constant != tt.want {
			t.Errorf("Queue constant = %q, want %q", tt.constant, tt.want)
		}
	}
}
