package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/specvital/web/src/backend/common/docs"
	"github.com/specvital/web/src/backend/common/health"
	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/common/ratelimit"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/client"
	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/internal/infra"
	analyzeradapter "github.com/specvital/web/src/backend/modules/analyzer/adapter"
	analyzerhandler "github.com/specvital/web/src/backend/modules/analyzer/handler"
	analyzerusecase "github.com/specvital/web/src/backend/modules/analyzer/usecase"
	authadapter "github.com/specvital/web/src/backend/modules/auth/adapter"
	authhandler "github.com/specvital/web/src/backend/modules/auth/handler"
	authusecase "github.com/specvital/web/src/backend/modules/auth/usecase"
	ghappadapter "github.com/specvital/web/src/backend/modules/github-app/adapter"
	ghapphandler "github.com/specvital/web/src/backend/modules/github-app/handler"
	ghappusecase "github.com/specvital/web/src/backend/modules/github-app/usecase"
	githubadapter "github.com/specvital/web/src/backend/modules/github/adapter"
	githubhandler "github.com/specvital/web/src/backend/modules/github/handler"
	githubusecase "github.com/specvital/web/src/backend/modules/github/usecase"
	useradapter "github.com/specvital/web/src/backend/modules/user/adapter"
	userhandler "github.com/specvital/web/src/backend/modules/user/handler"
	activetaskuc "github.com/specvital/web/src/backend/modules/user/usecase/activetask"
	bookmarkuc "github.com/specvital/web/src/backend/modules/user/usecase/bookmark"
	historyuc "github.com/specvital/web/src/backend/modules/user/usecase/history"

	specviewadapter "github.com/specvital/web/src/backend/modules/spec-view/adapter"
	specviewhandler "github.com/specvital/web/src/backend/modules/spec-view/handler"
	specviewusecase "github.com/specvital/web/src/backend/modules/spec-view/usecase"
	subscriptionadapter "github.com/specvital/web/src/backend/modules/subscription/adapter"
	subscriptionhandler "github.com/specvital/web/src/backend/modules/subscription/handler"
	subscriptionusecase "github.com/specvital/web/src/backend/modules/subscription/usecase"
	usageadapter "github.com/specvital/web/src/backend/modules/usage/adapter"
	usagehandler "github.com/specvital/web/src/backend/modules/usage/handler"
	usageusecase "github.com/specvital/web/src/backend/modules/usage/usecase"
)

type Handlers struct {
	API     api.StrictServerInterface
	Docs    *docs.Handler
	Health  *health.Handler
	Webhook api.WebhookHandlers
}

type App struct {
	AuthMiddleware *middleware.AuthMiddleware
	Handlers       *Handlers
	closers        []io.Closer
	infra          *infra.Container
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := infra.ConfigFromEnv()
	container, err := infra.NewContainer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init infra: %w", err)
	}

	handlers, closers, err := initHandlers(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("init handlers: %w", err)
	}
	authMiddleware := middleware.NewAuthMiddleware(
		container.JWTManager,
		authhandler.AccessCookieName,
		authhandler.RefreshCookieName,
	)

	return &App{
		AuthMiddleware: authMiddleware,
		Handlers:       handlers,
		closers:        closers,
		infra:          container,
	}, nil
}

func initHandlers(ctx context.Context, container *infra.Container) (*Handlers, []io.Closer, error) {
	log := logger.New()
	var closers []io.Closer

	queries := db.New(container.DB)

	authRepo := authadapter.NewPostgresRepository(queries)
	refreshTokenRepo := authadapter.NewRefreshTokenPostgresRepository(container.DB, queries)
	stateStore := authadapter.NewMemoryStateStore()

	getCurrentUserUC := authusecase.NewGetCurrentUserUseCase(authRepo)
	getUserGitHubTokenUC := authusecase.NewGetUserGitHubTokenUseCase(container.Encryptor, authRepo)
	tokenProvider := authadapter.NewTokenProviderAdapter(getUserGitHubTokenUC)

	subscriptionConfig := subscriptionadapter.ConfigFromEnv()
	subscriptionRepo := subscriptionadapter.NewPostgresRepository(queries)
	assignDefaultPlanUC := subscriptionusecase.NewAssignDefaultPlanUseCase(subscriptionConfig.DefaultPlanTier, subscriptionRepo)

	usageRepo := usageadapter.NewPostgresRepository(queries)
	reservationRepo := usageadapter.NewQuotaReservationPostgresRepository(queries)

	checkQuotaUC := usageusecase.NewCheckQuotaUseCase(subscriptionRepo, usageRepo, reservationRepo)
	getCurrentUsageUC := usageusecase.NewGetCurrentUsageUseCase(subscriptionRepo, usageRepo, reservationRepo)

	handleOAuthCallbackUC := authusecase.NewHandleOAuthCallbackUseCase(
		container.Encryptor,
		container.GitHubOAuth,
		refreshTokenRepo,
		authRepo,
		stateStore,
		assignDefaultPlanUC,
		container.JWTManager,
	)
	initiateOAuthUC := authusecase.NewInitiateOAuthUseCase(container.GitHubOAuth, stateStore)
	refreshTokenUC := authusecase.NewRefreshTokenUseCase(refreshTokenRepo, authRepo, container.JWTManager)

	var devLoginUC *authusecase.DevLoginUseCase
	if container.Environment != "production" {
		devLoginUC = authusecase.NewDevLoginUseCase(
			true,
			refreshTokenRepo,
			authRepo,
			container.JWTManager,
		)
	}

	authHandler, err := authhandler.NewHandler(&authhandler.HandlerConfig{
		CookieDomain:        container.CookieDomain,
		CookieSecure:        container.SecureCookie,
		DevLogin:            devLoginUC,
		FrontendURL:         container.FrontendURL,
		GetCurrentUser:      getCurrentUserUC,
		HandleOAuthCallback: handleOAuthCallbackUC,
		InitiateOAuth:       initiateOAuthUC,
		Logger:              log,
		RefreshToken:        refreshTokenUC,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create auth handler: %w", err)
	}

	bookmarkRepo := useradapter.NewBookmarkRepository(queries)
	historyRepo := useradapter.NewHistoryRepository(queries)
	activeTaskRepo := useradapter.NewActiveTaskRepository(queries)

	addAnalyzedRepoUC := historyuc.NewAddAnalyzedRepoUseCase(historyRepo)
	addBookmarkUC := bookmarkuc.NewAddBookmarkUseCase(bookmarkRepo)
	getBookmarksUC := bookmarkuc.NewGetBookmarksUseCase(bookmarkRepo)
	removeBookmarkUC := bookmarkuc.NewRemoveBookmarkUseCase(bookmarkRepo)
	getAnalyzedReposUC := historyuc.NewGetAnalyzedReposUseCase(historyRepo)
	getUserActiveTasksUC := activetaskuc.NewGetUserActiveTasksUseCase(activeTaskRepo)

	userHandler, err := userhandler.NewHandler(&userhandler.HandlerConfig{
		AddAnalyzedRepo:    addAnalyzedRepoUC,
		AddBookmark:        addBookmarkUC,
		GetAnalyzedRepos:   getAnalyzedReposUC,
		GetBookmarks:       getBookmarksUC,
		GetUserActiveTasks: getUserActiveTasksUC,
		Logger:             log,
		RemoveBookmark:     removeBookmarkUC,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create user handler: %w", err)
	}

	analyzerRepo := analyzeradapter.NewPostgresRepository(queries)
	analyzerQueue := analyzeradapter.NewRiverQueueService(container.River.Client(), analyzerRepo)
	analyzerGitClient := analyzeradapter.NewGitClientAdapter(container.GitClient)
	systemConfig := analyzeradapter.NewSystemConfigPostgres(queries)

	analyzeRepositoryUC := analyzerusecase.NewAnalyzeRepositoryUseCase(analyzerGitClient, analyzerQueue, analyzerRepo, systemConfig, tokenProvider, container.DB, reservationRepo)
	getAnalysisUC := analyzerusecase.NewGetAnalysisUseCase(analyzerQueue, analyzerRepo)
	getAnalysisHistoryUC := analyzerusecase.NewGetAnalysisHistoryUseCase(analyzerRepo)
	listRepositoryCardsUC := analyzerusecase.NewListRepositoryCardsUseCase(analyzerGitClient, analyzerRepo, tokenProvider)
	getUpdateStatusUC := analyzerusecase.NewGetUpdateStatusUseCase(analyzerGitClient, analyzerRepo, systemConfig, tokenProvider)
	getRepositoryStatsUC := analyzerusecase.NewGetRepositoryStatsUseCase(analyzerRepo)
	reanalyzeRepositoryUC := analyzerusecase.NewReanalyzeRepositoryUseCase(analyzerGitClient, analyzerQueue, analyzerRepo, tokenProvider)

	anonymousRateLimiter := ratelimit.NewIPRateLimiter(10, time.Minute)
	closers = append(closers, anonymousRateLimiter)

	tierLookup := subscriptionadapter.NewTierLookupAdapter(subscriptionRepo)

	analyzerHandler := analyzerhandler.NewHandler(
		log,
		analyzeRepositoryUC,
		getAnalysisUC,
		getAnalysisHistoryUC,
		listRepositoryCardsUC,
		getUpdateStatusUC,
		getRepositoryStatsUC,
		reanalyzeRepositoryUC,
		historyRepo,
		anonymousRateLimiter,
		tierLookup,
	)

	githubRepo := githubadapter.NewPostgresRepository(container.DB, queries)
	githubClientFactory := githubadapter.NewGitHubClientFactory(client.NewGitHubClientFactory())

	ghAppRepo := ghappadapter.NewPostgresRepository(queries)
	installationLookup := githubadapter.NewInstallationLookupAdapter(ghAppRepo)
	getInstallationTokenUC := ghappusecase.NewGetInstallationTokenUseCase(container.GitHubAppClient, ghAppRepo)
	installationTokenProvider := githubadapter.NewInstallationTokenProviderAdapter(getInstallationTokenUC)

	listUserReposUC := githubusecase.NewListUserReposUseCase(githubClientFactory, githubRepo, tokenProvider)
	listUserOrgsUC := githubusecase.NewListUserOrgsUseCase(githubClientFactory, githubRepo, tokenProvider, installationLookup)
	listOrgReposUC := githubusecase.NewListOrgReposUseCase(githubClientFactory, githubRepo, tokenProvider, installationLookup, installationTokenProvider)

	githubHandler, err := githubhandler.NewHandler(&githubhandler.HandlerConfig{
		ListOrgRepos:  listOrgReposUC,
		ListUserOrgs:  listUserOrgsUC,
		ListUserRepos: listUserReposUC,
		Logger:        log,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create github handler: %w", err)
	}

	handleWebhookUC := ghappusecase.NewHandleWebhookUseCase(ghAppRepo)
	webhookVerifier, err := ghappadapter.NewWebhookVerifier(container.GitHubAppWebhookSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("create webhook verifier: %w", err)
	}

	webhookHandler, err := ghapphandler.NewHandler(&ghapphandler.HandlerConfig{
		HandleWebhook: handleWebhookUC,
		Logger:        log,
		Verifier:      webhookVerifier,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create github-app handler: %w", err)
	}

	listInstallationsUC := ghappusecase.NewListInstallationsUseCase(ghAppRepo)
	getInstallURLUC := ghappusecase.NewGetInstallURLUseCase(container.GitHubAppClient)

	ghAppAPIHandler, err := ghapphandler.NewAPIHandler(&ghapphandler.APIHandlerConfig{
		GetInstallURL:     getInstallURLUC,
		ListInstallations: listInstallationsUC,
		Logger:            log,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create github-app api handler: %w", err)
	}

	specViewRepo := specviewadapter.NewPostgresRepository(queries)
	specViewQueue := specviewadapter.NewRiverQueueService(container.River.Client())

	getSpecDocumentUC := specviewusecase.NewGetSpecDocumentUseCase(specViewRepo)
	requestGenerationUC := specviewusecase.NewRequestGenerationUseCase(specViewRepo, specViewQueue, checkQuotaUC, container.DB, reservationRepo)
	getGenerationStatusUC := specviewusecase.NewGetGenerationStatusUseCase(specViewRepo)
	getVersionsUC := specviewusecase.NewGetVersionsUseCase(specViewRepo)
	getCacheAvailabilityUC := specviewusecase.NewGetCacheAvailabilityUseCase(specViewRepo)
	getCachePredictionUC := specviewusecase.NewGetCachePredictionUseCase(specViewRepo)
	getSpecByRepositoryUC := specviewusecase.NewGetSpecByRepositoryUseCase(specViewRepo)
	getVersionHistoryByRepoUC := specviewusecase.NewGetVersionHistoryByRepositoryUseCase(specViewRepo)

	specViewHandler, err := specviewhandler.NewHandler(&specviewhandler.HandlerConfig{
		GetCacheAvailability:    getCacheAvailabilityUC,
		GetCachePrediction:      getCachePredictionUC,
		GetGenerationStatus:     getGenerationStatusUC,
		GetSpecByRepository:     getSpecByRepositoryUC,
		GetSpecDocument:         getSpecDocumentUC,
		GetVersionHistoryByRepo: getVersionHistoryByRepoUC,
		GetVersions:             getVersionsUC,
		Logger:                  log,
		RequestGeneration:       requestGenerationUC,
		TierLookup:              tierLookup,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create spec-view handler: %w", err)
	}

	usageHandler, err := usagehandler.NewHandler(&usagehandler.HandlerConfig{
		CheckQuota:      checkQuotaUC,
		GetCurrentUsage: getCurrentUsageUC,
		Logger:          log,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create usage handler: %w", err)
	}

	getPricingPlansUC := subscriptionusecase.NewGetPricingPlansUseCase(subscriptionRepo)
	getUserSubscriptionUC := subscriptionusecase.NewGetUserSubscriptionUseCase(subscriptionRepo)
	subscriptionHandler, err := subscriptionhandler.NewHandler(&subscriptionhandler.HandlerConfig{
		GetPricingPlans:     getPricingPlansUC,
		GetUserSubscription: getUserSubscriptionUC,
		Logger:              log,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create subscription handler: %w", err)
	}

	apiHandlers := api.NewAPIHandlers(analyzerHandler, userHandler, authHandler, userHandler, githubHandler, ghAppAPIHandler, subscriptionHandler, analyzerHandler, specViewHandler, subscriptionHandler, usageHandler, userHandler, webhookHandler)

	return &Handlers{
		API:     apiHandlers,
		Docs:    docs.NewHandler(),
		Health:  health.NewHandler(log),
		Webhook: webhookHandler,
	}, closers, nil
}

func (a *App) APIHandler() api.StrictServerInterface {
	return a.Handlers.API
}

func (a *App) RouteRegistrars() []RouteRegistrar {
	return []RouteRegistrar{
		a.Handlers.Docs,
		a.Handlers.Health,
	}
}

func (a *App) WebhookHandler() api.WebhookHandlers {
	return a.Handlers.Webhook
}

func (a *App) Close() error {
	var errs []error

	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if a.infra != nil {
		if err := a.infra.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}
