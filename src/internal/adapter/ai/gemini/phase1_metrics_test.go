package gemini

import "testing"

func TestPhase1MetricsCollector_EmptyResults(t *testing.T) {
	collector := NewPhase1MetricsCollector()
	metrics := collector.Collect()

	if metrics.TotalTests != 0 {
		t.Errorf("expected TotalTests=0, got %d", metrics.TotalTests)
	}
	if metrics.ClassificationRate != 0 {
		t.Errorf("expected ClassificationRate=0, got %f", metrics.ClassificationRate)
	}
	if metrics.UncategorizedRate != 0 {
		t.Errorf("expected UncategorizedRate=0, got %f", metrics.UncategorizedRate)
	}
	if len(metrics.DomainDistribution) != 0 {
		t.Errorf("expected empty DomainDistribution, got %v", metrics.DomainDistribution)
	}
}

func TestPhase1MetricsCollector_SingleBatch(t *testing.T) {
	collector := NewPhase1MetricsCollector()

	results := []v3BatchResult{
		{Domain: "Authentication", Feature: "Login"},
		{Domain: "Authentication", Feature: "Logout"},
		{Domain: "Navigation", Feature: "Routing"},
	}
	collector.RecordBatch(results, 0, 0)

	metrics := collector.Collect()

	if metrics.TotalTests != 3 {
		t.Errorf("expected TotalTests=3, got %d", metrics.TotalTests)
	}
	if metrics.ClassificationRate != 1.0 {
		t.Errorf("expected ClassificationRate=1.0, got %f", metrics.ClassificationRate)
	}
	if metrics.UncategorizedRate != 0 {
		t.Errorf("expected UncategorizedRate=0, got %f", metrics.UncategorizedRate)
	}
	if metrics.DomainDistribution["Authentication"] != 2 {
		t.Errorf("expected Authentication=2, got %d", metrics.DomainDistribution["Authentication"])
	}
	if metrics.DomainDistribution["Navigation"] != 1 {
		t.Errorf("expected Navigation=1, got %d", metrics.DomainDistribution["Navigation"])
	}
}

func TestPhase1MetricsCollector_MultipleBatches(t *testing.T) {
	collector := NewPhase1MetricsCollector()

	batch1 := []v3BatchResult{
		{Domain: "Authentication", Feature: "Login"},
		{Domain: "Authentication", Feature: "Logout"},
	}
	collector.RecordBatch(batch1, 1, 0)

	batch2 := []v3BatchResult{
		{Domain: "Navigation", Feature: "Routing"},
		{Domain: "Forms", Feature: "Validation"},
	}
	collector.RecordBatch(batch2, 0, 1)

	metrics := collector.Collect()

	if metrics.TotalTests != 4 {
		t.Errorf("expected TotalTests=4, got %d", metrics.TotalTests)
	}
	if metrics.RetryRate != 0.5 {
		t.Errorf("expected RetryRate=0.5 (1 retry / 2 batches), got %f", metrics.RetryRate)
	}
	if metrics.FallbackRate != 0.25 {
		t.Errorf("expected FallbackRate=0.25 (1 fallback / 4 tests), got %f", metrics.FallbackRate)
	}
	if metrics.ClassificationRate != 0.75 {
		t.Errorf("expected ClassificationRate=0.75 (3 successful / 4 total), got %f", metrics.ClassificationRate)
	}
}

func TestPhase1MetricsCollector_UncategorizedDetection(t *testing.T) {
	collector := NewPhase1MetricsCollector()

	results := []v3BatchResult{
		{Domain: "Authentication", Feature: "Login"},
		{Domain: "Uncategorized", Feature: "General"},
		{Domain: "Navigation", Feature: "Other"},
		{Domain: "Forms", Feature: "Validation"},
	}
	collector.RecordBatch(results, 0, 0)

	metrics := collector.Collect()

	if metrics.TotalTests != 4 {
		t.Errorf("expected TotalTests=4, got %d", metrics.TotalTests)
	}
	// 2 uncategorized: "Uncategorized/General" and "Navigation/Other"
	if metrics.UncategorizedRate != 0.5 {
		t.Errorf("expected UncategorizedRate=0.5 (2 uncategorized / 4 total), got %f", metrics.UncategorizedRate)
	}
}

func TestIsUncategorized(t *testing.T) {
	tests := []struct {
		domain   string
		feature  string
		expected bool
	}{
		{"Uncategorized", "General", true},
		{"uncategorized", "general", true},
		{"UNCATEGORIZED", "GENERAL", true},
		{"General", "Something", true},
		{"Other", "Feature", true},
		{"Misc", "Test", true},
		{"Miscellaneous", "Something", true},
		{"Authentication", "General", true},
		{"Navigation", "Other", true},
		{"Forms", "Misc", true},
		{"Authentication", "Login", false},
		{"Navigation", "Routing", false},
		{"Forms", "Validation", false},
		{" Uncategorized ", " General ", true},
		{"  general  ", "feature", true},
	}

	for _, tc := range tests {
		t.Run(tc.domain+"/"+tc.feature, func(t *testing.T) {
			result := isUncategorized(tc.domain, tc.feature)
			if result != tc.expected {
				t.Errorf("isUncategorized(%q, %q) = %v, expected %v",
					tc.domain, tc.feature, result, tc.expected)
			}
		})
	}
}

func TestSafeRate(t *testing.T) {
	tests := []struct {
		numerator   int
		denominator int
		expected    float64
	}{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
		{1, 2, 0.5},
		{3, 4, 0.75},
		{5, 5, 1.0},
	}

	for _, tc := range tests {
		result := safeRate(tc.numerator, tc.denominator)
		if result != tc.expected {
			t.Errorf("safeRate(%d, %d) = %f, expected %f",
				tc.numerator, tc.denominator, result, tc.expected)
		}
	}
}

func TestPhase1MetricsCollector_DomainDistribution(t *testing.T) {
	collector := NewPhase1MetricsCollector()

	batch1 := []v3BatchResult{
		{Domain: "Authentication", Feature: "Login"},
		{Domain: "Authentication", Feature: "Logout"},
		{Domain: "Authentication", Feature: "Session"},
	}
	collector.RecordBatch(batch1, 0, 0)

	batch2 := []v3BatchResult{
		{Domain: "Navigation", Feature: "Routing"},
		{Domain: "Authentication", Feature: "OAuth"},
	}
	collector.RecordBatch(batch2, 0, 0)

	metrics := collector.Collect()

	if metrics.DomainDistribution["Authentication"] != 4 {
		t.Errorf("expected Authentication=4 across batches, got %d",
			metrics.DomainDistribution["Authentication"])
	}
	if metrics.DomainDistribution["Navigation"] != 1 {
		t.Errorf("expected Navigation=1, got %d",
			metrics.DomainDistribution["Navigation"])
	}
	if len(metrics.DomainDistribution) != 2 {
		t.Errorf("expected 2 domains, got %d", len(metrics.DomainDistribution))
	}
}
