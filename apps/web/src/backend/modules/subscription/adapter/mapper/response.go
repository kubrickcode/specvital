package mapper

import (
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

func ToUserSubscriptionResponse(sub *entity.SubscriptionWithPlan) api.UserSubscriptionResponse {
	return api.UserSubscriptionResponse{
		Plan:               toPlanInfo(&sub.Plan),
		CurrentPeriodStart: sub.CurrentPeriodStart,
		CurrentPeriodEnd:   sub.CurrentPeriodEnd,
	}
}

func toPlanInfo(plan *entity.Plan) api.PlanInfo {
	var specviewLimit, analysisLimit, retentionDays *int

	if plan.SpecviewMonthlyLimit != nil {
		l := int(*plan.SpecviewMonthlyLimit)
		specviewLimit = &l
	}
	if plan.AnalysisMonthlyLimit != nil {
		l := int(*plan.AnalysisMonthlyLimit)
		analysisLimit = &l
	}
	if plan.RetentionDays != nil {
		d := int(*plan.RetentionDays)
		retentionDays = &d
	}

	return api.PlanInfo{
		Tier:                 api.PlanTier(plan.Tier.String()),
		SpecviewMonthlyLimit: specviewLimit,
		AnalysisMonthlyLimit: analysisLimit,
		RetentionDays:        retentionDays,
	}
}

func ToPricingResponse(plans []entity.PricingPlan) api.PricingResponse {
	data := make([]api.PricingPlan, len(plans))
	for i, p := range plans {
		data[i] = toPricingPlan(&p)
	}
	return api.PricingResponse{Data: data}
}

func toPricingPlan(p *entity.PricingPlan) api.PricingPlan {
	var price, specview, analysis, retention *int

	if p.MonthlyPrice != nil {
		v := int(*p.MonthlyPrice)
		price = &v
	}
	if p.SpecviewMonthlyLimit != nil {
		v := int(*p.SpecviewMonthlyLimit)
		specview = &v
	}
	if p.AnalysisMonthlyLimit != nil {
		v := int(*p.AnalysisMonthlyLimit)
		analysis = &v
	}
	if p.RetentionDays != nil {
		v := int(*p.RetentionDays)
		retention = &v
	}

	return api.PricingPlan{
		Tier:                 api.PlanTier(p.Tier.String()),
		MonthlyPrice:         price,
		SpecviewMonthlyLimit: specview,
		AnalysisMonthlyLimit: analysis,
		RetentionDays:        retention,
	}
}
