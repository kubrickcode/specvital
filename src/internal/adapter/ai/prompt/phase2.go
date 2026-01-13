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
func BuildPhase2UserPrompt(input specview.Phase2Input, language specview.Language) string {
	var sb strings.Builder

	sb.WriteString("Convert the following test names to user-friendly descriptions.\n\n")
	sb.WriteString("Context:\n")
	sb.WriteString(fmt.Sprintf("- Domain: %s\n", input.DomainContext))
	sb.WriteString(fmt.Sprintf("- Feature: %s\n", input.FeatureName))
	sb.WriteString(fmt.Sprintf("- Target Language: %s\n\n", language))
	sb.WriteString("<tests>\n")

	for _, test := range input.Tests {
		sb.WriteString(fmt.Sprintf("%d|%s\n", test.Index, test.Name))
	}

	sb.WriteString("</tests>\n\n")
	sb.WriteString(fmt.Sprintf("Convert all %d tests. Output JSON only.", len(input.Tests)))

	return sb.String()
}
