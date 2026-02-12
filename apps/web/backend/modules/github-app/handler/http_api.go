package handler

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/kubrickcode/specvital/apps/web/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/backend/common/middleware"
	"github.com/kubrickcode/specvital/apps/web/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/adapter/mapper"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/github-app/usecase"
)

type APIHandler struct {
	getInstallURL     *usecase.GetInstallURLUseCase
	listInstallations *usecase.ListInstallationsUseCase
	logger            *logger.Logger
}

type APIHandlerConfig struct {
	GetInstallURL     *usecase.GetInstallURLUseCase
	ListInstallations *usecase.ListInstallationsUseCase
	Logger            *logger.Logger
}

var _ api.GitHubAppHandlers = (*APIHandler)(nil)

func NewAPIHandler(cfg *APIHandlerConfig) (*APIHandler, error) {
	if cfg == nil {
		return nil, errors.New("handler config is required")
	}
	if cfg.GetInstallURL == nil {
		return nil, errors.New("GetInstallURL usecase is required")
	}
	if cfg.ListInstallations == nil {
		return nil, errors.New("ListInstallations usecase is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}

	return &APIHandler{
		getInstallURL:     cfg.GetInstallURL,
		listInstallations: cfg.ListInstallations,
		logger:            cfg.Logger,
	}, nil
}

func (h *APIHandler) GetGitHubAppInstallURL(ctx context.Context, _ api.GetGitHubAppInstallURLRequestObject) (api.GetGitHubAppInstallURLResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetGitHubAppInstallURL401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	installURL := h.getInstallURL.Execute()

	return api.GetGitHubAppInstallURL200JSONResponse{
		InstallURL: installURL,
	}, nil
}

func (h *APIHandler) GetUserGitHubAppInstallations(ctx context.Context, _ api.GetUserGitHubAppInstallationsRequestObject) (api.GetUserGitHubAppInstallationsResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserGitHubAppInstallations401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	installations, err := h.listInstallations.Execute(ctx, usecase.ListInstallationsInput{
		UserID: userID,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to list installations", "error", err)
		return api.GetUserGitHubAppInstallations500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to list installations"),
		}, nil
	}

	return api.GetUserGitHubAppInstallations200JSONResponse{
		Data: mapper.ToAPIInstallations(installations),
	}, nil
}
