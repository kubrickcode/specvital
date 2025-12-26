package usecase

import "github.com/specvital/web/src/backend/modules/analyzer/domain/entity"

type AnalyzeResult struct {
	Analysis *entity.Analysis
	Progress *entity.AnalysisProgress
}
