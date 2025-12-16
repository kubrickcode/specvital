package infra

import (
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
)

type AsynqComponents struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewAsynqComponents(redisURL string) (*AsynqComponents, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := asynq.NewClient(opt)
	inspector := asynq.NewInspector(opt)

	_, err = inspector.Servers()
	if err != nil {
		client.Close()
		inspector.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &AsynqComponents{
		Client:    client,
		Inspector: inspector,
	}, nil
}

func (c *AsynqComponents) Close() error {
	var errs []error
	if err := c.Client.Close(); err != nil {
		errs = append(errs, fmt.Errorf("client: %w", err))
	}
	if err := c.Inspector.Close(); err != nil {
		errs = append(errs, fmt.Errorf("inspector: %w", err))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
