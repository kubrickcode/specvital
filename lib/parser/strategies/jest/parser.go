package jest

import (
	"context"

	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/strategies/shared/jstest"
)

const frameworkName = "jest"

func parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	return jstest.Parse(ctx, source, filename, frameworkName)
}
