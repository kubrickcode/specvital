package matchers

import (
	"context"
	"regexp"

	"github.com/specvital/core/pkg/parser/framework"
)

// ContentMatcher matches framework-specific content patterns in test files.
// For example: "test.describe(" for Playwright, "describe(" for Jest/Vitest with globals.
type ContentMatcher struct {
	// Patterns is a list of regex patterns to match in file content.
	// Patterns should match framework-specific test function calls or syntax.
	Patterns []*regexp.Regexp
}

// NewContentMatcher creates a ContentMatcher for the given regex patterns.
func NewContentMatcher(patterns ...*regexp.Regexp) *ContentMatcher {
	return &ContentMatcher{Patterns: patterns}
}

// NewContentMatcherFromStrings creates a ContentMatcher from regex pattern strings.
// Panics if any pattern fails to compile (use during initialization only).
func NewContentMatcherFromStrings(patterns ...string) *ContentMatcher {
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexps[i] = regexp.MustCompile(pattern)
	}
	return &ContentMatcher{Patterns: regexps}
}

// Match evaluates if a signal contains matching content patterns.
func (m *ContentMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileContent {
		return framework.NoMatch()
	}

	content, ok := signal.Context.([]byte)
	if !ok {
		return framework.NoMatch()
	}

	var evidence []string
	for _, pattern := range m.Patterns {
		if pattern.Match(content) {
			evidence = append(evidence, "pattern: "+pattern.String())
		}
	}

	if len(evidence) > 0 {
		return framework.PartialMatch(40, evidence...)
	}

	return framework.NoMatch()
}
