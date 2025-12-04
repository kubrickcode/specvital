package analyzer

import (
	"testing"

	"github.com/specvital/web/src/backend/common/clients/github"
)

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		// JavaScript/TypeScript test files
		{"src/components/Button.test.ts", true},
		{"src/components/Button.test.tsx", true},
		{"src/components/Button.spec.ts", true},
		{"src/components/Button.spec.tsx", true},
		{"src/components/Button.test.js", true},
		{"src/components/Button.spec.js", true},
		{"src/__tests__/Button.ts", true},
		{"src/__tests__/Button.tsx", true},
		{"__tests__/utils.js", true},

		// Go test files
		{"src/service/user_test.go", true},
		{"internal/handler_test.go", true},

		// Non-test files
		{"src/components/Button.ts", false},
		{"src/components/Button.tsx", false},
		{"src/utils/helper.js", false},
		{"main.go", false},
		{"package.json", false},
		{"README.md", false},

		// Files in skip directories
		{"node_modules/package/file.test.ts", false},
		{"vendor/module/file_test.go", false},
		{".git/hooks/pre-commit", false},
		{"dist/bundle.test.js", false},
		{"build/output.test.ts", false},
		{"coverage/report.test.ts", false},
		{"fixtures/sample.test.ts", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("isTestFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestFilterTestFiles(t *testing.T) {
	files := []github.FileInfo{
		{Path: "src/app.ts", Type: "file"},
		{Path: "src/app.test.ts", Type: "file"},
		{Path: "src/components", Type: "dir"},
		{Path: "src/components/Button.tsx", Type: "file"},
		{Path: "src/components/Button.test.tsx", Type: "file"},
		{Path: "node_modules/package/file.test.ts", Type: "file"},
		{Path: "internal/handler_test.go", Type: "file"},
	}

	result := filterTestFiles(files)

	if len(result) != 3 {
		t.Errorf("filterTestFiles() returned %d files, want 3", len(result))
	}

	expectedPaths := map[string]bool{
		"src/app.test.ts":                true,
		"src/components/Button.test.tsx": true,
		"internal/handler_test.go":       true,
	}

	for _, f := range result {
		if !expectedPaths[f.Path] {
			t.Errorf("unexpected file in result: %s", f.Path)
		}
	}
}

func TestShouldSkipPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"node_modules/package/file.ts", true},
		{"src/node_modules/package/file.ts", true},
		{"vendor/module/file.go", true},
		{".git/config", true},
		{"dist/bundle.js", true},
		{"build/output.js", true},
		{"coverage/report.html", true},
		{"fixtures/sample.json", true},

		{"src/components/Button.tsx", false},
		{"internal/handler.go", false},
		{"e2e/login.spec.ts", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldSkipPath(tt.path)
			if result != tt.expected {
				t.Errorf("shouldSkipPath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestMapFramework(t *testing.T) {
	tests := []struct {
		input    string
		expected Framework
	}{
		{"jest", FrameworkJest},
		{"Jest", FrameworkJest},
		{"JEST", FrameworkJest},
		{"vitest", FrameworkVitest},
		{"Vitest", FrameworkVitest},
		{"playwright", FrameworkPlaywright},
		{"Playwright", FrameworkPlaywright},
		{"gotesting", FrameworkGo},
		{"go", FrameworkGo},
		{"Go", FrameworkGo},
		{"unknown", FrameworkJest},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapFramework(tt.input)
			if result != tt.expected {
				t.Errorf("mapFramework(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalculateSummary(t *testing.T) {
	suites := []TestSuite{
		{
			FilePath:  "src/app.test.ts",
			Framework: FrameworkJest,
			Tests: []TestCase{
				{Status: TestStatusActive},
				{Status: TestStatusActive},
				{Status: TestStatusSkipped},
			},
		},
		{
			FilePath:  "src/utils.test.ts",
			Framework: FrameworkVitest,
			Tests: []TestCase{
				{Status: TestStatusActive},
				{Status: TestStatusTodo},
			},
		},
	}

	summary := calculateSummary(suites)

	if summary.Total != 5 {
		t.Errorf("summary.Total = %d, want 5", summary.Total)
	}
	if summary.Active != 3 {
		t.Errorf("summary.Active = %d, want 3", summary.Active)
	}
	if summary.Skipped != 1 {
		t.Errorf("summary.Skipped = %d, want 1", summary.Skipped)
	}
	if summary.Todo != 1 {
		t.Errorf("summary.Todo = %d, want 1", summary.Todo)
	}
	if len(summary.Frameworks) != 2 {
		t.Errorf("len(summary.Frameworks) = %d, want 2", len(summary.Frameworks))
	}
}

func TestCalculateSummaryEmpty(t *testing.T) {
	summary := calculateSummary([]TestSuite{})

	if summary.Total != 0 {
		t.Errorf("summary.Total = %d, want 0", summary.Total)
	}
	if len(summary.Frameworks) != 0 {
		t.Errorf("len(summary.Frameworks) = %d, want 0", len(summary.Frameworks))
	}
}
