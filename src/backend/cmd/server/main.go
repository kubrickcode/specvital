package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/specvital/web/src/backend/common/config"
	"github.com/specvital/web/src/backend/common/health"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/modules/analyzer"
)

const (
	apiTimeout      = 10 * time.Minute
	shutdownTimeout = 10 * time.Second
)

func main() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			slog.Debug(".env files not found, using system environment variables")
		}
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if os.Getenv("GITHUB_TOKEN") == "" {
		slog.Warn("GITHUB_TOKEN not set")
	}

	analyzer.InitializeParserStrategies()

	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	origins, err := middleware.GetAllowedOrigins()
	if err != nil {
		return fmt.Errorf("failed to get allowed origins: %w", err)
	}

	router := newRouter(origins)
	return startServer(router)
}

func newRouter(origins []string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger())
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(origins))
	r.Use(chimiddleware.Timeout(apiTimeout))
	r.Use(middleware.Compress())

	analyzer.RegisterRoutes(r)
	health.RegisterRoutes(r)

	return r
}

func startServer(handler http.Handler) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = config.DefaultPort
	}

	addr := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server failed: %w", err)
	case <-quit:
	}

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("server exited")
	return nil
}
