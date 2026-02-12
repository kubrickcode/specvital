package usecase

import "github.com/kubrickcode/specvital/apps/web/backend/modules/analyzer/domain/entity"

type AnalyzeResult struct {
	Analysis *entity.Analysis
	Progress *entity.AnalysisProgress
}
