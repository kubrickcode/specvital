package adapter

import (
	"log/slog"
	"os"

	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

type Config struct {
	DefaultPlanTier entity.PlanTier
}

func ConfigFromEnv() *Config {
	tier := os.Getenv("PROMO_DEFAULT_PLAN")
	if tier == "" {
		tier = string(entity.PlanTierFree)
	}

	planTier, valid := entity.ParsePlanTier(tier)
	if !valid {
		slog.Warn("invalid PROMO_DEFAULT_PLAN, falling back to free",
			"provided", tier,
			"valid_options", "free, pro, pro_plus, enterprise")
		planTier = entity.PlanTierFree
	}

	return &Config{
		DefaultPlanTier: planTier,
	}
}
