package server

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/common/docs"
	"github.com/specvital/web/src/backend/common/health"
	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
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
	githubadapter "github.com/specvital/web/src/backend/modules/github/adapter"
	githubhandler "github.com/specvital/web/src/backend/modules/github/handler"
	githubusecase "github.com/specvital/web/src/backend/modules/github/usecase"
	"github.com/specvital/web/src/backend/modules/user"
)

type Handlers struct {
	API    api.StrictServerInterface
	Docs   *docs.Handler
	Health *health.Handler
}

type App struct {
	AuthMiddleware *middleware.AuthMiddleware
	Handlers       *Handlers
	infra          *infra.Container
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := infra.ConfigFromEnv()
	container, err := infra.NewContainer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init infra: %w", err)
	}

	handlers, err := initHandlers(container)
	if err != nil {
		return nil, fmt.Errorf("init handlers: %w", err)
	}
	authMiddleware := middleware.NewAuthMiddleware(container.JWTManager, authhandler.CookieName)

	return &App{
		AuthMiddleware: authMiddleware,
		Handlers:       handlers,
		infra:          container,
	}, nil
}

func initHandlers(container *infra.Container) (*Handlers, error) {
	log := logger.New()

	queries := db.New(container.DB)

	authRepo := authadapter.NewPostgresRepository(queries)
	stateStore := authadapter.NewMemoryStateStore()

	getCurrentUserUC := authusecase.NewGetCurrentUserUseCase(authRepo)
	getUserGitHubTokenUC := authusecase.NewGetUserGitHubTokenUseCase(container.Encryptor, authRepo)
	tokenProvider := authadapter.NewTokenProviderAdapter(getUserGitHubTokenUC)
	handleOAuthCallbackUC := authusecase.NewHandleOAuthCallbackUseCase(
		container.Encryptor,
		container.GitHubOAuth,
		authRepo,
		stateStore,
		container.JWTManager,
	)
	initiateOAuthUC := authusecase.NewInitiateOAuthUseCase(container.GitHubOAuth, stateStore)

	authHandler, err := authhandler.NewHandler(&authhandler.HandlerConfig{
		CookieDomain:        container.CookieDomain,
		CookieSecure:        container.SecureCookie,
		FrontendURL:         container.FrontendURL,
		GetCurrentUser:      getCurrentUserUC,
		HandleOAuthCallback: handleOAuthCallbackUC,
		InitiateOAuth:       initiateOAuthUC,
		Logger:              log,
	})
	if err != nil {
		return nil, fmt.Errorf("create auth handler: %w", err)
	}

	bookmarkRepo := user.NewBookmarkRepository(queries)
	bookmarkService := user.NewBookmarkService(bookmarkRepo)
	bookmarkHandler, err := user.NewBookmarkHandler(&user.BookmarkHandlerConfig{
		Logger:  log,
		Service: bookmarkService,
	})
	if err != nil {
		return nil, fmt.Errorf("create bookmark handler: %w", err)
	}

	analysisHistoryRepo := user.NewAnalysisHistoryRepository(queries)
	analysisHistoryService := user.NewAnalysisHistoryService(analysisHistoryRepo)
	analysisHistoryHandler, err := user.NewAnalysisHistoryHandler(&user.AnalysisHistoryHandlerConfig{
		Logger:  log,
		Service: analysisHistoryService,
	})
	if err != nil {
		return nil, fmt.Errorf("create analysis history handler: %w", err)
	}

	analyzerRepo := analyzeradapter.NewPostgresRepository(queries)
	analyzerQueue := analyzeradapter.NewRiverQueueService(container.River.Client(), analyzerRepo)
	analyzerGitClient := analyzeradapter.NewGitClientAdapter(container.GitClient)

	analyzeRepositoryUC := analyzerusecase.NewAnalyzeRepositoryUseCase(analyzerGitClient, analyzerQueue, analyzerRepo, tokenProvider)
	getAnalysisUC := analyzerusecase.NewGetAnalysisUseCase(analyzerQueue, analyzerRepo)
	listRepositoryCardsUC := analyzerusecase.NewListRepositoryCardsUseCase(analyzerRepo)
	getUpdateStatusUC := analyzerusecase.NewGetUpdateStatusUseCase(analyzerGitClient, analyzerRepo, tokenProvider)
	getRepositoryStatsUC := analyzerusecase.NewGetRepositoryStatsUseCase(analyzerRepo)
	reanalyzeRepositoryUC := analyzerusecase.NewReanalyzeRepositoryUseCase(analyzerGitClient, analyzerQueue, analyzerRepo, tokenProvider)

	analyzerHandler := analyzerhandler.NewHandler(
		log,
		analyzeRepositoryUC,
		getAnalysisUC,
		listRepositoryCardsUC,
		getUpdateStatusUC,
		getRepositoryStatsUC,
		reanalyzeRepositoryUC,
	)

	githubRepo := githubadapter.NewPostgresRepository(container.DB, queries)
	githubClientFactory := githubadapter.NewGitHubClientFactory(client.NewGitHubClientFactory())

	listUserReposUC := githubusecase.NewListUserReposUseCase(githubClientFactory, githubRepo, tokenProvider)
	listUserOrgsUC := githubusecase.NewListUserOrgsUseCase(githubClientFactory, githubRepo, tokenProvider)
	listOrgReposUC := githubusecase.NewListOrgReposUseCase(githubClientFactory, githubRepo, tokenProvider)

	githubHandler, err := githubhandler.NewHandler(&githubhandler.HandlerConfig{
		ListOrgRepos:  listOrgReposUC,
		ListUserOrgs:  listUserOrgsUC,
		ListUserRepos: listUserReposUC,
		Logger:        log,
	})
	if err != nil {
		return nil, fmt.Errorf("create github handler: %w", err)
	}

	apiHandlers := api.NewAPIHandlers(analyzerHandler, analysisHistoryHandler, authHandler, bookmarkHandler, githubHandler, analyzerHandler)

	return &Handlers{
		API:    apiHandlers,
		Docs:   docs.NewHandler(),
		Health: health.NewHandler(log),
	}, nil
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

func (a *App) Close() error {
	if a.infra != nil {
		return a.infra.Close()
	}
	return nil
}
