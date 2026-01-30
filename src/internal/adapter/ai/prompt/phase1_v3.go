package prompt

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/specvital/worker/internal/domain/specview"
)

//go:embed templates/phase1_v3_system.md
var Phase1V3SystemPrompt string

// DomainSummary provides context about existing domains for anchor propagation.
// Used to guide AI toward consistent domain naming across batches.
type DomainSummary struct {
	Description string
	Features    []string
	Name        string
}

// BuildV3BatchUserPrompt builds the user prompt for V3 batch classification.
// It includes test list and optional existing domains for anchor propagation.
//
// Output format requirement: Exactly len(tests) items in response array.
func BuildV3BatchUserPrompt(tests []specview.TestForAssignment, existingDomains []DomainSummary, lang specview.Language) string {
	var sb strings.Builder

	sb.WriteString("Classify the following tests into domain/feature pairs.\n\n")
	fmt.Fprintf(&sb, "Target Language: %s\n\n", lang)

	writeExistingDomains(&sb, existingDomains)
	writeV3TestsSection(&sb, tests)

	return sb.String()
}

func writeExistingDomains(sb *strings.Builder, domains []DomainSummary) {
	if len(domains) == 0 {
		sb.WriteString("<existing-domains>\n")
		sb.WriteString("(none - create new domains as needed)\n")
		sb.WriteString("</existing-domains>\n\n")
		return
	}

	sb.WriteString("<existing-domains>\n")
	sb.WriteString("Prefer assigning to these existing domains when appropriate:\n\n")

	for _, domain := range domains {
		fmt.Fprintf(sb, "- %s", domain.Name)
		if domain.Description != "" {
			fmt.Fprintf(sb, ": %s", domain.Description)
		}
		sb.WriteString("\n")

		for _, feature := range domain.Features {
			fmt.Fprintf(sb, "  - %s\n", feature)
		}
	}

	sb.WriteString("</existing-domains>\n\n")
}

func writeV3TestsSection(sb *strings.Builder, tests []specview.TestForAssignment) {
	sb.WriteString("<tests>\n")

	if len(tests) == 0 {
		sb.WriteString("</tests>\n\nTotal: 0 tests. No tests to classify.")
		return
	}

	for i, test := range tests {
		writeV3TestEntry(sb, i, test)
	}

	sb.WriteString("</tests>\n\n")

	totalTests := len(tests)
	if totalTests == 1 {
		sb.WriteString("Total: 1 test. Return exactly 1 classification.")
	} else {
		fmt.Fprintf(sb, "Total: %d tests. Return exactly %d classifications in the same order.", totalTests, totalTests)
	}
}

func writeV3TestEntry(sb *strings.Builder, idx int, test specview.TestForAssignment) {
	fmt.Fprintf(sb, "[%d] %s: %s", idx, test.FilePath, test.Name)

	if test.SuitePath != "" {
		fmt.Fprintf(sb, " (suite: %s)", test.SuitePath)
	}

	sb.WriteString("\n")
}
