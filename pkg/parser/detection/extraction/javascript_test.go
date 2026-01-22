package extraction

import (
	"context"
	"regexp"
	"testing"
)

func TestMatchPatternExcludingComments(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`globals\s*:\s*true`)

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "pattern found without comments",
			content: `{ test: { globals: true } }`,
			want:    true,
		},
		{
			name:    "pattern in single-line comment",
			content: `{ test: { // globals: true } }`,
			want:    false,
		},
		{
			name:    "pattern in multi-line comment",
			content: `{ test: { /* globals: true */ } }`,
			want:    false,
		},
		{
			name:    "regression: glob pattern with /* in string",
			content: `{ test: { include: ["**/*.ts"], globals: true } }`,
			want:    true,
		},
		{
			name:    "regression: glob pattern with */ in string",
			content: `{ test: { exclude: ["view/**/*"], globals: true } }`,
			want:    true,
		},
		{
			name: "regression: multiple glob patterns",
			content: `{
				test: {
					include: ["extension/**/*.ts", "src/**/*.ts"],
					exclude: ["**/node_modules/**", "**/dist/**", "view/**/*"],
					globals: true,
				}
			}`,
			want: true,
		},
		{
			name:    "pattern after single-line comment",
			content: "{ test: { // comment\nglobals: true } }",
			want:    true,
		},
		{
			name:    "pattern after multi-line comment",
			content: `{ test: { /* comment */ globals: true } }`,
			want:    true,
		},
		{
			name:    "empty content",
			content: ``,
			want:    false,
		},
		{
			name:    "no pattern match",
			content: `{ test: { environment: 'node' } }`,
			want:    false,
		},
		// Fast path tests (no comment markers)
		{
			name:    "fast path: no comment markers",
			content: `{ test: { globals: true } }`,
			want:    true,
		},
		{
			name:    "fast path: no comment markers, no match",
			content: `{ test: { globals: false } }`,
			want:    false,
		},
		// Fallback path tests (malformed content)
		{
			name:    "fallback: malformed JS with pattern",
			content: `{ test: { globals: true } // unclosed`,
			want:    true,
		},
		{
			name:    "fallback: malformed JS pattern in comment",
			content: `{ test: { // globals: true`,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MatchPatternExcludingComments(context.Background(), []byte(tt.content), pattern)
			if got != tt.want {
				t.Errorf("MatchPatternExcludingComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasJSTestPatterns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name: "vitest globals mode test file with describe/test",
			content: `
describe('Calculator', () => {
  test('adds two numbers', () => {
    expect(1 + 1).toBe(2);
  });
});
`,
			want: true,
		},
		{
			name: "vitest globals mode test file with describe/it",
			content: `
describe('Calculator', () => {
  it('adds two numbers', () => {
    expect(1 + 1).toBe(2);
  });
});
`,
			want: true,
		},
		{
			name: "test file with beforeEach/afterEach",
			content: `
describe('UserService', () => {
  beforeEach(() => {
    // setup
  });
  afterEach(() => {
    // cleanup
  });
  test('creates user', () => {});
});
`,
			want: true,
		},
		{
			name: "test file with beforeAll/afterAll",
			content: `
describe('Database', () => {
  beforeAll(() => {
    // connect
  });
  afterAll(() => {
    // disconnect
  });
  it('queries data', () => {});
});
`,
			want: true,
		},
		{
			name: "regular source file without test patterns",
			content: `
export function add(a, b) {
  return a + b;
}

export function multiply(a, b) {
  return a * b;
}
`,
			want: false,
		},
		{
			name: "source file with test-like variable names",
			content: `
const testConfig = { enabled: true };
const describe = 'This is a description';
const it = 42;
`,
			want: false,
		},
		{
			name: "test patterns in comments should be ignored",
			content: `
// describe('Calculator', () => {
//   test('adds two numbers', () => {
//     expect(1 + 1).toBe(2);
//   });
// });

export function calculator() {
  return { add: (a, b) => a + b };
}
`,
			want: false,
		},
		{
			name: "test patterns in multi-line comments should be ignored",
			content: `
/*
describe('Calculator', () => {
  test('adds two numbers', () => {
    expect(1 + 1).toBe(2);
  });
});
*/

export const VERSION = '1.0.0';
`,
			want: false,
		},
		{
			name: "top-level test without describe",
			content: `
test('standalone test', () => {
  expect(true).toBe(true);
});
`,
			want: true,
		},
		{
			name: "top-level it without describe",
			content: `
it('standalone it', () => {
  expect(true).toBe(true);
});
`,
			want: true,
		},
		{
			name: "empty file",
			content: ``,
			want:    false,
		},
		{
			name: "vitest test file with explicit import (also has patterns)",
			content: `
import { describe, test, expect } from 'vitest';

describe('Calculator', () => {
  test('adds two numbers', () => {
    expect(1 + 1).toBe(2);
  });
});
`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := HasJSTestPatterns(context.Background(), []byte(tt.content))
			if got != tt.want {
				t.Errorf("HasJSTestPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}
