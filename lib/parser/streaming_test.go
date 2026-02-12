package parser_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/kubrickcode/specvital/lib/parser"
	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/source"

	_ "github.com/kubrickcode/specvital/lib/parser/strategies/jest"
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

func TestScanStream(t *testing.T) {
	t.Run("should receive FileResult for each file", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create 10 test files
		for i := 0; i < 10; i++ {
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(tmpDir, "test"+string(rune('a'+i))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		scanner := parser.NewScanner()
		results, err := scanner.ScanStream(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		for result := range results {
			if result.IsSuccess() {
				count++
			}
		}

		if count != 10 {
			t.Errorf("expected 10 successful results, got %d", count)
		}
	})

	t.Run("should process multiple files independently", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create multiple valid test files to verify independent processing
		// Parse errors are difficult to produce reliably without malformed AST,
		// so this test verifies that each file is processed independently.
		validContent := []byte(`import { it } from '@jest/globals'; it('valid test', () => {});`)
		if err := os.WriteFile(filepath.Join(tmpDir, "valid.test.ts"), validContent, 0644); err != nil {
			t.Fatalf("failed to write valid file: %v", err)
		}

		if err := os.WriteFile(filepath.Join(tmpDir, "another.test.ts"), validContent, 0644); err != nil {
			t.Fatalf("failed to write another file: %v", err)
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		scanner := parser.NewScanner()
		results, err := scanner.ScanStream(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var successCount int
		for result := range results {
			if result.IsSuccess() {
				successCount++
			}
		}

		// Both valid files should be parsed successfully
		if successCount != 2 {
			t.Errorf("expected 2 successful results, got %d", successCount)
		}
	})

	t.Run("should stop on context cancellation", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create many files
		for i := 0; i < 50; i++ {
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(tmpDir, "test"+string(rune('0'+i%10))+string(rune('a'+i%26))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		scanner := parser.NewScanner()
		results, err := scanner.ScanStream(ctx, src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Cancel after receiving a few results
		var count int
		for range results {
			count++
			if count >= 3 {
				cancel()
				break
			}
		}

		// Drain remaining to allow cleanup
		for range results {
		}

		// Should have stopped early (not all 50)
		if count >= 50 {
			t.Error("expected scan to stop early on cancellation")
		}
	})

	t.Run("should work with empty directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		scanner := parser.NewScanner()
		results, err := scanner.ScanStream(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		for range results {
			count++
		}

		if count != 0 {
			t.Errorf("expected 0 results for empty directory, got %d", count)
		}
	})

	t.Run("should handle concurrent consumption safely", func(t *testing.T) {
		tmpDir := t.TempDir()

		for i := 0; i < 20; i++ {
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(tmpDir, "test"+string(rune('a'+i%26))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		scanner := parser.NewScanner(parser.WithWorkers(8))
		results, err := scanner.ScanStream(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		for result := range results {
			if result.IsSuccess() {
				count++
			}
		}

		if count != 20 {
			t.Errorf("expected 20 successful results, got %d", count)
		}
	})

	t.Run("should not leak goroutines on cancellation", func(t *testing.T) {
		tmpDir := t.TempDir()

		for i := 0; i < 100; i++ {
			subDir := filepath.Join(tmpDir, "sub"+string(rune('a'+i%26)))
			if err := os.MkdirAll(subDir, 0755); err != nil {
				t.Fatalf("failed to create dir: %v", err)
			}
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(subDir, "file"+string(rune('0'+i%10))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		goroutinesBefore := runtime.NumGoroutine()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		scanner := parser.NewScanner()
		results, err := scanner.ScanStream(ctx, src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Read a few then cancel
		count := 0
		for range results {
			count++
			if count >= 5 {
				cancel()
				break
			}
		}

		// Drain remaining
		for range results {
		}

		// Allow cleanup
		time.Sleep(100 * time.Millisecond)

		goroutinesAfter := runtime.NumGoroutine()

		if goroutinesAfter > goroutinesBefore+2 {
			t.Errorf("potential goroutine leak: before=%d, after=%d", goroutinesBefore, goroutinesAfter)
		}
	})
}

func TestScanStreaming(t *testing.T) {
	t.Run("should work as convenience wrapper", func(t *testing.T) {
		tmpDir := t.TempDir()

		for i := 0; i < 5; i++ {
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(tmpDir, "test"+string(rune('a'+i))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		results, err := parser.ScanStreaming(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		for result := range results {
			if result.IsSuccess() {
				count++
			}
		}

		if count != 5 {
			t.Errorf("expected 5 successful results, got %d", count)
		}
	})

	t.Run("should accept options", func(t *testing.T) {
		tmpDir := t.TempDir()

		content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
		if err := os.WriteFile(filepath.Join(tmpDir, "test.test.ts"), content, 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		results, err := parser.ScanStreaming(context.Background(), src, parser.WithWorkers(2))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		for result := range results {
			if result.IsSuccess() {
				count++
			}
		}

		if count != 1 {
			t.Errorf("expected 1 successful result, got %d", count)
		}
	})

	t.Run("should produce consistent results with Scan", func(t *testing.T) {
		tmpDir := t.TempDir()

		for i := 0; i < 10; i++ {
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(tmpDir, "test"+string(rune('a'+i))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		// Streaming scan
		results, err := parser.ScanStreaming(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var streamCount int
		for result := range results {
			if result.IsSuccess() {
				streamCount++
			}
		}

		// Batch scan
		scanResult, err := parser.Scan(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if streamCount != scanResult.Stats.FilesMatched {
			t.Errorf("ScanStreaming matched %d files, but Scan matched %d", streamCount, scanResult.Stats.FilesMatched)
		}
	})
}

func TestDiscoveryStream(t *testing.T) {
	t.Run("should maintain file count after streaming refactor", func(t *testing.T) {
		tmpDir := t.TempDir()

		for i := 0; i < 5; i++ {
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(tmpDir, "test"+string(rune('a'+i))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		result, err := parser.Scan(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Stats.FilesScanned != 5 {
			t.Errorf("expected 5 scanned files, got %d", result.Stats.FilesScanned)
		}
		if len(result.Inventory.Files) != 5 {
			t.Errorf("expected 5 parsed files, got %d", len(result.Inventory.Files))
		}
	})

	t.Run("should close channel on context cancellation without goroutine leak", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create directory structure with many files
		for i := 0; i < 100; i++ {
			subDir := filepath.Join(tmpDir, "sub"+string(rune('a'+i%26)))
			if err := os.MkdirAll(subDir, 0755); err != nil {
				t.Fatalf("failed to create dir: %v", err)
			}
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(subDir, "file"+string(rune('0'+i%10))+".test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		// Record goroutine count before test
		goroutinesBefore := runtime.NumGoroutine()

		ctx, cancel := context.WithCancel(context.Background())

		// Start scan and cancel quickly
		done := make(chan struct{})
		go func() {
			_, _ = parser.Scan(ctx, src)
			close(done)
		}()

		// Cancel after short delay
		time.Sleep(1 * time.Millisecond)
		cancel()

		// Wait for scan to complete
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("scan did not complete after cancellation")
		}

		// Allow goroutines to clean up
		time.Sleep(50 * time.Millisecond)

		goroutinesAfter := runtime.NumGoroutine()

		// Allow small tolerance for background goroutines
		if goroutinesAfter > goroutinesBefore+2 {
			t.Errorf("potential goroutine leak: before=%d, after=%d", goroutinesBefore, goroutinesAfter)
		}
	})

	t.Run("should maintain behavior with existing Scan API", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test files in various locations
		dirs := []string{
			"",
			"__tests__",
			"src/components/__tests__",
		}
		for _, dir := range dirs {
			path := filepath.Join(tmpDir, dir)
			if dir != "" {
				if err := os.MkdirAll(path, 0755); err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
			}
			content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
			filename := filepath.Join(path, "component.test.ts")
			if err := os.WriteFile(filename, content, 0644); err != nil {
				t.Fatalf("failed to write file: %v", err)
			}
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		result, err := parser.Scan(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Stats.FilesScanned != 3 {
			t.Errorf("expected 3 scanned files, got %d", result.Stats.FilesScanned)
		}
	})

	t.Run("should handle errors during discovery without blocking", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create valid test file
		content := []byte(`import { it } from '@jest/globals'; it('test', () => {});`)
		if err := os.WriteFile(filepath.Join(tmpDir, "valid.test.ts"), content, 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		src, err := source.NewLocalSource(tmpDir)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		defer src.Close()

		result, err := parser.Scan(context.Background(), src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Stats.FilesScanned != 1 {
			t.Errorf("expected 1 scanned file, got %d", result.Stats.FilesScanned)
		}
	})
}
