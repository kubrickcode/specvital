package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/specvital/worker/internal/adapter/ai/prompt"
	"github.com/specvital/worker/internal/adapter/ai/reliability"
	"github.com/specvital/worker/internal/domain/specview"
)

// mockV3Provider simulates a V3 Provider for integration testing.
// Allows injecting custom batch processor behavior for various scenarios.
type mockV3Provider struct {
	phase1Model    string
	batchProcessor func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error)
	phase1Retry    *reliability.Retryer
	phase1CB       *reliability.CircuitBreaker
}

func newMockV3Provider(batchProcessor func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error)) *mockV3Provider {
	return &mockV3Provider{
		phase1Model:    "test-model",
		batchProcessor: batchProcessor,
		phase1Retry:    reliability.NewRetryer(reliability.DefaultPhase1RetryConfig()),
		phase1CB:       reliability.NewCircuitBreaker(reliability.DefaultPhase1CircuitConfig()),
	}
}

// classifyDomainsV3Integration performs V3 classification using mock batch processor.
func (m *mockV3Provider) classifyDomainsV3Integration(ctx context.Context, input specview.Phase1Input, lang specview.Language) (*specview.Phase1Output, *specview.TokenUsage, error) {
	tests := flattenTests(input.Files)
	if len(tests) == 0 {
		return nil, nil, fmt.Errorf("%w: no tests to classify", specview.ErrInvalidInput)
	}

	batches := splitIntoBatches(tests, v3BatchSize)
	totalUsage := &specview.TokenUsage{Model: m.phase1Model}
	totalRetries := 0
	totalFallbacks := 0

	var allResults [][]v3BatchResult
	var existingDomains []prompt.DomainSummary

	for batchIdx, batch := range batches {
		if err := ctx.Err(); err != nil {
			return nil, totalUsage, err
		}

		results, usage, retries, fallbacks, err := m.processV3BatchWithRetry(ctx, batch, existingDomains, lang)
		accumulateTokenUsage(totalUsage, usage)
		totalRetries += retries
		totalFallbacks += fallbacks

		if err != nil {
			return nil, totalUsage, fmt.Errorf("batch %d/%d failed: %w", batchIdx+1, len(batches), err)
		}

		allResults = append(allResults, results)
		existingDomains = extractDomainSummaries(results, existingDomains)
	}

	_ = totalRetries
	_ = totalFallbacks

	output := mergeV3Results(allResults, input)
	return output, totalUsage, nil
}

func (m *mockV3Provider) processV3BatchWithRetry(
	ctx context.Context,
	tests []specview.TestForAssignment,
	existingDomains []prompt.DomainSummary,
	lang specview.Language,
) ([]v3BatchResult, *specview.TokenUsage, int, int, error) {
	totalUsage := &specview.TokenUsage{Model: m.phase1Model}
	retryCount := 0
	fallbackCount := 0

	if len(tests) == 0 {
		return []v3BatchResult{}, totalUsage, 0, 0, nil
	}

	// Try batch processing with retries
	for attempt := 1; attempt <= v3MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, totalUsage, retryCount, fallbackCount, err
		}

		results, usage, err := m.batchProcessor(ctx, tests, existingDomains)
		if usage != nil {
			accumulateTokenUsage(totalUsage, usage)
		}

		if err == nil && len(results) == len(tests) {
			return results, totalUsage, retryCount, fallbackCount, nil
		}

		retryCount++
	}

	// Batch processing failed - try splitting
	if len(tests) >= v3MinBatchSizeForSplit {
		results, splitUsage, splitRetries, splitFallbacks, err := m.processV3SplitBatch(ctx, tests, existingDomains, lang)
		accumulateTokenUsage(totalUsage, splitUsage)
		retryCount += splitRetries
		fallbackCount += splitFallbacks

		if err == nil {
			return results, totalUsage, retryCount, fallbackCount, nil
		}
	}

	// Fall back to individual processing
	results, indivUsage, indivFallbacks := m.processV3Individual(ctx, tests, existingDomains)
	accumulateTokenUsage(totalUsage, indivUsage)
	fallbackCount += indivFallbacks

	return results, totalUsage, retryCount, fallbackCount, nil
}

func (m *mockV3Provider) processV3SplitBatch(
	ctx context.Context,
	tests []specview.TestForAssignment,
	existingDomains []prompt.DomainSummary,
	lang specview.Language,
) ([]v3BatchResult, *specview.TokenUsage, int, int, error) {
	left, right := splitBatch(tests)
	totalUsage := &specview.TokenUsage{Model: m.phase1Model}
	totalRetries := 0
	totalFallbacks := 0

	leftResults, leftUsage, leftRetries, leftFallbacks, leftErr := m.processV3BatchWithRetry(ctx, left, existingDomains, lang)
	accumulateTokenUsage(totalUsage, leftUsage)
	totalRetries += leftRetries
	totalFallbacks += leftFallbacks

	if leftErr != nil {
		return nil, totalUsage, totalRetries, totalFallbacks, fmt.Errorf("left split failed: %w", leftErr)
	}

	rightResults, rightUsage, rightRetries, rightFallbacks, rightErr := m.processV3BatchWithRetry(ctx, right, existingDomains, lang)
	accumulateTokenUsage(totalUsage, rightUsage)
	totalRetries += rightRetries
	totalFallbacks += rightFallbacks

	if rightErr != nil {
		return nil, totalUsage, totalRetries, totalFallbacks, fmt.Errorf("right split failed: %w", rightErr)
	}

	results := make([]v3BatchResult, 0, len(leftResults)+len(rightResults))
	results = append(results, leftResults...)
	results = append(results, rightResults...)

	return results, totalUsage, totalRetries, totalFallbacks, nil
}

func (m *mockV3Provider) processV3Individual(
	ctx context.Context,
	tests []specview.TestForAssignment,
	existingDomains []prompt.DomainSummary,
) ([]v3BatchResult, *specview.TokenUsage, int) {
	totalUsage := &specview.TokenUsage{Model: m.phase1Model}
	results := make([]v3BatchResult, 0, len(tests))
	fallbackCount := 0

	for _, test := range tests {
		if err := ctx.Err(); err != nil {
			return results, totalUsage, fallbackCount
		}

		singleTest := []specview.TestForAssignment{test}
		result, usage, err := m.batchProcessor(ctx, singleTest, existingDomains)
		if usage != nil {
			accumulateTokenUsage(totalUsage, usage)
		}

		if err != nil || len(result) != 1 {
			domain, feature := deriveDomainFromPath(test.FilePath)
			results = append(results, v3BatchResult{
				Domain:     domain,
				DomainDesc: "Derived from file path",
				Feature:    feature,
			})
			fallbackCount++
			continue
		}

		results = append(results, result[0])
	}

	return results, totalUsage, fallbackCount
}

// Helper to create test input
func createV3TestInput(fileCount, testsPerFile int) specview.Phase1Input {
	files := make([]specview.FileInfo, fileCount)
	testIndex := 0

	for i := 0; i < fileCount; i++ {
		tests := make([]specview.TestInfo, testsPerFile)
		for j := 0; j < testsPerFile; j++ {
			tests[j] = specview.TestInfo{
				Index: testIndex,
				Name:  fmt.Sprintf("test_%d_%d", i, j),
			}
			testIndex++
		}
		files[i] = specview.FileInfo{
			Path:  fmt.Sprintf("src/file_%d_test.go", i),
			Tests: tests,
		}
	}

	return specview.Phase1Input{
		AnalysisID: "test-analysis-id",
		Files:      files,
		Language:   "Korean",
	}
}

// --- Integration Tests ---

func TestV3Integration_SmallInput_SingleBatch(t *testing.T) {
	// 15 tests = 1 batch (under v3BatchSize of 20)
	input := createV3TestInput(3, 5)
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		callCount.Add(1)
		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "Authentication", Feature: "Login"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100, CandidatesTokens: 50}, nil
	})

	ctx := context.Background()
	output, usage, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount.Load() != 1 {
		t.Errorf("expected 1 batch call, got %d", callCount.Load())
	}

	if len(output.Domains) != 1 {
		t.Errorf("expected 1 domain, got %d", len(output.Domains))
	}

	if output.Domains[0].Name != "Authentication" {
		t.Errorf("expected domain 'Authentication', got %q", output.Domains[0].Name)
	}

	if usage.PromptTokens != 100 {
		t.Errorf("expected 100 prompt tokens, got %d", usage.PromptTokens)
	}
}

func TestV3Integration_MediumInput_MultipleBatches(t *testing.T) {
	// 50 tests = 3 batches (20 + 20 + 10)
	input := createV3TestInput(10, 5)
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		callCount.Add(1)
		results := make([]v3BatchResult, len(tests))
		// Use existing domains if available
		domain := "Domain1"
		if len(existingDomains) > 0 {
			domain = existingDomains[0].Name
		}
		for i := range tests {
			results[i] = v3BatchResult{Domain: domain, Feature: "Feature1"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100}, nil
	})

	ctx := context.Background()
	output, usage, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 batches expected
	if callCount.Load() != 3 {
		t.Errorf("expected 3 batch calls, got %d", callCount.Load())
	}

	// All tests should be classified
	totalTests := 0
	for _, domain := range output.Domains {
		for _, feature := range domain.Features {
			totalTests += len(feature.TestIndices)
		}
	}
	if totalTests != 50 {
		t.Errorf("expected 50 tests classified, got %d", totalTests)
	}

	// Token usage accumulated
	if usage.PromptTokens != 300 {
		t.Errorf("expected 300 accumulated prompt tokens, got %d", usage.PromptTokens)
	}
}

func TestV3Integration_LargeInput_ManyBatches(t *testing.T) {
	// 100 tests = 5 batches
	input := createV3TestInput(20, 5)
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		batchNum := callCount.Add(1)
		results := make([]v3BatchResult, len(tests))

		// Simulate varying domains per batch
		domain := fmt.Sprintf("Domain%d", batchNum)
		for i := range tests {
			results[i] = v3BatchResult{Domain: domain, Feature: "Feature1"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100, CandidatesTokens: 50}, nil
	})

	ctx := context.Background()
	output, usage, err := provider.classifyDomainsV3Integration(ctx, input, "English")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount.Load() != 5 {
		t.Errorf("expected 5 batch calls, got %d", callCount.Load())
	}

	// Should have multiple domains
	if len(output.Domains) < 2 {
		t.Errorf("expected multiple domains from different batches, got %d", len(output.Domains))
	}

	if usage.PromptTokens != 500 {
		t.Errorf("expected 500 accumulated prompt tokens, got %d", usage.PromptTokens)
	}
}

func TestV3Integration_AnchorDomainPropagation(t *testing.T) {
	// Verify that domains from previous batches are passed to subsequent batches
	input := createV3TestInput(8, 5) // 40 tests = 2 batches
	var receivedDomains [][]prompt.DomainSummary

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		receivedDomains = append(receivedDomains, existingDomains)

		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "Auth", Feature: "Login"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100}, nil
	})

	ctx := context.Background()
	_, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First batch should have no existing domains
	if len(receivedDomains) < 2 {
		t.Fatalf("expected at least 2 batches, got %d", len(receivedDomains))
	}

	if len(receivedDomains[0]) != 0 {
		t.Errorf("first batch should have 0 existing domains, got %d", len(receivedDomains[0]))
	}

	// Second batch should have domains from first batch
	if len(receivedDomains[1]) == 0 {
		t.Error("second batch should have existing domains from first batch")
	}

	if receivedDomains[1][0].Name != "Auth" {
		t.Errorf("expected anchor domain 'Auth', got %q", receivedDomains[1][0].Name)
	}
}

func TestV3Integration_RetryThenSuccess(t *testing.T) {
	input := createV3TestInput(2, 5) // 10 tests = 1 batch
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		count := callCount.Add(1)
		// Fail first 2 attempts
		if count <= 2 {
			return nil, &specview.TokenUsage{PromptTokens: 50}, fmt.Errorf("transient error")
		}

		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "Payment", Feature: "Checkout"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100}, nil
	})

	ctx := context.Background()
	output, usage, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 calls: 2 failures + 1 success
	if callCount.Load() != 3 {
		t.Errorf("expected 3 calls, got %d", callCount.Load())
	}

	if len(output.Domains) != 1 {
		t.Errorf("expected 1 domain, got %d", len(output.Domains))
	}

	// Token usage: 2 * 50 (failures) + 100 (success) = 200
	if usage.PromptTokens != 200 {
		t.Errorf("expected 200 prompt tokens, got %d", usage.PromptTokens)
	}
}

func TestV3Integration_SplitAfterRetryExhausted(t *testing.T) {
	// 8 tests, batch fails at size > 4, succeeds at 4 or less
	input := createV3TestInput(2, 4) // 8 tests
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		callCount.Add(1)

		if len(tests) > 4 {
			return nil, &specview.TokenUsage{PromptTokens: 10}, fmt.Errorf("batch too large")
		}

		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "User", Feature: "Profile"}
		}
		return results, &specview.TokenUsage{PromptTokens: 50}, nil
	})

	ctx := context.Background()
	output, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Original batch (8) fails 3 times, then splits to 4+4, both succeed
	// Calls: 3 (retries for 8) + 1 (4) + 1 (4) = 5
	if callCount.Load() != 5 {
		t.Errorf("expected 5 calls (3 retries + 2 split batches), got %d", callCount.Load())
	}

	// Verify all tests classified
	totalTests := 0
	for _, domain := range output.Domains {
		for _, feature := range domain.Features {
			totalTests += len(feature.TestIndices)
		}
	}
	if totalTests != 8 {
		t.Errorf("expected 8 tests classified, got %d", totalTests)
	}
}

func TestV3Integration_IndividualFallback(t *testing.T) {
	// 3 tests (too small to split), batch fails, falls back to individual
	input := createV3TestInput(1, 3)
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		callCount.Add(1)

		// Batch (>1 test) fails, individual succeeds
		if len(tests) > 1 {
			return nil, &specview.TokenUsage{PromptTokens: 10}, fmt.Errorf("batch fails")
		}

		return []v3BatchResult{{Domain: "Individual", Feature: "Success"}}, &specview.TokenUsage{PromptTokens: 20}, nil
	})

	ctx := context.Background()
	output, usage, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 retries for batch + 3 individual calls
	if callCount.Load() != 6 {
		t.Errorf("expected 6 calls, got %d", callCount.Load())
	}

	// All classified to Individual domain
	if len(output.Domains) != 1 || output.Domains[0].Name != "Individual" {
		t.Errorf("expected all tests in 'Individual' domain, got %v", output.Domains)
	}

	// Token usage: 3 * 10 (batch fails) + 3 * 20 (individual) = 90
	if usage.PromptTokens != 90 {
		t.Errorf("expected 90 prompt tokens, got %d", usage.PromptTokens)
	}
}

func TestV3Integration_PathBasedFallback(t *testing.T) {
	// All API calls fail - tests should be assigned path-based domains
	input := createV3TestInput(1, 2)
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		callCount.Add(1)
		return nil, &specview.TokenUsage{PromptTokens: 10}, fmt.Errorf("always fails")
	})

	ctx := context.Background()
	output, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All tests should be path-based
	if len(output.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(output.Domains))
	}

	// Input path is "src/file_0_test.go" -> src is the only directory, becomes "Src"
	if output.Domains[0].Name == uncategorizedDomainName {
		t.Errorf("expected path-based domain, got Uncategorized")
	}
	// When src is the only significant part (and it's skipped), fallback to last part of allParts
	if output.Domains[0].Name != "Src" {
		t.Errorf("expected domain 'Src' (from path), got %q", output.Domains[0].Name)
	}
}

func TestV3Integration_EmptyInput(t *testing.T) {
	input := specview.Phase1Input{
		Files:    []specview.FileInfo{},
		Language: "Korean",
	}

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		t.Error("batch processor should not be called for empty input")
		return nil, nil, nil
	})

	ctx := context.Background()
	_, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestV3Integration_ContextCancellation(t *testing.T) {
	input := createV3TestInput(10, 5) // 50 tests = 3 batches
	var callCount atomic.Int32

	ctx, cancel := context.WithCancel(context.Background())

	provider := newMockV3Provider(func(innerCtx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		callCount.Add(1)

		// Cancel after first batch
		if callCount.Load() == 1 {
			cancel()
		}

		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "Domain", Feature: "Feature"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100}, nil
	})

	_, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err == nil {
		t.Error("expected error from context cancellation")
	}

	// Should stop after first batch (or during second batch's context check)
	if callCount.Load() > 2 {
		t.Errorf("expected at most 2 batch calls before cancellation, got %d", callCount.Load())
	}
}

func TestV3Integration_MergeResultsPreservesTestIndices(t *testing.T) {
	// Verify that test indices are correctly preserved through the pipeline
	input := createV3TestInput(3, 5) // 15 tests with indices 0-14
	expectedIndices := make(map[int]bool)
	for i := 0; i < 15; i++ {
		expectedIndices[i] = true
	}

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "Domain", Feature: "Feature"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100}, nil
	})

	ctx := context.Background()
	output, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Collect all test indices from output
	foundIndices := make(map[int]bool)
	for _, domain := range output.Domains {
		for _, feature := range domain.Features {
			for _, idx := range feature.TestIndices {
				foundIndices[idx] = true
			}
		}
	}

	// Verify all expected indices are present
	for idx := range expectedIndices {
		if !foundIndices[idx] {
			t.Errorf("missing test index %d in output", idx)
		}
	}

	// Verify no extra indices
	for idx := range foundIndices {
		if !expectedIndices[idx] {
			t.Errorf("unexpected test index %d in output", idx)
		}
	}
}

func TestV3Integration_ResponseParsingWithRealJSON(t *testing.T) {
	// Test with realistic JSON response format
	input := createV3TestInput(1, 3)

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		// Simulate real API response parsing
		jsonResponse := `[{"d": "Authentication", "f": "Login"}, {"d": "Authentication", "f": "Logout"}, {"d": "Payment", "f": "Checkout"}]`

		var results []v3BatchResult
		if err := json.Unmarshal([]byte(jsonResponse), &results); err != nil {
			return nil, nil, err
		}

		return results, &specview.TokenUsage{PromptTokens: 100, CandidatesTokens: 50}, nil
	})

	ctx := context.Background()
	output, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 2 domains (Authentication, Payment)
	if len(output.Domains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(output.Domains))
	}

	// Find Authentication domain
	var authDomain *specview.DomainGroup
	for i := range output.Domains {
		if output.Domains[i].Name == "Authentication" {
			authDomain = &output.Domains[i]
			break
		}
	}

	if authDomain == nil {
		t.Fatal("Authentication domain not found")
	}

	// Should have 2 features (Login, Logout)
	if len(authDomain.Features) != 2 {
		t.Errorf("expected 2 features in Authentication, got %d", len(authDomain.Features))
	}
}

func TestV3Integration_CountMismatchTriggersRetry(t *testing.T) {
	input := createV3TestInput(1, 5) // 5 tests
	var callCount atomic.Int32

	provider := newMockV3Provider(func(ctx context.Context, tests []specview.TestForAssignment, existingDomains []prompt.DomainSummary) ([]v3BatchResult, *specview.TokenUsage, error) {
		count := callCount.Add(1)

		// Return wrong count for first 2 attempts
		if count <= 2 {
			// Return 3 results instead of 5 - simulates count mismatch
			return []v3BatchResult{
				{Domain: "D", Feature: "F"},
				{Domain: "D", Feature: "F"},
				{Domain: "D", Feature: "F"},
			}, &specview.TokenUsage{PromptTokens: 50}, nil
		}

		// Return correct count on 3rd attempt
		results := make([]v3BatchResult, len(tests))
		for i := range tests {
			results[i] = v3BatchResult{Domain: "Success", Feature: "Feature"}
		}
		return results, &specview.TokenUsage{PromptTokens: 100}, nil
	})

	ctx := context.Background()
	output, _, err := provider.classifyDomainsV3Integration(ctx, input, "Korean")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 calls expected
	if callCount.Load() != 3 {
		t.Errorf("expected 3 calls, got %d", callCount.Load())
	}

	// Final result should be from successful attempt
	if len(output.Domains) != 1 || output.Domains[0].Name != "Success" {
		t.Errorf("expected 'Success' domain, got %v", output.Domains)
	}
}
