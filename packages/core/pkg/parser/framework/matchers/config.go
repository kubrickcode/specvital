package matchers

import (
	"context"
	"path/filepath"

	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework"
)

// ConfigMatcher matches framework-specific configuration files.
// For example: "jest.config.js", "vitest.config.ts"
type ConfigMatcher struct {
	// Patterns is a list of config file name patterns to match.
	// These should be exact file names (e.g., "jest.config.js").
	Patterns []string
}

// NewConfigMatcher creates a ConfigMatcher for the given config file patterns.
func NewConfigMatcher(patterns ...string) *ConfigMatcher {
	return &ConfigMatcher{Patterns: patterns}
}

// Match evaluates if a signal contains a matching config file.
func (m *ConfigMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalConfigFile {
		return framework.NoMatch()
	}

	filename := signal.Value
	base := filepath.Base(filename)

	for _, pattern := range m.Patterns {
		if base == pattern {
			return framework.DefiniteMatch("config: " + base)
		}
	}

	return framework.NoMatch()
}
