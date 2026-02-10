package parser

import (
	"context"
	"fmt"

	coreparser "github.com/kubrickcode/specvital/packages/core/pkg/parser"
	"github.com/kubrickcode/specvital/packages/core/pkg/source"
	"github.com/kubrickcode/specvital/apps/worker/src/internal/adapter/mapping"
	"github.com/kubrickcode/specvital/apps/worker/src/internal/domain/analysis"
)

// CoreParser implements analysis.Parser using specvital/core's parser package.
type CoreParser struct{}

// NewCoreParser creates a new CoreParser.
func NewCoreParser() *CoreParser {
	return &CoreParser{}
}

// coreSourceProvider is implemented by sources that can provide
// the underlying source.Source for the core parser.
type coreSourceProvider interface {
	CoreSource() source.Source
}

// Scan implements analysis.Parser by delegating to the core parser
// and converting the result to domain types.
func (p *CoreParser) Scan(ctx context.Context, src analysis.Source) (*analysis.Inventory, error) {
	provider, ok := src.(coreSourceProvider)
	if !ok {
		return nil, fmt.Errorf("source does not implement coreSourceProvider interface")
	}

	result, err := coreparser.Scan(ctx, provider.CoreSource())
	if err != nil {
		return nil, fmt.Errorf("core parser scan: %w", err)
	}

	return mapping.ConvertCoreToDomainInventory(result.Inventory), nil
}

// ScanStream implements analysis.StreamingParser by delegating to the core parser's
// ScanStreaming and converting results to domain types.
func (p *CoreParser) ScanStream(ctx context.Context, src analysis.Source) (<-chan analysis.FileResult, error) {
	provider, ok := src.(coreSourceProvider)
	if !ok {
		return nil, fmt.Errorf("source does not implement coreSourceProvider interface")
	}

	coreCh, err := coreparser.ScanStreaming(ctx, provider.CoreSource())
	if err != nil {
		return nil, fmt.Errorf("core parser scan stream: %w", err)
	}

	domainCh := make(chan analysis.FileResult)
	go func() {
		defer close(domainCh)
		for coreResult := range coreCh {
			select {
			case <-ctx.Done():
				return
			case domainCh <- mapping.ConvertCoreFileResult(coreResult):
			}
		}
	}()

	return domainCh, nil
}
