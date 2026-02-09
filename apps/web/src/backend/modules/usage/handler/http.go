package handler

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	subscriptiondomain "github.com/specvital/web/src/backend/modules/subscription/domain"
	"github.com/specvital/web/src/backend/modules/usage/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/usage/domain/entity"
	"github.com/specvital/web/src/backend/modules/usage/usecase"
)

type Handler struct {
	checkQuota      *usecase.CheckQuotaUseCase
	getCurrentUsage *usecase.GetCurrentUsageUseCase
	logger          *logger.Logger
}

var _ api.UsageHandlers = (*Handler)(nil)

type HandlerConfig struct {
	CheckQuota      *usecase.CheckQuotaUseCase
	GetCurrentUsage *usecase.GetCurrentUsageUseCase
	Logger          *logger.Logger
}

func NewHandler(cfg *HandlerConfig) (*Handler, error) {
	if cfg.CheckQuota == nil {
		return nil, errors.New("CheckQuota usecase is required")
	}
	if cfg.GetCurrentUsage == nil {
		return nil, errors.New("GetCurrentUsage usecase is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("Logger is required")
	}

	return &Handler{
		checkQuota:      cfg.CheckQuota,
		getCurrentUsage: cfg.GetCurrentUsage,
		logger:          cfg.Logger,
	}, nil
}

func (h *Handler) CheckQuota(ctx context.Context, request api.CheckQuotaRequestObject) (api.CheckQuotaResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.CheckQuota401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	if request.Body == nil {
		return api.CheckQuota400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("request body is required"),
		}, nil
	}

	eventType := entity.EventType(request.Body.EventType)
	if !eventType.IsValid() {
		return api.CheckQuota400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid event type"),
		}, nil
	}

	amount := 1
	if request.Body.Amount != nil {
		amount = *request.Body.Amount
	}
	if amount <= 0 {
		return api.CheckQuota400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("amount must be at least 1"),
		}, nil
	}

	result, err := h.checkQuota.Execute(ctx, usecase.CheckQuotaInput{
		UserID:    userID,
		EventType: eventType,
		Amount:    amount,
	})
	if err != nil {
		if errors.Is(err, subscriptiondomain.ErrNoActiveSubscription) {
			return api.CheckQuota404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("no active subscription found"),
			}, nil
		}
		h.logger.Error(ctx, "failed to check quota", "error", err)
		return api.CheckQuota500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to check quota"),
		}, nil
	}

	return api.CheckQuota200JSONResponse(mapper.ToCheckQuotaResponse(result)), nil
}

func (h *Handler) GetCurrentUsage(ctx context.Context, _ api.GetCurrentUsageRequestObject) (api.GetCurrentUsageResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetCurrentUsage401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	result, err := h.getCurrentUsage.Execute(ctx, usecase.GetCurrentUsageInput{
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, subscriptiondomain.ErrNoActiveSubscription) {
			return api.GetCurrentUsage404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("no active subscription found"),
			}, nil
		}
		h.logger.Error(ctx, "failed to get current usage", "error", err)
		return api.GetCurrentUsage500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get current usage"),
		}, nil
	}

	return api.GetCurrentUsage200JSONResponse(mapper.ToUsageStatusResponse(result)), nil
}
