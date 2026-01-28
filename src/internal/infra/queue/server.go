package queue

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
)

const (
	DefaultConcurrency     = 5
	DefaultShutdownTimeout = 30 * time.Second
)

// QueueAllocation defines worker count for a specific queue.
type QueueAllocation struct {
	Name       string
	MaxWorkers int
}

type ServerConfig struct {
	Middleware      []rivertype.WorkerMiddleware
	Pool            *pgxpool.Pool
	Queues          []QueueAllocation
	ShutdownTimeout time.Duration
	Workers         *river.Workers
}

type Server struct {
	client          *river.Client[pgx.Tx]
	shutdownTimeout time.Duration
}

func NewServer(ctx context.Context, cfg ServerConfig) (*Server, error) {
	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = DefaultShutdownTimeout
	}

	queues := buildQueueConfig(cfg)

	riverConfig := &river.Config{
		Queues:  queues,
		Workers: cfg.Workers,
	}
	if len(cfg.Middleware) > 0 {
		// Ensure WorkerMiddleware implements Middleware at compile time
		var _ rivertype.Middleware = (rivertype.WorkerMiddleware)(nil)

		// WorkerMiddleware embeds Middleware, assign to interface slice for River config
		middleware := make([]rivertype.Middleware, len(cfg.Middleware))
		for i, m := range cfg.Middleware {
			middleware[i] = m
		}
		riverConfig.Middleware = middleware
	}

	client, err := river.NewClient(riverpgxv5.New(cfg.Pool), riverConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		client:          client,
		shutdownTimeout: shutdownTimeout,
	}, nil
}

// buildQueueConfig creates River queue configuration from ServerConfig.
func buildQueueConfig(cfg ServerConfig) map[string]river.QueueConfig {
	queues := make(map[string]river.QueueConfig, len(cfg.Queues))
	for _, q := range cfg.Queues {
		name := q.Name
		if name == "" {
			name = river.QueueDefault
		}
		maxWorkers := q.MaxWorkers
		if maxWorkers <= 0 {
			maxWorkers = DefaultConcurrency
		}
		queues[name] = river.QueueConfig{MaxWorkers: maxWorkers}
	}
	return queues
}

func (s *Server) Start(ctx context.Context) error {
	return s.client.Start(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()
	return s.client.Stop(ctx)
}

func (s *Server) Client() *river.Client[pgx.Tx] {
	return s.client
}
