package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	DefaultConcurrency = 5
)

type ServerConfig struct {
	Concurrency int
	RedisURL    string
}

func NewServer(cfg ServerConfig) (*asynq.Server, error) {
	concurrency := cfg.Concurrency
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	opt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URI: %w", err)
	}

	return asynq.NewServer(opt, asynq.Config{
		Concurrency: concurrency,
	}), nil
}

func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}
