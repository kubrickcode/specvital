package mapper

import (
	"github.com/kubrickcode/specvital/apps/web/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/usage/usecase"
)

func ToCheckQuotaResponse(output *usecase.CheckQuotaOutput) api.CheckQuotaResponse {
	var limit *int
	if output.Limit != nil {
		l := int(*output.Limit)
		limit = &l
	}

	return api.CheckQuotaResponse{
		IsAllowed: output.IsAllowed,
		Used:      int(output.Used),
		Reserved:  int(output.Reserved),
		Limit:     limit,
		ResetAt:   output.ResetAt,
	}
}

func ToUsageStatusResponse(output *usecase.GetCurrentUsageOutput) api.UsageStatusResponse {
	return api.UsageStatusResponse{
		Specview: toUsageMetric(output.Specview),
		Analysis: toUsageMetric(output.Analysis),
		ResetAt:  output.ResetAt,
		Plan:     toPlanInfo(output.Plan),
	}
}

func toUsageMetric(metric usecase.UsageMetricOutput) api.UsageMetric {
	var limit *int
	if metric.Limit != nil {
		l := int(*metric.Limit)
		limit = &l
	}

	return api.UsageMetric{
		Used:       int(metric.Used),
		Reserved:   int(metric.Reserved),
		Limit:      limit,
		Percentage: metric.Percentage,
	}
}

func toPlanInfo(plan usecase.PlanInfoOutput) api.PlanInfo {
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
		Tier:                 api.PlanTier(plan.Tier),
		SpecviewMonthlyLimit: specviewLimit,
		AnalysisMonthlyLimit: analysisLimit,
		RetentionDays:        retentionDays,
	}
}
