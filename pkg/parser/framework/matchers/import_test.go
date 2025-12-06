package matchers_test

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/parser/framework"
	"github.com/specvital/core/pkg/parser/framework/matchers"
)

func TestImportMatcher_ExactMatch(t *testing.T) {
	m := matchers.NewImportMatcher("vitest", "@playwright/test")

	tests := []struct {
		name       string
		importPath string
		wantMatch  bool
	}{
		{"exact match vitest", "vitest", true},
		{"exact match playwright", "@playwright/test", true},
		{"no match", "jest", false},
		{"prefix not matched without slash", "vitest-extra", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalImport,
				Value: tt.importPath,
			}

			result := m.Match(context.Background(), signal)

			if tt.wantMatch {
				if result.Confidence != 100 {
					t.Errorf("expected definite match, got confidence %d", result.Confidence)
				}
			} else {
				if result.Confidence != 0 {
					t.Errorf("expected no match, got confidence %d", result.Confidence)
				}
			}
		})
	}
}

func TestImportMatcher_PrefixMatch(t *testing.T) {
	m := matchers.NewImportMatcher("vitest/", "@jest/")

	tests := []struct {
		name       string
		importPath string
		wantMatch  bool
	}{
		{"prefix match vitest", "vitest/config", true},
		{"prefix match jest", "@jest/globals", true},
		{"no match - different package", "jest", false},
		{"no match - not a prefix", "vitestify", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalImport,
				Value: tt.importPath,
			}

			result := m.Match(context.Background(), signal)

			if tt.wantMatch {
				if result.Confidence != 100 {
					t.Errorf("expected definite match, got confidence %d", result.Confidence)
				}
			} else {
				if result.Confidence != 0 {
					t.Errorf("expected no match, got confidence %d", result.Confidence)
				}
			}
		})
	}
}

func TestImportMatcher_WrongSignalType(t *testing.T) {
	m := matchers.NewImportMatcher("vitest")

	signal := framework.Signal{
		Type:  framework.SignalConfigFile,
		Value: "vitest.config.js",
	}

	result := m.Match(context.Background(), signal)
	if result.Confidence != 0 {
		t.Errorf("expected no match for wrong signal type, got confidence %d", result.Confidence)
	}
}
