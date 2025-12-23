package auth

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/mapper"
)

type BookmarkHandler struct {
	logger  *logger.Logger
	service BookmarkService
}

type BookmarkHandlerConfig struct {
	Logger  *logger.Logger
	Service BookmarkService
}

func NewBookmarkHandler(cfg *BookmarkHandlerConfig) (*BookmarkHandler, error) {
	if cfg == nil {
		return nil, errors.New("bookmark handler config is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Service == nil {
		return nil, errors.New("service is required")
	}
	return &BookmarkHandler{
		logger:  cfg.Logger,
		service: cfg.Service,
	}, nil
}

func (h *BookmarkHandler) AddBookmark(ctx context.Context, request api.AddBookmarkRequestObject) (api.AddBookmarkResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.AddBookmark401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	isBookmarked, err := h.service.AddBookmark(ctx, userID, request.Owner, request.Repo)
	if err != nil {
		if errors.Is(err, domain.ErrCodebaseNotFound) {
			return api.AddBookmark404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("repository not found"),
			}, nil
		}

		h.logger.Error(ctx, "failed to add bookmark", "error", err, "owner", request.Owner, "repo", request.Repo)
		return api.AddBookmark500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to add bookmark"),
		}, nil
	}

	return api.AddBookmark200JSONResponse(mapper.ToBookmarkResponse(isBookmarked)), nil
}

func (h *BookmarkHandler) GetUserBookmarks(ctx context.Context, _ api.GetUserBookmarksRequestObject) (api.GetUserBookmarksResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserBookmarks401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	repos, err := h.service.GetUserBookmarkedRepos(ctx, userID)
	if err != nil {
		h.logger.Error(ctx, "failed to get bookmarks", "error", err, "userID", userID)
		return api.GetUserBookmarks500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get bookmarks"),
		}, nil
	}

	return api.GetUserBookmarks200JSONResponse(mapper.ToBookmarkedRepositoriesResponse(repos)), nil
}

func (h *BookmarkHandler) RemoveBookmark(ctx context.Context, request api.RemoveBookmarkRequestObject) (api.RemoveBookmarkResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.RemoveBookmark401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	isBookmarked, err := h.service.RemoveBookmark(ctx, userID, request.Owner, request.Repo)
	if err != nil {
		if errors.Is(err, domain.ErrCodebaseNotFound) {
			return api.RemoveBookmark404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("repository not found"),
			}, nil
		}

		h.logger.Error(ctx, "failed to remove bookmark", "error", err, "owner", request.Owner, "repo", request.Repo)
		return api.RemoveBookmark500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to remove bookmark"),
		}, nil
	}

	return api.RemoveBookmark200JSONResponse(mapper.ToBookmarkResponse(isBookmarked)), nil
}
