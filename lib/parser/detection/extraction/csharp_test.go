package extraction

import (
	"context"
	"testing"
)

func TestExtractCSharpUsings(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "context cancellation returns nil",
			content:  `using Xunit;`,
			expected: nil, // Handled separately with cancelled context
		},
		{
			name: "simple using",
			content: `
using Xunit;

public class CalculatorTests {}
`,
			expected: []string{"Xunit"},
		},
		{
			name: "multiple usings",
			content: `
using System;
using Xunit;
using FluentAssertions;

public class CalculatorTests {}
`,
			expected: []string{"System", "Xunit", "FluentAssertions"},
		},
		{
			name: "nested namespace",
			content: `
using Microsoft.VisualStudio.TestTools.UnitTesting;

public class CalculatorTests {}
`,
			expected: []string{"Microsoft.VisualStudio.TestTools.UnitTesting"},
		},
		{
			name: "no usings",
			content: `
public class Calculator {}
`,
			expected: nil,
		},
		{
			name: "NUnit framework",
			content: `
using NUnit.Framework;

public class CalculatorTests {}
`,
			expected: []string{"NUnit.Framework"},
		},
		{
			name: "using static",
			content: `
using static System.Math;

public class CalculatorTests {}
`,
			expected: []string{"System.Math"},
		},
		{
			name: "global using",
			content: `
global using Xunit;

public class CalculatorTests {}
`,
			expected: []string{"Xunit"},
		},
		{
			name: "global using static",
			content: `
global using static System.Console;

public class CalculatorTests {}
`,
			expected: []string{"System.Console"},
		},
		{
			name: "using alias should be ignored",
			content: `
using MyAlias = System.Collections.Generic;
using Xunit;

public class CalculatorTests {}
`,
			expected: []string{"Xunit"},
		},
		{
			name: "global using alias should be ignored",
			content: `
global using Project = PC.MyCompany.Project;
using Xunit;
using NUnit.Framework;

public class CalculatorTests {}
`,
			expected: []string{"Xunit", "NUnit.Framework"},
		},
		{
			name: "mixed aliases and regular usings",
			content: `
using System;
using Dict = System.Collections.Generic.Dictionary<string, int>;
using Xunit;
using Assert = Xunit.Assert;
using FluentAssertions;

public class CalculatorTests {}
`,
			expected: []string{"System", "Xunit", "FluentAssertions"},
		},
		{
			name:     "malformed double dot",
			content:  `using A..B;`,
			expected: nil,
		},
		{
			name:     "malformed leading dot",
			content:  `using .A;`,
			expected: nil,
		},
		{
			name:     "malformed trailing dot",
			content:  `using A.;`,
			expected: nil,
		},
		{
			name: "whitespace variations",
			content: `
using   System  ;
using	Xunit	;
`,
			expected: []string{"System", "Xunit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip context cancellation test in main loop
			if tt.name == "context cancellation returns nil" {
				return
			}

			ctx := context.Background()
			result := ExtractCSharpUsings(ctx, []byte(tt.content))

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d usings, got %d: %v", len(tt.expected), len(result), result)
				return
			}

			for i, exp := range tt.expected {
				if result[i] != exp {
					t.Errorf("expected result[%d]='%s', got '%s'", i, exp, result[i])
				}
			}
		})
	}

	// Dedicated test for context cancellation
	t.Run("context cancellation returns nil", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result := ExtractCSharpUsings(ctx, []byte(`using Xunit;`))
		if result != nil {
			t.Errorf("expected nil for cancelled context, got %v", result)
		}
	})
}
