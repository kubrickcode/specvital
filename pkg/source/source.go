// Package source provides abstractions for reading files from various data sources.
// It defines a unified interface that supports local filesystem, Git repositories,
// and potentially other sources like S3 or GitLab in the future.
package source

import (
	"context"
	"errors"
	"io"
	"io/fs"
)

// Source defines the interface for reading files from a data source.
// Implementations must be safe for concurrent use.
type Source interface {
	// Root returns the root path of the source.
	// For LocalSource, this is the absolute path to the directory.
	// For GitSource, this is the path to the cloned repository.
	Root() string

	// Open opens the file at the given path for reading.
	// The path should be relative to the source root.
	// Callers must close the returned ReadCloser when done.
	Open(ctx context.Context, path string) (io.ReadCloser, error)

	// Stat returns file info for the given path.
	// The path should be relative to the source root.
	Stat(ctx context.Context, path string) (fs.FileInfo, error)

	// Close releases any resources held by the source.
	// For GitSource, this removes the cloned repository.
	// Close is idempotent; calling it multiple times has no effect.
	Close() error
}

// Sentinel errors for source operations.
var (
	// ErrInvalidPath indicates the provided path is invalid or inaccessible.
	ErrInvalidPath = errors.New("source: invalid path")

	// ErrGitCloneFailed indicates a git clone operation failed.
	ErrGitCloneFailed = errors.New("source: git clone failed")

	// ErrRepositoryNotFound indicates the repository does not exist or is inaccessible.
	ErrRepositoryNotFound = errors.New("source: repository not found")
)
