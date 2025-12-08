package extraction

import (
	"context"
	"regexp"
)

// Ruby require patterns:
// - require 'foo'
// - require "foo"
// - require('foo')
// - require_relative 'foo'
// - require_relative('../foo')
//
// Pattern uses possessive-like matching with explicit character class
// to avoid ReDoS vulnerabilities. Module names are alphanumeric with / _ - .

var rbRequirePattern = regexp.MustCompile(`(?m)^\s*require(?:_relative)?\s*\(?['"]([a-zA-Z0-9_./-]+)['"]\)?`)

// ExtractRubyRequires extracts module names from Ruby require statements.
func ExtractRubyRequires(_ context.Context, content []byte) []string {
	matches := rbRequirePattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	requires := make([]string, 0, len(matches))

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		mod := string(match[1])
		if mod == "" {
			continue
		}

		if _, ok := seen[mod]; ok {
			continue
		}

		seen[mod] = struct{}{}
		requires = append(requires, mod)
	}

	return requires
}
