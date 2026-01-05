package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

type AIProvider interface {
	ConvertTestNames(ctx context.Context, input ConvertInput) (map[string]string, error)
	ModelID() string
}

type ConvertInput struct {
	FilePath string
	Language entity.Language
	Suites   []SuiteInput
}

type SuiteInput struct {
	Hierarchy string
	Tests     []string
}
