package analyzer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/detection"
	"github.com/specvital/core/pkg/parser/framework"
	"github.com/specvital/core/pkg/parser/framework/matchers"
	"github.com/specvital/web/src/backend/common/clients/github"
	"golang.org/x/sync/errgroup"

	// Framework registration via blank imports
	_ "github.com/specvital/core/pkg/parser/strategies/gotesting"
	_ "github.com/specvital/core/pkg/parser/strategies/jest"
	_ "github.com/specvital/core/pkg/parser/strategies/junit5"
	_ "github.com/specvital/core/pkg/parser/strategies/playwright"
	_ "github.com/specvital/core/pkg/parser/strategies/pytest"
	_ "github.com/specvital/core/pkg/parser/strategies/unittest"
	_ "github.com/specvital/core/pkg/parser/strategies/vitest"
)

// ErrRateLimitTooLow indicates GitHub API rate limit is insufficient for analysis.
var ErrRateLimitTooLow = errors.New("rate limit too low")

const (
	minRateLimitRemaining = 10
	maxConcurrentFetches  = 10
	maxFilesToProcess     = 200
)

var testFilePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\.test\.[jt]sx?$`),
	regexp.MustCompile(`\.spec\.[jt]sx?$`),
	regexp.MustCompile(`_test\.go$`),
	regexp.MustCompile(`(^|/)__tests__/.+\.[jt]sx?$`),
	regexp.MustCompile(`(^|/)test_[^/]+\.py$`),
	regexp.MustCompile(`_test\.py$`),
	regexp.MustCompile(`(^|/)conftest\.py$`),
	// Java/JUnit patterns
	regexp.MustCompile(`Test\.java$`),
	regexp.MustCompile(`Tests\.java$`),
	regexp.MustCompile(`(^|/)Test[A-Z][^/]*\.java$`),
}

var (
	configPatternsOnce sync.Once
	configPatterns     map[string]bool
)

func getConfigPatterns() map[string]bool {
	configPatternsOnce.Do(func() {
		configPatterns = make(map[string]bool)
		for _, def := range framework.DefaultRegistry().All() {
			for _, m := range def.Matchers {
				if cm, ok := m.(*matchers.ConfigMatcher); ok {
					for _, p := range cm.Patterns {
						configPatterns[filepath.Base(p)] = true
					}
				}
			}
		}
	})
	return configPatterns
}

type Service struct {
	gitHost GitHost
}

func NewService(gitHost GitHost) *Service {
	return &Service{
		gitHost: gitHost,
	}
}

func (s *Service) Analyze(ctx context.Context, owner, repo string) (*AnalysisResult, error) {
	rateLimit := s.gitHost.GetRateLimit()
	if rateLimit.Remaining > 0 && rateLimit.Remaining < minRateLimitRemaining {
		return nil, fmt.Errorf("%w: %d remaining", ErrRateLimitTooLow, rateLimit.Remaining)
	}

	commitInfo, err := s.gitHost.GetLatestCommit(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	files, err := s.gitHost.ListFiles(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	testFiles := filterTestFiles(files)
	if len(testFiles) == 0 {
		return &AnalysisResult{
			AnalyzedAt: time.Now().UTC().Format(time.RFC3339),
			CommitSHA:  commitInfo.SHA,
			Owner:      owner,
			Repo:       repo,
			Suites:     []TestSuite{},
			Summary:    Summary{Frameworks: []FrameworkSummary{}},
		}, nil
	}

	configFiles := filterConfigFiles(files)
	projectScope := s.buildProjectScope(ctx, owner, repo, configFiles)

	suites, err := s.parseTestFiles(ctx, owner, repo, testFiles, projectScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test files: %w", err)
	}

	summary := calculateSummary(suites)

	return &AnalysisResult{
		AnalyzedAt: time.Now().UTC().Format(time.RFC3339),
		CommitSHA:  commitInfo.SHA,
		Owner:      owner,
		Repo:       repo,
		Suites:     suites,
		Summary:    summary,
	}, nil
}

func filterTestFiles(files []github.FileInfo) []github.FileInfo {
	var testFiles []github.FileInfo
	for _, file := range files {
		if file.Type != "file" {
			continue
		}
		if isTestFile(file.Path) {
			testFiles = append(testFiles, file)
		}
	}
	return testFiles
}

func isTestFile(path string) bool {
	if shouldSkipPath(path) {
		return false
	}
	for _, pattern := range testFilePatterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

func shouldSkipPath(path string) bool {
	skipDirs := []string{"node_modules", "vendor", ".git", "dist", "build", "coverage", "fixtures"}
	for _, dir := range skipDirs {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") {
			return true
		}
	}
	return false
}

func filterConfigFiles(files []github.FileInfo) []github.FileInfo {
	patterns := getConfigPatterns()
	var configFiles []github.FileInfo
	for _, file := range files {
		if file.Type != "file" {
			continue
		}
		baseName := filepath.Base(file.Path)
		if patterns[baseName] {
			configFiles = append(configFiles, file)
		}
	}
	return configFiles
}

func matchesConfigPattern(def *framework.Definition, baseName string) bool {
	for _, m := range def.Matchers {
		if cm, ok := m.(*matchers.ConfigMatcher); ok {
			for _, p := range cm.Patterns {
				if filepath.Base(p) == baseName {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) buildProjectScope(ctx context.Context, owner, repo string, configFiles []github.FileInfo) *framework.AggregatedProjectScope {
	scope := framework.NewProjectScope()
	registry := framework.DefaultRegistry()

	loadedCount := 0
	for _, cf := range configFiles {
		content, err := s.gitHost.GetFileContent(ctx, owner, repo, cf.Path)
		if err != nil {
			slog.Warn("failed to get config file content", "path", cf.Path, "error", err)
			continue
		}

		baseName := filepath.Base(cf.Path)

		for _, def := range registry.All() {
			if def.ConfigParser == nil {
				continue
			}

			// Only try ConfigParser if file matches this framework's patterns
			if !matchesConfigPattern(def, baseName) {
				continue
			}

			configScope, err := def.ConfigParser.Parse(ctx, cf.Path, []byte(content))
			if err != nil {
				continue
			}
			if configScope != nil {
				scope.AddConfig(cf.Path, configScope)
				loadedCount++
				break
			}
		}
	}

	if len(configFiles) > 0 && loadedCount == 0 {
		slog.Warn("all config files failed to load", "total", len(configFiles))
	}

	return scope
}

func (s *Service) parseTestFiles(ctx context.Context, owner, repo string, files []github.FileInfo, projectScope *framework.AggregatedProjectScope) ([]TestSuite, error) {
	if len(files) > maxFilesToProcess {
		slog.Warn("truncating test files", "total", len(files), "limit", maxFilesToProcess)
		files = files[:maxFilesToProcess]
	}

	registry := framework.DefaultRegistry()
	detector := detection.NewDetector(registry)
	detector.SetProjectScope(projectScope)

	var mu sync.Mutex
	suites := []TestSuite{}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentFetches)

	for _, file := range files {
		g.Go(func() error {
			content, err := s.gitHost.GetFileContent(ctx, owner, repo, file.Path)
			if err != nil {
				slog.Warn("failed to get file content", "path", file.Path, "error", err)
				return nil
			}

			testFile, err := parseTestFile(ctx, registry, detector, []byte(content), file.Path)
			if err != nil {
				slog.Warn("failed to parse test file", "path", file.Path, "error", err)
				return nil
			}
			if testFile == nil {
				return nil
			}

			suite := convertToTestSuite(file.Path, testFile)
			if len(suite.Tests) > 0 {
				mu.Lock()
				suites = append(suites, suite)
				mu.Unlock()
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("failed to parse test files: %w", err)
	}

	slog.Info("parseTestFiles completed", "suiteCount", len(suites))
	return suites, nil
}

func parseTestFile(ctx context.Context, registry *framework.Registry, detector *detection.Detector, source []byte, filename string) (*domain.TestFile, error) {
	result := detector.Detect(ctx, filename, source)
	if result.Framework == "" {
		return nil, nil
	}

	def := registry.Find(result.Framework)
	if def == nil || def.Parser == nil {
		return nil, nil
	}

	return def.Parser.Parse(ctx, source, filename)
}

func convertToTestSuite(filePath string, testFile *domain.TestFile) TestSuite {
	framework := mapFramework(testFile.Framework)
	var tests []TestCase

	for _, test := range testFile.Tests {
		tests = append(tests, TestCase{
			FilePath:  filePath,
			Framework: framework,
			Line:      test.Location.StartLine,
			Name:      test.Name,
			Status:    mapStatus(test.Status),
		})
	}

	for _, suite := range testFile.Suites {
		tests = append(tests, extractTestsFromSuite(filePath, framework, suite, "")...)
	}

	return TestSuite{
		FilePath:  filePath,
		Framework: framework,
		Tests:     tests,
	}
}

func extractTestsFromSuite(filePath string, framework Framework, suite domain.TestSuite, prefix string) []TestCase {
	var tests []TestCase
	suiteName := suite.Name
	if prefix != "" {
		suiteName = prefix + " > " + suiteName
	}

	for _, test := range suite.Tests {
		name := test.Name
		if suiteName != "" {
			name = suiteName + " > " + name
		}
		tests = append(tests, TestCase{
			FilePath:  filePath,
			Framework: framework,
			Line:      test.Location.StartLine,
			Name:      name,
			Status:    mapStatus(test.Status),
		})
	}

	for _, nested := range suite.Suites {
		tests = append(tests, extractTestsFromSuite(filePath, framework, nested, suiteName)...)
	}

	return tests
}

// mapFramework normalizes core framework names for display.
// Only maps names that need normalization, others pass through as-is.
var frameworkDisplayNames = map[string]string{
	"go-testing": "go",
}

func mapFramework(name string) Framework {
	if display, ok := frameworkDisplayNames[name]; ok {
		return display
	}
	return name
}

func mapStatus(status domain.TestStatus) TestStatus {
	switch status {
	case domain.TestStatusSkipped:
		return TestStatusSkipped
	case domain.TestStatusFixme:
		return TestStatusTodo
	case domain.TestStatusPending:
		// Go testing returns "pending" for all tests (no skip/only modifiers)
		// Treat as active since Go tests don't have a "pending" concept
		return TestStatusActive
	case domain.TestStatusOnly, "":
		return TestStatusActive
	default:
		slog.Warn("unknown test status", "status", status)
		return TestStatusActive
	}
}

func (s *Service) CheckRateLimit() (int, error) {
	rateLimit := s.gitHost.GetRateLimit()
	if rateLimit.Remaining > 0 && rateLimit.Remaining < minRateLimitRemaining {
		return rateLimit.Remaining, fmt.Errorf("rate limit too low: %d remaining", rateLimit.Remaining)
	}
	return rateLimit.Remaining, nil
}

func (s *Service) GetRateLimit() github.RateLimitInfo {
	return s.gitHost.GetRateLimit()
}

func calculateSummary(suites []TestSuite) Summary {
	frameworkStats := make(map[Framework]*FrameworkSummary)
	var totalActive, totalSkipped, totalTodo int

	for _, suite := range suites {
		stats, exists := frameworkStats[suite.Framework]
		if !exists {
			stats = &FrameworkSummary{Framework: suite.Framework}
			frameworkStats[suite.Framework] = stats
		}

		for _, test := range suite.Tests {
			stats.Total++
			switch test.Status {
			case TestStatusActive:
				stats.Active++
				totalActive++
			case TestStatusSkipped:
				stats.Skipped++
				totalSkipped++
			case TestStatusTodo:
				stats.Todo++
				totalTodo++
			}
		}
	}

	frameworks := make([]FrameworkSummary, 0, len(frameworkStats))
	for _, stats := range frameworkStats {
		frameworks = append(frameworks, *stats)
	}

	return Summary{
		Active:     totalActive,
		Frameworks: frameworks,
		Skipped:    totalSkipped,
		Todo:       totalTodo,
		Total:      totalActive + totalSkipped + totalTodo,
	}
}
