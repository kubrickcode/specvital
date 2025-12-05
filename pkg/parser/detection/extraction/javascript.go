package extraction

import (
	"context"
	"regexp"
)

var jsImportPattern = regexp.MustCompile(`(?:import\s+.*?\s+from|require\()\s*['"]([^'"]+)['"]`)

func ExtractJSImports(_ context.Context, content []byte) []string {
	matches := jsImportPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	imports := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			imports = append(imports, string(match[1]))
		}
	}
	return imports
}

var commentStripRegex = regexp.MustCompile(`//.*|/\*[\s\S]*?\*/`)

// MatchPatternExcludingComments checks if pattern matches content after stripping comments.
// Handles both single-line (//) and multi-line (/* */) comments.
// Limitation: Does not handle comments inside string literals.
func MatchPatternExcludingComments(content []byte, pattern *regexp.Regexp) bool {
	noComments := commentStripRegex.ReplaceAll(content, []byte{})
	return pattern.Match(noComments)
}
