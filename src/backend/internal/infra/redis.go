package infra

import (
	"fmt"

	"github.com/hibiken/asynq"
)

func NewAsynqClient(redisURL string) (*asynq.Client, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := asynq.NewClient(opt)

	inspector := asynq.NewInspector(opt)
	_, err = inspector.Servers()
	inspector.Close()

	if err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
