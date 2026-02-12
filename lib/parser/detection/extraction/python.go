package extraction

import (
	"context"
	"regexp"
)

var pyImportPattern = regexp.MustCompile(`(?m)^(?:import\s+(\S+)|from\s+(\S+)\s+import\s+(\S+))`)

var (
	pyUnittestMockOnlyPattern      = regexp.MustCompile(`(?m)^from\s+unittest\s+import\s+mock\b`)
	pyUnittestMockSubmodulePattern = regexp.MustCompile(`(?m)^from\s+unittest\.mock\s+import`)
	pyUnittestMockDirectPattern    = regexp.MustCompile(`(?m)^import\s+unittest\.mock\b`)
	pyUnittestDirectPattern        = regexp.MustCompile(`(?m)^import\s+unittest\s*$`)
	pyUnittestFromPattern          = regexp.MustCompile(`(?m)^from\s+unittest\s+import\s+(\w+)`)
)

// ExtractPythonImports parses Python import statements and returns module names.
func ExtractPythonImports(_ context.Context, content []byte) []string {
	hasUnittestMockOnly := pyUnittestMockOnlyPattern.Match(content)
	hasUnittestMockSubmodule := pyUnittestMockSubmodulePattern.Match(content)
	hasUnittestMockDirect := pyUnittestMockDirectPattern.Match(content)
	hasUnittestMock := hasUnittestMockOnly || hasUnittestMockSubmodule || hasUnittestMockDirect

	hasUnittestDirect := pyUnittestDirectPattern.Match(content)

	hasUnittestFromNonMock := false
	fromMatches := pyUnittestFromPattern.FindAllSubmatch(content, -1)
	for _, m := range fromMatches {
		if len(m) > 1 && string(m[1]) != "mock" {
			hasUnittestFromNonMock = true
			break
		}
	}
	hasRealUnittest := hasUnittestDirect || hasUnittestFromNonMock
	usesOnlyUnittestMock := hasUnittestMock && !hasRealUnittest

	matches := pyImportPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	imports := make([]string, 0, len(matches))

	for _, match := range matches {
		var mod string
		if len(match) > 1 && len(match[1]) > 0 {
			mod = string(match[1])
		} else if len(match) > 2 && len(match[2]) > 0 {
			mod = string(match[2])
			if mod == "unittest" && len(match) > 3 && string(match[3]) == "mock" {
				mod = "unittest.mock"
			}
		}

		if mod == "" {
			continue
		}

		if mod == "unittest" && usesOnlyUnittestMock {
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
