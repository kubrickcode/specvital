package handler

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/kubrickcode/specvital/apps/web/src/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/src/backend/common/middleware"
	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/adapter/mapper"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/entity"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/domain/port"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/spec-view/usecase"
	subscriptionDomain "github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain"
	subscription "github.com/kubrickcode/specvital/apps/web/src/backend/modules/subscription/domain/entity"
)

type Handler struct {
	getCacheAvailability    *usecase.GetCacheAvailabilityUseCase
	getCachePrediction      *usecase.GetCachePredictionUseCase
	getGenerationStatus     *usecase.GetGenerationStatusUseCase
	getSpecByRepository     *usecase.GetSpecByRepositoryUseCase
	getSpecDocument         *usecase.GetSpecDocumentUseCase
	getVersionHistoryByRepo *usecase.GetVersionHistoryByRepositoryUseCase
	getVersions             *usecase.GetVersionsUseCase
	logger                  *logger.Logger
	requestGeneration       *usecase.RequestGenerationUseCase
	tierLookup              port.TierLookup
}

var _ api.SpecViewHandlers = (*Handler)(nil)

type HandlerConfig struct {
	GetCacheAvailability    *usecase.GetCacheAvailabilityUseCase
	GetCachePrediction      *usecase.GetCachePredictionUseCase
	GetGenerationStatus     *usecase.GetGenerationStatusUseCase
	GetSpecByRepository     *usecase.GetSpecByRepositoryUseCase
	GetSpecDocument         *usecase.GetSpecDocumentUseCase
	GetVersionHistoryByRepo *usecase.GetVersionHistoryByRepositoryUseCase
	GetVersions             *usecase.GetVersionsUseCase
	Logger                  *logger.Logger
	RequestGeneration       *usecase.RequestGenerationUseCase
	// TierLookup is optional. If nil, all requests use default queue.
	TierLookup port.TierLookup
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
	if cfg.GetVersions == nil {
		return nil, errors.New("GetVersions usecase is required")
	}
	if cfg.GetCacheAvailability == nil {
		return nil, errors.New("GetCacheAvailability usecase is required")
	}
	if cfg.GetCachePrediction == nil {
		return nil, errors.New("GetCachePrediction usecase is required")
	}
	if cfg.GetSpecByRepository == nil {
		return nil, errors.New("GetSpecByRepository usecase is required")
	}
	if cfg.GetVersionHistoryByRepo == nil {
		return nil, errors.New("GetVersionHistoryByRepo usecase is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("Logger is required")
	}

	return &Handler{
		getCacheAvailability:    cfg.GetCacheAvailability,
		getCachePrediction:      cfg.GetCachePrediction,
		getGenerationStatus:     cfg.GetGenerationStatus,
		getSpecByRepository:     cfg.GetSpecByRepository,
		getSpecDocument:         cfg.GetSpecDocument,
		getVersionHistoryByRepo: cfg.GetVersionHistoryByRepo,
		getVersions:             cfg.GetVersions,
		logger:                  cfg.Logger,
		requestGeneration:       cfg.RequestGeneration,
		tierLookup:              cfg.TierLookup,
	}, nil
}

func (h *Handler) GetSpecDocument(ctx context.Context, request api.GetSpecDocumentRequestObject) (api.GetSpecDocumentResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	analysisID := request.AnalysisID.String()

	language := ""
	if request.Params.Language != nil {
		language = string(*request.Params.Language)
	}

	version := 0
	if request.Params.Version != nil {
		version = *request.Params.Version
	}

	result, err := h.getSpecDocument.Execute(ctx, usecase.GetSpecDocumentInput{
		AnalysisID: analysisID,
		Language:   language,
		UserID:     userID,
		Version:    version,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetSpecDocument401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return api.GetSpecDocument403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("access denied to this resource"),
			}, nil
		case errors.Is(err, domain.ErrDocumentNotFound):
			return api.GetSpecDocument404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("spec document not found"),
			}, nil
		case errors.Is(err, domain.ErrInvalidAnalysisID):
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

	mode := entity.GenerationModeInitial
	if request.Body.GenerationMode != nil {
		mode = entity.GenerationMode(*request.Body.GenerationMode)
	}

	tier := h.lookupUserTier(ctx, userID)

	result, err := h.requestGeneration.Execute(ctx, usecase.RequestGenerationInput{
		AnalysisID: request.Body.AnalysisID.String(),
		Language:   language,
		Mode:       mode,
		Tier:       tier,
		UserID:     userID,
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
		case errors.Is(err, subscriptionDomain.ErrNoActiveSubscription):
			return api.RequestSpecGeneration403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("Active subscription required for spec generation. Please subscribe to a plan."),
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
	userID := middleware.GetUserID(ctx)
	analysisID := request.AnalysisID.String()

	input := usecase.GetGenerationStatusInput{
		AnalysisID: analysisID,
		UserID:     userID,
	}
	if request.Params.Language != nil {
		input.Language = *request.Params.Language
	}

	result, err := h.getGenerationStatus.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetSpecGenerationStatus401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return api.GetSpecGenerationStatus403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("access denied to this resource"),
			}, nil
		case errors.Is(err, domain.ErrInvalidAnalysisID):
			return api.GetSpecGenerationStatus404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("invalid analysis ID"),
			}, nil
		case errors.Is(err, domain.ErrInvalidLanguage):
			return api.GetSpecGenerationStatus400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid language"),
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

func (h *Handler) GetSpecVersions(ctx context.Context, request api.GetSpecVersionsRequestObject) (api.GetSpecVersionsResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	analysisID := request.AnalysisID.String()
	language := string(request.Params.Language)

	result, err := h.getVersions.Execute(ctx, usecase.GetVersionsInput{
		AnalysisID: analysisID,
		Language:   language,
		UserID:     userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetSpecVersions401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return api.GetSpecVersions403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("access denied to this resource"),
			}, nil
		case errors.Is(err, domain.ErrInvalidAnalysisID):
			return api.GetSpecVersions400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid analysis ID format"),
			}, nil
		case errors.Is(err, domain.ErrInvalidLanguage):
			return api.GetSpecVersions400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid language"),
			}, nil
		case errors.Is(err, domain.ErrDocumentNotFound):
			return api.GetSpecVersions404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("no versions found for this language"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get versions", "error", err)
		return api.GetSpecVersions500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get versions"),
		}, nil
	}

	resp := mapper.ToVersionHistoryResponse(&mapper.VersionHistoryInput{
		Language:      result.Language,
		LatestVersion: result.LatestVersion,
		Versions:      result.Versions,
	})
	return api.GetSpecVersions200JSONResponse(resp), nil
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

func (h *Handler) GetSpecCacheAvailability(ctx context.Context, request api.GetSpecCacheAvailabilityRequestObject) (api.GetSpecCacheAvailabilityResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	analysisID := request.AnalysisID.String()

	result, err := h.getCacheAvailability.Execute(ctx, usecase.GetCacheAvailabilityInput{
		AnalysisID: analysisID,
		UserID:     userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetSpecCacheAvailability401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrInvalidAnalysisID), errors.Is(err, domain.ErrAnalysisNotFound):
			return api.GetSpecCacheAvailability404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get cache availability", "error", err)
		return api.GetSpecCacheAvailability500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get cache availability"),
		}, nil
	}

	return api.GetSpecCacheAvailability200JSONResponse{
		Languages: result.Languages,
	}, nil
}

func (h *Handler) GetSpecCachePrediction(ctx context.Context, request api.GetSpecCachePredictionRequestObject) (api.GetSpecCachePredictionResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	analysisID := request.AnalysisID.String()
	language := string(request.Params.Language)

	result, err := h.getCachePrediction.Execute(ctx, usecase.GetCachePredictionInput{
		AnalysisID: analysisID,
		Language:   language,
		UserID:     userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetSpecCachePrediction401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrInvalidAnalysisID), errors.Is(err, domain.ErrAnalysisNotFound):
			return api.GetSpecCachePrediction404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found"),
			}, nil
		case errors.Is(err, domain.ErrInvalidLanguage):
			return api.GetSpecCachePrediction400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid language"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get cache prediction", "error", err)
		return api.GetSpecCachePrediction500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get cache prediction"),
		}, nil
	}

	return api.GetSpecCachePrediction200JSONResponse{
		TotalBehaviors:     result.TotalBehaviors,
		CacheableBehaviors: result.CacheableBehaviors,
		NewBehaviors:       result.NewBehaviors,
		EstimatedCost:      result.EstimatedCost,
	}, nil
}

func (h *Handler) lookupUserTier(ctx context.Context, userID string) subscription.PlanTier {
	if userID == "" || h.tierLookup == nil {
		return ""
	}
	tierStr, err := h.tierLookup.GetUserTier(ctx, userID)
	if err != nil {
		h.logger.Warn(ctx, "failed to lookup user tier, using default", "error", err)
		return ""
	}
	return subscription.PlanTier(tierStr)
}

func (h *Handler) GetSpecDocumentByRepository(ctx context.Context, request api.GetSpecDocumentByRepositoryRequestObject) (api.GetSpecDocumentByRepositoryResponseObject, error) {
	userID := middleware.GetUserID(ctx)

	language := ""
	if request.Params.Language != nil {
		language = string(*request.Params.Language)
	}

	version := 0
	if request.Params.Version != nil {
		version = *request.Params.Version
	}

	documentID := ""
	if request.Params.DocumentID != nil {
		documentID = request.Params.DocumentID.String()
	}

	result, err := h.getSpecByRepository.Execute(ctx, usecase.GetSpecByRepositoryInput{
		DocumentID: documentID,
		Language:   language,
		Name:       request.Repo,
		Owner:      request.Owner,
		UserID:     userID,
		Version:    version,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetSpecDocumentByRepository401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return api.GetSpecDocumentByRepository403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("access denied to this resource"),
			}, nil
		case errors.Is(err, domain.ErrDocumentNotFound), errors.Is(err, domain.ErrCodebaseNotFound):
			return api.GetSpecDocumentByRepository404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("spec document not found"),
			}, nil
		case errors.Is(err, domain.ErrInvalidRepository):
			return api.GetSpecDocumentByRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid repository"),
			}, nil
		case errors.Is(err, domain.ErrInvalidLanguage):
			return api.GetSpecDocumentByRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid language"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get spec document by repository", "error", err)
		return api.GetSpecDocumentByRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get spec document"),
		}, nil
	}

	resp, err := mapper.ToRepoSpecDocumentResponse(result.Document)
	if err != nil {
		h.logger.Error(ctx, "failed to marshal repo spec document", "error", err)
		return api.GetSpecDocumentByRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}

	return api.GetSpecDocumentByRepository200JSONResponse(resp), nil
}

func (h *Handler) GetVersionHistoryByRepository(ctx context.Context, request api.GetVersionHistoryByRepositoryRequestObject) (api.GetVersionHistoryByRepositoryResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	language := string(request.Params.Language)

	result, err := h.getVersionHistoryByRepo.Execute(ctx, usecase.GetVersionHistoryByRepositoryInput{
		Language: language,
		Name:     request.Repo,
		Owner:    request.Owner,
		UserID:   userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnauthorized):
			return api.GetVersionHistoryByRepository401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return api.GetVersionHistoryByRepository403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("access denied to this resource"),
			}, nil
		case errors.Is(err, domain.ErrCodebaseNotFound):
			return api.GetVersionHistoryByRepository404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("repository not found"),
			}, nil
		case errors.Is(err, domain.ErrInvalidRepository):
			return api.GetVersionHistoryByRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid repository"),
			}, nil
		case errors.Is(err, domain.ErrInvalidLanguage):
			return api.GetVersionHistoryByRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid language"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get version history by repository", "error", err)
		return api.GetVersionHistoryByRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get version history"),
		}, nil
	}

	resp, err := mapper.ToRepoVersionHistoryResponse(result.Language, result.Versions)
	if err != nil {
		h.logger.Error(ctx, "failed to marshal version history response", "error", err)
		return api.GetVersionHistoryByRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}

	return api.GetVersionHistoryByRepository200JSONResponse(resp), nil
}
