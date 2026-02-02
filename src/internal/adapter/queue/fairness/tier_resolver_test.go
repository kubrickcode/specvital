package fairness

import (
	"context"
	"testing"

	"github.com/specvital/worker/internal/infra/db"
)

func TestConvertDBTier(t *testing.T) {
	tests := []struct {
		name   string
		dbTier db.PlanTier
		want   PlanTier
	}{
		{"free tier", db.PlanTierFree, TierFree},
		{"pro tier", db.PlanTierPro, TierPro},
		{"pro plus tier", db.PlanTierProPlus, TierProPlus},
		{"enterprise tier", db.PlanTierEnterprise, TierEnterprise},
		{"unknown tier defaults to free", db.PlanTier("unknown"), TierFree},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertDBTier(tt.dbTier)
			if got != tt.want {
				t.Errorf("convertDBTier(%q) = %q, want %q", tt.dbTier, got, tt.want)
			}
		})
	}
}

func TestDBTierResolver_ResolveTier_EmptyUserID(t *testing.T) {
	resolver := NewDBTierResolver(nil)
	tier := resolver.ResolveTier(context.Background(), "")
	if tier != TierFree {
		t.Errorf("ResolveTier with empty userID = %q, want %q", tier, TierFree)
	}
}

func TestDBTierResolver_ResolveTier_InvalidUUID(t *testing.T) {
	resolver := NewDBTierResolver(nil)
	tier := resolver.ResolveTier(context.Background(), "not-a-valid-uuid")
	if tier != TierFree {
		t.Errorf("ResolveTier with invalid UUID = %q, want %q", tier, TierFree)
	}
}
