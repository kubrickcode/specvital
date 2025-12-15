package api

import "context"

type AnalyzerHandlers interface {
	AnalyzeRepository(ctx context.Context, request AnalyzeRepositoryRequestObject) (AnalyzeRepositoryResponseObject, error)
	GetAnalysisStatus(ctx context.Context, request GetAnalysisStatusRequestObject) (GetAnalysisStatusResponseObject, error)
}

type APIHandlers struct {
	analyzer AnalyzerHandlers
}

var _ StrictServerInterface = (*APIHandlers)(nil)

func NewAPIHandlers(analyzer AnalyzerHandlers) *APIHandlers {
	return &APIHandlers{
		analyzer: analyzer,
	}
}

func (h *APIHandlers) AnalyzeRepository(ctx context.Context, request AnalyzeRepositoryRequestObject) (AnalyzeRepositoryResponseObject, error) {
	return h.analyzer.AnalyzeRepository(ctx, request)
}

func (h *APIHandlers) GetAnalysisStatus(ctx context.Context, request GetAnalysisStatusRequestObject) (GetAnalysisStatusResponseObject, error) {
	return h.analyzer.GetAnalysisStatus(ctx, request)
}
