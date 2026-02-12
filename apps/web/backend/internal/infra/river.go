package infra

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type RiverClient struct {
	client *river.Client[pgx.Tx]
}

func NewRiverClient(pool *pgxpool.Pool) (*RiverClient, error) {
	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		return nil, err
	}

	return &RiverClient{client: riverClient}, nil
}

func (r *RiverClient) Client() *river.Client[pgx.Tx] {
	return r.client
}
