package parser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/detection"
	domain_hints "github.com/kubrickcode/specvital/packages/core/pkg/parser/domain_hints"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/strategies/shared/dotnetast"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/strategies/shared/kotlinast"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/strategies/shared/swiftast"
	"github.com/kubrickcode/specvital/packages/core/pkg/source"
	"golang.org/x/sync/semaphore"
)

const (
	// DefaultWorkers indicates that the scanner should use GOMAXPROCS as the worker count.
	DefaultWorkers = 0
	// DefaultTimeout is the default scan timeout duration.
	DefaultTimeout = 5 * time.Minute
	// MaxWorkers is the maximum number of concurrent workers allowed.
	MaxWorkers = 1024
	// DefaultMaxFileSize is the default maximum file size for scanning (10MB).
	DefaultMaxFileSize = 10 * 1024 * 1024
)

// DefaultSkipPatterns contains directory names that are skipped by default during scanning.
var DefaultSkipPatterns = []string{
	"node_modules",
	".git",
	"vendor",
	"dist",
	".next",
	"__pycache__",
	"coverage",
	".cache",
}

var (
	// ErrScanCancelled is returned when scanning is cancelled via context.
	ErrScanCancelled = errors.New("scanner: scan cancelled")
	// ErrScanTimeout is returned when scanning exceeds the timeout duration.
	ErrScanTimeout = errors.New("scanner: scan timeout")
)

// Scanner performs framework detection and test file parsing.
// It integrates framework.Registry and detection.Detector for improved accuracy.
type Scanner struct {
	registry     *framework.Registry
	detector     *detection.Detector
	projectScope *framework.AggregatedProjectScope
	options      *ScanOptions
}

// ScanResult contains the outcome of a scan operation.
type ScanResult struct {
	// Inventory contains all parsed test files.
	Inventory *domain.Inventory

	// Errors contains non-fatal errors encountered during scanning.
	Errors []ScanError

	// Stats provides scan statistics including confidence distribution.
	Stats ScanStats
}

// ScanError represents an error that occurred during a specific phase of scanning.
type ScanError struct {
	// Err is the underlying error.
	Err error

	// Path is the file path where the error occurred (may be empty for non-file errors).
	Path string

	// Phase indicates which phase the error occurred in.
	// Values: "discovery", "config-parse", "detection", "parsing"
	Phase string
}

// Error implements the error interface.
func (e ScanError) Error() string {
	if e.Path == "" {
		return fmt.Sprintf("[%s] %v", e.Phase, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %v", e.Phase, e.Path, e.Err)
}

// ScanStats provides statistics about the scan operation.
type ScanStats struct {
	// FilesScanned is the total number of test file candidates discovered.
	FilesScanned int

	// FilesMatched is the number of files that were successfully parsed.
	FilesMatched int

	// FilesFailed is the number of files that failed to parse.
	FilesFailed int

	// FilesSkipped is the number of files skipped due to low confidence or other reasons.
	FilesSkipped int

	// ConfidenceDist tracks detection confidence distribution.
	// Keys: "definite", "moderate", "weak", "unknown"
	ConfidenceDist map[string]int

	// ConfigsFound is the number of config files discovered and parsed.
	ConfigsFound int

	// Duration is the total scan duration.
	Duration time.Duration
}

// NewScanner creates a new scanner with the given options.
func NewScanner(opts ...ScanOption) *Scanner {
	options := newDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	applyDefaults(&options)

	detector := detection.NewDetector(options.Registry)

	return &Scanner{
		registry:     options.Registry,
		detector:     detector,
		projectScope: nil,
		options:      &options,
	}
}

// SetProjectScope sets pre-parsed config information for remote sources.
// This is useful when scanning from sources without filesystem access (e.g., GitHub API).
func (s *Scanner) SetProjectScope(scope *framework.AggregatedProjectScope) {
	s.projectScope = scope
	s.detector.SetProjectScope(scope)
}

// Scan performs the complete scanning process:
//  1. Discover and parse config files
//  2. Build project scope
//  3. Discover test files
//  4. Detect framework for each file
//  5. Parse test files in parallel
//
// Internally uses streaming scan for unified implementation.
// The caller is responsible for calling src.Close() when done.
// For GitSource, failure to close will leak temporary directories.
func (s *Scanner) Scan(ctx context.Context, src source.Source) (*ScanResult, error) {
	startTime := time.Now()

	rootPath := src.Root()

	result := &ScanResult{
		Inventory: &domain.Inventory{
			RootPath: rootPath,
			Files:    []domain.TestFile{},
		},
		Errors: []ScanError{},
		Stats: ScanStats{
			ConfidenceDist: make(map[string]int),
		},
	}

	// ScanStream handles timeout and config parsing internally
	resultCh, err := s.ScanStream(ctx, src)
	if err != nil {
		result.Stats.Duration = time.Since(startTime)
		if errors.Is(err, context.DeadlineExceeded) {
			return result, ErrScanTimeout
		}
		if errors.Is(err, context.Canceled) {
			return result, ErrScanCancelled
		}
		return result, err
	}

	// Retrieve config stats after ScanStream initialization
	if s.projectScope != nil {
		result.Stats.ConfigsFound = len(s.projectScope.Configs)
	}

	// Collect all streaming results
	var files []domain.TestFile
	for fileResult := range resultCh {
		result.Stats.FilesScanned++

		if fileResult.Confidence != "" {
			result.Stats.ConfidenceDist[fileResult.Confidence]++
		}

		if fileResult.Err != nil {
			phase := "parsing"
			if fileResult.Path == "" {
				phase = "discovery"
			}
			result.Errors = append(result.Errors, ScanError{
				Err:   fileResult.Err,
				Path:  fileResult.Path,
				Phase: phase,
			})
			continue
		}

		if fileResult.File != nil {
			files = append(files, *fileResult.File)
		}
	}

	// Sort by path for deterministic output order.
	// Parallel goroutines complete in variable order based on file size and parsing complexity.
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	result.Inventory.Files = files
	result.Stats.FilesMatched = len(files)
	result.Stats.FilesFailed = len(result.Errors)
	result.Stats.FilesSkipped = result.Stats.FilesScanned - result.Stats.FilesMatched - result.Stats.FilesFailed
	result.Stats.Duration = time.Since(startTime)

	// Check for timeout or cancellation after processing
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return result, ErrScanTimeout
		}
		if errors.Is(err, context.Canceled) {
			return result, ErrScanCancelled
		}
	}

	return result, nil
}

// ScanFiles scans specific files (for incremental/watch mode).
// This bypasses file discovery and directly scans the provided file paths.
// Internally uses streaming scan for unified implementation.
//
// The caller is responsible for calling src.Close() when done.
func (s *Scanner) ScanFiles(ctx context.Context, src source.Source, files []string) (*ScanResult, error) {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(ctx, s.options.Timeout)
	defer cancel()

	result := &ScanResult{
		Inventory: &domain.Inventory{
			RootPath: src.Root(),
			Files:    []domain.TestFile{},
		},
		Errors: []ScanError{},
		Stats: ScanStats{
			FilesScanned:   len(files),
			ConfidenceDist: make(map[string]int),
		},
	}

	if len(files) == 0 {
		result.Stats.Duration = time.Since(startTime)
		return result, nil
	}

	// Collect all streaming results
	var parsedFiles []domain.TestFile
	for fileResult := range s.scanFilesStream(ctx, src, files) {
		if fileResult.Confidence != "" {
			result.Stats.ConfidenceDist[fileResult.Confidence]++
		}

		if fileResult.Err != nil {
			result.Errors = append(result.Errors, ScanError{
				Err:   fileResult.Err,
				Path:  fileResult.Path,
				Phase: "parsing",
			})
			continue
		}

		if fileResult.File != nil {
			parsedFiles = append(parsedFiles, *fileResult.File)
		}
	}

	// Sort by path for deterministic output order
	sort.Slice(parsedFiles, func(i, j int) bool {
		return parsedFiles[i].Path < parsedFiles[j].Path
	})

	result.Inventory.Files = parsedFiles
	result.Stats.FilesMatched = len(parsedFiles)
	result.Stats.FilesFailed = len(result.Errors)
	result.Stats.FilesSkipped = result.Stats.FilesScanned - result.Stats.FilesMatched - result.Stats.FilesFailed
	result.Stats.Duration = time.Since(startTime)

	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return result, ErrScanTimeout
		}
		if errors.Is(err, context.Canceled) {
			return result, ErrScanCancelled
		}
	}

	return result, nil
}

// discoverConfigFiles walks the source root to find framework config files.
// Returns relative paths from the source root for consistent Source.Open() usage.
func (s *Scanner) discoverConfigFiles(ctx context.Context, src source.Source) []string {
	patterns := []string{
		"jest.config.js",
		"jest.config.ts",
		"jest.config.mjs",
		"jest.config.cjs",
		"jest.config.json",
		"vitest.config.js",
		"vitest.config.ts",
		"vitest.config.mjs",
		"vitest.config.cjs",
		"playwright.config.js",
		"playwright.config.ts",
		"cypress.config.cjs",
		"cypress.config.js",
		"cypress.config.mjs",
		"cypress.config.mts",
		"cypress.config.ts",
		"pytest.ini",
		"pyproject.toml",
		"conftest.py",
		".rspec",
		"spec_helper.rb",
		"rails_helper.rb",
		"phpunit.xml",
		"phpunit.xml.dist",
		"phpunit.dist.xml",
		".mocharc.cjs",
		".mocharc.js",
		".mocharc.json",
		".mocharc.jsonc",
		".mocharc.mjs",
		".mocharc.yaml",
		".mocharc.yml",
		"mocha.opts",
	}

	rootPath := src.Root()
	skipSet := buildSkipSet(s.options.ExcludePatterns)
	var configFiles []string

	_ = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, walkErr error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if walkErr != nil {
			return nil
		}

		if d.IsDir() {
			if shouldSkipDir(path, rootPath, skipSet) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}

		filename := filepath.Base(path)
		for _, pattern := range patterns {
			if filename == pattern {
				relPath, err := filepath.Rel(rootPath, path)
				if err == nil {
					configFiles = append(configFiles, relPath)
				}
				break
			}
		}

		return nil
	})

	return configFiles
}

func (s *Scanner) parseConfigFiles(ctx context.Context, src source.Source, files []string, errors *[]ScanError) *framework.AggregatedProjectScope {
	scope := framework.NewProjectScope()

	for _, file := range files {
		if ctx.Err() != nil {
			break
		}

		content, err := readFileFromSource(ctx, src, file)
		if err != nil {
			*errors = append(*errors, ScanError{
				Err:   err,
				Path:  file,
				Phase: "config-parse",
			})
			continue
		}

		filename := filepath.Base(file)
		parsed := false

		for _, def := range s.registry.All() {
			if def.ConfigParser == nil {
				continue
			}

			signal := framework.Signal{
				Type:  framework.SignalConfigFile,
				Value: filename,
			}

			matched := false
			for _, matcher := range def.Matchers {
				result := matcher.Match(ctx, signal)
				if result.Confidence > 0 {
					matched = true
					break
				}
			}

			if matched {
				// Use absolute path for config parsing to ensure correct BaseDir resolution
				absConfigPath := filepath.Join(src.Root(), file)
				configScope, err := def.ConfigParser.Parse(ctx, absConfigPath, content)
				if err != nil {
					*errors = append(*errors, ScanError{
						Err:   err,
						Path:  file,
						Phase: "config-parse",
					})
				} else {
					scope.AddConfig(absConfigPath, configScope)
					parsed = true
				}
				break
			}
		}

		if !parsed {
			*errors = append(*errors, ScanError{
				Err:   fmt.Errorf("no matching framework config parser"),
				Path:  file,
				Phase: "config-parse",
			})
		}
	}

	return scope
}

// DiscoveryResult contains a discovered file path or an error encountered during discovery.
type DiscoveryResult struct {
	// Path is the relative file path from source root (empty when Err is set).
	Path string

	// Err is the error encountered during discovery (nil on success).
	Err error
}

// discoverTestFilesStream walks the source root to find test file candidates,
// sending each discovered path through a channel for incremental processing.
// The channel is closed when discovery completes or context is cancelled.
// Callers must consume all values from the channel to avoid goroutine leaks.
func (s *Scanner) discoverTestFilesStream(ctx context.Context, src source.Source) <-chan DiscoveryResult {
	out := make(chan DiscoveryResult)

	go func() {
		defer close(out)

		rootPath := src.Root()
		skipSet := buildSkipSet(append(DefaultSkipPatterns, s.options.ExcludePatterns...))

		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, walkErr error) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if walkErr != nil {
				select {
				case out <- DiscoveryResult{Err: fmt.Errorf("access error at %s: %w", path, walkErr)}:
				case <-ctx.Done():
					return ctx.Err()
				}
				return nil
			}

			if d.IsDir() {
				if shouldSkipDir(path, rootPath, skipSet) {
					return filepath.SkipDir
				}
				return nil
			}

			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}

			relPath, err := filepath.Rel(rootPath, path)
			if err != nil {
				select {
				case out <- DiscoveryResult{Err: fmt.Errorf("compute relative path for %s: %w", path, err)}:
				case <-ctx.Done():
					return ctx.Err()
				}
				return nil
			}

			if !isTestFileCandidate(relPath) {
				return nil
			}

			if len(s.options.Patterns) > 0 {
				if !matchesAnyPattern(path, rootPath, s.options.Patterns) {
					return nil
				}
			}

			if s.options.MaxFileSize > 0 {
				info, err := d.Info()
				if err != nil {
					select {
					case out <- DiscoveryResult{Err: fmt.Errorf("failed to get file info for %s: %w", path, err)}:
					case <-ctx.Done():
						return ctx.Err()
					}
					return nil
				}
				if info.Size() > s.options.MaxFileSize {
					return nil
				}
			}

			select {
			case out <- DiscoveryResult{Path: relPath}:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})

		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			select {
			case out <- DiscoveryResult{Err: fmt.Errorf("walk directory: %w", err)}:
			case <-ctx.Done():
			}
		}
	}()

	return out
}

// scanFilesStream parses a list of files and streams results through a channel.
// This is the internal implementation used by ScanFiles.
// The caller is responsible for consuming all results from the channel.
func (s *Scanner) scanFilesStream(ctx context.Context, src source.Source, files []string) <-chan *FileResult {
	out := make(chan *FileResult)

	go func() {
		defer close(out)

		workers := s.options.Workers
		if workers <= 0 {
			workers = runtime.GOMAXPROCS(0)
		}
		if workers > MaxWorkers {
			workers = MaxWorkers
		}

		sem := semaphore.NewWeighted(int64(workers))
		var wg sync.WaitGroup

		for _, file := range files {
			if ctx.Err() != nil {
				break
			}

			if err := sem.Acquire(ctx, 1); err != nil {
				break
			}

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				defer sem.Release(1)

				result := s.parseFileToResult(ctx, src, filePath)

				select {
				case out <- result:
				case <-ctx.Done():
				}
			}(file)
		}

		wg.Wait()
	}()

	return out
}

func (s *Scanner) parseFile(ctx context.Context, src source.Source, path string) (*domain.TestFile, *ScanError, string) {
	if err := ctx.Err(); err != nil {
		return nil, &ScanError{
			Err:   err,
			Path:  path,
			Phase: "parsing",
		}, ""
	}

	content, err := readFileFromSource(ctx, src, path)
	if err != nil {
		return nil, &ScanError{
			Err:   err,
			Path:  path,
			Phase: "parsing",
		}, ""
	}

	// Use absolute path for detection to match config scope paths
	absPath := filepath.Join(src.Root(), path)
	detectionResult := s.detector.Detect(ctx, absPath, content)

	if !detectionResult.IsDetected() {
		return nil, nil, "unknown"
	}

	def := s.registry.Find(detectionResult.Framework)
	if def == nil || def.Parser == nil {
		return nil, &ScanError{
			Err:   fmt.Errorf("no parser for framework %s", detectionResult.Framework),
			Path:  path,
			Phase: "detection",
		}, string(detectionResult.Source)
	}

	testFile, err := def.Parser.Parse(ctx, content, path)
	if err != nil {
		return nil, &ScanError{
			Err:   fmt.Errorf("parse: %w", err),
			Path:  path,
			Phase: "parsing",
		}, string(detectionResult.Source)
	}

	if s.options.ExtractDomainHints {
		if extractor := domain_hints.GetExtractor(testFile.Language); extractor != nil {
			testFile.DomainHints = extractor.Extract(ctx, content)
		}
	}

	return testFile, nil, string(detectionResult.Source)
}

// readFileFromSource reads a file from source using relative path.
// The relPath must be relative to src.Root().
func readFileFromSource(ctx context.Context, src source.Source, relPath string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	reader, err := src.Open(ctx, relPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", relPath, err)
	}

	return content, nil
}

func buildSkipSet(patterns []string) map[string]bool {
	skipSet := make(map[string]bool, len(patterns))
	for _, p := range patterns {
		skipSet[p] = true
	}
	return skipSet
}

func shouldSkipDir(path, rootPath string, skipSet map[string]bool) bool {
	if path == rootPath {
		return false
	}

	base := filepath.Base(path)

	if base == "coverage" {
		parent := filepath.Dir(path)
		return parent == rootPath
	}

	return skipSet[base]
}

func isTestFileCandidate(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs":
		return isJSTestFile(path)
	case ".go":
		return isGoTestFile(path)
	case ".java":
		return isJavaTestFile(path)
	case ".kt", ".kts":
		return isKotlinTestFile(path)
	case ".py":
		return isPythonTestFile(path)
	case ".cs":
		return isCSharpTestFile(path)
	case ".rb":
		return isRubyTestFile(path)
	case ".rs":
		return isRustTestFile(path)
	case ".cc", ".cpp", ".cxx":
		return isCppTestFile(path)
	case ".php":
		return isPHPTestFile(path)
	case ".swift":
		return isSwiftTestFile(path)
	default:
		return false
	}
}

func isGoTestFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasSuffix(base, "_test.go")
}

func isJavaTestFile(path string) bool {
	normalizedPath := filepath.ToSlash(path)

	// Exclude src/main/ (production code in Maven/Gradle structure)
	if strings.Contains(normalizedPath, "/src/main/") || strings.HasPrefix(normalizedPath, "src/main/") {
		return false
	}

	base := filepath.Base(path)
	name := strings.TrimSuffix(base, ".java")

	// JUnit conventions: *Test.java, *Tests.java, Test*.java
	if strings.HasSuffix(name, "Test") || strings.HasSuffix(name, "Tests") {
		return true
	}
	if strings.HasPrefix(name, "Test") {
		return true
	}

	// Files in test/ or tests/ directory
	if strings.Contains(normalizedPath, "/test/") || strings.Contains(normalizedPath, "/tests/") {
		return true
	}
	if strings.Contains(normalizedPath, "/src/test/") {
		return true
	}

	return false
}

func isKotlinTestFile(path string) bool {
	return kotlinast.IsKotlinTestFile(path)
}

func isJSTestFile(path string) bool {
	base := filepath.Base(path)
	lowerBase := strings.ToLower(base)

	if strings.Contains(lowerBase, ".test.") || strings.Contains(lowerBase, ".spec.") {
		return true
	}

	// Playwright setup/teardown files: *.setup.{js,ts,jsx,tsx}
	ext := filepath.Ext(lowerBase)
	if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" {
		nameWithoutExt := strings.TrimSuffix(lowerBase, ext)
		if strings.HasSuffix(nameWithoutExt, ".setup") || strings.HasSuffix(nameWithoutExt, ".teardown") {
			return true
		}
	}

	// Cypress E2E test files: *.cy.{js,ts,jsx,tsx}
	if strings.Contains(lowerBase, ".cy.") {
		return true
	}

	normalizedPath := filepath.ToSlash(path)

	// Exclude fixture and mock directories (not actual test files)
	if strings.Contains(normalizedPath, "/__fixtures__/") || strings.HasPrefix(normalizedPath, "__fixtures__/") ||
		strings.Contains(normalizedPath, "/__mocks__/") || strings.HasPrefix(normalizedPath, "__mocks__/") {
		return false
	}

	if strings.Contains(normalizedPath, "/__tests__/") || strings.HasPrefix(normalizedPath, "__tests__/") {
		return true
	}

	// Cypress e2e/ and component/ directories
	if strings.Contains(normalizedPath, "/cypress/e2e/") || strings.Contains(normalizedPath, "/cypress/component/") {
		return true
	}

	// Files in test/ or tests/ directory (common convention like other languages)
	if strings.Contains(normalizedPath, "/test/") || strings.HasPrefix(normalizedPath, "test/") {
		return true
	}
	if strings.Contains(normalizedPath, "/tests/") || strings.HasPrefix(normalizedPath, "tests/") {
		return true
	}

	return false
}

func isPythonTestFile(path string) bool {
	base := filepath.Base(path)

	// pytest conventions: test_*.py or *_test.py
	if strings.HasPrefix(base, "test_") && strings.HasSuffix(base, ".py") {
		return true
	}
	if strings.HasSuffix(base, "_test.py") {
		return true
	}

	// conftest.py is a pytest configuration/fixture file, not a test file.
	// It's discovered as a config file but doesn't contain tests.
	if base == "conftest.py" {
		return false
	}

	// Files in tests/ directory
	normalizedPath := filepath.ToSlash(path)
	if strings.Contains(normalizedPath, "/tests/") || strings.HasPrefix(normalizedPath, "tests/") {
		return strings.HasSuffix(base, ".py")
	}

	return false
}

func isCSharpTestFile(path string) bool {
	// Check filename pattern first
	if dotnetast.IsCSharpTestFileName(path) {
		return true
	}

	// Check directory patterns for .NET project conventions
	normalizedPath := filepath.ToSlash(path)
	if strings.Contains(normalizedPath, "/test/") || strings.Contains(normalizedPath, "/tests/") {
		return true
	}
	if strings.Contains(normalizedPath, ".Tests/") || strings.Contains(normalizedPath, ".Test/") {
		return true
	}
	if strings.Contains(normalizedPath, ".Specs/") || strings.Contains(normalizedPath, ".Spec/") {
		return true
	}
	// "test/", "tests/", "Tests/" as project folder patterns
	if strings.HasPrefix(normalizedPath, "test/") || strings.HasPrefix(normalizedPath, "tests/") ||
		strings.HasPrefix(normalizedPath, "Tests/") || strings.Contains(normalizedPath, "/Tests/") {
		return true
	}

	return false
}

func isRubyTestFile(path string) bool {
	base := filepath.Base(path)

	// Exclude config/helper files
	if base == "spec_helper.rb" || base == "rails_helper.rb" {
		return false
	}

	// RSpec convention: *_spec.rb
	if strings.HasSuffix(base, "_spec.rb") {
		return true
	}

	// Minitest convention: *_test.rb
	if strings.HasSuffix(base, "_test.rb") {
		return true
	}

	normalizedPath := filepath.ToSlash(path)

	// Files in spec/ directory (excluding spec/support/ subdirectory)
	if strings.Contains(normalizedPath, "/spec/") || strings.HasPrefix(normalizedPath, "spec/") {
		// Exclude spec/support/ directory (helpers, not tests)
		if strings.Contains(normalizedPath, "/spec/support/") || strings.HasPrefix(normalizedPath, "spec/support/") {
			return false
		}
		return strings.HasSuffix(base, ".rb")
	}

	// Minitest convention: Files in test/ directory
	if strings.Contains(normalizedPath, "/test/") || strings.HasPrefix(normalizedPath, "test/") {
		return strings.HasSuffix(base, ".rb")
	}

	return false
}

func isRustTestFile(path string) bool {
	base := filepath.Base(path)

	// Rust test file convention: *_test.rs
	if strings.HasSuffix(base, "_test.rs") {
		return true
	}

	normalizedPath := filepath.ToSlash(path)

	// tests/ directory: all .rs files are candidates (content matcher filters non-tests)
	if strings.Contains(normalizedPath, "/tests/") || strings.HasPrefix(normalizedPath, "tests/") {
		return strings.HasSuffix(base, ".rs")
	}

	// src/ directory: Rust places unit tests inline with #[cfg(test)] modules
	if strings.Contains(normalizedPath, "/src/") || strings.HasPrefix(normalizedPath, "src/") {
		return strings.HasSuffix(base, ".rs")
	}

	// crates/ directory: Cargo workspaces often use crates/ for sub-crates
	// Each sub-crate may have inline tests with #[cfg(test)] modules
	if strings.Contains(normalizedPath, "/crates/") || strings.HasPrefix(normalizedPath, "crates/") {
		return strings.HasSuffix(base, ".rs")
	}

	return false
}

func isCppTestFile(path string) bool {
	base := filepath.Base(path)
	baseLower := strings.ToLower(base)

	// Strip extension to check name patterns
	ext := filepath.Ext(baseLower)
	name := strings.TrimSuffix(baseLower, ext)

	// Google Test conventions: *_test, *_unittest
	if strings.HasSuffix(name, "_test") || strings.HasSuffix(name, "_unittest") {
		return true
	}

	// *Test pattern (e.g., DatabaseTest.cc) - uppercase T avoids false positives like "contest.cc"
	baseOriginal := filepath.Base(path)
	nameOriginal := strings.TrimSuffix(baseOriginal, filepath.Ext(baseOriginal))
	if strings.HasSuffix(nameOriginal, "Test") && len(nameOriginal) > 4 {
		return true
	}

	normalizedPath := filepath.ToSlash(path)

	// test/ or tests/ directory
	if strings.Contains(normalizedPath, "/test/") || strings.Contains(normalizedPath, "/tests/") {
		return true
	}
	if strings.HasPrefix(normalizedPath, "test/") || strings.HasPrefix(normalizedPath, "tests/") {
		return true
	}

	return false
}

func isPHPTestFile(path string) bool {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, ".php")

	if strings.HasSuffix(name, "Test") || strings.HasSuffix(name, "Tests") {
		return true
	}
	if strings.HasPrefix(name, "Test") {
		return true
	}

	normalizedPath := filepath.ToSlash(path)

	if strings.Contains(normalizedPath, "/test/") || strings.Contains(normalizedPath, "/tests/") {
		return true
	}
	if strings.HasPrefix(normalizedPath, "test/") || strings.HasPrefix(normalizedPath, "tests/") {
		return true
	}

	return false
}

func isSwiftTestFile(path string) bool {
	return swiftast.IsSwiftTestFile(path)
}

func matchesAnyPattern(path, rootPath string, patterns []string) bool {
	relPath, err := filepath.Rel(rootPath, path)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range patterns {
		matched, err := doublestar.Match(pattern, relPath)
		if err != nil {
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

// ScanStream performs streaming scan that returns results through a channel.
// Each file is sent as FileResult immediately after parsing completes.
//
// This method enables memory-efficient processing of large repositories by allowing
// consumers to process files incrementally rather than accumulating all results.
//
// Usage:
//
//	results, err := scanner.ScanStream(ctx, src)
//	if err != nil {
//	    return err
//	}
//	for result := range results {
//	    if result.Err != nil {
//	        log.Warn("parse error", "path", result.Path, "err", result.Err)
//	        continue
//	    }
//	    // Process result.File
//	}
//
// The channel is closed when:
//   - All files have been processed
//   - Context is cancelled
//   - Context deadline exceeded
//
// Parse errors are included in FileResult.Err rather than aborting the scan.
// The caller is responsible for calling src.Close() when done.
func (s *Scanner) ScanStream(ctx context.Context, src source.Source) (<-chan *FileResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.options.Timeout)

	// Config parsing first (projectScope required for detection)
	// Config errors are not propagated in streaming mode.
	// Use Scan() if config error reporting is required.
	if s.projectScope == nil {
		configFiles := s.discoverConfigFiles(ctx, src)
		var configErrors []ScanError
		s.projectScope = s.parseConfigFiles(ctx, src, configFiles, &configErrors)
		s.detector.SetProjectScope(s.projectScope)
	}

	// Check for early cancellation
	if err := ctx.Err(); err != nil {
		cancel()
		return nil, err
	}

	// Unbuffered channel for natural backpressure
	out := make(chan *FileResult)

	go func() {
		defer close(out)
		defer cancel()

		workers := s.options.Workers
		if workers <= 0 {
			workers = runtime.GOMAXPROCS(0)
		}
		if workers > MaxWorkers {
			workers = MaxWorkers
		}

		sem := semaphore.NewWeighted(int64(workers))
		var wg sync.WaitGroup

		for discoveryResult := range s.discoverTestFilesStream(ctx, src) {
			if ctx.Err() != nil {
				break
			}

			if discoveryResult.Err != nil {
				select {
				case out <- &FileResult{
					Err:  discoveryResult.Err,
					Path: "",
				}:
				case <-ctx.Done():
					return
				}
				continue
			}

			path := discoveryResult.Path

			if err := sem.Acquire(ctx, 1); err != nil {
				break
			}

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				defer sem.Release(1)

				result := s.parseFileToResult(ctx, src, filePath)

				select {
				case out <- result:
				case <-ctx.Done():
				}
			}(path)
		}

		wg.Wait()
	}()

	return out, nil
}

// parseFileToResult parses a single file and returns FileResult.
// This is the streaming-oriented version that wraps parseFile.
func (s *Scanner) parseFileToResult(ctx context.Context, src source.Source, path string) *FileResult {
	testFile, scanErr, confidence := s.parseFile(ctx, src, path)

	if scanErr != nil {
		return &FileResult{
			Err:        scanErr.Err,
			Path:       path,
			Confidence: confidence,
		}
	}

	if testFile == nil {
		// File was skipped (unknown framework, low confidence, etc.)
		return &FileResult{
			Path:       path,
			Confidence: confidence,
		}
	}

	return &FileResult{
		File:       testFile,
		Path:       path,
		Confidence: confidence,
	}
}

func Scan(ctx context.Context, src source.Source, opts ...ScanOption) (*ScanResult, error) {
	scanner := NewScanner(opts...)
	return scanner.Scan(ctx, src)
}

// ScanStreaming performs streaming scan that returns results through a channel.
// Each file is sent as FileResult immediately after parsing completes.
//
// This enables memory-efficient processing of large repositories by allowing
// consumers to process files incrementally rather than accumulating all results.
//
// Usage:
//
//	results, err := parser.ScanStreaming(ctx, src)
//	if err != nil {
//	    return err
//	}
//	batch := make([]*domain.TestFile, 0, 100)
//	for result := range results {
//	    if result.Err != nil {
//	        log.Warn("parse error", "path", result.Path, "err", result.Err)
//	        continue
//	    }
//	    if result.File == nil {
//	        continue // Skipped file (unknown framework)
//	    }
//	    batch = append(batch, result.File)
//	    if len(batch) >= 100 {
//	        processBatch(batch)
//	        batch = batch[:0]
//	    }
//	}
//	if len(batch) > 0 {
//	    processBatch(batch)
//	}
//
// The caller is responsible for calling src.Close() when done.
func ScanStreaming(ctx context.Context, src source.Source, opts ...ScanOption) (<-chan *FileResult, error) {
	scanner := NewScanner(opts...)
	return scanner.ScanStream(ctx, src)
}
