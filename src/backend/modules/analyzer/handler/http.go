package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/usecase"
)

const defaultRecentLimit = 10

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type Handler struct {
	analyzeRepository   *usecase.AnalyzeRepositoryUseCase
	getAnalysis         *usecase.GetAnalysisUseCase
	getRepositoryStats  *usecase.GetRepositoryStatsUseCase
	getUpdateStatus     *usecase.GetUpdateStatusUseCase
	listRepositoryCards *usecase.ListRepositoryCardsUseCase
	logger              *logger.Logger
	reanalyzeRepository *usecase.ReanalyzeRepositoryUseCase
}

var _ api.AnalyzerHandlers = (*Handler)(nil)
var _ api.RepositoryHandlers = (*Handler)(nil)

func NewHandler(
	logger *logger.Logger,
	analyzeRepository *usecase.AnalyzeRepositoryUseCase,
	getAnalysis *usecase.GetAnalysisUseCase,
	listRepositoryCards *usecase.ListRepositoryCardsUseCase,
	getUpdateStatus *usecase.GetUpdateStatusUseCase,
	getRepositoryStats *usecase.GetRepositoryStatsUseCase,
	reanalyzeRepository *usecase.ReanalyzeRepositoryUseCase,
) *Handler {
	return &Handler{
		analyzeRepository:   analyzeRepository,
		getAnalysis:         getAnalysis,
		getRepositoryStats:  getRepositoryStats,
		getUpdateStatus:     getUpdateStatus,
		listRepositoryCards: listRepositoryCards,
		logger:              logger,
		reanalyzeRepository: reanalyzeRepository,
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
	result, err := h.analyzeRepository.Execute(ctx, usecase.AnalyzeRepositoryInput{
		Owner:  owner,
		Repo:   repo,
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
		response, mapErr := mapper.ToCompletedResponse(result.Analysis)
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

func (h *Handler) GetAnalysisStatus(ctx context.Context, request api.GetAnalysisStatusRequestObject) (api.GetAnalysisStatusResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

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
		response, mapErr := mapper.ToCompletedResponse(result.Analysis)
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
	limit := defaultRecentLimit
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	cards, err := h.listRepositoryCards.Execute(ctx, usecase.ListRepositoryCardsInput{
		Limit: limit,
	})
	if err != nil {
		h.logger.Error(ctx, "failed to get recent repositories", "error", err)
		return api.GetRecentRepositories500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get recent repositories"),
		}, nil
	}

	return api.GetRecentRepositories200JSONResponse{
		Data: mapper.ToRepositoryCards(cards),
	}, nil
}

func (h *Handler) GetRepositoryStats(ctx context.Context, _ api.GetRepositoryStatsRequestObject) (api.GetRepositoryStatsResponseObject, error) {
	stats, err := h.getRepositoryStats.Execute(ctx)
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
	result, err := h.reanalyzeRepository.Execute(ctx, usecase.ReanalyzeRepositoryInput{
		Owner:  owner,
		Repo:   repo,
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
