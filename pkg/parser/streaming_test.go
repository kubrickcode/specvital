package parser_test

import (
	"errors"
	"testing"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser"
)

func TestFileResult(t *testing.T) {
	t.Run("IsSuccess returns true for successful parse", func(t *testing.T) {
		result := &parser.FileResult{
			File: &domain.TestFile{
				Path:      "example.test.ts",
				Framework: "jest",
				Language:  "typescript",
			},
			Err:  nil,
			Path: "example.test.ts",
		}

		if !result.IsSuccess() {
			t.Error("expected IsSuccess() to return true for successful parse")
		}
	})

	t.Run("IsSuccess returns false when Err is set", func(t *testing.T) {
		result := &parser.FileResult{
			File: nil,
			Err:  errors.New("parse error"),
			Path: "broken.test.ts",
		}

		if result.IsSuccess() {
			t.Error("expected IsSuccess() to return false when Err is set")
		}
	})

	t.Run("IsSuccess returns false when File is nil", func(t *testing.T) {
		result := &parser.FileResult{
			File: nil,
			Err:  nil,
			Path: "unknown.ts",
		}

		if result.IsSuccess() {
			t.Error("expected IsSuccess() to return false when File is nil")
		}
	})

	t.Run("Path is always accessible for error tracking", func(t *testing.T) {
		path := "src/components/Button.test.tsx"
		result := &parser.FileResult{
			File: nil,
			Err:  errors.New("syntax error at line 42"),
			Path: path,
		}

		if result.Path != path {
			t.Errorf("expected Path %q, got %q", path, result.Path)
		}
	})
}

func TestFileResult_ChannelUsage(t *testing.T) {
	t.Run("FileResult can be sent through channel", func(t *testing.T) {
		ch := make(chan *parser.FileResult, 3)

		// Simulate sending results through channel
		ch <- &parser.FileResult{
			File: &domain.TestFile{Path: "test1.ts", Framework: "jest"},
			Path: "test1.ts",
		}
		ch <- &parser.FileResult{
			Err:  errors.New("parse failed"),
			Path: "test2.ts",
		}
		ch <- &parser.FileResult{
			File: &domain.TestFile{Path: "test3.ts", Framework: "vitest"},
			Path: "test3.ts",
		}
		close(ch)

		var successCount, errorCount int
		for result := range ch {
			if result.IsSuccess() {
				successCount++
			} else {
				errorCount++
			}
		}

		if successCount != 2 {
			t.Errorf("expected 2 successful results, got %d", successCount)
		}
		if errorCount != 1 {
			t.Errorf("expected 1 error result, got %d", errorCount)
		}
	})
}
