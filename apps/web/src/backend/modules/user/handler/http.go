package handler

import (
	"context"
	"regexp"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/user/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/user/domain"
	"github.com/specvital/web/src/backend/modules/user/domain/entity"
	activetaskuc "github.com/specvital/web/src/backend/modules/user/usecase/activetask"
	bookmarkuc "github.com/specvital/web/src/backend/modules/user/usecase/bookmark"
	historyuc "github.com/specvital/web/src/backend/modules/user/usecase/history"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

type Handler struct {
	addAnalyzedRepo    *historyuc.AddAnalyzedRepoUseCase
	addBookmark        *bookmarkuc.AddBookmarkUseCase
	getAnalyzedRepos   *historyuc.GetAnalyzedReposUseCase
	getBookmarks       *bookmarkuc.GetBookmarksUseCase
	getUserActiveTasks *activetaskuc.GetUserActiveTasksUseCase
	logger             *logger.Logger
	removeBookmark     *bookmarkuc.RemoveBookmarkUseCase
}

type HandlerConfig struct {
	AddAnalyzedRepo    *historyuc.AddAnalyzedRepoUseCase
	AddBookmark        *bookmarkuc.AddBookmarkUseCase
	GetAnalyzedRepos   *historyuc.GetAnalyzedReposUseCase
	GetBookmarks       *bookmarkuc.GetBookmarksUseCase
	GetUserActiveTasks *activetaskuc.GetUserActiveTasksUseCase
	Logger             *logger.Logger
	RemoveBookmark     *bookmarkuc.RemoveBookmarkUseCase
}

var _ api.BookmarkHandlers = (*Handler)(nil)
var _ api.AnalysisHistoryHandlers = (*Handler)(nil)
var _ api.UserActiveTasksHandlers = (*Handler)(nil)

func NewHandler(cfg *HandlerConfig) (*Handler, error) {
	if cfg == nil {
		return nil, errors.New("handler config is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.AddAnalyzedRepo == nil {
		return nil, errors.New("addAnalyzedRepo usecase is required")
	}
	if cfg.AddBookmark == nil {
		return nil, errors.New("addBookmark usecase is required")
	}
	if cfg.GetBookmarks == nil {
		return nil, errors.New("getBookmarks usecase is required")
	}
	if cfg.RemoveBookmark == nil {
		return nil, errors.New("removeBookmark usecase is required")
	}
	if cfg.GetAnalyzedRepos == nil {
		return nil, errors.New("getAnalyzedRepos usecase is required")
	}
	if cfg.GetUserActiveTasks == nil {
		return nil, errors.New("getUserActiveTasks usecase is required")
	}
	return &Handler{
		addAnalyzedRepo:    cfg.AddAnalyzedRepo,
		addBookmark:        cfg.AddBookmark,
		getAnalyzedRepos:   cfg.GetAnalyzedRepos,
		getBookmarks:       cfg.GetBookmarks,
		getUserActiveTasks: cfg.GetUserActiveTasks,
		logger:             cfg.Logger,
		removeBookmark:     cfg.RemoveBookmark,
	}, nil
}

func (h *Handler) AddUserAnalyzedRepository(
	ctx context.Context,
	request api.AddUserAnalyzedRepositoryRequestObject,
) (api.AddUserAnalyzedRepositoryResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.AddUserAnalyzedRepository401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	if err := validateOwnerRepo(request.Body.Owner, request.Body.Repo); err != nil {
		return api.AddUserAnalyzedRepository400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	output, err := h.addAnalyzedRepo.Execute(ctx, historyuc.AddAnalyzedRepoInput{
		Owner:  request.Body.Owner,
		Repo:   request.Body.Repo,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAnalysisNotFound) {
			return api.AddUserAnalyzedRepository404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found for this repository"),
			}, nil
		}
		h.logger.Error(ctx, "failed to add analyzed repository", "error", err, "owner", request.Body.Owner, "repo", request.Body.Repo)
		return api.AddUserAnalyzedRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to add analyzed repository"),
		}, nil
	}

	return api.AddUserAnalyzedRepository200JSONResponse{
		AnalysisID: output.AnalysisID,
		UpdatedAt:  output.UpdatedAt,
	}, nil
}

func (h *Handler) AddBookmark(ctx context.Context, request api.AddBookmarkRequestObject) (api.AddBookmarkResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.AddBookmark401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	output, err := h.addBookmark.Execute(ctx, bookmarkuc.AddBookmarkInput{
		Owner:  request.Owner,
		Repo:   request.Repo,
		UserID: userID,
	})
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

	return api.AddBookmark200JSONResponse(mapper.ToBookmarkResponse(output.IsBookmarked)), nil
}

func (h *Handler) GetUserBookmarks(ctx context.Context, _ api.GetUserBookmarksRequestObject) (api.GetUserBookmarksResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserBookmarks401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	repos, err := h.getBookmarks.Execute(ctx, bookmarkuc.GetBookmarksInput{
		UserID: userID,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to get bookmarks", "error", err, "userID", userID)
		return api.GetUserBookmarks500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get bookmarks"),
		}, nil
	}

	return api.GetUserBookmarks200JSONResponse(mapper.ToBookmarkedRepositoriesResponse(repos)), nil
}

func (h *Handler) RemoveBookmark(ctx context.Context, request api.RemoveBookmarkRequestObject) (api.RemoveBookmarkResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.RemoveBookmark401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	output, err := h.removeBookmark.Execute(ctx, bookmarkuc.RemoveBookmarkInput{
		Owner:  request.Owner,
		Repo:   request.Repo,
		UserID: userID,
	})
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

	return api.RemoveBookmark200JSONResponse(mapper.ToBookmarkResponse(output.IsBookmarked)), nil
}

func (h *Handler) GetUserAnalyzedRepositories(
	ctx context.Context,
	request api.GetUserAnalyzedRepositoriesRequestObject,
) (api.GetUserAnalyzedRepositoriesResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserAnalyzedRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	limit := defaultPageLimit
	if request.Params.Limit != nil && *request.Params.Limit > 0 && *request.Params.Limit <= maxPageLimit {
		limit = *request.Params.Limit
	}

	ownership := entity.OwnershipAll
	if request.Params.Ownership != nil {
		ownership = entity.ParseOwnershipFilter(string(*request.Params.Ownership))
	}

	result, err := h.getAnalyzedRepos.Execute(ctx, historyuc.GetAnalyzedReposInput{
		Cursor:    request.Params.Cursor,
		Limit:     limit,
		Ownership: ownership,
		UserID:    userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCursor) {
			return api.GetUserAnalyzedRepositories400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid cursor format"),
			}, nil
		}
		h.logger.Error(ctx, "failed to get analyzed repositories", "error", err, "userID", userID)
		return api.GetUserAnalyzedRepositories500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to retrieve analyzed repositories"),
		}, nil
	}

	return api.GetUserAnalyzedRepositories200JSONResponse(mapper.ToUserAnalyzedRepositoriesResponse(result)), nil
}

func (h *Handler) GetUserActiveTasks(
	ctx context.Context,
	_ api.GetUserActiveTasksRequestObject,
) (api.GetUserActiveTasksResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetUserActiveTasks401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	result, err := h.getUserActiveTasks.Execute(ctx, activetaskuc.GetUserActiveTasksInput{
		UserID: userID,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to get user active tasks", "error", err, "userID", userID)
		return api.GetUserActiveTasks500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to retrieve active tasks"),
		}, nil
	}

	return api.GetUserActiveTasks200JSONResponse(mapper.ToActiveTasksResponse(result.Tasks)), nil
}

func validateOwnerRepo(owner, repo string) error {
	if owner == "" || repo == "" {
		return errors.New("owner and repo are required")
	}
	if !validNamePattern.MatchString(owner) {
		return errors.New("invalid owner format")
	}
	if !validNamePattern.MatchString(repo) {
		return errors.New("invalid repo format")
	}
	return nil
}
