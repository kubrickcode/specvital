package prompt

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/kubrickcode/specvital/apps/worker/src/internal/domain/specview"
)

//go:embed templates/phase2_system.md
var Phase2SystemPrompt string

// BuildPhase2UserPrompt builds the user prompt for Phase 2 conversion.
// Returns the prompt string and index mapping (0-based index → original index).
// AI receives 0-based indices matching the system prompt example format.
func BuildPhase2UserPrompt(input specview.Phase2Input, language specview.Language) (string, []int) {
	var sb strings.Builder

	sb.WriteString("Convert the following test names to user-friendly descriptions.\n\n")

	// Inject language-specific examples for few-shot learning
	examples := GetPhase2Examples(language)
	if len(examples) > 0 {
		sb.WriteString("## Examples:\n")
		for _, ex := range examples {
			sb.WriteString(fmt.Sprintf("- `%s` → \"%s\"\n", ex.Input, ex.Output))
		}
		sb.WriteString("\n")
	} else {
		// Fallback: use English examples with strong translation instruction
		fallbackExamples := GetPhase2Examples("English")
		if len(fallbackExamples) > 0 {
			sb.WriteString("## Reference Examples (English):\n")
			for _, ex := range fallbackExamples {
				sb.WriteString(fmt.Sprintf("- `%s` → \"%s\"\n", ex.Input, ex.Output))
			}
			sb.WriteString(fmt.Sprintf("\n**CRITICAL: The examples above are in English for reference only. You MUST translate ALL output to %s. Do NOT output in English.**\n\n", language))
		}
	}

	sb.WriteString("Context:\n")
	sb.WriteString(fmt.Sprintf("- Domain: %s\n", input.DomainContext))
	sb.WriteString(fmt.Sprintf("- Feature: %s\n", input.FeatureName))
	sb.WriteString(fmt.Sprintf("- Target Language: %s (ALL output MUST be in this language)\n\n", language))
	sb.WriteString("<tests>\n")

	// Build index mapping: position in slice → original test index
	indexMapping := make([]int, len(input.Tests))
	for i, test := range input.Tests {
		indexMapping[i] = test.Index
		// Send 0-based index to AI (matches system prompt example)
		sb.WriteString(fmt.Sprintf("%d|%s\n", i, test.Name))
	}

	sb.WriteString("</tests>\n\n")
	sb.WriteString(fmt.Sprintf("Convert all %d tests. Output JSON only.", len(input.Tests)))

	return sb.String(), indexMapping
}
