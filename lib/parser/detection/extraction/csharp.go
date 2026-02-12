package extraction

import (
	"context"
	"regexp"
)

// C# using patterns:
// - using Xunit;
// - using NUnit.Framework;
// - using Microsoft.VisualStudio.TestTools.UnitTesting;
//
// Explicitly excludes:
// - using aliases: using MyAlias = Some.Namespace;
// - using directives with aliases are skipped to avoid false matches

// csharpUsingPattern matches C# using statements with proper namespace validation.
// Pattern ensures dots are always followed by valid identifiers to prevent backtracking.
var csharpUsingPattern = regexp.MustCompile(`(?m)^(?:global\s+)?using\s+(?:static\s+)?([A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*)\s*;`)

// csharpUsingAliasPattern matches using alias declarations to filter them out.
var csharpUsingAliasPattern = regexp.MustCompile(`(?m)^(?:global\s+)?using\s+[A-Za-z_][A-Za-z0-9_]*\s*=`)

// ExtractCSharpUsings extracts namespace names from C# using statements.
// It filters out using alias declarations (e.g., using MyAlias = Some.Namespace).
func ExtractCSharpUsings(ctx context.Context, content []byte) []string {
	// Check context before expensive operations
	if err := ctx.Err(); err != nil {
		return nil
	}

	// Build alias position set using byte positions (O(n) vs O(n²) with line counting)
	aliasPositions := make(map[int]struct{})
	for _, loc := range csharpUsingAliasPattern.FindAllIndex(content, -1) {
		aliasPositions[loc[0]] = struct{}{}
	}

	matches := csharpUsingPattern.FindAllSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	usings := make([]string, 0, len(matches))

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		// Check if this match starts at an alias position
		if _, isAlias := aliasPositions[match[0]]; isAlias {
			continue
		}

		namespace := string(content[match[2]:match[3]])
		if namespace == "" {
			continue
		}

		if _, ok := seen[namespace]; ok {
			continue
		}

		seen[namespace] = struct{}{}
		usings = append(usings, namespace)
	}

	return usings
}
