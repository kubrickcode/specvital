package entity

type PlanTier string

const (
	PlanTierFree       PlanTier = "free"
	PlanTierPro        PlanTier = "pro"
	PlanTierProPlus    PlanTier = "pro_plus"
	PlanTierEnterprise PlanTier = "enterprise"
)

func (t PlanTier) String() string {
	return string(t)
}

func (t PlanTier) IsValid() bool {
	switch t {
	case PlanTierFree, PlanTierPro, PlanTierProPlus, PlanTierEnterprise:
		return true
	default:
		return false
	}
}

func (t PlanTier) IsUnlimited() bool {
	return t == PlanTierEnterprise
}

func ParsePlanTier(s string) (PlanTier, bool) {
	tier := PlanTier(s)
	return tier, tier.IsValid()
}
