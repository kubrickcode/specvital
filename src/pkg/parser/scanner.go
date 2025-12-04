package parser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/specvital/core/domain"
	"github.com/specvital/core/parser/strategies"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const (
	DefaultWorkers = 0
	DefaultTimeout = 5 * time.Minute
	MaxWorkers     = 1024
)

var (
	ErrScanCancelled = errors.New("scanner: scan cancelled")
	ErrScanTimeout   = errors.New("scanner: scan timeout")
)

type ScanResult struct {
	Errors    []ScanError
	Inventory *domain.Inventory
}

type ScanError struct {
	Err   error
	Path  string
	Phase string
}

func (e ScanError) Error() string {
	if e.Path == "" {
		return fmt.Sprintf("[%s] %v", e.Phase, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %v", e.Phase, e.Path, e.Err)
}

func Scan(ctx context.Context, rootPath string, opts ...ScanOption) (*ScanResult, error) {
	options := &ScanOptions{
		ExcludePatterns: nil,
		MaxFileSize:     DefaultMaxFileSize,
		Patterns:        nil,
		Timeout:         DefaultTimeout,
		Workers:         DefaultWorkers,
	}

	for _, opt := range opts {
		opt(options)
	}

	workers := options.Workers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	if workers > MaxWorkers {
		workers = MaxWorkers
	}

	timeout := options.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	detectorOpts := buildDetectorOpts(options)
	detectionResult, err := DetectTestFiles(ctx, rootPath, detectorOpts...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrScanTimeout
		}
		if errors.Is(err, context.Canceled) {
			return nil, ErrScanCancelled
		}
		return nil, fmt.Errorf("scanner: detection failed: %w", err)
	}

	scanResult := &ScanResult{
		Errors: make([]ScanError, 0),
		Inventory: &domain.Inventory{
			Files:    make([]domain.TestFile, 0),
			RootPath: rootPath,
		},
	}

	for _, detErr := range detectionResult.Errors {
		scanResult.Errors = append(scanResult.Errors, ScanError{
			Err:   detErr,
			Path:  "",
			Phase: "detection",
		})
	}

	if len(detectionResult.Files) == 0 {
		return scanResult, nil
	}

	files, errs := parseFilesParallel(ctx, detectionResult.Files, workers)

	scanResult.Inventory.Files = files
	scanResult.Errors = append(scanResult.Errors, errs...)

	return scanResult, nil
}

func buildDetectorOpts(options *ScanOptions) []DetectorOption {
	var detectorOpts []DetectorOption

	if len(options.ExcludePatterns) > 0 {
		merged := make([]string, 0, len(DefaultSkipPatterns)+len(options.ExcludePatterns))
		merged = append(merged, DefaultSkipPatterns...)
		merged = append(merged, options.ExcludePatterns...)
		detectorOpts = append(detectorOpts, WithSkipPatterns(merged))
	}

	if len(options.Patterns) > 0 {
		detectorOpts = append(detectorOpts, WithPatterns(options.Patterns))
	}

	if options.MaxFileSize > 0 {
		detectorOpts = append(detectorOpts, WithMaxFileSize(options.MaxFileSize))
	}

	return detectorOpts
}

func parseFilesParallel(ctx context.Context, files []string, workers int) ([]domain.TestFile, []ScanError) {
	sem := semaphore.NewWeighted(int64(workers))
	g, gCtx := errgroup.WithContext(ctx)

	var (
		mu         sync.Mutex
		results    = make([]domain.TestFile, 0, len(files))
		scanErrors = make([]ScanError, 0)
	)

	for _, file := range files {
		g.Go(func() error {
			if err := sem.Acquire(gCtx, 1); err != nil {
				return nil // Context cancelled
			}
			defer sem.Release(1)

			testFile, err := parseFile(gCtx, file)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				scanErrors = append(scanErrors, ScanError{
					Err:   err,
					Path:  file,
					Phase: "parsing",
				})
				return nil // Continue with other files
			}

			if testFile != nil {
				results = append(results, *testFile)
			}

			return nil
		})
	}

	_ = g.Wait() // Errors are collected in scanErrors

	return results, scanErrors
}

func parseFile(ctx context.Context, path string) (*domain.TestFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	strategy := strategies.FindStrategy(path, content)
	if strategy == nil {
		return nil, nil // No matching strategy
	}

	testFile, err := strategy.Parse(ctx, content, path)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	return testFile, nil
}
