package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/specvital/worker/internal/adapter/ai/prompt"
	"github.com/specvital/worker/internal/domain/specview"
)

// phase2Response represents the expected JSON response from Phase 2.
type phase2Response struct {
	Conversions []phase2Conversion `json:"conversions"`
}

type phase2Conversion struct {
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
	Index       int     `json:"index"`
}

// convertTestNames performs Phase 2: test name to behavior conversion.
func (p *Provider) convertTestNames(ctx context.Context, input specview.Phase2Input, lang specview.Language) (*specview.Phase2Output, error) {
	// Validate input
	if len(input.Tests) == 0 {
		return nil, fmt.Errorf("%w: no tests to convert", specview.ErrInvalidInput)
	}

	systemPrompt := prompt.Phase2SystemPrompt
	userPrompt := prompt.BuildPhase2UserPrompt(input, lang)

	var result string

	// Retry logic
	err := p.phase2Retry.Do(ctx, func() error {
		var innerErr error
		result, innerErr = p.generateContent(ctx, p.phase2Model, systemPrompt, userPrompt, p.phase2CB)
		return innerErr
	})
	if err != nil {
		return nil, fmt.Errorf("phase 2 conversion failed: %w", err)
	}

	// Parse JSON response
	output, err := parsePhase2Response(result)
	if err != nil {
		slog.WarnContext(ctx, "failed to parse phase 2 response",
			"error", err,
			"response", truncateForLog(result, 500),
		)
		return nil, fmt.Errorf("failed to parse phase 2 response: %w", err)
	}

	// Validate output
	if err := validatePhase2Output(ctx, output, input); err != nil {
		slog.WarnContext(ctx, "phase 2 output validation failed",
			"error", err,
		)
		return nil, fmt.Errorf("phase 2 output validation failed: %w", err)
	}

	return output, nil
}

// parsePhase2Response parses the JSON response into Phase2Output.
func parsePhase2Response(jsonStr string) (*specview.Phase2Output, error) {
	var resp phase2Response
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	output := &specview.Phase2Output{
		Behaviors: make([]specview.BehaviorSpec, 0, len(resp.Conversions)),
	}

	for _, c := range resp.Conversions {
		behavior := specview.BehaviorSpec{
			Confidence:  c.Confidence,
			Description: c.Description,
			TestIndex:   c.Index,
		}
		output.Behaviors = append(output.Behaviors, behavior)
	}

	return output, nil
}

// validatePhase2Output validates the Phase 2 output against input.
func validatePhase2Output(ctx context.Context, output *specview.Phase2Output, input specview.Phase2Input) error {
	if output == nil || len(output.Behaviors) == 0 {
		return fmt.Errorf("no behaviors in output")
	}

	// Collect all test indices from input
	expectedIndices := make(map[int]bool)
	for _, test := range input.Tests {
		expectedIndices[test.Index] = true
	}

	// Collect all test indices from output
	coveredIndices := make(map[int]bool)
	for _, behavior := range output.Behaviors {
		if behavior.Description == "" {
			return fmt.Errorf("behavior description is empty for test index %d", behavior.TestIndex)
		}
		if !expectedIndices[behavior.TestIndex] {
			return fmt.Errorf("unexpected test index %d in output", behavior.TestIndex)
		}
		coveredIndices[behavior.TestIndex] = true
	}

	// Check coverage
	if len(coveredIndices) < len(expectedIndices) {
		missing := len(expectedIndices) - len(coveredIndices)
		slog.WarnContext(ctx, "phase 2 output missing test indices",
			"expected", len(expectedIndices),
			"covered", len(coveredIndices),
			"missing", missing,
		)
	}

	return nil
}
