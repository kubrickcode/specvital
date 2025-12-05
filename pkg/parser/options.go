package parser

import (
	"time"

	"github.com/specvital/core/pkg/parser/detection"
)

// ScanOptions configures the behavior of [Scan].
type ScanOptions struct {
	// ExcludePatterns specifies directory names to skip during scanning.
	ExcludePatterns []string
	// MaxFileSize is the maximum file size in bytes to process.
	MaxFileSize int64
	// Patterns specifies glob patterns to match test files (e.g., "**/*.test.ts").
	Patterns []string
	// ProjectContext provides project-level metadata for source-agnostic detection.
	// When set, enables detection without filesystem access (e.g., GitHub API environment).
	ProjectContext *detection.ProjectContext
	// Timeout is the maximum duration for the entire scan operation.
	Timeout time.Duration
	// Workers is the number of concurrent file parsers.
	Workers int
}

// ScanOption is a functional option for configuring [Scan].
type ScanOption func(*ScanOptions)

// WithExclude returns a [ScanOption] that adds directory patterns to skip.
// These patterns are matched against directory base names.
func WithExclude(patterns []string) ScanOption {
	return func(o *ScanOptions) {
		o.ExcludePatterns = patterns
	}
}

// WithScanMaxFileSize returns a [ScanOption] that sets the maximum file size.
// Files larger than this size are skipped. Negative values are ignored.
func WithScanMaxFileSize(size int64) ScanOption {
	return func(o *ScanOptions) {
		if size < 0 {
			return
		}
		o.MaxFileSize = size
	}
}

// WithScanPatterns returns a [ScanOption] that filters files by glob patterns.
// Only files matching at least one pattern are processed.
// Uses doublestar syntax (e.g., "**/*.test.ts", "src/**/*.spec.js").
func WithScanPatterns(patterns []string) ScanOption {
	return func(o *ScanOptions) {
		o.Patterns = patterns
	}
}

// WithTimeout returns a [ScanOption] that sets the scan timeout.
// Negative values are ignored.
func WithTimeout(d time.Duration) ScanOption {
	return func(o *ScanOptions) {
		if d < 0 {
			return
		}
		o.Timeout = d
	}
}

// WithWorkers returns a [ScanOption] that sets the number of parallel workers.
// Zero uses GOMAXPROCS, negative values are ignored.
func WithWorkers(n int) ScanOption {
	return func(o *ScanOptions) {
		if n < 0 {
			return
		}
		o.Workers = n
	}
}

// WithProjectContext returns a [ScanOption] that sets the project context.
// This enables source-agnostic detection for environments without filesystem access
// (e.g., GitHub API). The ProjectContext should contain config file paths and
// their parsed contents.
func WithProjectContext(ctx *detection.ProjectContext) ScanOption {
	return func(o *ScanOptions) {
		o.ProjectContext = ctx
	}
}
