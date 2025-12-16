package server

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/common/docs"
	"github.com/specvital/web/src/backend/common/health"
	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/internal/infra"
	"github.com/specvital/web/src/backend/modules/analyzer"
	"github.com/specvital/web/src/backend/modules/auth"
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
	authMiddleware := middleware.NewAuthMiddleware(container.JWTManager, auth.CookieName)

	return &App{
		AuthMiddleware: authMiddleware,
		Handlers:       handlers,
		infra:          container,
	}, nil
}

func initHandlers(container *infra.Container) (*Handlers, error) {
	log := logger.New()

	queries := db.New(container.DB)

	authRepo := auth.NewRepository(queries)
	stateStore := auth.NewStateStore()
	authService := auth.NewService(&auth.ServiceConfig{
		Encryptor:    container.Encryptor,
		GitHubClient: container.GitHubOAuth,
		JWTManager:   container.JWTManager,
		Repository:   authRepo,
		StateStore:   stateStore,
	})
	authHandler, err := auth.NewHandler(&auth.HandlerConfig{
		CookieSecure: container.SecureCookie,
		FrontendURL:  container.FrontendURL,
		Logger:       log,
		Service:      authService,
	})
	if err != nil {
		return nil, fmt.Errorf("create auth handler: %w", err)
	}

	analyzerRepo := analyzer.NewRepository(queries)
	queueSvc := analyzer.NewQueueService(container.Asynq.Client, container.Asynq.Inspector)
	analyzerService := analyzer.NewAnalyzerService(log, analyzerRepo, queueSvc, container.GitClient, authService)
	analyzerHandler := analyzer.NewAnalyzerHandler(log, analyzerService)

	apiHandlers := api.NewAPIHandlers(analyzerHandler, authHandler)

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
