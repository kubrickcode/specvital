package matchers_test

import (
	"context"
	"testing"

	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework/matchers"
)

func TestConfigMatcher_Match(t *testing.T) {
	m := matchers.NewConfigMatcher("jest.config.js", "jest.config.ts")

	tests := []struct {
		name      string
		filename  string
		wantMatch bool
	}{
		{"exact match .js", "jest.config.js", true},
		{"exact match .ts", "jest.config.ts", true},
		{"match with path", "/project/jest.config.js", true},
		{"no match different config", "vitest.config.js", false},
		{"no match similar name", "jest.config.json", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:  framework.SignalConfigFile,
				Value: tt.filename,
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

func TestConfigMatcher_WrongSignalType(t *testing.T) {
	m := matchers.NewConfigMatcher("jest.config.js")

	signal := framework.Signal{
		Type:  framework.SignalImport,
		Value: "jest.config.js",
	}

	result := m.Match(context.Background(), signal)
	if result.Confidence != 0 {
		t.Errorf("expected no match for wrong signal type, got confidence %d", result.Confidence)
	}
}
