package analyzer

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/mapper"
)

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type AnalyzerHandler struct {
	logger  *logger.Logger
	service AnalyzerService
}

var _ api.AnalyzerHandlers = (*AnalyzerHandler)(nil)

func NewAnalyzerHandler(logger *logger.Logger, service AnalyzerService) *AnalyzerHandler {
	return &AnalyzerHandler{logger: logger, service: service}
}

func (h *AnalyzerHandler) AnalyzeRepository(ctx context.Context, request api.AnalyzeRepositoryRequestObject) (api.AnalyzeRepositoryResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.AnalyzeRepository400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	result, err := h.service.AnalyzeRepository(ctx, owner, repo)
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

		log.Error(ctx, "service error in AnalyzeRepository", "error", err)
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

func (h *AnalyzerHandler) GetAnalysisStatus(ctx context.Context, request api.GetAnalysisStatusRequestObject) (api.GetAnalysisStatusResponseObject, error) {
	owner, repo := request.Owner, request.Repo
	log := h.logger.With("owner", owner, "repo", repo)

	if err := validateOwnerRepo(owner, repo); err != nil {
		return api.GetAnalysisStatus400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest(err.Error()),
		}, nil
	}

	result, err := h.service.GetAnalysisStatus(ctx, owner, repo)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.GetAnalysisStatus404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: api.NewNotFound("analysis not found"),
			}, nil
		}
		log.Error(ctx, "service error in GetAnalysisStatus", "error", err)
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

// Union response helpers for discriminated union handling.
// The generated StrictServerInterface doesn't directly support union types,
// so we marshal to raw JSON and write using the Visitor pattern.

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
