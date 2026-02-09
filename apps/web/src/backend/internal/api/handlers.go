package api

import (
	"context"
	"net/http"
)

type AnalyzerHandlers interface {
	AnalyzeRepository(ctx context.Context, request AnalyzeRepositoryRequestObject) (AnalyzeRepositoryResponseObject, error)
	GetAnalysisHistory(ctx context.Context, request GetAnalysisHistoryRequestObject) (GetAnalysisHistoryResponseObject, error)
	GetAnalysisStatus(ctx context.Context, request GetAnalysisStatusRequestObject) (GetAnalysisStatusResponseObject, error)
}

type WebhookHandlers interface {
	HandleGitHubAppWebhookRaw(w http.ResponseWriter, r *http.Request)
}

type AuthHandlers interface {
	AuthCallback(ctx context.Context, request AuthCallbackRequestObject) (AuthCallbackResponseObject, error)
	AuthDevLogin(ctx context.Context, request AuthDevLoginRequestObject) (AuthDevLoginResponseObject, error)
	AuthLogin(ctx context.Context, request AuthLoginRequestObject) (AuthLoginResponseObject, error)
	AuthLogout(ctx context.Context, request AuthLogoutRequestObject) (AuthLogoutResponseObject, error)
	AuthMe(ctx context.Context, request AuthMeRequestObject) (AuthMeResponseObject, error)
	AuthRefresh(ctx context.Context, request AuthRefreshRequestObject) (AuthRefreshResponseObject, error)
}

type BookmarkHandlers interface {
	AddBookmark(ctx context.Context, request AddBookmarkRequestObject) (AddBookmarkResponseObject, error)
	GetUserBookmarks(ctx context.Context, request GetUserBookmarksRequestObject) (GetUserBookmarksResponseObject, error)
	RemoveBookmark(ctx context.Context, request RemoveBookmarkRequestObject) (RemoveBookmarkResponseObject, error)
}

type AnalysisHistoryHandlers interface {
	AddUserAnalyzedRepository(ctx context.Context, request AddUserAnalyzedRepositoryRequestObject) (AddUserAnalyzedRepositoryResponseObject, error)
	GetUserAnalyzedRepositories(ctx context.Context, request GetUserAnalyzedRepositoriesRequestObject) (GetUserAnalyzedRepositoriesResponseObject, error)
}

type UserActiveTasksHandlers interface {
	GetUserActiveTasks(ctx context.Context, request GetUserActiveTasksRequestObject) (GetUserActiveTasksResponseObject, error)
}

type GitHubHandlers interface {
	GetOrganizationRepositories(ctx context.Context, request GetOrganizationRepositoriesRequestObject) (GetOrganizationRepositoriesResponseObject, error)
	GetUserGitHubOrganizations(ctx context.Context, request GetUserGitHubOrganizationsRequestObject) (GetUserGitHubOrganizationsResponseObject, error)
	GetUserGitHubRepositories(ctx context.Context, request GetUserGitHubRepositoriesRequestObject) (GetUserGitHubRepositoriesResponseObject, error)
}

type RepositoryHandlers interface {
	GetRecentRepositories(ctx context.Context, request GetRecentRepositoriesRequestObject) (GetRecentRepositoriesResponseObject, error)
	GetRepositoryStats(ctx context.Context, request GetRepositoryStatsRequestObject) (GetRepositoryStatsResponseObject, error)
	GetUpdateStatus(ctx context.Context, request GetUpdateStatusRequestObject) (GetUpdateStatusResponseObject, error)
	ReanalyzeRepository(ctx context.Context, request ReanalyzeRepositoryRequestObject) (ReanalyzeRepositoryResponseObject, error)
}

type GitHubAppHandlers interface {
	GetGitHubAppInstallURL(ctx context.Context, request GetGitHubAppInstallURLRequestObject) (GetGitHubAppInstallURLResponseObject, error)
	GetUserGitHubAppInstallations(ctx context.Context, request GetUserGitHubAppInstallationsRequestObject) (GetUserGitHubAppInstallationsResponseObject, error)
}

type SpecViewHandlers interface {
	GetSpecCacheAvailability(ctx context.Context, request GetSpecCacheAvailabilityRequestObject) (GetSpecCacheAvailabilityResponseObject, error)
	GetSpecCachePrediction(ctx context.Context, request GetSpecCachePredictionRequestObject) (GetSpecCachePredictionResponseObject, error)
	GetSpecDocument(ctx context.Context, request GetSpecDocumentRequestObject) (GetSpecDocumentResponseObject, error)
	GetSpecDocumentByRepository(ctx context.Context, request GetSpecDocumentByRepositoryRequestObject) (GetSpecDocumentByRepositoryResponseObject, error)
	GetSpecGenerationStatus(ctx context.Context, request GetSpecGenerationStatusRequestObject) (GetSpecGenerationStatusResponseObject, error)
	GetSpecVersions(ctx context.Context, request GetSpecVersionsRequestObject) (GetSpecVersionsResponseObject, error)
	GetVersionHistoryByRepository(ctx context.Context, request GetVersionHistoryByRepositoryRequestObject) (GetVersionHistoryByRepositoryResponseObject, error)
	RequestSpecGeneration(ctx context.Context, request RequestSpecGenerationRequestObject) (RequestSpecGenerationResponseObject, error)
}

type UsageHandlers interface {
	CheckQuota(ctx context.Context, request CheckQuotaRequestObject) (CheckQuotaResponseObject, error)
	GetCurrentUsage(ctx context.Context, request GetCurrentUsageRequestObject) (GetCurrentUsageResponseObject, error)
}

type SubscriptionHandlers interface {
	GetUserSubscription(ctx context.Context, request GetUserSubscriptionRequestObject) (GetUserSubscriptionResponseObject, error)
}

type PricingHandlers interface {
	GetPricing(ctx context.Context, request GetPricingRequestObject) (GetPricingResponseObject, error)
}

type APIHandlers struct {
	analyzer        AnalyzerHandlers
	analysisHistory AnalysisHistoryHandlers
	auth            AuthHandlers
	bookmark        BookmarkHandlers
	github          GitHubHandlers
	githubApp       GitHubAppHandlers
	pricing         PricingHandlers
	repository      RepositoryHandlers
	specView        SpecViewHandlers
	subscription    SubscriptionHandlers
	usage           UsageHandlers
	userActiveTasks UserActiveTasksHandlers
	webhook         WebhookHandlers
}

var _ StrictServerInterface = (*APIHandlers)(nil)

func NewAPIHandlers(
	analyzer AnalyzerHandlers,
	analysisHistory AnalysisHistoryHandlers,
	auth AuthHandlers,
	bookmark BookmarkHandlers,
	github GitHubHandlers,
	githubApp GitHubAppHandlers,
	pricing PricingHandlers,
	repository RepositoryHandlers,
	specView SpecViewHandlers,
	subscription SubscriptionHandlers,
	usage UsageHandlers,
	userActiveTasks UserActiveTasksHandlers,
	webhook WebhookHandlers,
) *APIHandlers {
	return &APIHandlers{
		analyzer:        analyzer,
		analysisHistory: analysisHistory,
		auth:            auth,
		bookmark:        bookmark,
		github:          github,
		githubApp:       githubApp,
		pricing:         pricing,
		repository:      repository,
		specView:        specView,
		subscription:    subscription,
		usage:           usage,
		userActiveTasks: userActiveTasks,
		webhook:         webhook,
	}
}

func (h *APIHandlers) AnalyzeRepository(ctx context.Context, request AnalyzeRepositoryRequestObject) (AnalyzeRepositoryResponseObject, error) {
	return h.analyzer.AnalyzeRepository(ctx, request)
}

func (h *APIHandlers) GetAnalysisHistory(ctx context.Context, request GetAnalysisHistoryRequestObject) (GetAnalysisHistoryResponseObject, error) {
	return h.analyzer.GetAnalysisHistory(ctx, request)
}

func (h *APIHandlers) GetAnalysisStatus(ctx context.Context, request GetAnalysisStatusRequestObject) (GetAnalysisStatusResponseObject, error) {
	return h.analyzer.GetAnalysisStatus(ctx, request)
}

func (h *APIHandlers) AuthCallback(ctx context.Context, request AuthCallbackRequestObject) (AuthCallbackResponseObject, error) {
	return h.auth.AuthCallback(ctx, request)
}

func (h *APIHandlers) AuthLogin(ctx context.Context, request AuthLoginRequestObject) (AuthLoginResponseObject, error) {
	return h.auth.AuthLogin(ctx, request)
}

func (h *APIHandlers) AuthLogout(ctx context.Context, request AuthLogoutRequestObject) (AuthLogoutResponseObject, error) {
	return h.auth.AuthLogout(ctx, request)
}

func (h *APIHandlers) AuthMe(ctx context.Context, request AuthMeRequestObject) (AuthMeResponseObject, error) {
	return h.auth.AuthMe(ctx, request)
}

func (h *APIHandlers) AuthRefresh(ctx context.Context, request AuthRefreshRequestObject) (AuthRefreshResponseObject, error) {
	return h.auth.AuthRefresh(ctx, request)
}

func (h *APIHandlers) AuthDevLogin(ctx context.Context, request AuthDevLoginRequestObject) (AuthDevLoginResponseObject, error) {
	return h.auth.AuthDevLogin(ctx, request)
}

func (h *APIHandlers) AddBookmark(ctx context.Context, request AddBookmarkRequestObject) (AddBookmarkResponseObject, error) {
	return h.bookmark.AddBookmark(ctx, request)
}

func (h *APIHandlers) GetRecentRepositories(ctx context.Context, request GetRecentRepositoriesRequestObject) (GetRecentRepositoriesResponseObject, error) {
	return h.repository.GetRecentRepositories(ctx, request)
}

func (h *APIHandlers) GetRepositoryStats(ctx context.Context, request GetRepositoryStatsRequestObject) (GetRepositoryStatsResponseObject, error) {
	return h.repository.GetRepositoryStats(ctx, request)
}

func (h *APIHandlers) GetUpdateStatus(ctx context.Context, request GetUpdateStatusRequestObject) (GetUpdateStatusResponseObject, error) {
	return h.repository.GetUpdateStatus(ctx, request)
}

func (h *APIHandlers) ReanalyzeRepository(ctx context.Context, request ReanalyzeRepositoryRequestObject) (ReanalyzeRepositoryResponseObject, error) {
	return h.repository.ReanalyzeRepository(ctx, request)
}

func (h *APIHandlers) GetUserBookmarks(ctx context.Context, request GetUserBookmarksRequestObject) (GetUserBookmarksResponseObject, error) {
	return h.bookmark.GetUserBookmarks(ctx, request)
}

func (h *APIHandlers) RemoveBookmark(ctx context.Context, request RemoveBookmarkRequestObject) (RemoveBookmarkResponseObject, error) {
	return h.bookmark.RemoveBookmark(ctx, request)
}

func (h *APIHandlers) AddUserAnalyzedRepository(ctx context.Context, request AddUserAnalyzedRepositoryRequestObject) (AddUserAnalyzedRepositoryResponseObject, error) {
	return h.analysisHistory.AddUserAnalyzedRepository(ctx, request)
}

func (h *APIHandlers) GetUserAnalyzedRepositories(ctx context.Context, request GetUserAnalyzedRepositoriesRequestObject) (GetUserAnalyzedRepositoriesResponseObject, error) {
	return h.analysisHistory.GetUserAnalyzedRepositories(ctx, request)
}

func (h *APIHandlers) GetUserActiveTasks(ctx context.Context, request GetUserActiveTasksRequestObject) (GetUserActiveTasksResponseObject, error) {
	return h.userActiveTasks.GetUserActiveTasks(ctx, request)
}

func (h *APIHandlers) GetOrganizationRepositories(ctx context.Context, request GetOrganizationRepositoriesRequestObject) (GetOrganizationRepositoriesResponseObject, error) {
	return h.github.GetOrganizationRepositories(ctx, request)
}

func (h *APIHandlers) GetUserGitHubOrganizations(ctx context.Context, request GetUserGitHubOrganizationsRequestObject) (GetUserGitHubOrganizationsResponseObject, error) {
	return h.github.GetUserGitHubOrganizations(ctx, request)
}

func (h *APIHandlers) GetUserGitHubRepositories(ctx context.Context, request GetUserGitHubRepositoriesRequestObject) (GetUserGitHubRepositoriesResponseObject, error) {
	return h.github.GetUserGitHubRepositories(ctx, request)
}

func (h *APIHandlers) HandleGitHubAppWebhook(_ context.Context, _ HandleGitHubAppWebhookRequestObject) (HandleGitHubAppWebhookResponseObject, error) {
	return HandleGitHubAppWebhook500ApplicationProblemPlusJSONResponse{
		InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("use raw handler instead"),
	}, nil
}

func (h *APIHandlers) GetGitHubAppInstallURL(ctx context.Context, request GetGitHubAppInstallURLRequestObject) (GetGitHubAppInstallURLResponseObject, error) {
	if h.githubApp == nil {
		return GetGitHubAppInstallURL500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("GitHub App not configured"),
		}, nil
	}
	return h.githubApp.GetGitHubAppInstallURL(ctx, request)
}

func (h *APIHandlers) GetUserGitHubAppInstallations(ctx context.Context, request GetUserGitHubAppInstallationsRequestObject) (GetUserGitHubAppInstallationsResponseObject, error) {
	if h.githubApp == nil {
		return GetUserGitHubAppInstallations500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("GitHub App not configured"),
		}, nil
	}
	return h.githubApp.GetUserGitHubAppInstallations(ctx, request)
}

func (h *APIHandlers) WebhookHandler() WebhookHandlers {
	return h.webhook
}

func (h *APIHandlers) GetSpecDocument(ctx context.Context, request GetSpecDocumentRequestObject) (GetSpecDocumentResponseObject, error) {
	if h.specView == nil {
		return GetSpecDocument500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetSpecDocument(ctx, request)
}

func (h *APIHandlers) GetSpecGenerationStatus(ctx context.Context, request GetSpecGenerationStatusRequestObject) (GetSpecGenerationStatusResponseObject, error) {
	if h.specView == nil {
		return GetSpecGenerationStatus500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetSpecGenerationStatus(ctx, request)
}

func (h *APIHandlers) RequestSpecGeneration(ctx context.Context, request RequestSpecGenerationRequestObject) (RequestSpecGenerationResponseObject, error) {
	if h.specView == nil {
		return RequestSpecGeneration500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.RequestSpecGeneration(ctx, request)
}

func (h *APIHandlers) GetSpecVersions(ctx context.Context, request GetSpecVersionsRequestObject) (GetSpecVersionsResponseObject, error) {
	if h.specView == nil {
		return GetSpecVersions500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetSpecVersions(ctx, request)
}

func (h *APIHandlers) GetSpecCacheAvailability(ctx context.Context, request GetSpecCacheAvailabilityRequestObject) (GetSpecCacheAvailabilityResponseObject, error) {
	if h.specView == nil {
		return GetSpecCacheAvailability500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetSpecCacheAvailability(ctx, request)
}

func (h *APIHandlers) GetSpecCachePrediction(ctx context.Context, request GetSpecCachePredictionRequestObject) (GetSpecCachePredictionResponseObject, error) {
	if h.specView == nil {
		return GetSpecCachePrediction500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetSpecCachePrediction(ctx, request)
}

func (h *APIHandlers) CheckQuota(ctx context.Context, request CheckQuotaRequestObject) (CheckQuotaResponseObject, error) {
	if h.usage == nil {
		return CheckQuota500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Usage feature not configured"),
		}, nil
	}
	return h.usage.CheckQuota(ctx, request)
}

func (h *APIHandlers) GetCurrentUsage(ctx context.Context, request GetCurrentUsageRequestObject) (GetCurrentUsageResponseObject, error) {
	if h.usage == nil {
		return GetCurrentUsage500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Usage feature not configured"),
		}, nil
	}
	return h.usage.GetCurrentUsage(ctx, request)
}

func (h *APIHandlers) GetUserSubscription(ctx context.Context, request GetUserSubscriptionRequestObject) (GetUserSubscriptionResponseObject, error) {
	if h.subscription == nil {
		return GetUserSubscription500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Subscription feature not configured"),
		}, nil
	}
	return h.subscription.GetUserSubscription(ctx, request)
}

func (h *APIHandlers) GetPricing(ctx context.Context, request GetPricingRequestObject) (GetPricingResponseObject, error) {
	if h.pricing == nil {
		return GetPricing500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Pricing feature not configured"),
		}, nil
	}
	return h.pricing.GetPricing(ctx, request)
}

func (h *APIHandlers) GetSpecDocumentByRepository(ctx context.Context, request GetSpecDocumentByRepositoryRequestObject) (GetSpecDocumentByRepositoryResponseObject, error) {
	if h.specView == nil {
		return GetSpecDocumentByRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetSpecDocumentByRepository(ctx, request)
}

func (h *APIHandlers) GetVersionHistoryByRepository(ctx context.Context, request GetVersionHistoryByRepositoryRequestObject) (GetVersionHistoryByRepositoryResponseObject, error) {
	if h.specView == nil {
		return GetVersionHistoryByRepository500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: NewInternalError("Spec View feature not configured"),
		}, nil
	}
	return h.specView.GetVersionHistoryByRepository(ctx, request)
}
