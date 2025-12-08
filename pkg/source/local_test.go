package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestNewLocalSource(t *testing.T) {
	t.Run("should create source with valid directory path", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()

		// When
		src, err := NewLocalSource(tmpDir)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if src == nil {
			t.Fatal("expected source to be non-nil")
		}
		if src.Root() != tmpDir {
			t.Errorf("expected root %q, got %q", tmpDir, src.Root())
		}
	})

	t.Run("should convert relative path to absolute path", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to change directory: %v", err)
		}
		defer func() { _ = os.Chdir(originalDir) }()

		// When
		src, err := NewLocalSource(".")

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !filepath.IsAbs(src.Root()) {
			t.Errorf("expected absolute path, got %q", src.Root())
		}
	})

	t.Run("should fail with non-existent path", func(t *testing.T) {
		// Given
		nonExistentPath := "/path/that/does/not/exist/12345"

		// When
		src, err := NewLocalSource(nonExistentPath)

		// Then
		if err == nil {
			t.Fatal("expected error for non-existent path")
		}
		if src != nil {
			t.Error("expected source to be nil")
		}
		if !isInvalidPathError(err) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})

	t.Run("should fail when path is a file not a directory", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "file.txt")
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// When
		src, err := NewLocalSource(filePath)

		// Then
		if err == nil {
			t.Fatal("expected error for file path")
		}
		if src != nil {
			t.Error("expected source to be nil")
		}
		if !isInvalidPathError(err) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})
}

func TestLocalSource_Open(t *testing.T) {
	t.Run("should open existing file", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		testContent := []byte("hello world")
		testFile := "test.txt"
		if err := os.WriteFile(filepath.Join(tmpDir, testFile), testContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, testFile)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read content: %v", err)
		}
		if string(content) != string(testContent) {
			t.Errorf("expected %q, got %q", testContent, content)
		}
	})

	t.Run("should open file in nested directory", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		nestedDir := filepath.Join(tmpDir, "nested", "dir")
		if err := os.MkdirAll(nestedDir, 0755); err != nil {
			t.Fatalf("failed to create nested dir: %v", err)
		}
		testContent := []byte("nested content")
		testFile := filepath.Join("nested", "dir", "file.txt")
		if err := os.WriteFile(filepath.Join(tmpDir, testFile), testContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, testFile)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read content: %v", err)
		}
		if string(content) != string(testContent) {
			t.Errorf("expected %q, got %q", testContent, content)
		}
	})

	t.Run("should fail for non-existent file", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, "non-existent.txt")

		// Then
		if err == nil {
			t.Fatal("expected error for non-existent file")
		}
		if reader != nil {
			t.Error("expected reader to be nil")
		}
	})

	t.Run("should prevent path traversal attacks", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		traversalPaths := []string{
			"../etc/passwd",
			"../../etc/passwd",
			"subdir/../../etc/passwd",
		}

		for _, path := range traversalPaths {
			// When
			reader, err := src.Open(ctx, path)

			// Then
			if err == nil {
				reader.Close()
				t.Errorf("expected error for path traversal: %s", path)
			}
			if !isInvalidPathError(err) {
				t.Errorf("expected ErrInvalidPath for %s, got: %v", path, err)
			}
		}
	})
}

func TestLocalSource_Stat(t *testing.T) {
	t.Run("should return file info for existing file", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		testContent := []byte("hello world")
		testFile := "test.txt"
		if err := os.WriteFile(filepath.Join(tmpDir, testFile), testContent, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		info, err := src.Stat(ctx, testFile)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Name() != testFile {
			t.Errorf("expected name %q, got %q", testFile, info.Name())
		}
		if info.Size() != int64(len(testContent)) {
			t.Errorf("expected size %d, got %d", len(testContent), info.Size())
		}
		if info.IsDir() {
			t.Error("expected IsDir to be false")
		}
	})

	t.Run("should return info for directory", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		subDir := "subdir"
		if err := os.Mkdir(filepath.Join(tmpDir, subDir), 0755); err != nil {
			t.Fatalf("failed to create subdir: %v", err)
		}
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		info, err := src.Stat(ctx, subDir)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected IsDir to be true")
		}
	})

	t.Run("should fail for non-existent path", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		info, err := src.Stat(ctx, "non-existent")

		// Then
		if err == nil {
			t.Fatal("expected error for non-existent path")
		}
		if info != nil {
			t.Error("expected info to be nil")
		}
	})
}

func TestLocalSource_Close(t *testing.T) {
	t.Run("should be idempotent", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		src, _ := NewLocalSource(tmpDir)

		// When
		err1 := src.Close()
		err2 := src.Close()
		err3 := src.Close()

		// Then
		if err1 != nil {
			t.Errorf("first close returned error: %v", err1)
		}
		if err2 != nil {
			t.Errorf("second close returned error: %v", err2)
		}
		if err3 != nil {
			t.Errorf("third close returned error: %v", err3)
		}
	})
}

func TestLocalSource_ReadTestdata(t *testing.T) {
	t.Run("should read existing testdata files", func(t *testing.T) {
		// Given
		wd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}

		repoRoot := filepath.Dir(filepath.Dir(wd))
		testdataPath := filepath.Join(repoRoot, "tests", "integration", "testdata")

		if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
			t.Skip("testdata directory not found")
		}

		src, err := NewLocalSource(testdataPath)
		if err != nil {
			t.Fatalf("failed to create source: %v", err)
		}
		ctx := context.Background()

		// When
		info, err := src.Stat(ctx, "cache")

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !info.IsDir() {
			t.Error("expected cache to be a directory")
		}
	})
}

func TestLocalSource_Open_SymlinkAttack(t *testing.T) {
	t.Run("should prevent symlink-based path traversal", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		outsideDir := t.TempDir()
		secretFile := filepath.Join(outsideDir, "secret.txt")
		if err := os.WriteFile(secretFile, []byte("secret"), 0644); err != nil {
			t.Fatalf("failed to create secret file: %v", err)
		}

		symlinkPath := filepath.Join(tmpDir, "evil_link")
		if err := os.Symlink(secretFile, symlinkPath); err != nil {
			t.Skipf("symlink not supported: %v", err)
		}

		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, "evil_link")

		// Then
		if err == nil {
			reader.Close()
			t.Fatal("expected error for symlink escape")
		}
		if !isInvalidPathError(err) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})

	t.Run("should prevent symlink directory escape", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		outsideDir := t.TempDir()
		secretFile := filepath.Join(outsideDir, "secret.txt")
		if err := os.WriteFile(secretFile, []byte("secret"), 0644); err != nil {
			t.Fatalf("failed to create secret file: %v", err)
		}

		symlinkDir := filepath.Join(tmpDir, "evil_dir")
		if err := os.Symlink(outsideDir, symlinkDir); err != nil {
			t.Skipf("symlink not supported: %v", err)
		}

		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, filepath.Join("evil_dir", "secret.txt"))

		// Then
		if err == nil {
			reader.Close()
			t.Fatal("expected error for symlink directory escape")
		}
		if !isInvalidPathError(err) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})

	t.Run("should allow symlink within root", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		realFile := filepath.Join(tmpDir, "real.txt")
		if err := os.WriteFile(realFile, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create real file: %v", err)
		}

		symlinkPath := filepath.Join(tmpDir, "link.txt")
		if err := os.Symlink(realFile, symlinkPath); err != nil {
			t.Skipf("symlink not supported: %v", err)
		}

		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, "link.txt")

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer reader.Close()

		content, _ := io.ReadAll(reader)
		if string(content) != "content" {
			t.Errorf("expected 'content', got %q", content)
		}
	})
}

func TestLocalSource_ConcurrentAccess(t *testing.T) {
	t.Run("should be safe for concurrent use", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		testContent := []byte("concurrent")
		for i := 0; i < 10; i++ {
			filename := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
			if err := os.WriteFile(filename, testContent, 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}
		}
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		var wg sync.WaitGroup
		errChan := make(chan error, 30)

		for i := 0; i < 10; i++ {
			for j := 0; j < 3; j++ {
				wg.Add(1)
				go func(fileNum int) {
					defer wg.Done()
					reader, err := src.Open(ctx, fmt.Sprintf("file%d.txt", fileNum))
					if err != nil {
						errChan <- err
						return
					}
					defer reader.Close()
					_, err = io.ReadAll(reader)
					if err != nil {
						errChan <- err
					}
				}(i)
			}
		}

		wg.Wait()
		close(errChan)

		// Then
		for err := range errChan {
			t.Errorf("concurrent access error: %v", err)
		}
	})
}

func TestLocalSource_Open_EmptyPath(t *testing.T) {
	t.Run("should reject empty path", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, "")

		// Then
		if err == nil {
			reader.Close()
			t.Fatal("expected error for empty path")
		}
		if !isInvalidPathError(err) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})

	t.Run("should reject dot path", func(t *testing.T) {
		// Given
		tmpDir := t.TempDir()
		src, _ := NewLocalSource(tmpDir)
		ctx := context.Background()

		// When
		reader, err := src.Open(ctx, ".")

		// Then
		if err == nil {
			reader.Close()
			t.Fatal("expected error for dot path")
		}
		if !isInvalidPathError(err) {
			t.Errorf("expected ErrInvalidPath, got: %v", err)
		}
	})
}

func isInvalidPathError(err error) bool {
	return errors.Is(err, ErrInvalidPath)
}
