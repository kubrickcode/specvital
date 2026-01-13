package prompt

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/specvital/worker/internal/domain/specview"
)

//go:embed templates/phase1_system.md
var Phase1SystemPrompt string

// BuildPhase1UserPrompt builds the user prompt for Phase 1 classification.
func BuildPhase1UserPrompt(input specview.Phase1Input, language specview.Language) string {
	var sb strings.Builder

	sb.WriteString("Classify the following tests into business domains and features.\n\n")
	sb.WriteString(fmt.Sprintf("Target Language: %s\n\n", language))
	sb.WriteString("<files>\n")

	totalTests := 0
	for fileIdx, file := range input.Files {
		sb.WriteString(fmt.Sprintf("[%d] %s", fileIdx, file.Path))
		if file.Framework != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", file.Framework))
		}
		sb.WriteString("\n")

		// Domain hints (imports and calls)
		if file.DomainHints != nil {
			if len(file.DomainHints.Imports) > 0 {
				sb.WriteString(fmt.Sprintf("  imports: %s\n", strings.Join(file.DomainHints.Imports, ", ")))
			}
			if len(file.DomainHints.Calls) > 0 {
				sb.WriteString(fmt.Sprintf("  calls: %s\n", strings.Join(file.DomainHints.Calls, ", ")))
			}
		}

		// Tests
		sb.WriteString("  tests:\n")
		for _, test := range file.Tests {
			if test.SuitePath != "" {
				sb.WriteString(fmt.Sprintf("    %d|%s|%s\n", test.Index, test.SuitePath, test.Name))
			} else {
				sb.WriteString(fmt.Sprintf("    %d|%s\n", test.Index, test.Name))
			}
			totalTests++
		}
	}

	sb.WriteString("</files>\n\n")
	sb.WriteString(fmt.Sprintf("Total: %d tests (indices 0-%d). Assign ALL to exactly one feature.", totalTests, totalTests-1))

	return sb.String()
}

