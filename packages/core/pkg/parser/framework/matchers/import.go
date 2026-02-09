// Package matchers provides reusable matcher implementations for framework detection.
package matchers

import (
	"context"
	"strings"

	"github.com/specvital/core/pkg/parser/framework"
)

// ImportMatcher matches framework-specific import statements.
// For example: "import { test } from 'vitest'" matches Vitest.
type ImportMatcher struct {
	// Patterns is a list of import path patterns to match.
	// Supports exact matches and prefix matching.
	// Examples: ["vitest", "vitest/"], ["@playwright/test", "@playwright/test/"]
	Patterns []string
}

// NewImportMatcher creates an ImportMatcher for the given import patterns.
func NewImportMatcher(patterns ...string) *ImportMatcher {
	return &ImportMatcher{Patterns: patterns}
}

// Match evaluates if a signal contains a matching import statement.
func (m *ImportMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalImport {
		return framework.NoMatch()
	}

	importPath := signal.Value
	for _, pattern := range m.Patterns {
		if m.matchesPattern(importPath, pattern) {
			return framework.DefiniteMatch("import: " + importPath)
		}
	}

	return framework.NoMatch()
}

func (m *ImportMatcher) matchesPattern(importPath, pattern string) bool {
	if importPath == pattern {
		return true
	}
	if strings.HasSuffix(pattern, "/") && strings.HasPrefix(importPath, pattern) {
		return true
	}
	return false
}
