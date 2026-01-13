package prompt

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/specvital/worker/internal/domain/specview"
)

//go:embed templates/phase2_system.md
var Phase2SystemPrompt string

// BuildPhase2UserPrompt builds the user prompt for Phase 2 conversion.
// Returns the prompt string and index mapping (0-based index → original index).
// AI receives 0-based indices matching the system prompt example format.
func BuildPhase2UserPrompt(input specview.Phase2Input, language specview.Language) (string, []int) {
	var sb strings.Builder

	sb.WriteString("Convert the following test names to user-friendly descriptions.\n\n")
	sb.WriteString("Context:\n")
	sb.WriteString(fmt.Sprintf("- Domain: %s\n", input.DomainContext))
	sb.WriteString(fmt.Sprintf("- Feature: %s\n", input.FeatureName))
	sb.WriteString(fmt.Sprintf("- Target Language: %s\n\n", language))
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
