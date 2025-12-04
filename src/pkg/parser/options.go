package parser

import "time"

type ScanOptions struct {
	ExcludePatterns []string
	MaxFileSize     int64
	Patterns        []string
	Timeout         time.Duration
	Workers         int
}

type ScanOption func(*ScanOptions)

func WithExclude(patterns []string) ScanOption {
	return func(o *ScanOptions) {
		o.ExcludePatterns = patterns
	}
}

func WithScanMaxFileSize(size int64) ScanOption {
	return func(o *ScanOptions) {
		if size < 0 {
			return
		}
		o.MaxFileSize = size
	}
}

func WithScanPatterns(patterns []string) ScanOption {
	return func(o *ScanOptions) {
		o.Patterns = patterns
	}
}

func WithTimeout(d time.Duration) ScanOption {
	return func(o *ScanOptions) {
		if d < 0 {
			return
		}
		o.Timeout = d
	}
}

func WithWorkers(n int) ScanOption {
	return func(o *ScanOptions) {
		if n < 0 {
			return
		}
		o.Workers = n
	}
}
