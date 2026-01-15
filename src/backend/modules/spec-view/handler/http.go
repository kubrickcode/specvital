package handler

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/spec-view/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/usecase"
)

type Handler struct {
	getGenerationStatus *usecase.GetGenerationStatusUseCase
	getSpecDocument     *usecase.GetSpecDocumentUseCase
	logger              *logger.Logger
	requestGeneration   *usecase.RequestGenerationUseCase
}

var _ api.SpecViewHandlers = (*Handler)(nil)

type HandlerConfig struct {
	GetGenerationStatus *usecase.GetGenerationStatusUseCase
	GetSpecDocument     *usecase.GetSpecDocumentUseCase
	Logger              *logger.Logger
	RequestGeneration   *usecase.RequestGenerationUseCase
}

func NewHandler(cfg *HandlerConfig) (*Handler, error) {
	if cfg.GetSpecDocument == nil {
		return nil, errors.New("GetSpecDocument usecase is required")
	}
	if cfg.RequestGeneration == nil {
		return nil, errors.New("RequestGeneration usecase is required")
	}
	if cfg.GetGenerationStatus == nil {
		return nil, errors.New("GetGenerationStatus usecase is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("Logger is required")
	}

	return &Handler{
		getGenerationStatus: cfg.GetGenerationStatus,
		getSpecDocument:     cfg.GetSpecDocument,
		logger:              cfg.Logger,
		requestGeneration:   cfg.RequestGeneration,
	}, nil
}

func (h *Handler) GetSpecDocument(ctx context.Context, request api.GetSpecDocumentRequestObject) (api.GetSpecDocumentResponseObject, error) {
	analysisID := request.AnalysisID.String()

	result, err := h.getSpecDocument.Execute(ctx, usecase.GetSpecDocumentInput{
		AnalysisID: analysisID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			return api.GetSpecDocument404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("spec document not found"),
			}, nil
		}
		if errors.Is(err, domain.ErrInvalidAnalysisID) {
			return api.GetSpecDocument404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("invalid analysis ID"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get spec document", "error", err)
		return api.GetSpecDocument500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get spec document"),
		}, nil
	}

	if result.Document != nil {
		data, err := mapper.ToSpecDocumentResponse(result.Document)
		if err != nil {
			h.logger.Error(ctx, "failed to marshal spec document", "error", err)
			return api.GetSpecDocument500ApplicationProblemPlusJSONResponse{
				InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
			}, nil
		}
		return newSpecDocument200Response(data), nil
	}

	data, err := mapper.ToGeneratingResponse(result.GenerationStatus)
	if err != nil {
		h.logger.Error(ctx, "failed to marshal generation status", "error", err)
		return api.GetSpecDocument500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}
	return newSpecDocument200Response(data), nil
}

func (h *Handler) RequestSpecGeneration(ctx context.Context, request api.RequestSpecGenerationRequestObject) (api.RequestSpecGenerationResponseObject, error) {
	if request.Body == nil {
		return api.RequestSpecGeneration400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("request body is required"),
		}, nil
	}

	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.RequestSpecGeneration401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required for spec generation"),
		}, nil
	}

	language := "English"
	if request.Body.Language != nil {
		language = string(*request.Body.Language)
	}

	isForce := false
	if request.Body.IsForceRegenerate != nil {
		isForce = *request.Body.IsForceRegenerate
	}

	result, err := h.requestGeneration.Execute(ctx, usecase.RequestGenerationInput{
		AnalysisID:        request.Body.AnalysisID.String(),
		IsForceRegenerate: isForce,
		Language:          language,
		UserID:            userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAnalysisNotFound):
			return api.RequestSpecGeneration404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found"),
			}, nil
		case errors.Is(err, domain.ErrAlreadyExists):
			return api.RequestSpecGeneration409ApplicationProblemPlusJSONResponse(
				api.ProblemDetail{
					Detail: "spec document already exists",
					Status: 409,
					Title:  "Conflict",
				},
			), nil
		case errors.Is(err, domain.ErrGenerationPending), errors.Is(err, domain.ErrGenerationRunning):
			return api.RequestSpecGeneration409ApplicationProblemPlusJSONResponse(
				api.ProblemDetail{
					Detail: "generation already in progress",
					Status: 409,
					Title:  "Conflict",
				},
			), nil
		case errors.Is(err, domain.ErrInvalidAnalysisID):
			return api.RequestSpecGeneration400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid analysis ID"),
			}, nil
		case errors.Is(err, domain.ErrQuotaExceeded):
			return api.RequestSpecGeneration429ApplicationProblemPlusJSONResponse{
				TooManyRequestsApplicationProblemPlusJSONResponse: api.TooManyRequestsApplicationProblemPlusJSONResponse{
					Detail: "Monthly SpecView generation quota exceeded. Upgrade your plan or wait for the next billing period.",
					Status: 429,
					Title:  "Quota Exceeded",
					Type:   ptr("quota_exceeded"),
				},
			}, nil
		}

		h.logger.Error(ctx, "failed to request spec generation", "error", err)
		return api.RequestSpecGeneration500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to request spec generation"),
		}, nil
	}

	resp, err := mapper.ToRequestGenerationResponse(&entity.SpecGenerationStatus{
		AnalysisID: result.AnalysisID,
		Status:     result.Status,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to map generation response", "error", err)
		return api.RequestSpecGeneration500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}

	return api.RequestSpecGeneration202JSONResponse(resp), nil
}

func (h *Handler) GetSpecGenerationStatus(ctx context.Context, request api.GetSpecGenerationStatusRequestObject) (api.GetSpecGenerationStatusResponseObject, error) {
	analysisID := request.AnalysisID.String()

	result, err := h.getGenerationStatus.Execute(ctx, usecase.GetGenerationStatusInput{
		AnalysisID: analysisID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAnalysisID) {
			return api.GetSpecGenerationStatus404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("invalid analysis ID"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get generation status", "error", err)
		return api.GetSpecGenerationStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get generation status"),
		}, nil
	}

	resp, err := mapper.ToGenerationStatusResponse(result.Status)
	if err != nil {
		h.logger.Error(ctx, "failed to map generation status response", "error", err)
		return api.GetSpecGenerationStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}

	return api.GetSpecGenerationStatus200JSONResponse(resp), nil
}

type specDocument200Response struct {
	data []byte
}

func newSpecDocument200Response(data []byte) api.GetSpecDocumentResponseObject {
	return specDocument200Response{data: data}
}

func (r specDocument200Response) VisitGetSpecDocumentResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := w.Write(r.data)
	return err
}

func ptr(s string) *string {
	return &s
}
