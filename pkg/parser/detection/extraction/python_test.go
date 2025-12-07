package extraction

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractPythonImports(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "import statement",
			content: `
import pytest
import os
`,
			expected: []string{"pytest", "os"},
		},
		{
			name: "from import statement",
			content: `
from pytest import fixture
from os import path
`,
			expected: []string{"pytest", "os"},
		},
		{
			name: "mixed imports",
			content: `
import pytest
from os import path
import sys
from collections import defaultdict
`,
			expected: []string{"pytest", "os", "sys", "collections"},
		},
		{
			name: "dotted module names",
			content: `
import pytest.mark
from pytest.mark import parametrize
`,
			expected: []string{"pytest.mark"},
		},
		{
			name: "no imports",
			content: `
def test_something():
    pass
`,
			expected: nil,
		},
		{
			name: "dedup imports",
			content: `
import pytest
from pytest import fixture
import pytest
`,
			expected: []string{"pytest"},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPythonImports(ctx, []byte(tt.content))
			assert.Equal(t, tt.expected, result)
		})
	}
}
