package batch

import (
	"fmt"

	"google.golang.org/genai"

	"github.com/specvital/worker/internal/adapter/ai/prompt"
	"github.com/specvital/worker/internal/domain/specview"
)

const (
	defaultSeedForRequest = int32(42) // Fixed seed for deterministic output
	maxOutputTokensValue  = int32(65536)
)

// CreateClassificationJob creates a batch job request for Phase 1 classification.
// Converts Phase1Input into a single InlinedRequest with proper configuration.
func (p *Provider) CreateClassificationJob(input specview.Phase1Input) (BatchRequest, error) {
	if len(input.Files) == 0 {
		return BatchRequest{}, fmt.Errorf("no files to classify")
	}

	// Build prompts using existing prompt builder
	systemPrompt := prompt.Phase1SystemPrompt
	userPrompt := prompt.BuildPhase1UserPrompt(input, input.Language)

	// Create InlinedRequest
	request := &genai.InlinedRequest{
		Contents: []*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: userPrompt},
				},
				Role: "user",
			},
		},
		Config: &genai.GenerateContentConfig{
			Temperature:      genai.Ptr(float32(0.0)), // Deterministic output
			Seed:             genai.Ptr(defaultSeedForRequest),
			MaxOutputTokens:  maxOutputTokensValue,
			ResponseMIMEType: "application/json",
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemPrompt}},
			},
			// Disable thinking to reduce processing time
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: genai.Ptr(int32(0)),
			},
		},
	}

	return BatchRequest{
		AnalysisID: input.AnalysisID,
		Model:      p.config.Phase1Model,
		Requests:   []*genai.InlinedRequest{request},
	}, nil
}
