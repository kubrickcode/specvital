package analyzer

import (
	"context"
	"testing"

	"github.com/specvital/web/src/backend/common/clients/github"
)

// Test helpers

type mockGitHost struct {
	files        map[string]string
	latestCommit *github.CommitInfo
	fileList     []github.FileInfo
	rateLimit    github.RateLimitInfo
}

func (m *mockGitHost) GetLatestCommit(ctx context.Context, owner, repo string) (*github.CommitInfo, error) {
	return m.latestCommit, nil
}

func (m *mockGitHost) ListFiles(ctx context.Context, owner, repo string) ([]github.FileInfo, error) {
	return m.fileList, nil
}

func (m *mockGitHost) GetFileContent(ctx context.Context, owner, repo, path string) (string, error) {
	if content, ok := m.files[path]; ok {
		return content, nil
	}
	return "", github.ErrNotFound
}

func (m *mockGitHost) GetRateLimit() github.RateLimitInfo {
	return m.rateLimit
}

func assertFilePathsMatch(t *testing.T, got []github.FileInfo, wantPaths []string) {
	t.Helper()
	if len(got) != len(wantPaths) {
		t.Errorf("got %d files, want %d", len(got), len(wantPaths))
		return
	}

	wantMap := make(map[string]bool)
	for _, p := range wantPaths {
		wantMap[p] = true
	}

	for _, f := range got {
		if !wantMap[f.Path] {
			t.Errorf("unexpected file: %s", f.Path)
		}
	}
}

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

		// Python test files
		{"tests/test_user.py", true},
		{"test_main.py", true},
		{"src/user_test.py", true},
		{"tests/conftest.py", true},
		{"conftest.py", true},

		// Java test files
		{"src/test/java/UserTest.java", true},
		{"src/test/java/UserTests.java", true},
		{"src/test/java/TestUserService.java", true},

		// C#/xUnit test files
		{"Tests/UserTests.cs", true},
		{"Tests/UserTest.cs", true},
		{"Tests/UserSpecs.cs", true},
		{"Tests/UserSpec.cs", true},
		{"Tests/FluentAssertions.Specs/BasicSpecs.cs", true},
		{"src/MyProject.Tests/UserTests.cs", true},

		// Non-test files
		{"src/components/Button.ts", false},
		{"src/components/Button.tsx", false},
		{"src/utils/helper.js", false},
		{"main.go", false},
		{"package.json", false},
		{"README.md", false},
		{"src/User.cs", false},
		{"src/Services/UserService.cs", false},

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
		// Normalized names
		{"go-testing", "go"},
		// Pass-through (core names used as-is)
		{"jest", "jest"},
		{"vitest", "vitest"},
		{"playwright", "playwright"},
		{"pytest", "pytest"},
		{"unknown", "unknown"},
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
			Framework: "jest",
			Tests: []TestCase{
				{Status: TestStatusActive},
				{Status: TestStatusActive},
				{Status: TestStatusSkipped},
			},
		},
		{
			FilePath:  "src/utils.test.ts",
			Framework: "vitest",
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

func TestGetConfigPatterns(t *testing.T) {
	patterns := getConfigPatterns()

	// Should return non-empty map from core's matchers
	if len(patterns) == 0 {
		t.Error("getConfigPatterns() returned empty map, expected patterns from core")
	}

	// Should contain known vitest config patterns
	expectedPatterns := []string{
		"vitest.config.ts",
		"vitest.config.js",
		"jest.config.ts",
		"jest.config.js",
		"playwright.config.ts",
	}

	for _, p := range expectedPatterns {
		if !patterns[p] {
			t.Errorf("expected pattern %q not found in config patterns", p)
		}
	}

	// Should be idempotent (sync.Once)
	patterns2 := getConfigPatterns()
	if len(patterns) != len(patterns2) {
		t.Error("getConfigPatterns() not idempotent")
	}
}

func TestBuildProjectScope(t *testing.T) {
	tests := []struct {
		name         string
		configFiles  []github.FileInfo
		fileContents map[string]string
		wantNil      bool
	}{
		{
			name:         "empty config files returns non-nil scope",
			configFiles:  []github.FileInfo{},
			fileContents: map[string]string{},
			wantNil:      false,
		},
		{
			name: "vitest config with globals",
			configFiles: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
			},
			fileContents: map[string]string{
				"vitest.config.ts": `import { defineConfig } from 'vitest/config'
export default defineConfig({
  test: {
    globals: true,
  },
})`,
			},
			wantNil: false,
		},
		{
			name: "multiple config files",
			configFiles: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
				{Path: "playwright.config.ts", Type: "file"},
			},
			fileContents: map[string]string{
				"vitest.config.ts":     `export default { test: { globals: true } }`,
				"playwright.config.ts": `export default { use: { headless: true } }`,
			},
			wantNil: false,
		},
		{
			name: "missing config file content gracefully handled",
			configFiles: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
			},
			fileContents: map[string]string{},
			wantNil:      false,
		},
		{
			name: "malformed config content handled gracefully",
			configFiles: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
			},
			fileContents: map[string]string{
				"vitest.config.ts": `this is { invalid javascript syntax`,
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockGitHost{
				files: tt.fileContents,
			}
			service := NewService(mock)

			result := service.buildProjectScope(context.Background(), "owner", "repo", tt.configFiles)

			if tt.wantNil && result != nil {
				t.Error("expected nil ProjectScope, got non-nil")
			}
			if !tt.wantNil && result == nil {
				t.Error("expected non-nil ProjectScope, got nil")
			}
		})
	}
}

func TestFilterConfigFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    []github.FileInfo
		expected []string
	}{
		{
			name: "filters vitest config files",
			files: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
				{Path: "src/app.ts", Type: "file"},
				{Path: "src/app.test.ts", Type: "file"},
			},
			expected: []string{"vitest.config.ts"},
		},
		{
			name: "filters multiple config files",
			files: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
				{Path: "jest.config.js", Type: "file"},
				{Path: "playwright.config.ts", Type: "file"},
				{Path: "src/index.ts", Type: "file"},
			},
			expected: []string{"vitest.config.ts", "jest.config.js", "playwright.config.ts"},
		},
		{
			name: "filters nested config files",
			files: []github.FileInfo{
				{Path: "packages/app/vitest.config.ts", Type: "file"},
				{Path: "packages/lib/jest.config.js", Type: "file"},
			},
			expected: []string{"packages/app/vitest.config.ts", "packages/lib/jest.config.js"},
		},
		{
			name: "ignores directories",
			files: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "dir"},
				{Path: "jest.config.js", Type: "file"},
			},
			expected: []string{"jest.config.js"},
		},
		{
			name: "returns empty for no config files",
			files: []github.FileInfo{
				{Path: "src/app.ts", Type: "file"},
				{Path: "package.json", Type: "file"},
			},
			expected: []string{},
		},
		{
			name: "handles all vitest config extensions",
			files: []github.FileInfo{
				{Path: "vitest.config.ts", Type: "file"},
				{Path: "vitest.config.js", Type: "file"},
				{Path: "vitest.config.mts", Type: "file"},
				{Path: "vitest.config.mjs", Type: "file"},
			},
			expected: []string{"vitest.config.ts", "vitest.config.js", "vitest.config.mts", "vitest.config.mjs"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterConfigFiles(tt.files)
			assertFilePathsMatch(t, result, tt.expected)
		})
	}
}
