package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/specvital/worker/internal/adapter/ai/prompt"
	"github.com/specvital/worker/internal/domain/specview"
)

// phase1Response represents the expected JSON response from Phase 1.
type phase1Response struct {
	Domains []phase1Domain `json:"domains"`
}

type phase1Domain struct {
	Confidence  float64         `json:"confidence"`
	Description string          `json:"description"`
	Features    []phase1Feature `json:"features"`
	Name        string          `json:"name"`
}

type phase1Feature struct {
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
	TestIndices []int   `json:"test_indices"`
}

// classifyDomains performs Phase 1: domain and feature classification.
func (p *Provider) classifyDomains(ctx context.Context, input specview.Phase1Input, lang specview.Language) (*specview.Phase1Output, *specview.TokenUsage, error) {
	// Validate input
	if len(input.Files) == 0 {
		return nil, nil, fmt.Errorf("%w: no files to classify", specview.ErrInvalidInput)
	}

	systemPrompt := prompt.Phase1SystemPrompt
	userPrompt := prompt.BuildPhase1UserPrompt(input, lang)

	var result string
	var usage *specview.TokenUsage

	// Retry logic
	err := p.phase1Retry.Do(ctx, func() error {
		var innerErr error
		result, usage, innerErr = p.generateContent(ctx, p.phase1Model, systemPrompt, userPrompt, p.phase1CB)
		return innerErr
	})
	if err != nil {
		return nil, nil, fmt.Errorf("phase 1 classification failed: %w", err)
	}

	// Parse JSON response
	output, err := parsePhase1Response(result)
	if err != nil {
		slog.WarnContext(ctx, "failed to parse phase 1 response",
			"error", err,
			"response", truncateForLog(result, 500),
		)
		return nil, nil, fmt.Errorf("failed to parse phase 1 response: %w", err)
	}

	// Validate output
	if err := validatePhase1Output(ctx, output, input); err != nil {
		slog.WarnContext(ctx, "phase 1 output validation failed",
			"error", err,
		)
		return nil, nil, fmt.Errorf("phase 1 output validation failed: %w", err)
	}

	return output, usage, nil
}

// parsePhase1Response parses the JSON response into Phase1Output.
func parsePhase1Response(jsonStr string) (*specview.Phase1Output, error) {
	var resp phase1Response
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	output := &specview.Phase1Output{
		Domains: make([]specview.DomainGroup, 0, len(resp.Domains)),
	}

	for _, d := range resp.Domains {
		domain := specview.DomainGroup{
			Confidence:  d.Confidence,
			Description: d.Description,
			Features:    make([]specview.FeatureGroup, 0, len(d.Features)),
			Name:        d.Name,
		}

		for _, f := range d.Features {
			feature := specview.FeatureGroup{
				Confidence:  f.Confidence,
				Description: f.Description,
				Name:        f.Name,
				TestIndices: f.TestIndices,
			}
			domain.Features = append(domain.Features, feature)
		}

		output.Domains = append(output.Domains, domain)
	}

	return output, nil
}

// validatePhase1Output validates the Phase 1 output against input.
func validatePhase1Output(ctx context.Context, output *specview.Phase1Output, input specview.Phase1Input) error {
	if output == nil || len(output.Domains) == 0 {
		return fmt.Errorf("no domains in output")
	}

	// Collect all test indices from input
	expectedIndices := make(map[int]bool)
	for _, file := range input.Files {
		for _, test := range file.Tests {
			expectedIndices[test.Index] = true
		}
	}

	// Collect all test indices from output
	coveredIndices := make(map[int]bool)
	for _, domain := range output.Domains {
		if domain.Name == "" {
			return fmt.Errorf("domain name is empty")
		}
		for _, feature := range domain.Features {
			if feature.Name == "" {
				return fmt.Errorf("feature name is empty in domain %q", domain.Name)
			}
			for _, idx := range feature.TestIndices {
				if !expectedIndices[idx] {
					return fmt.Errorf("unexpected test index %d in feature %q", idx, feature.Name)
				}
				coveredIndices[idx] = true
			}
		}
	}

	// Check coverage
	if len(coveredIndices) < len(expectedIndices) {
		missing := len(expectedIndices) - len(coveredIndices)
		// Log warning but don't fail
		slog.WarnContext(ctx, "phase 1 output missing test indices",
			"expected", len(expectedIndices),
			"covered", len(coveredIndices),
			"missing", missing,
		)
	}

	return nil
}

// truncateForLog truncates a string for logging purposes.
func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
