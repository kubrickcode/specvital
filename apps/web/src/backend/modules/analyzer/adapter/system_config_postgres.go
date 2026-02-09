package adapter

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
	"github.com/specvital/web/src/backend/modules/analyzer/domain/port"
)

const keyParserVersion = "parser_version"

// SystemConfigQuerier defines the query methods needed for system config operations.
type SystemConfigQuerier interface {
	GetSystemConfigValue(ctx context.Context, key string) (string, error)
}

var _ port.SystemConfigReader = (*SystemConfigPostgres)(nil)

type SystemConfigPostgres struct {
	queries SystemConfigQuerier
}

func NewSystemConfigPostgres(queries SystemConfigQuerier) *SystemConfigPostgres {
	return &SystemConfigPostgres{queries: queries}
}

func (s *SystemConfigPostgres) GetParserVersion(ctx context.Context) (string, error) {
	value, err := s.queries.GetSystemConfigValue(ctx, keyParserVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrParserVersionNotConfigured
		}
		return "", fmt.Errorf("get parser version from system_config: %w", err)
	}
	return value, nil
}
