package prompt

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/kubrickcode/specvital/apps/worker/src/internal/domain/specview"
)

//go:embed templates/placement_system.md
var PlacementSystemPrompt string

// BuildPlacementUserPrompt builds the user prompt for new test placement.
// Converts existing domain/feature structure and new tests to a compact format
// optimized for minimal token usage (~300-500 tokens for typical case).
func BuildPlacementUserPrompt(input specview.PlacementInput) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Target Language: %s\n\n", input.Language)

	// Existing structure (compact format)
	sb.WriteString("<structure>\n")
	for _, domain := range input.ExistingStructure.Domains {
		fmt.Fprintf(&sb, "D:%s\n", domain.Name)
		for _, feature := range domain.Features {
			fmt.Fprintf(&sb, "  F:%s\n", feature.Name)
		}
	}
	sb.WriteString("</structure>\n\n")

	// New tests to place
	sb.WriteString("<new_tests>\n")
	for i, test := range input.NewTests {
		if test.SuitePath != "" {
			fmt.Fprintf(&sb, "%d|%s|%s\n", i, test.SuitePath, test.Name)
		} else {
			fmt.Fprintf(&sb, "%d|%s\n", i, test.Name)
		}
	}
	sb.WriteString("</new_tests>\n\n")

	fmt.Fprintf(&sb, "Place %d tests into existing structure.", len(input.NewTests))

	return sb.String()
}
