package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/common/ratelimit"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/entity"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
	"github.com/specvital/web/src/backend/modules/analyzer/usecase"
	subscription "github.com/specvital/web/src/backend/modules/subscription/domain/entity"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type Handler struct {
	analyzeRepository    *usecase.AnalyzeRepositoryUseCase
	anonymousRateLimiter *ratelimit.IPRateLimiter
	getAnalysis          *usecase.GetAnalysisUseCase
	getAnalysisHistory   *usecase.GetAnalysisHistoryUseCase
	getRepositoryStats   *usecase.GetRepositoryStatsUseCase
	getUpdateStatus      *usecase.GetUpdateStatusUseCase
	historyChecker       port.HistoryChecker
	listRepositoryCards  *usecase.ListRepositoryCardsUseCase
	logger               *logger.Logger
	reanalyzeRepository  *usecase.ReanalyzeRepositoryUseCase
	tierLookup           port.TierLookup
}

var _ api.AnalyzerHandlers = (*Handler)(nil)
var _ api.RepositoryHandlers = (*Handler)(nil)

func NewHandler(
	logger *logger.Logger,
	analyzeRepository *usecase.AnalyzeRepositoryUseCase,
	getAnalysis *usecase.GetAnalysisUseCase,
	getAnalysisHistory *usecase.GetAnalysisHistoryUseCase,
	listRepositoryCards *usecase.ListRepositoryCardsUseCase,
	getUpdateStatus *usecase.GetUpdateStatusUseCase,
	getRepositoryStats *usecase.GetRepositoryStatsUseCase,
	reanalyzeRepository *usecase.ReanalyzeRepositoryUseCase,
	historyChecker port.HistoryChecker,
	anonymousRateLimiter *ratelimit.IPRateLimiter,
	tierLookup port.TierLookup,
) *Handler {
	return &Handler{
		analyzeRepository:    analyzeRepository,
		anonymousRateLimiter: anonymousRateLimiter,
		getAnalysis:          getAnalysis,
		getAnalysisHistory:   getAnalysisHistory,
		getRepositoryStats:   getRepositoryStats,
		getUpdateStatus:      getUpdateStatus,
		historyChecker:       historyChecker,
		listRepositoryCards:  listRepositoryCards,
		logger:               logger,
		reanalyzeRepository:  reanalyzeRepository,
		tierLookup:           tierLookup,
	}
}

func (h *Handler) AnalyzeRepository(ctx context.Context, request api.AnalyzeRepositoryRequestObject) (api.AnalyzeRepositoryResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.AnalyzeRepository400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	userID := middleware.GetUserID(ctx)

	// Specific commit query - use getAnalysis usecase
	if request.Params.Commit != nil {
		if err := validateCommitSHA(*request.Params.Commit); err != nil {
			return api.AnalyzeRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
			}, nil
		}
		return h.analyzeRepositoryByCommit(ctx, owner, repo, *request.Params.Commit, userID, log)
	}

	if userID == "" && h.anonymousRateLimiter != nil {
		clientIP := middleware.GetClientIP(ctx)
		if clientIP == "" {
			clientIP = "unknown"
		}
		if !h.anonymousRateLimiter.Allow(clientIP) {
			log.Warn(ctx, "rate limit exceeded for anonymous user", "client_ip", clientIP)
			return api.AnalyzeRepository429ApplicationProblemPlusJSONResponse{
				TooManyRequestsApplicationProblemPlusJSONResponse: api.TooManyRequestsApplicationProblemPlusJSONResponse{
					Detail: "Rate limit exceeded. Please sign in for higher limits or try again later.",
					Status: 429,
					Title:  "Too Many Requests",
				},
			}, nil
		}
	}

	tier := h.lookupUserTier(ctx, log, userID)

	result, err := h.analyzeRepository.Execute(ctx, usecase.AnalyzeRepositoryInput{
		Owner:  owner,
		Repo:   repo,
		Tier:   tier,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, client.ErrRepoNotFound) {
			return api.AnalyzeRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("repository not found"),
			}, nil
		}
		if errors.Is(err, client.ErrForbidden) {
			return api.AnalyzeRepository400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("repository access forbidden"),
			}, nil
		}

		log.Error(ctx, "usecase error in AnalyzeRepository", "error", err)
		return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process analysis request"),
		}, nil
	}

	if result.Analysis == nil && result.Progress == nil {
		log.Error(ctx, "invalid result: neither analysis nor progress is set")
		return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("internal error"),
		}, nil
	}

	if result.Analysis != nil {
		opts := h.buildHistoryOptions(ctx, userID, owner, repo)
		response, mapErr := mapper.ToCompletedResponse(result.Analysis, opts)
		if mapErr != nil {
			log.Error(ctx, "failed to map completed response", "error", mapErr)
			return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
				InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
			}, nil
		}
		completed, err := response.AsCompletedResponse()
		if err != nil {
			log.Error(ctx, "failed to unmarshal completed response", "error", err)
			return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
				InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
			}, nil
		}
		return api.AnalyzeRepository200JSONResponse(completed), nil
	}

	response, mapErr := mapper.ToStatusResponse(result.Progress)
	if mapErr != nil {
		log.Error(ctx, "failed to map status response", "error", mapErr)
		return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}
	return newAnalyze202Response(response)
}

func (h *Handler) analyzeRepositoryByCommit(ctx context.Context, owner, repo, commitSHA, userID string, log *logger.Logger) (api.AnalyzeRepositoryResponseObject, error) {
	result, err := h.getAnalysis.Execute(ctx, usecase.GetAnalysisInput{
		CommitSHA: commitSHA,
		Owner:     owner,
		Repo:      repo,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.AnalyzeRepository404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found for commit"),
			}, nil
		}
		log.Error(ctx, "usecase error in AnalyzeRepository by commit", "error", err, "commit", commitSHA)
		return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get analysis"),
		}, nil
	}

	if result.Analysis == nil {
		return api.AnalyzeRepository404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found for commit"),
		}, nil
	}

	opts := h.buildHistoryOptions(ctx, userID, owner, repo)
	response, mapErr := mapper.ToCompletedResponse(result.Analysis, opts)
	if mapErr != nil {
		log.Error(ctx, "failed to map completed response", "error", mapErr)
		return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}
	completed, err := response.AsCompletedResponse()
	if err != nil {
		log.Error(ctx, "failed to unmarshal completed response", "error", err)
		return api.AnalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}
	return api.AnalyzeRepository200JSONResponse(completed), nil
}

func (h *Handler) GetAnalysisHistory(ctx context.Context, request api.GetAnalysisHistoryRequestObject) (api.GetAnalysisHistoryResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.GetAnalysisHistory400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	result, err := h.getAnalysisHistory.Execute(ctx, usecase.GetAnalysisHistoryInput{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			return api.GetAnalysisHistory400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
			}, nil
		}
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetAnalysisHistory404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("repository not found"),
			}, nil
		}
		log.Error(ctx, "usecase error in GetAnalysisHistory", "error", err)
		return api.GetAnalysisHistory500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get analysis history"),
		}, nil
	}

	if len(result.Items) == 0 {
		return api.GetAnalysisHistory404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("no completed analyses found"),
		}, nil
	}

	items := make([]mapper.AnalysisHistoryInput, len(result.Items))
	for i, item := range result.Items {
		items[i] = mapper.AnalysisHistoryInput{
			BranchName:  item.BranchName,
			CommitSHA:   item.CommitSHA,
			CommittedAt: item.CommittedAt,
			CompletedAt: item.CompletedAt,
			ID:          item.ID,
			TotalTests:  item.TotalTests,
		}
	}

	headCommitSHA := result.Items[0].CommitSHA

	response, err := mapper.ToAnalysisHistoryResponse(items, headCommitSHA)
	if err != nil {
		log.Error(ctx, "mapper error in GetAnalysisHistory", "error", err)
		return api.GetAnalysisHistory500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to build response"),
		}, nil
	}

	return api.GetAnalysisHistory200JSONResponse(response), nil
}

func (h *Handler) GetAnalysisStatus(ctx context.Context, request api.GetAnalysisStatusRequestObject) (api.GetAnalysisStatusResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)
	userID := middleware.GetUserID(ctx)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.GetAnalysisStatus400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	result, err := h.getAnalysis.Execute(ctx, usecase.GetAnalysisInput{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetAnalysisStatus404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found"),
			}, nil
		}
		log.Error(ctx, "usecase error in GetAnalysisStatus", "error", err)
		return api.GetAnalysisStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get analysis status"),
		}, nil
	}

	if result.Analysis == nil && result.Progress == nil {
		log.Error(ctx, "invalid result: neither analysis nor progress is set")
		return api.GetAnalysisStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("internal error"),
		}, nil
	}

	if result.Analysis != nil {
		opts := h.buildHistoryOptions(ctx, userID, owner, repo)
		response, mapErr := mapper.ToCompletedResponse(result.Analysis, opts)
		if mapErr != nil {
			log.Error(ctx, "failed to map completed response", "error", mapErr)
			return api.GetAnalysisStatus500ApplicationProblemPlusJSONResponse{
				InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
			}, nil
		}
		return newStatus200Response(response)
	}

	response, mapErr := mapper.ToStatusResponse(result.Progress)
	if mapErr != nil {
		log.Error(ctx, "failed to map status response", "error", mapErr)
		return api.GetAnalysisStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
		}, nil
	}
	return newStatus200Response(response)
}

func (h *Handler) GetRecentRepositories(ctx context.Context, request api.GetRecentRepositoriesRequestObject) (api.GetRecentRepositoriesResponseObject, error) {
	params := request.Params
	userID := middleware.GetUserID(ctx)

	if err := validateRecentRepositoriesAuth(userID, params.View, params.Ownership); err != nil {
		return api.GetRecentRepositories401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized(err.Error()),
		}, nil
	}

	input := usecase.ListRepositoryCardsPaginatedInput{
		UserID: userID,
	}

	if params.Cursor != nil {
		input.Cursor = *params.Cursor
	}
	if params.Limit != nil {
		input.Limit = *params.Limit
	}
	if params.SortBy != nil {
		input.SortBy = entity.ParseSortBy(string(*params.SortBy))
	}
	if params.SortOrder != nil {
		input.SortOrder = entity.ParseSortOrder(string(*params.SortOrder))
	}
	if params.View != nil {
		input.View = entity.ParseViewFilter(string(*params.View))
	}
	if params.Ownership != nil {
		input.Ownership = entity.ParseOwnershipFilter(string(*params.Ownership))
	}

	result, err := h.listRepositoryCards.ExecutePaginated(ctx, input)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCursor) {
			return api.GetRecentRepositories400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid cursor"),
			}, nil
		}
		h.logger.Error(ctx, "failed to get recent repositories", "error", err)
		return api.GetRecentRepositories500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get recent repositories"),
		}, nil
	}

	return api.GetRecentRepositories200JSONResponse(mapper.ToPaginatedRepositoriesResponse(result)), nil
}

func (h *Handler) GetRepositoryStats(ctx context.Context, _ api.GetRepositoryStatsRequestObject) (api.GetRepositoryStatsResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.GetRepositoryStats401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	stats, err := h.getRepositoryStats.Execute(ctx, usecase.GetRepositoryStatsInput{
		UserID: userID,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to get repository stats", "error", err)
		return api.GetRepositoryStats500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get repository stats"),
		}, nil
	}

	return api.GetRepositoryStats200JSONResponse(mapper.ToRepositoryStatsResponse(stats)), nil
}

func (h *Handler) GetUpdateStatus(ctx context.Context, request api.GetUpdateStatusRequestObject) (api.GetUpdateStatusResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.GetUpdateStatus400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	userID := middleware.GetUserID(ctx)
	result, err := h.getUpdateStatus.Execute(ctx, usecase.GetUpdateStatusInput{
		Owner:  owner,
		Repo:   repo,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetUpdateStatus404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("repository not found"),
			}, nil
		}
		log.Error(ctx, "failed to get update status", "error", err)
		return api.GetUpdateStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get update status"),
		}, nil
	}

	return api.GetUpdateStatus200JSONResponse(mapper.ToUpdateStatusResponse(result)), nil
}

func (h *Handler) ReanalyzeRepository(ctx context.Context, request api.ReanalyzeRepositoryRequestObject) (api.ReanalyzeRepositoryResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.ReanalyzeRepository400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	userID := middleware.GetUserID(ctx)
	tier := h.lookupUserTier(ctx, log, userID)

	result, err := h.reanalyzeRepository.Execute(ctx, usecase.ReanalyzeRepositoryInput{
		Owner:  owner,
		Repo:   repo,
		Tier:   tier,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.ReanalyzeRepository404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("repository not found"),
			}, nil
		}
		log.Error(ctx, "failed to trigger reanalysis", "error", err)
		return api.ReanalyzeRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to trigger reanalysis"),
		}, nil
	}

	if result.Progress != nil {
		response, mapErr := mapper.ToStatusResponse(result.Progress)
		if mapErr != nil {
			log.Error(ctx, "failed to map status response", "error", mapErr)
			return api.ReanalyzeRepository500ApplicationProblemPlusJSONResponse{
				InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to process response"),
			}, nil
		}
		return newReanalyze202Response(response)
	}

	return api.ReanalyzeRepository500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("unexpected response state"),
	}, nil
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

func validateCommitSHA(sha string) error {
	if len(sha) < 7 || len(sha) > 40 {
		return errors.New("commit SHA must be between 7 and 40 characters")
	}
	for _, c := range sha {
		if !((c >= 'a' && c <= 'f') || (c >= '0' && c <= '9')) {
			return errors.New("commit SHA must contain only lowercase hex characters")
		}
	}
	return nil
}

func validateRecentRepositoriesAuth(userID string, view *api.ViewFilterParam, ownership *api.OwnershipFilterParam) error {
	if userID != "" {
		return nil
	}

	if view != nil && *view == api.ViewFilterParamMy {
		return errors.New("authentication required to view your analyzed repositories")
	}

	if ownership != nil && *ownership != api.OwnershipFilterParamAll {
		return errors.New("authentication required to filter by ownership")
	}

	return nil
}

type analyze202Response struct {
	union json.RawMessage
}

func newAnalyze202Response(response api.AnalysisResponse) (api.AnalyzeRepositoryResponseObject, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return analyze202Response{union: data}, nil
}

func (r analyze202Response) VisitAnalyzeRepositoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write(r.union)
	return err
}

type status200Response struct {
	union json.RawMessage
}

func newStatus200Response(response api.AnalysisResponse) (api.GetAnalysisStatusResponseObject, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return status200Response{union: data}, nil
}

func (r status200Response) VisitGetAnalysisStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(r.union)
	return err
}

type reanalyze202Response struct {
	union json.RawMessage
}

func newReanalyze202Response(response api.AnalysisResponse) (api.ReanalyzeRepositoryResponseObject, error) {
	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return reanalyze202Response{union: data}, nil
}

func (r reanalyze202Response) VisitReanalyzeRepositoryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write(r.union)
	return err
}

func (h *Handler) buildHistoryOptions(ctx context.Context, userID, owner, repo string) mapper.CompletedResponseOptions {
	if userID == "" || h.historyChecker == nil {
		return mapper.CompletedResponseOptions{}
	}

	exists, err := h.historyChecker.CheckUserHistoryExists(ctx, userID, owner, repo)
	if err != nil {
		h.logger.Warn(ctx, "failed to check user history", "error", err)
		return mapper.CompletedResponseOptions{}
	}

	return mapper.CompletedResponseOptions{IsInMyHistory: &exists}
}

func (h *Handler) lookupUserTier(ctx context.Context, log *logger.Logger, userID string) subscription.PlanTier {
	if userID == "" || h.tierLookup == nil {
		return ""
	}
	tierStr, err := h.tierLookup.GetUserTier(ctx, userID)
	if err != nil {
		log.Warn(ctx, "failed to lookup user tier, using default", "error", err)
		return ""
	}
	return subscription.PlanTier(tierStr)
}
