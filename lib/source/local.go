package source

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// LocalSource implements Source for local filesystem access.
type LocalSource struct {
	root string
}

// NewLocalSource creates a new LocalSource for the given root path.
// The path must be an existing directory. Relative paths are converted
// to absolute paths.
func NewLocalSource(rootPath string) (*LocalSource, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to resolve path: %v", ErrInvalidPath, err)
	}
	absPath = filepath.Clean(absPath)

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: path does not exist: %s", ErrInvalidPath, absPath)
		}
		return nil, fmt.Errorf("%w: failed to stat path: %v", ErrInvalidPath, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%w: path is not a directory: %s", ErrInvalidPath, absPath)
	}

	return &LocalSource{root: absPath}, nil
}

// Root returns the absolute path to the source root directory.
func (s *LocalSource) Root() string {
	return s.root
}

// Open opens the file at the given path for reading.
// The path must be relative to the source root. Paths attempting to escape
// the root directory will return ErrInvalidPath.
func (s *LocalSource) Open(_ context.Context, path string) (io.ReadCloser, error) {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Stat returns file info for the given path.
// The path must be relative to the source root. Paths attempting to escape
// the root directory will return ErrInvalidPath.
func (s *LocalSource) Stat(_ context.Context, path string) (fs.FileInfo, error) {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	return info, nil
}

// resolvePath validates and resolves a relative path within the source root.
// Returns ErrInvalidPath if the path escapes the root directory.
func (s *LocalSource) resolvePath(path string) (string, error) {
	cleaned := filepath.Clean(path)
	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("%w: empty or current directory path not allowed", ErrInvalidPath)
	}

	fullPath := filepath.Join(s.root, cleaned)

	rootWithSep := s.root + string(filepath.Separator)
	if !strings.HasPrefix(fullPath, rootWithSep) && fullPath != s.root {
		return "", fmt.Errorf("%w: path escapes root directory", ErrInvalidPath)
	}

	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err == nil {
		resolvedRoot, rootErr := filepath.EvalSymlinks(s.root)
		if rootErr == nil {
			resolvedRootWithSep := resolvedRoot + string(filepath.Separator)
			if !strings.HasPrefix(resolvedPath, resolvedRootWithSep) && resolvedPath != resolvedRoot {
				return "", fmt.Errorf("%w: path escapes root directory via symlink", ErrInvalidPath)
			}
		}
	}

	return fullPath, nil
}

// Close is a no-op for LocalSource as there are no resources to release.
func (s *LocalSource) Close() error {
	return nil
}
