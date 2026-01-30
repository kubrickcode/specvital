package gemini

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"unicode"

	"github.com/specvital/worker/internal/domain/specview"
)

const (
	// waveSize is the number of concurrent batches per wave.
	// 10 concurrent requests balance throughput vs rate limiting.
	waveSize = 10

	// defaultClassificationConfidence is the confidence score for AI-classified domains/features.
	defaultClassificationConfidence = 1.0
)

// classifyDomainsV2 performs Phase 1 using two-stage architecture.
// Stage 1: Extract taxonomy from file metadata (single API call).
// Stage 2: Assign tests to fixed taxonomy (parallel batch processing).
func (p *Provider) classifyDomainsV2(ctx context.Context, input specview.Phase1Input, lang specview.Language) (*specview.Phase1Output, *specview.TokenUsage, error) {
	if len(input.Files) == 0 {
		return nil, nil, fmt.Errorf("%w: no files to classify", specview.ErrInvalidInput)
	}

	slog.InfoContext(ctx, "starting phase 1 v2 two-stage classification",
		"file_count", len(input.Files),
		"test_count", countTests(input.Files),
	)

	totalUsage := &specview.TokenUsage{Model: p.phase1Model}

	// Stage 1: Taxonomy Extraction (with cache)
	cacheKey := TaxonomyCacheKey{
		AnalysisID: input.AnalysisID,
		Language:   lang,
		ModelID:    p.phase1Model,
	}

	taxonomy := p.taxonomyCache.Get(cacheKey)
	if taxonomy != nil {
		slog.InfoContext(ctx, "taxonomy cache hit, skipping Stage 1 API call",
			"analysis_id", input.AnalysisID,
			"domain_count", len(taxonomy.Domains),
		)
	} else {
		taxonomyInput := prepareTaxonomyInput(input)
		var stage1Usage *specview.TokenUsage
		var err error
		taxonomy, stage1Usage, err = p.extractTaxonomy(ctx, taxonomyInput)
		if err != nil {
			slog.WarnContext(ctx, "stage 1 taxonomy extraction failed, using heuristic fallback",
				"error", err,
			)
			// Intentionally NOT caching heuristic fallback - retry API on next call
			taxonomy = generateHeuristicTaxonomy(input.Files)
		} else {
			aggregateUsage(totalUsage, stage1Usage)
			// Cache successful taxonomy extraction
			p.taxonomyCache.Set(cacheKey, taxonomy)
		}
	}

	slog.InfoContext(ctx, "stage 1 complete",
		"domain_count", len(taxonomy.Domains),
	)

	// Stage 2: Test Assignment
	batches := createAssignmentBatches(input, assignmentBatchSize)
	if len(batches) == 0 {
		return mergeAssignmentsToPhase1Output(taxonomy, nil, input), totalUsage, nil
	}

	assignments, stage2Usage, err := p.processAssignmentWaves(ctx, batches, taxonomy, lang)
	if err != nil {
		return nil, nil, fmt.Errorf("stage 2 assignment failed: %w", err)
	}
	aggregateUsage(totalUsage, stage2Usage)

	slog.InfoContext(ctx, "stage 2 complete",
		"batch_count", len(batches),
		"assignment_count", len(assignments),
	)

	output := mergeAssignmentsToPhase1Output(taxonomy, assignments, input)

	slog.InfoContext(ctx, "phase 1 v2 classification complete",
		"domain_count", len(output.Domains),
		"total_prompt_tokens", totalUsage.PromptTokens,
		"total_output_tokens", totalUsage.CandidatesTokens,
	)

	return output, totalUsage, nil
}

// processAssignmentWaves processes batches in waves of concurrent requests.
// Each wave contains up to waveSize batches running in parallel.
func (p *Provider) processAssignmentWaves(
	ctx context.Context,
	batches []specview.AssignmentBatch,
	taxonomy *specview.TaxonomyOutput,
	lang specview.Language,
) ([]specview.AssignmentOutput, *specview.TokenUsage, error) {
	totalUsage := &specview.TokenUsage{Model: p.phase2Model}
	allAssignments := make([]specview.AssignmentOutput, 0, len(batches))

	for waveStart := 0; waveStart < len(batches); waveStart += waveSize {
		if err := ctx.Err(); err != nil {
			return nil, nil, fmt.Errorf("assignment cancelled: %w", err)
		}

		waveEnd := min(waveStart+waveSize, len(batches))
		waveBatches := batches[waveStart:waveEnd]

		slog.InfoContext(ctx, "processing assignment wave",
			"wave_start", waveStart,
			"wave_end", waveEnd,
			"wave_size", len(waveBatches),
		)

		waveAssignments, waveUsage, err := p.processWave(ctx, waveBatches, taxonomy, lang)
		if err != nil {
			return nil, nil, fmt.Errorf("wave %d-%d failed: %w", waveStart, waveEnd, err)
		}

		allAssignments = append(allAssignments, waveAssignments...)
		aggregateUsage(totalUsage, waveUsage)
	}

	return allAssignments, totalUsage, nil
}

// batchResult holds the result of a single batch assignment.
type batchResult struct {
	batchIndex int
	err        error
	output     *specview.AssignmentOutput
	usage      *specview.TokenUsage
}

// processWave processes a single wave of batches concurrently.
// Uses context cancellation to stop remaining goroutines on first error.
func (p *Provider) processWave(
	ctx context.Context,
	batches []specview.AssignmentBatch,
	taxonomy *specview.TaxonomyOutput,
	lang specview.Language,
) ([]specview.AssignmentOutput, *specview.TokenUsage, error) {
	waveCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan batchResult, len(batches))
	var wg sync.WaitGroup

	for _, batch := range batches {
		wg.Add(1)
		go func(b specview.AssignmentBatch) {
			defer wg.Done()

			output, usage, err := p.assignTestsBatch(waveCtx, b, taxonomy, lang)
			results <- batchResult{
				batchIndex: b.BatchIndex,
				err:        err,
				output:     output,
				usage:      usage,
			}
		}(batch)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results, cancel on first error
	resultMap := make(map[int]batchResult)
	var firstErr error
	for result := range results {
		if result.err != nil && firstErr == nil {
			firstErr = fmt.Errorf("batch %d failed: %w", result.batchIndex, result.err)
			cancel()
		}
		if firstErr == nil {
			resultMap[result.batchIndex] = result
		}
	}

	if firstErr != nil {
		return nil, nil, firstErr
	}

	assignments := make([]specview.AssignmentOutput, 0, len(batches))
	totalUsage := &specview.TokenUsage{Model: p.phase2Model}

	for _, batch := range batches {
		result := resultMap[batch.BatchIndex]
		if result.output != nil {
			assignments = append(assignments, *result.output)
		}
		aggregateUsage(totalUsage, result.usage)
	}

	return assignments, totalUsage, nil
}

// mergeAssignmentsToPhase1Output converts taxonomy and assignments to Phase1Output.
// Maps test indices to domain/feature structure expected by downstream processing.
func mergeAssignmentsToPhase1Output(
	taxonomy *specview.TaxonomyOutput,
	assignments []specview.AssignmentOutput,
	input specview.Phase1Input,
) *specview.Phase1Output {
	// Handle nil or empty taxonomy
	if taxonomy == nil || len(taxonomy.Domains) == 0 {
		return &specview.Phase1Output{Domains: []specview.DomainGroup{
			{
				Confidence:  defaultClassificationConfidence,
				Description: "Files that do not fit into specific domains",
				Features: []specview.FeatureGroup{
					{
						Confidence:  defaultClassificationConfidence,
						Name:        uncategorizedFeatureName,
						TestIndices: collectAllTestIndices(input),
					},
				},
				Name: uncategorizedDomainName,
			},
		}}
	}

	// Build domain/feature → test indices map from assignments
	testMap := make(map[string]map[string][]int)
	for _, assignment := range assignments {
		for _, a := range assignment.Assignments {
			if testMap[a.Domain] == nil {
				testMap[a.Domain] = make(map[string][]int)
			}
			testMap[a.Domain][a.Feature] = append(testMap[a.Domain][a.Feature], a.TestIndices...)
		}
	}

	// Build output domains from taxonomy structure
	domains := make([]specview.DomainGroup, 0, len(taxonomy.Domains))
	for _, td := range taxonomy.Domains {
		domain := specview.DomainGroup{
			Confidence:  defaultClassificationConfidence,
			Description: td.Description,
			Name:        td.Name,
		}

		featureTestMap := testMap[td.Name]
		features := make([]specview.FeatureGroup, 0, len(td.Features))

		for _, tf := range td.Features {
			indices := featureTestMap[tf.Name]
			if len(indices) == 0 {
				continue
			}

			feature := specview.FeatureGroup{
				Confidence:  defaultClassificationConfidence,
				Name:        tf.Name,
				TestIndices: indices,
			}
			features = append(features, feature)
		}

		if len(features) > 0 {
			domain.Features = features
			domains = append(domains, domain)
		}
	}

	// Ensure Uncategorized domain exists if there are unassigned tests
	if len(domains) == 0 {
		domains = append(domains, specview.DomainGroup{
			Confidence:  defaultClassificationConfidence,
			Description: "Files that do not fit into specific domains",
			Features: []specview.FeatureGroup{
				{
					Confidence:  defaultClassificationConfidence,
					Name:        uncategorizedFeatureName,
					TestIndices: collectAllTestIndices(input),
				},
			},
			Name: uncategorizedDomainName,
		})
	}

	return &specview.Phase1Output{Domains: domains}
}

// collectAllTestIndices extracts all test indices from input.
func collectAllTestIndices(input specview.Phase1Input) []int {
	var indices []int
	for _, file := range input.Files {
		for _, test := range file.Tests {
			indices = append(indices, test.Index)
		}
	}
	return indices
}

// aggregateUsage adds usage to total, handling nil safely.
func aggregateUsage(total, usage *specview.TokenUsage) {
	if usage == nil {
		return
	}
	total.CandidatesTokens += usage.CandidatesTokens
	total.PromptTokens += usage.PromptTokens
	total.TotalTokens += usage.TotalTokens
}

// generateHeuristicTaxonomy creates a fallback taxonomy from file paths.
// Used when Stage 1 AI extraction fails.
func generateHeuristicTaxonomy(files []specview.FileInfo) *specview.TaxonomyOutput {
	// Group files by top-level directory
	dirMap := make(map[string][]int)
	for i, file := range files {
		dir := extractTopLevelDir(file.Path)
		dirMap[dir] = append(dirMap[dir], i)
	}

	// Sort directories for deterministic output
	dirs := make([]string, 0, len(dirMap))
	for dir := range dirMap {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	domains := make([]specview.TaxonomyDomain, 0, len(dirMap))
	for _, dir := range dirs {
		indices := dirMap[dir]
		name := humanizeDirName(dir)
		domains = append(domains, specview.TaxonomyDomain{
			Description: fmt.Sprintf("Tests from %s directory", dir),
			Features: []specview.TaxonomyFeature{
				{
					FileIndices: indices,
					Name:        uncategorizedFeatureName,
				},
			},
			Name: name,
		})
	}

	if len(domains) == 0 {
		allIndices := make([]int, len(files))
		for i := range files {
			allIndices[i] = i
		}
		domains = []specview.TaxonomyDomain{
			{
				Description: "All test files",
				Features: []specview.TaxonomyFeature{
					{
						FileIndices: allIndices,
						Name:        uncategorizedFeatureName,
					},
				},
				Name: uncategorizedDomainName,
			},
		}
	}

	return &specview.TaxonomyOutput{Domains: domains}
}

// extractTopLevelDir extracts the first directory component from a path.
func extractTopLevelDir(path string) string {
	// Skip leading slash
	start := 0
	if len(path) > 0 && path[0] == '/' {
		start = 1
	}

	for i := start; i < len(path); i++ {
		if path[i] == '/' {
			return path[start:i]
		}
	}

	// No slash found - file at root
	return "root"
}

// humanizeDirName converts directory names to human-readable format.
func humanizeDirName(dir string) string {
	nameMap := map[string]string{
		"auth":     "Authentication",
		"api":      "API",
		"db":       "Database",
		"payment":  "Payment",
		"user":     "User Management",
		"users":    "User Management",
		"admin":    "Administration",
		"config":   "Configuration",
		"core":     "Core",
		"lib":      "Library",
		"pkg":      "Package",
		"internal": "Internal",
		"src":      "Source",
		"test":     "Testing",
		"tests":    "Testing",
		"util":     "Utilities",
		"utils":    "Utilities",
		"root":     uncategorizedDomainName,
	}

	if name, ok := nameMap[dir]; ok {
		return name
	}

	// Capitalize first letter safely (handles unicode and non-lowercase)
	if len(dir) > 0 {
		runes := []rune(dir)
		runes[0] = unicode.ToUpper(runes[0])
		return string(runes)
	}
	return uncategorizedDomainName
}
