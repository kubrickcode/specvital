package gemini

import (
	"log/slog"
	"slices"
	"strings"
)

// Phase1Metrics holds quality metrics for Phase 1 classification.
type Phase1Metrics struct {
	// ClassificationRate is the percentage of tests successfully classified (0.0-1.0).
	ClassificationRate float64
	// DomainDistribution maps domain names to their test counts.
	DomainDistribution map[string]int
	// FallbackRate is the percentage of tests that required individual processing fallback (0.0-1.0).
	FallbackRate float64
	// RetryRate is the percentage of batches that required retries (0.0-1.0).
	RetryRate float64
	// TotalTests is the total number of tests processed.
	TotalTests int
	// UncategorizedRate is the percentage of tests classified as Uncategorized/General (0.0-1.0).
	UncategorizedRate float64
}

// Phase1MetricsCollector accumulates metrics during Phase 1 processing.
type Phase1MetricsCollector struct {
	batchCount     int
	fallbackCount  int
	retryCount     int
	results        []v3BatchResult
	totalTestCount int
}

// NewPhase1MetricsCollector creates a new metrics collector.
func NewPhase1MetricsCollector() *Phase1MetricsCollector {
	return &Phase1MetricsCollector{
		results: make([]v3BatchResult, 0),
	}
}

// RecordBatch records results from a single batch processing.
func (c *Phase1MetricsCollector) RecordBatch(results []v3BatchResult, retries, fallbacks int) {
	c.results = append(c.results, results...)
	c.batchCount++
	c.retryCount += retries
	c.fallbackCount += fallbacks
	c.totalTestCount += len(results)
}

// Collect calculates final metrics from all recorded batches.
func (c *Phase1MetricsCollector) Collect() Phase1Metrics {
	if c.totalTestCount == 0 {
		return Phase1Metrics{
			ClassificationRate: 0,
			DomainDistribution: make(map[string]int),
			FallbackRate:       0,
			RetryRate:          0,
			TotalTests:         0,
			UncategorizedRate:  0,
		}
	}

	uncategorizedCount := 0
	domainDist := make(map[string]int)

	for _, r := range c.results {
		domainDist[r.Domain]++
		if isUncategorized(r.Domain, r.Feature) {
			uncategorizedCount++
		}
	}

	return Phase1Metrics{
		ClassificationRate: float64(c.totalTestCount-c.fallbackCount) / float64(c.totalTestCount),
		DomainDistribution: domainDist,
		FallbackRate:       float64(c.fallbackCount) / float64(c.totalTestCount),
		RetryRate:          safeRate(c.retryCount, c.batchCount),
		TotalTests:         c.totalTestCount,
		UncategorizedRate:  float64(uncategorizedCount) / float64(c.totalTestCount),
	}
}

// LogMetrics logs the collected metrics.
// Logs at error level if UncategorizedRate > 0 (should be eliminated).
func (c *Phase1MetricsCollector) LogMetrics(analysisID string) {
	metrics := c.Collect()

	logAttrs := []any{
		"analysis_id", analysisID,
		"total_tests", metrics.TotalTests,
		"classification_rate", metrics.ClassificationRate,
		"uncategorized_rate", metrics.UncategorizedRate,
		"retry_rate", metrics.RetryRate,
		"fallback_rate", metrics.FallbackRate,
		"domain_count", len(metrics.DomainDistribution),
	}

	if metrics.UncategorizedRate > 0 {
		slog.Error("phase 1 quality violation: uncategorized tests detected", logAttrs...)
		return
	}

	slog.Info("phase 1 metrics", logAttrs...)
}

// isUncategorized checks if a domain/feature pair represents an uncategorized classification.
func isUncategorized(domain, feature string) bool {
	normalizedDomain := strings.ToLower(strings.TrimSpace(domain))
	normalizedFeature := strings.ToLower(strings.TrimSpace(feature))

	uncategorizedDomains := []string{"uncategorized", "general", "other", "misc", "miscellaneous"}
	uncategorizedFeatures := []string{"general", "other", "misc", "miscellaneous", "uncategorized"}

	return slices.Contains(uncategorizedDomains, normalizedDomain) ||
		slices.Contains(uncategorizedFeatures, normalizedFeature)
}

// safeRate calculates a rate avoiding division by zero.
func safeRate(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}
