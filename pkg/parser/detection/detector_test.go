package detection

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/framework"
	"github.com/specvital/core/pkg/parser/framework/matchers"
)

// TestDetector_EarlyReturn_ImportWins tests that import detection returns immediately.
func TestDetector_EarlyReturn_ImportWins(t *testing.T) {
	registry := framework.NewRegistry()
	registry.Register(&framework.Definition{
		Name:      "vitest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("vitest", "vitest/"),
		},
	})

	detector := NewDetector(registry)

	// Setup scope for jest (should be ignored due to import)
	projectScope := framework.NewProjectScope()
	jestScope := &framework.ConfigScope{
		ConfigPath:  "/project/jest.config.js",
		BaseDir:     "/project",
		Framework:   "jest",
		GlobalsMode: true,
	}
	projectScope.AddConfig("/project/jest.config.js", jestScope)
	detector.SetProjectScope(projectScope)

	content := []byte(`
import { describe, it, expect } from 'vitest';

describe('test suite', () => {
  it('should work', () => {
    expect(true).toBe(true);
  });
});
`)

	result := detector.Detect(context.Background(), "/project/test.spec.ts", content)

	if result.Framework != "vitest" {
		t.Errorf("expected framework 'vitest', got '%s'", result.Framework)
	}

	if result.Source != SourceImport {
		t.Errorf("expected source 'import', got '%s'", result.Source)
	}
}

// TestDetector_EarlyReturn_FallbackToScope tests scope detection when no imports found.
func TestDetector_EarlyReturn_FallbackToScope(t *testing.T) {
	registry := framework.NewRegistry()
	registry.Register(&framework.Definition{
		Name:      "jest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("@jest/globals", "@jest/globals/"),
		},
	})

	detector := NewDetector(registry)

	projectScope := framework.NewProjectScope()
	jestScope := &framework.ConfigScope{
		ConfigPath:  "/project/jest.config.js",
		BaseDir:     "/project",
		Framework:   "jest",
		GlobalsMode: true,
	}
	projectScope.AddConfig("/project/jest.config.js", jestScope)
	detector.SetProjectScope(projectScope)

	// No imports - should fall back to scope
	content := []byte(`
describe('test suite', () => {
  it('should work', () => {
    expect(true).toBe(true);
  });
});
`)

	result := detector.Detect(context.Background(), "/project/src/test.spec.ts", content)

	if result.Framework != "jest" {
		t.Errorf("expected framework 'jest', got '%s'", result.Framework)
	}

	if result.Source != SourceConfigScope {
		t.Errorf("expected source 'config-scope', got '%s'", result.Source)
	}

	if result.Scope == nil {
		t.Error("expected scope to be set")
	}
}

// TestDetector_EarlyReturn_FallbackToContent tests content detection as last resort.
func TestDetector_EarlyReturn_FallbackToContent(t *testing.T) {
	registry := framework.NewRegistry()
	registry.Register(&framework.Definition{
		Name:      "playwright",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("@playwright/test", "@playwright/test/"),
			matchers.NewContentMatcherFromStrings(`test\.describe\(`),
		},
	})

	detector := NewDetector(registry)

	// No imports, no scope - should use content pattern
	content := []byte(`
test.describe('test suite', () => {
  test('should work', async ({ page }) => {
    await page.goto('https://example.com');
  });
});
`)

	result := detector.Detect(context.Background(), "/project/e2e/test.spec.ts", content)

	if result.Framework != "playwright" {
		t.Errorf("expected framework 'playwright', got '%s'", result.Framework)
	}

	if result.Source != SourceContentPattern {
		t.Errorf("expected source 'content-pattern', got '%s'", result.Source)
	}
}

// TestDetector_EarlyReturn_Unknown tests unknown result when nothing matches.
func TestDetector_EarlyReturn_Unknown(t *testing.T) {
	registry := framework.NewRegistry()
	detector := NewDetector(registry)

	content := []byte(`
console.log('hello world');
`)

	result := detector.Detect(context.Background(), "/project/test.ts", content)

	if result.Framework != "" {
		t.Errorf("expected no framework, got '%s'", result.Framework)
	}

	if result.Source != SourceUnknown {
		t.Errorf("expected source 'unknown', got '%s'", result.Source)
	}

	if result.IsDetected() {
		t.Error("expected IsDetected() to return false")
	}
}

// TestDetector_GoFileNotMatchedByJSFramework tests that Go files are not detected by JS frameworks.
func TestDetector_GoFileNotMatchedByJSFramework(t *testing.T) {
	registry := framework.NewRegistry()

	registry.Register(&framework.Definition{
		Name:      "vitest",
		Languages: []domain.Language{domain.LanguageTypeScript, domain.LanguageJavaScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("vitest", "vitest/"),
		},
	})

	registry.Register(&framework.Definition{
		Name:      "go-testing",
		Languages: []domain.Language{domain.LanguageGo},
		Matchers:  []framework.Matcher{},
	})

	detector := NewDetector(registry)

	// Setup vitest config scope that covers the entire project
	projectScope := framework.NewProjectScope()
	vitestScope := &framework.ConfigScope{
		ConfigPath:  "/project/vitest.config.ts",
		BaseDir:     "/project",
		Framework:   "vitest",
		GlobalsMode: false,
	}
	projectScope.AddConfig("/project/vitest.config.ts", vitestScope)
	detector.SetProjectScope(projectScope)

	content := []byte(`
package env

import "testing"

func TestEnv(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		// test code
	})
}
`)

	result := detector.Detect(context.Background(), "/project/src/go/libs/env/env_test.go", content)

	// Go file should NOT be detected as vitest
	if result.Framework == "vitest" {
		t.Errorf("Go file should not be detected as vitest, got framework '%s'", result.Framework)
	}
}

// TestDetector_ImportPriorityOverScope tests that import takes priority over scope.
func TestDetector_ImportPriorityOverScope(t *testing.T) {
	registry := framework.NewRegistry()

	registry.Register(&framework.Definition{
		Name:      "vitest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("vitest", "vitest/"),
		},
	})

	registry.Register(&framework.Definition{
		Name:      "playwright",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("@playwright/test", "@playwright/test/"),
		},
	})

	detector := NewDetector(registry)

	// Setup vitest scope that covers entire project
	projectScope := framework.NewProjectScope()
	vitestScope := &framework.ConfigScope{
		ConfigPath:  "/project/vitest.config.ts",
		BaseDir:     "/project",
		Framework:   "vitest",
		GlobalsMode: false,
	}
	projectScope.AddConfig("/project/vitest.config.ts", vitestScope)
	detector.SetProjectScope(projectScope)

	// Playwright test file with @playwright/test import
	content := []byte(`
import { expect, test } from "@playwright/test";

test("should work", async ({ page }) => {
	await page.goto("/");
	await expect(page).toHaveTitle("Test");
});
`)

	result := detector.Detect(context.Background(), "/project/src/view/e2e/test.spec.ts", content)

	// Import should win over scope
	if result.Framework != "playwright" {
		t.Errorf("expected framework 'playwright', got '%s'", result.Framework)
	}

	if result.Source != SourceImport {
		t.Errorf("expected source 'import', got '%s'", result.Source)
	}
}

// TestDetector_DeeperScopePriority tests that more specific scope wins.
func TestDetector_DeeperScopePriority(t *testing.T) {
	registry := framework.NewRegistry()

	registry.Register(&framework.Definition{
		Name:      "vitest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("vitest", "vitest/"),
		},
	})

	registry.Register(&framework.Definition{
		Name:      "playwright",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("@playwright/test", "@playwright/test/"),
		},
	})

	detector := NewDetector(registry)

	projectScope := framework.NewProjectScope()
	vitestScope := &framework.ConfigScope{
		ConfigPath:  "/project/vitest.config.ts",
		BaseDir:     "/project", // Shallow
		Framework:   "vitest",
		GlobalsMode: true,
	}
	playwrightScope := &framework.ConfigScope{
		ConfigPath:  "/project/e2e/playwright.config.ts",
		BaseDir:     "/project/e2e", // Deeper
		Framework:   "playwright",
		GlobalsMode: false,
	}
	projectScope.AddConfig("/project/vitest.config.ts", vitestScope)
	projectScope.AddConfig("/project/e2e/playwright.config.ts", playwrightScope)
	detector.SetProjectScope(projectScope)

	// No imports - should use deeper scope
	content := []byte(`
test("should work", async ({ page }) => {
	await page.goto("/");
});
`)

	result := detector.Detect(context.Background(), "/project/e2e/login.spec.ts", content)

	if result.Framework != "playwright" {
		t.Errorf("expected framework 'playwright' (deeper scope), got '%s'", result.Framework)
	}

	if result.Source != SourceConfigScope {
		t.Errorf("expected source 'config-scope', got '%s'", result.Source)
	}
}

// TestDetector_UnsupportedLanguage tests detection for unsupported file types.
func TestDetector_UnsupportedLanguage(t *testing.T) {
	registry := framework.NewRegistry()
	registry.Register(&framework.Definition{
		Name:      "jest",
		Languages: []domain.Language{domain.LanguageTypeScript},
		Matchers:  []framework.Matcher{},
	})

	detector := NewDetector(registry)

	content := []byte(`print("hello world")`)

	result := detector.Detect(context.Background(), "/project/test.py", content)

	if result.Framework != "" {
		t.Errorf("expected no framework for Python file, got '%s'", result.Framework)
	}

	if result.Source != SourceUnknown {
		t.Errorf("expected source 'unknown', got '%s'", result.Source)
	}
}

// TestDetectLanguage tests language detection from file extensions.
func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		path string
		want domain.Language
	}{
		{"/project/test.ts", domain.LanguageTypeScript},
		{"/project/test.tsx", domain.LanguageTypeScript},
		{"/project/test.js", domain.LanguageJavaScript},
		{"/project/test.jsx", domain.LanguageJavaScript},
		{"/project/test.mjs", domain.LanguageJavaScript},
		{"/project/test.cjs", domain.LanguageJavaScript},
		{"/project/test.go", domain.LanguageGo},
		{"/project/test.py", ""},
		{"/project/test.txt", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := detectLanguage(tt.path); got != tt.want {
				t.Errorf("detectLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestResult_String tests Result string representation.
func TestResult_String(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   string
	}{
		{
			name:   "detected",
			result: Confirmed("jest", SourceImport),
			want:   "jest (source: import)",
		},
		{
			name:   "unknown",
			result: Unknown(),
			want:   "no framework detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestResult_IsDetected tests IsDetected method.
func TestResult_IsDetected(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   bool
	}{
		{"import", Confirmed("jest", SourceImport), true},
		{"scope", Confirmed("vitest", SourceConfigScope), true},
		{"content", Confirmed("playwright", SourceContentPattern), true},
		{"unknown", Unknown(), false},
		{"empty framework", Result{Framework: "", Source: SourceImport}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.IsDetected(); got != tt.want {
				t.Errorf("IsDetected() = %v, want %v", got, tt.want)
			}
		})
	}
}
