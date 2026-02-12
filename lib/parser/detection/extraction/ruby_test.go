package extraction

import (
	"context"
	"reflect"
	"testing"
)

func TestExtractRubyRequires(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "single quotes",
			content: `
require 'rspec'
require 'rspec/core'
`,
			expected: []string{"rspec", "rspec/core"},
		},
		{
			name: "double quotes",
			content: `
require "rspec"
require "rspec/expectations"
`,
			expected: []string{"rspec", "rspec/expectations"},
		},
		{
			name: "require_relative",
			content: `
require_relative 'spec_helper'
require_relative '../support/helpers'
`,
			expected: []string{"spec_helper", "../support/helpers"},
		},
		{
			name: "mixed requires",
			content: `
require 'rspec'
require_relative 'spec_helper'
require "json"
`,
			expected: []string{"rspec", "spec_helper", "json"},
		},
		{
			name: "with leading whitespace",
			content: `
  require 'rspec'
	require 'json'
`,
			expected: []string{"rspec", "json"},
		},
		{
			name: "no requires",
			content: `
class User
  def initialize
  end
end
`,
			expected: nil,
		},
		{
			name: "dedup requires",
			content: `
require 'rspec'
require 'json'
require 'rspec'
`,
			expected: []string{"rspec", "json"},
		},
		{
			name: "ignore comments",
			content: `
# require 'not_loaded'
require 'rspec'
`,
			expected: []string{"rspec"},
		},
		{
			name: "require with parentheses",
			content: `
require('rspec')
require_relative('../spec_helper')
`,
			expected: []string{"rspec", "../spec_helper"},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractRubyRequires(ctx, []byte(tt.content))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ExtractRubyRequires() = %v, want %v", got, tt.expected)
			}
		})
	}
}
