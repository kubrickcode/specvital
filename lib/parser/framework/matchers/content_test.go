package matchers_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/kubrickcode/specvital/lib/parser/framework"
	"github.com/kubrickcode/specvital/lib/parser/framework/matchers"
)

func TestContentMatcher_Match(t *testing.T) {
	m := matchers.NewContentMatcherFromStrings(
		`\bdescribe\s*\(`,
		`\btest\s*\(`,
		`\bit\s*\(`,
	)

	tests := []struct {
		name      string
		content   string
		wantMatch bool
	}{
		{
			name:      "matches describe",
			content:   `describe("test suite", () => {})`,
			wantMatch: true,
		},
		{
			name:      "matches test",
			content:   `test("test case", () => {})`,
			wantMatch: true,
		},
		{
			name:      "matches it",
			content:   `it("should work", () => {})`,
			wantMatch: true,
		},
		{
			name:      "no match",
			content:   `const foo = "bar";`,
			wantMatch: false,
		},
		{
			name:      "matches with whitespace variations",
			content:   `describe  ("test", () => {})`,
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := framework.Signal{
				Type:    framework.SignalFileContent,
				Value:   "",
				Context: []byte(tt.content),
			}

			result := m.Match(context.Background(), signal)

			if tt.wantMatch {
				if result.Confidence == 0 {
					t.Error("expected match, got no match")
				}
				if len(result.Evidence) == 0 {
					t.Error("expected evidence for match")
				}
			} else {
				if result.Confidence != 0 {
					t.Errorf("expected no match, got confidence %d", result.Confidence)
				}
			}
		})
	}
}

func TestContentMatcher_WrongSignalType(t *testing.T) {
	m := matchers.NewContentMatcher(regexp.MustCompile(`test`))

	signal := framework.Signal{
		Type:  framework.SignalImport,
		Value: "test content",
	}

	result := m.Match(context.Background(), signal)
	if result.Confidence != 0 {
		t.Errorf("expected no match for wrong signal type, got confidence %d", result.Confidence)
	}
}

func TestContentMatcher_InvalidContext(t *testing.T) {
	m := matchers.NewContentMatcher(regexp.MustCompile(`test`))

	signal := framework.Signal{
		Type:    framework.SignalFileContent,
		Value:   "",
		Context: "invalid type", // Should be []byte
	}

	result := m.Match(context.Background(), signal)
	if result.Confidence != 0 {
		t.Errorf("expected no match for invalid context type, got confidence %d", result.Confidence)
	}
}
