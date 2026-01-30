package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/specvital/worker/internal/adapter/ai/prompt"
	"github.com/specvital/worker/internal/adapter/ai/reliability"
	"github.com/specvital/worker/internal/domain/specview"
)

const (
	// v3BatchSize is the number of tests processed per batch in V3 architecture.
	// 20 tests produce ~600 output tokens, safely under token limits.
	v3BatchSize = 20
)

// v3BatchResult represents a single classification result from V3 batch processing.
// Uses compact field names to minimize output tokens.
type v3BatchResult struct {
	Domain  string `json:"d"`
	Feature string `json:"f"`
}

// processV3Batch processes a single batch of tests and returns classifications.
// Returns error if API call fails or response validation fails.
func (p *Provider) processV3Batch(
	ctx context.Context,
	tests []specview.TestForAssignment,
	existingDomains []prompt.DomainSummary,
	lang specview.Language,
) ([]v3BatchResult, *specview.TokenUsage, error) {
	if len(tests) == 0 {
		return []v3BatchResult{}, nil, nil
	}

	systemPrompt := prompt.Phase1V3SystemPrompt
	userPrompt := prompt.BuildV3BatchUserPrompt(tests, existingDomains, lang)

	var results []v3BatchResult
	var usage *specview.TokenUsage

	err := p.phase1Retry.Do(ctx, func() error {
		result, innerUsage, innerErr := p.generateContent(ctx, p.phase1Model, systemPrompt, userPrompt, p.phase1CB)
		if innerErr != nil {
			return innerErr
		}
		usage = innerUsage

		parsed, parseErr := parseV3BatchResponse(result)
		if parseErr != nil {
			slog.WarnContext(ctx, "failed to parse v3 batch response, will retry",
				"error", parseErr,
				"response", truncateForLog(result, 500),
			)
			return &reliability.RetryableError{Err: parseErr}
		}

		if err := validateV3BatchCount(parsed, len(tests)); err != nil {
			slog.WarnContext(ctx, "v3 batch count validation failed, will retry",
				"error", err,
				"expected", len(tests),
				"got", len(parsed),
			)
			return &reliability.RetryableError{Err: err}
		}

		results = parsed
		return nil
	})
	if err != nil {
		return nil, usage, fmt.Errorf("v3 batch processing failed: %w", err)
	}

	return results, usage, nil
}

// parseV3BatchResponse parses the JSON array response from V3 batch API.
func parseV3BatchResponse(jsonStr string) ([]v3BatchResult, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("empty response")
	}

	var results []v3BatchResult
	if err := json.Unmarshal([]byte(jsonStr), &results); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if results == nil {
		return nil, fmt.Errorf("null response array")
	}

	for i, r := range results {
		if r.Domain == "" || r.Feature == "" {
			return nil, fmt.Errorf("empty domain or feature at index %d", i)
		}
	}

	return results, nil
}

// validateV3BatchCount validates that the response count matches expected count.
func validateV3BatchCount(results []v3BatchResult, expectedCount int) error {
	if len(results) != expectedCount {
		return fmt.Errorf("count mismatch: expected %d, got %d", expectedCount, len(results))
	}
	return nil
}
