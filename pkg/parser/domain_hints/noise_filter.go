package domain_hints

import "strings"

// ShouldFilterImportNoise filters out universal noise patterns from import paths.
// Removes: empty strings, bare relative markers (., ..)
func ShouldFilterImportNoise(importPath string) bool {
	if importPath == "" {
		return true
	}
	// Bare relative markers provide no domain classification value
	// Note: "../utils" or "./helper" are preserved (only bare "." or ".." filtered)
	if importPath == "." || importPath == ".." {
		return true
	}
	return false
}

// ShouldFilterNoise filters out universal noise patterns from domain hints.
// Removes: empty strings, malformed identifiers, and framework-specific noise.
func ShouldFilterNoise(call string) bool {
	if call == "" {
		return true
	}

	// Malformed patterns: starts with "[" (e.g., "[." from spread array handling)
	if call[0] == '[' {
		return true
	}

	// Malformed patterns: starts with "(" (e.g., "(0.", "(1." from decimal literals)
	if call[0] == '(' {
		return true
	}

	// String literal method calls: starts with quote (e.g., "str".encode, 'str'.upper)
	// These are parser artifacts from Python/JS where string literals call methods
	if call[0] == '"' || call[0] == '\'' {
		return true
	}

	// Function arguments leaked into call: contains "=" (e.g., func(key="value"))
	// This indicates parser captured arguments along with function name
	if strings.Contains(call, "=") {
		return true
	}

	// URL patterns leaked into call (e.g., requests.Request("GET","http://example"))
	// This indicates parser captured URL arguments along with function name
	if strings.Contains(call, "http://") || strings.Contains(call, "https://") {
		return true
	}

	// Cheerio/jQuery selector: single "$" is noise
	if call == "$" {
		return true
	}

	// Generic callback variable name: no domain signal
	if call == "fn" {
		return true
	}

	// JavaScript/C-style inline comments leaked into call
	// e.g., "res.json()//Byspec,theruntimecanonly..." from parser including trailing comment
	if strings.Contains(call, "//") {
		return true
	}

	// Short standalone calls (1-2 chars): no domain signal regardless of validity
	// Variables like x, y, f, cb, fn are all generic and don't contribute to domain classification
	// Exception: calls with dots like "io.Reader" are preserved as they indicate package usage
	// But single dot "." or ".." are still noise
	if len(call) <= 2 {
		// Check if it's a valid package.method pattern (has dot with content on both sides)
		dotIdx := strings.Index(call, ".")
		if dotIdx == -1 || dotIdx == 0 || dotIdx == len(call)-1 {
			return true
		}
	}

	// Unbalanced parentheses: parser artifact from method chaining
	// e.g., "json.NewDecoder(w" from "json.NewDecoder(w).Decode(...)"
	// e.g., "expect(mockChromeStorage.session" from "expect(...).toBe(...)"
	if strings.Count(call, "(") != strings.Count(call, ")") {
		return true
	}

	return false
}
