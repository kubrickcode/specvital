package parser

import "github.com/specvital/core/pkg/domain"

// FileResult represents the result of parsing a single test file in streaming mode.
// It is sent through a channel for each file processed by ScanStream.
type FileResult struct {
	// File contains the parsed test file on success.
	// nil when Err is non-nil.
	File *domain.TestFile

	// Err contains the error encountered during parsing.
	// nil indicates successful parsing.
	Err error

	// Path is the original file path (relative to source root).
	// Always populated, allowing error tracking even on parse failure.
	Path string
}

// IsSuccess returns true if the file was parsed successfully.
func (r *FileResult) IsSuccess() bool {
	return r.Err == nil && r.File != nil
}
