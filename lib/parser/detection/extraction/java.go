package extraction

import (
	"context"
	"regexp"
)

// Java import patterns:
// - import org.junit.jupiter.api.Test;
// - import static org.junit.jupiter.api.Assertions.*;

var javaImportPattern = regexp.MustCompile(`(?m)^import\s+(?:static\s+)?([a-zA-Z_][a-zA-Z0-9_.]*(?:\.\*)?);`)

// ExtractJavaImports extracts package names from Java import statements.
func ExtractJavaImports(_ context.Context, content []byte) []string {
	matches := javaImportPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	imports := make([]string, 0, len(matches))

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		importPath := string(match[1])
		if importPath == "" {
			continue
		}

		if _, ok := seen[importPath]; ok {
			continue
		}

		seen[importPath] = struct{}{}
		imports = append(imports, importPath)
	}

	return imports
}
