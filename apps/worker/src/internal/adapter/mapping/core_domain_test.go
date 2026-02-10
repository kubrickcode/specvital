package mapping

import (
	"errors"
	"testing"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	coreparser "github.com/kubrickcode/specvital/packages/core/pkg/parser"
	"github.com/kubrickcode/specvital/apps/worker/src/internal/domain/analysis"
)

func TestConvertCoreToDomainInventory_Nil(t *testing.T) {
	result := ConvertCoreToDomainInventory(nil)
	if result != nil {
		t.Errorf("expected nil for nil input, got %v", result)
	}
}

func TestConvertCoreToDomainInventory_Empty(t *testing.T) {
	coreInv := &domain.Inventory{
		Files:    []domain.TestFile{},
		RootPath: "/test",
	}

	result := ConvertCoreToDomainInventory(coreInv)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Files) != 0 {
		t.Errorf("expected empty files, got %d files", len(result.Files))
	}
}

func TestConvertCoreTestStatus(t *testing.T) {
	tests := []struct {
		name       string
		coreStatus domain.TestStatus
		expected   analysis.TestStatus
	}{
		{
			name:       "active",
			coreStatus: domain.TestStatusActive,
			expected:   analysis.TestStatusActive,
		},
		{
			name:       "focused",
			coreStatus: domain.TestStatusFocused,
			expected:   analysis.TestStatusFocused,
		},
		{
			name:       "skipped",
			coreStatus: domain.TestStatusSkipped,
			expected:   analysis.TestStatusSkipped,
		},
		{
			name:       "todo",
			coreStatus: domain.TestStatusTodo,
			expected:   analysis.TestStatusTodo,
		},
		{
			name:       "xfail",
			coreStatus: domain.TestStatusXfail,
			expected:   analysis.TestStatusXfail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertCoreTestStatus(tt.coreStatus)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertCoreTestFile(t *testing.T) {
	coreFile := domain.TestFile{
		Path:      "test.ts",
		Framework: "jest",
		Language:  domain.LanguageTypeScript,
		Suites:    []domain.TestSuite{},
		Tests: []domain.Test{
			{
				Name: "test 1",
				Location: domain.Location{
					StartLine: 10,
					EndLine:   15,
				},
				Status: domain.TestStatusActive,
			},
		},
	}

	result := convertCoreTestFile(coreFile)

	if result.Path != "test.ts" {
		t.Errorf("expected path 'test.ts', got %s", result.Path)
	}
	if result.Framework != "jest" {
		t.Errorf("expected framework 'jest', got %s", result.Framework)
	}
	if len(result.Tests) != 1 {
		t.Errorf("expected 1 test, got %d", len(result.Tests))
	}
	if result.Tests[0].Name != "test 1" {
		t.Errorf("expected test name 'test 1', got %s", result.Tests[0].Name)
	}
}

func TestConvertCoreTestFile_WithDomainHints(t *testing.T) {
	coreFile := domain.TestFile{
		Path:      "auth.test.ts",
		Framework: "jest",
		Language:  domain.LanguageTypeScript,
		DomainHints: &domain.DomainHints{
			Calls:   []string{"authService.validateToken", "userRepo.findById"},
			Imports: []string{"@nestjs/jwt", "@nestjs/testing"},
		},
		Suites: []domain.TestSuite{},
		Tests:  []domain.Test{},
	}

	result := convertCoreTestFile(coreFile)

	if result.DomainHints == nil {
		t.Fatal("expected DomainHints to be non-nil")
	}
	if len(result.DomainHints.Calls) != 2 {
		t.Errorf("expected 2 calls, got %d", len(result.DomainHints.Calls))
	}
	if result.DomainHints.Calls[0] != "authService.validateToken" {
		t.Errorf("expected first call 'authService.validateToken', got %s", result.DomainHints.Calls[0])
	}
	if len(result.DomainHints.Imports) != 2 {
		t.Errorf("expected 2 imports, got %d", len(result.DomainHints.Imports))
	}
}

func TestConvertDomainHints(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		result := convertDomainHints(nil)
		if result != nil {
			t.Errorf("expected nil for nil input, got %v", result)
		}
	})

	t.Run("converts hints correctly", func(t *testing.T) {
		hints := &domain.DomainHints{
			Calls:   []string{"service.method"},
			Imports: []string{"package/module"},
		}

		result := convertDomainHints(hints)

		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if len(result.Calls) != 1 || result.Calls[0] != "service.method" {
			t.Errorf("unexpected Calls: %v", result.Calls)
		}
		if len(result.Imports) != 1 || result.Imports[0] != "package/module" {
			t.Errorf("unexpected Imports: %v", result.Imports)
		}
	})

	t.Run("empty slices", func(t *testing.T) {
		hints := &domain.DomainHints{
			Calls:   []string{},
			Imports: []string{},
		}

		result := convertDomainHints(hints)

		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if len(result.Calls) != 0 {
			t.Errorf("expected empty Calls, got %v", result.Calls)
		}
		if len(result.Imports) != 0 {
			t.Errorf("expected empty Imports, got %v", result.Imports)
		}
	})
}

func TestConvertCoreTestSuite(t *testing.T) {
	coreSuite := domain.TestSuite{
		Name: "suite 1",
		Location: domain.Location{
			StartLine: 5,
			EndLine:   20,
		},
		Suites: []domain.TestSuite{
			{
				Name: "nested suite",
				Location: domain.Location{
					StartLine: 10,
					EndLine:   15,
				},
			},
		},
		Tests: []domain.Test{
			{
				Name: "test in suite",
				Location: domain.Location{
					StartLine: 12,
					EndLine:   13,
				},
				Status: domain.TestStatusSkipped,
			},
		},
	}

	result := convertCoreTestSuite(coreSuite)

	if result.Name != "suite 1" {
		t.Errorf("expected name 'suite 1', got %s", result.Name)
	}
	if result.Location.StartLine != 5 {
		t.Errorf("expected start line 5, got %d", result.Location.StartLine)
	}
	if len(result.Suites) != 1 {
		t.Errorf("expected 1 nested suite, got %d", len(result.Suites))
	}
	if len(result.Tests) != 1 {
		t.Errorf("expected 1 test, got %d", len(result.Tests))
	}
	if result.Tests[0].Status != analysis.TestStatusSkipped {
		t.Errorf("expected status skipped, got %v", result.Tests[0].Status)
	}
}

func TestConvertCoreFileResult(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := ConvertCoreFileResult(nil)
		if result.Err != nil {
			t.Errorf("expected nil error, got %v", result.Err)
		}
		if result.File != nil {
			t.Errorf("expected nil file, got %v", result.File)
		}
	})

	t.Run("error result", func(t *testing.T) {
		parseErr := errors.New("parse failed")
		coreResult := &coreparser.FileResult{
			Err:  parseErr,
			Path: "broken.ts",
		}

		result := ConvertCoreFileResult(coreResult)

		if result.Err != parseErr {
			t.Errorf("expected error %v, got %v", parseErr, result.Err)
		}
		if result.File != nil {
			t.Errorf("expected nil file on error, got %v", result.File)
		}
	})

	t.Run("nil file (skipped)", func(t *testing.T) {
		coreResult := &coreparser.FileResult{
			Err:  nil,
			File: nil,
			Path: "unknown.xyz",
		}

		result := ConvertCoreFileResult(coreResult)

		if result.Err != nil {
			t.Errorf("expected nil error, got %v", result.Err)
		}
		if result.File != nil {
			t.Errorf("expected nil file, got %v", result.File)
		}
	})

	t.Run("success result", func(t *testing.T) {
		coreResult := &coreparser.FileResult{
			Err: nil,
			File: &domain.TestFile{
				Path:      "app.test.ts",
				Framework: "jest",
				Language:  domain.LanguageTypeScript,
				Tests: []domain.Test{
					{
						Name:   "should work",
						Status: domain.TestStatusActive,
					},
				},
			},
			Path:       "app.test.ts",
			Confidence: "scope",
		}

		result := ConvertCoreFileResult(coreResult)

		if result.Err != nil {
			t.Errorf("expected nil error, got %v", result.Err)
		}
		if result.File == nil {
			t.Fatal("expected non-nil file")
		}
		if result.File.Path != "app.test.ts" {
			t.Errorf("expected path 'app.test.ts', got %s", result.File.Path)
		}
		if result.File.Framework != "jest" {
			t.Errorf("expected framework 'jest', got %s", result.File.Framework)
		}
		if len(result.File.Tests) != 1 {
			t.Errorf("expected 1 test, got %d", len(result.File.Tests))
		}
	})
}
