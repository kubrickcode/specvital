package extraction

import (
	"context"
	"testing"
)

func TestExtractJavaImports(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "junit jupiter imports",
			content: `
package com.example.project;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.DisplayName;
import static org.junit.jupiter.api.Assertions.assertEquals;

class CalculatorTest {
}
`,
			expected: []string{
				"org.junit.jupiter.api.Test",
				"org.junit.jupiter.api.DisplayName",
				"org.junit.jupiter.api.Assertions.assertEquals",
			},
		},
		{
			name: "wildcard import",
			content: `
import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;
`,
			expected: []string{
				"org.junit.jupiter.api.*",
				"org.junit.jupiter.api.Assertions.*",
			},
		},
		{
			name: "no imports",
			content: `
package com.example;

class Simple {
}
`,
			expected: nil,
		},
		{
			name: "parameterized test imports",
			content: `
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;
import org.junit.jupiter.params.provider.CsvSource;
`,
			expected: []string{
				"org.junit.jupiter.params.ParameterizedTest",
				"org.junit.jupiter.params.provider.ValueSource",
				"org.junit.jupiter.params.provider.CsvSource",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractJavaImports(context.Background(), []byte(tt.content))

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d imports, got %d: %v", len(tt.expected), len(result), result)
				return
			}

			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("import[%d]: expected %q, got %q", i, exp, result[i])
				}
			}
		})
	}
}
