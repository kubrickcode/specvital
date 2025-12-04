package parser

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestDetectTestFiles(t *testing.T) {
	ctx := context.Background()

	t.Run("should return error for non-existent path", func(t *testing.T) {
		_, err := DetectTestFiles(ctx, "/non/existent/path")

		if !errors.Is(err, ErrInvalidRootPath) {
			t.Errorf("expected ErrInvalidRootPath, got %v", err)
		}
	})

	t.Run("should return error for file path instead of directory", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		_, err = DetectTestFiles(ctx, tmpFile.Name())

		if !errors.Is(err, ErrInvalidRootPath) {
			t.Errorf("expected ErrInvalidRootPath, got %v", err)
		}
	})

	t.Run("should detect JavaScript test files", func(t *testing.T) {
		tmpDir := t.TempDir()
		createTestFile(t, tmpDir, "user.test.ts", "")
		createTestFile(t, tmpDir, "auth.spec.js", "")
		createTestFile(t, tmpDir, "utils.ts", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 2 {
			t.Errorf("expected 2 files, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should detect Go test files", func(t *testing.T) {
		tmpDir := t.TempDir()
		createTestFile(t, tmpDir, "user_test.go", "")
		createTestFile(t, tmpDir, "user.go", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should detect files in __tests__ directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		testsDir := filepath.Join(tmpDir, "__tests__")
		if err := os.MkdirAll(testsDir, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, testsDir, "user.ts", "")
		createTestFile(t, testsDir, "auth.tsx", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 2 {
			t.Errorf("expected 2 files, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should skip node_modules directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		nodeModules := filepath.Join(tmpDir, "node_modules", "some-package")
		if err := os.MkdirAll(nodeModules, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, nodeModules, "index.test.ts", "")
		createTestFile(t, tmpDir, "app.test.ts", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should skip .git directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		gitDir := filepath.Join(tmpDir, ".git", "hooks")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, gitDir, "pre-commit.test.ts", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 0 {
			t.Errorf("expected 0 files, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should skip vendor directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		vendorDir := filepath.Join(tmpDir, "vendor", "github.com")
		if err := os.MkdirAll(vendorDir, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, vendorDir, "lib_test.go", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 0 {
			t.Errorf("expected 0 files, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should detect files in nested directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		nested := filepath.Join(tmpDir, "src", "components", "user")
		if err := os.MkdirAll(nested, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, nested, "user.test.tsx", "")

		result, err := DetectTestFiles(ctx, tmpDir)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("should stop on context cancellation", func(t *testing.T) {
		tmpDir := t.TempDir()
		createTestFile(t, tmpDir, "app.test.ts", "")

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := DetectTestFiles(canceledCtx, tmpDir)

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("with options/should detect all files when patterns is empty slice", func(t *testing.T) {
		tmpDir := t.TempDir()
		createTestFile(t, tmpDir, "a.test.ts", "")
		createTestFile(t, tmpDir, "b.test.ts", "")

		result, err := DetectTestFiles(ctx, tmpDir, WithPatterns([]string{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 2 {
			t.Errorf("expected 2 files when patterns is empty, got %d", len(result.Files))
		}
	})

	t.Run("with options/should use custom skip patterns", func(t *testing.T) {
		tmpDir := t.TempDir()
		customSkip := filepath.Join(tmpDir, "build")
		if err := os.MkdirAll(customSkip, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, customSkip, "output.test.ts", "")
		createTestFile(t, tmpDir, "app.test.ts", "")

		result, err := DetectTestFiles(ctx, tmpDir, WithSkipPatterns([]string{"build"}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("with options/should filter by glob patterns", func(t *testing.T) {
		tmpDir := t.TempDir()
		srcDir := filepath.Join(tmpDir, "src")
		libDir := filepath.Join(tmpDir, "lib")
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(libDir, 0755); err != nil {
			t.Fatal(err)
		}
		createTestFile(t, srcDir, "app.test.ts", "")
		createTestFile(t, libDir, "utils.test.ts", "")

		result, err := DetectTestFiles(ctx, tmpDir, WithPatterns([]string{"src/**/*.test.ts"}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d: %v", len(result.Files), result.Files)
		}
		if len(result.Files) > 0 && filepath.Base(result.Files[0]) != "app.test.ts" {
			t.Errorf("expected app.test.ts, got %s", filepath.Base(result.Files[0]))
		}
	})

	t.Run("with options/should skip files exceeding max size", func(t *testing.T) {
		tmpDir := t.TempDir()
		smallFile := filepath.Join(tmpDir, "small.test.ts")
		largeFile := filepath.Join(tmpDir, "large.test.ts")

		if err := os.WriteFile(smallFile, []byte("small"), 0644); err != nil {
			t.Fatal(err)
		}
		largeContent := make([]byte, 1024)
		if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
			t.Fatal(err)
		}

		result, err := DetectTestFiles(ctx, tmpDir, WithMaxFileSize(100))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file, got %d: %v", len(result.Files), result.Files)
		}
	})

	t.Run("with options/should not limit file size when max size is 0", func(t *testing.T) {
		tmpDir := t.TempDir()
		largeFile := filepath.Join(tmpDir, "large.test.ts")
		largeContent := make([]byte, 1024*1024) // 1MB
		if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
			t.Fatal(err)
		}

		result, err := DetectTestFiles(ctx, tmpDir, WithMaxFileSize(0))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Files) != 1 {
			t.Errorf("expected 1 file when maxSize=0 (no limit), got %d", len(result.Files))
		}
	})
}

func TestIsTestFileCandidate(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"TypeScript test file", "user.test.ts", true},
		{"TypeScript spec file", "user.spec.ts", true},
		{"JavaScript test file", "user.test.js", true},
		{"TSX test file", "Component.test.tsx", true},
		{"JSX spec file", "Component.spec.jsx", true},
		{"Go test file", "user_test.go", true},
		{"Regular TypeScript file", "user.ts", false},
		{"Regular Go file", "user.go", false},
		{"__tests__ directory file", "__tests__/user.ts", true},
		{"Python file", "user_test.py", false},
		{"Unsupported extension", "user.test.rb", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTestFileCandidate(tt.path)
			if result != tt.expected {
				t.Errorf("isTestFileCandidate(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsJSTestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"JS test file", "user.test.js", true},
		{"TS spec file", "user.spec.ts", true},
		{"Case insensitive test", "User.TEST.ts", true},
		{"Case insensitive spec", "User.SPEC.ts", true},
		{"__tests__ directory", "src/__tests__/user.ts", true},
		{"__tests__ at root", "__tests__/user.ts", true},
		{"Regular file", "user.ts", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("isJSTestFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsGoTestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Go test file", "user_test.go", true},
		{"Go regular file", "user.go", false},
		{"Nested path", "pkg/user_test.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGoTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("isGoTestFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestBuildSkipSet(t *testing.T) {
	patterns := []string{"node_modules", ".git", "vendor"}

	result := buildSkipSet(patterns)

	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
	for _, p := range patterns {
		if !result[p] {
			t.Errorf("expected %q to be in skip set", p)
		}
	}
}

func TestShouldSkipDir(t *testing.T) {
	skipSet := map[string]bool{
		"node_modules": true,
		".git":         true,
	}

	tests := []struct {
		name     string
		path     string
		rootPath string
		expected bool
	}{
		{"root path", "/project", "/project", false},
		{"node_modules", "/project/node_modules", "/project", true},
		{".git", "/project/.git", "/project", true},
		{"src directory", "/project/src", "/project", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipDir(tt.path, tt.rootPath, skipSet)
			if result != tt.expected {
				t.Errorf("shouldSkipDir(%q, %q) = %v, want %v", tt.path, tt.rootPath, result, tt.expected)
			}
		})
	}
}

func TestMatchesAnyPattern(t *testing.T) {
	rootPath := "/project"

	tests := []struct {
		name     string
		path     string
		patterns []string
		expected bool
	}{
		{
			"match single pattern",
			"/project/src/user.test.ts",
			[]string{"src/**/*.test.ts"},
			true,
		},
		{
			"match one of multiple patterns",
			"/project/lib/utils.spec.js",
			[]string{"src/**/*.test.ts", "lib/**/*.spec.js"},
			true,
		},
		{
			"no match",
			"/project/other/file.test.ts",
			[]string{"src/**/*.test.ts"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesAnyPattern(tt.path, rootPath, tt.patterns)
			if result != tt.expected {
				t.Errorf("matchesAnyPattern(%q, %q, %v) = %v, want %v",
					tt.path, rootPath, tt.patterns, result, tt.expected)
			}
		})
	}
}

func TestDefaultSkipPatterns(t *testing.T) {
	expectedPatterns := []string{
		"node_modules",
		".git",
		"vendor",
		"dist",
	}

	for _, p := range expectedPatterns {
		if !slices.Contains(DefaultSkipPatterns, p) {
			t.Errorf("expected %q to be in DefaultSkipPatterns", p)
		}
	}
}

func createTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", path, err)
	}
}
