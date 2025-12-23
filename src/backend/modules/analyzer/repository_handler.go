package analyzer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/mapper"
)

const defaultRecentLimit = 10

type RepositoryHandler struct {
	analyzerService AnalyzerService
	logger          *logger.Logger
	repoService     RepositoryService
}

var _ api.RepositoryHandlers = (*RepositoryHandler)(nil)

func NewRepositoryHandler(logger *logger.Logger, repoService RepositoryService, analyzerService AnalyzerService) *RepositoryHandler {
	return &RepositoryHandler{
		analyzerService: analyzerService,
		logger:          logger,
		repoService:     repoService,
	}
}

func (h *RepositoryHandler) GetRecentRepositories(ctx context.Context, request api.GetRecentRepositoriesRequestObject) (api.GetRecentRepositoriesResponseObject, error) {
	limit := defaultRecentLimit
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	cards, err := h.repoService.GetRecentRepositories(ctx, limit)
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

func (h *RepositoryHandler) GetRepositoryStats(ctx context.Context, _ api.GetRepositoryStatsRequestObject) (api.GetRepositoryStatsResponseObject, error) {
	stats, err := h.repoService.GetRepositoryStats(ctx)
	if err != nil {
		h.logger.Error(ctx, "failed to get repository stats", "error", err)
		return api.GetRepositoryStats500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to get repository stats"),
		}, nil
	}

	return api.GetRepositoryStats200JSONResponse(mapper.ToRepositoryStatsResponse(stats)), nil
}

func (h *RepositoryHandler) GetUpdateStatus(ctx context.Context, request api.GetUpdateStatusRequestObject) (api.GetUpdateStatusResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.GetUpdateStatus400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	result, err := h.repoService.GetUpdateStatus(ctx, owner, repo)
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

func (h *RepositoryHandler) ReanalyzeRepository(ctx context.Context, request api.ReanalyzeRepositoryRequestObject) (api.ReanalyzeRepositoryResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.ReanalyzeRepository400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	result, err := h.analyzerService.AnalyzeRepository(ctx, owner, repo)
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
