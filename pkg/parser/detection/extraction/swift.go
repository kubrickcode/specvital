package extraction

import (
	"context"
	"regexp"
)

// Swift import patterns:
// - import XCTest
// - import Foundation
// - @testable import MyApp

var swiftImportPattern = regexp.MustCompile(`(?m)^(?:@\w+\s+)?import\s+([A-Za-z_][A-Za-z0-9_]*)`)

// ExtractSwiftImports extracts module names from Swift import statements.
func ExtractSwiftImports(_ context.Context, content []byte) []string {
	matches := swiftImportPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	imports := make([]string, 0, len(matches))

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		moduleName := string(match[1])
		if moduleName == "" {
			continue
		}

		if _, ok := seen[moduleName]; ok {
			continue
		}

		seen[moduleName] = struct{}{}
		imports = append(imports, moduleName)
	}

	return imports
}
