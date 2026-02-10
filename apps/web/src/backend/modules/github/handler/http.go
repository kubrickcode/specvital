package handler

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/kubrickcode/specvital/apps/web/src/backend/common/logger"
	"github.com/kubrickcode/specvital/apps/web/src/backend/common/middleware"
	"github.com/kubrickcode/specvital/apps/web/src/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/adapter/mapper"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/domain"
	"github.com/kubrickcode/specvital/apps/web/src/backend/modules/github/usecase"
)

type Handler struct {
	listOrgRepos  *usecase.ListOrgReposUseCase
	listUserOrgs  *usecase.ListUserOrgsUseCase
	listUserRepos *usecase.ListUserReposUseCase
	logger        *logger.Logger
}

type HandlerConfig struct {
	ListOrgRepos  *usecase.ListOrgReposUseCase
	ListUserOrgs  *usecase.ListUserOrgsUseCase
	ListUserRepos *usecase.ListUserReposUseCase
	Logger        *logger.Logger
}

var _ api.GitHubHandlers = (*Handler)(nil)

func NewHandler(cfg *HandlerConfig) (*Handler, error) {
	if cfg == nil {
		return nil, errors.New("handler config is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.ListUserRepos == nil {
		return nil, errors.New("listUserRepos usecase is required")
	}
	if cfg.ListUserOrgs == nil {
		return nil, errors.New("listUserOrgs usecase is required")
	}
	if cfg.ListOrgRepos == nil {
		return nil, errors.New("listOrgRepos usecase is required")
	}
	return &Handler{
		listOrgRepos:  cfg.ListOrgRepos,
		listUserOrgs:  cfg.ListUserOrgs,
		listUserRepos: cfg.ListUserRepos,
		logger:        cfg.Logger,
	}, nil
}

func (h *Handler) GetUserGitHubRepositories(ctx context.Context, request api.GetUserGitHubRepositoriesRequestObject) (api.GetUserGitHubRepositoriesResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserGitHubRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	refresh := request.Params.Refresh != nil && *request.Params.Refresh
	repos, err := h.listUserRepos.Execute(ctx, usecase.ListUserReposInput{
		Refresh: refresh,
		UserID:  userID,
	})
	if err != nil {
		return h.handleReposError(ctx, err, "list user repositories")
	}

	return api.GetUserGitHubRepositories200JSONResponse{
		Data: mapper.ToGitHubRepositories(repos),
	}, nil
}

func (h *Handler) GetUserGitHubOrganizations(ctx context.Context, request api.GetUserGitHubOrganizationsRequestObject) (api.GetUserGitHubOrganizationsResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserGitHubOrganizations401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	refresh := request.Params.Refresh != nil && *request.Params.Refresh
	orgs, err := h.listUserOrgs.Execute(ctx, usecase.ListUserOrgsInput{
		Refresh: refresh,
		UserID:  userID,
	})
	if err != nil {
		return h.handleOrgsError(ctx, err, "list user organizations")
	}

	return api.GetUserGitHubOrganizations200JSONResponse{
		Data: mapper.ToGitHubOrganizations(orgs),
	}, nil
}

func (h *Handler) GetOrganizationRepositories(ctx context.Context, request api.GetOrganizationRepositoriesRequestObject) (api.GetOrganizationRepositoriesResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetOrganizationRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	refresh := request.Params.Refresh != nil && *request.Params.Refresh
	repos, err := h.listOrgRepos.Execute(ctx, usecase.ListOrgReposInput{
		OrgLogin: request.Org,
		Refresh:  refresh,
		UserID:   userID,
	})
	if err != nil {
		return h.handleOrgReposError(ctx, err, "list organization repositories")
	}

	return api.GetOrganizationRepositories200JSONResponse{
		Data: mapper.ToGitHubRepositories(repos),
	}, nil
}

func (h *Handler) classifyError(ctx context.Context, err error, operation string) (status int, message string) {
	switch {
	case errors.Is(err, domain.ErrNoGitHubToken):
		return http.StatusUnauthorized, "github token not available"
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, "github token expired or invalid"
	case errors.Is(err, domain.ErrInsufficientScope):
		return http.StatusUnauthorized, "github token lacks required permissions - please re-authenticate with additional scopes"
	case domain.IsRateLimitError(err):
		return http.StatusTooManyRequests, err.Error()
	case errors.Is(err, domain.ErrOrganizationNotFound):
		return http.StatusNotFound, "organization not found"
	default:
		h.logger.Error(ctx, "failed to "+operation, "error", err)
		return http.StatusInternalServerError, "internal error"
	}
}

func (h *Handler) handleReposError(ctx context.Context, err error, operation string) (api.GetUserGitHubRepositoriesResponseObject, error) {
	status, msg := h.classifyError(ctx, err, operation)
	switch status {
	case http.StatusUnauthorized:
		return api.GetUserGitHubRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized(msg),
		}, nil
	case http.StatusTooManyRequests:
		return api.GetUserGitHubRepositories429ApplicationProblemPlusJSONResponse{
			TooManyRequestsApplicationProblemPlusJSONResponse: api.NewTooManyRequests(msg),
		}, nil
	default:
		return api.GetUserGitHubRepositories500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to retrieve repositories"),
		}, nil
	}
}

func (h *Handler) handleOrgsError(ctx context.Context, err error, operation string) (api.GetUserGitHubOrganizationsResponseObject, error) {
	status, msg := h.classifyError(ctx, err, operation)
	switch status {
	case http.StatusUnauthorized:
		return api.GetUserGitHubOrganizations401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized(msg),
		}, nil
	case http.StatusTooManyRequests:
		return api.GetUserGitHubOrganizations429ApplicationProblemPlusJSONResponse{
			TooManyRequestsApplicationProblemPlusJSONResponse: api.NewTooManyRequests(msg),
		}, nil
	default:
		return api.GetUserGitHubOrganizations500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to retrieve organizations"),
		}, nil
	}
}

func (h *Handler) handleOrgReposError(ctx context.Context, err error, operation string) (api.GetOrganizationRepositoriesResponseObject, error) {
	status, msg := h.classifyError(ctx, err, operation)
	switch status {
	case http.StatusUnauthorized:
		return api.GetOrganizationRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized(msg),
		}, nil
	case http.StatusTooManyRequests:
		return api.GetOrganizationRepositories429ApplicationProblemPlusJSONResponse{
			TooManyRequestsApplicationProblemPlusJSONResponse: api.NewTooManyRequests(msg),
		}, nil
	case http.StatusNotFound:
		return api.GetOrganizationRepositories404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound(msg),
		}, nil
	default:
		return api.GetOrganizationRepositories500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to retrieve organization repositories"),
		}, nil
	}
}
