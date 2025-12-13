package server

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/common/docs"
	"github.com/specvital/web/src/backend/common/health"
	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/internal/infra"
	"github.com/specvital/web/src/backend/modules/analyzer"
)

type Handlers struct {
	API    api.StrictServerInterface
	Docs   *docs.Handler
	Health *health.Handler
}

type App struct {
	infra    *infra.Container
	Handlers *Handlers
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := infra.ConfigFromEnv()
	container, err := infra.NewContainer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init infra: %w", err)
	}

	handlers := initHandlers(container)

	return &App{
		infra:    container,
		Handlers: handlers,
	}, nil
}

func initHandlers(infra *infra.Container) *Handlers {
	log := logger.New()

	queries := db.New(infra.DB)
	repo := analyzer.NewRepository(queries)
	queueSvc := analyzer.NewQueueService(infra.Queue)

	analyzerService := analyzer.NewAnalyzerService(log, repo, queueSvc, infra.GitClient)
	analyzerHandler := analyzer.NewAnalyzerHandler(log, analyzerService)

	return &Handlers{
		API:    analyzerHandler,
		Docs:   docs.NewHandler(),
		Health: health.NewHandler(log),
	}
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
