package extraction

import (
	"context"
	"regexp"
)

// Python import patterns:
// - import foo
// - import foo.bar
// - from foo import bar
// - from foo.bar import baz

var pyImportPattern = regexp.MustCompile(`(?m)^(?:import\s+(\S+)|from\s+(\S+)\s+import)`)

// ExtractPythonImports extracts module names from Python import statements.
func ExtractPythonImports(_ context.Context, content []byte) []string {
	matches := pyImportPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	imports := make([]string, 0, len(matches))

	for _, match := range matches {
		var mod string
		if len(match) > 1 && len(match[1]) > 0 {
			mod = string(match[1]) // import foo
		} else if len(match) > 2 && len(match[2]) > 0 {
			mod = string(match[2]) // from foo import bar
		}

		if mod == "" {
			continue
		}

		if _, ok := seen[mod]; ok {
			continue
		}

		seen[mod] = struct{}{}
		imports = append(imports, mod)
	}

	return imports
}
