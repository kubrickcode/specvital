package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/ai/prompt"
	"github.com/kubrickcode/specvital/apps/worker/internal/adapter/ai/reliability"
	"github.com/kubrickcode/specvital/apps/worker/internal/domain/specview"
)

// placementResponse represents the expected JSON response from placement API call.
type placementResponse struct {
	Placements []placementItem `json:"placements"`
}

type placementItem struct {
	Domain    string `json:"domain"`
	Feature   string `json:"feature"`
	TestIndex int    `json:"test_index"`
}

// placeNewTests places new tests into an existing domain/feature structure.
// Uses Phase 1 circuit breaker and retry logic since it's a classification task.
func (p *Provider) placeNewTests(ctx context.Context, input specview.PlacementInput) (*specview.PlacementOutput, *specview.TokenUsage, error) {
	if len(input.NewTests) == 0 {
		return &specview.PlacementOutput{Placements: []specview.TestPlacement{}}, nil, nil
	}

	if input.ExistingStructure == nil || len(input.ExistingStructure.Domains) == 0 {
		return nil, nil, fmt.Errorf("%w: existing structure is required for placement", specview.ErrInvalidInput)
	}

	systemPrompt := prompt.PlacementSystemPrompt
	userPrompt := prompt.BuildPlacementUserPrompt(input)

	var output *specview.PlacementOutput
	var usage *specview.TokenUsage

	err := p.phase1Retry.Do(ctx, func() error {
		result, innerUsage, innerErr := p.generateContent(ctx, p.phase1Model, systemPrompt, userPrompt, p.phase1CB)
		if innerErr != nil {
			return innerErr
		}
		usage = innerUsage

		var parseErr error
		output, parseErr = parsePlacementResponse(result, len(input.NewTests))
		if parseErr != nil {
			slog.WarnContext(ctx, "failed to parse placement response, will retry",
				"error", parseErr,
				"response", truncateForLog(result, 500),
			)
			return &reliability.RetryableError{Err: parseErr}
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("placement failed: %w", err)
	}

	return output, usage, nil
}

// parsePlacementResponse parses the JSON response into PlacementOutput.
func parsePlacementResponse(jsonStr string, expectedCount int) (*specview.PlacementOutput, error) {
	var resp placementResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	output := &specview.PlacementOutput{
		Placements: make([]specview.TestPlacement, 0, len(resp.Placements)),
	}

	seenIndices := make(map[int]bool)
	for _, p := range resp.Placements {
		if p.TestIndex < 0 || p.TestIndex >= expectedCount {
			return nil, fmt.Errorf("invalid test_index %d (expected 0-%d)", p.TestIndex, expectedCount-1)
		}
		if seenIndices[p.TestIndex] {
			return nil, fmt.Errorf("duplicate test_index %d", p.TestIndex)
		}
		seenIndices[p.TestIndex] = true

		if p.Domain == "" || p.Feature == "" {
			return nil, fmt.Errorf("empty domain or feature for test_index %d", p.TestIndex)
		}

		output.Placements = append(output.Placements, specview.TestPlacement{
			DomainName:  p.Domain,
			FeatureName: p.Feature,
			TestIndex:   p.TestIndex,
		})
	}

	if len(output.Placements) != expectedCount {
		return nil, fmt.Errorf("placement count mismatch: got %d, expected %d", len(output.Placements), expectedCount)
	}

	return output, nil
}
