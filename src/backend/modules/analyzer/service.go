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
	"github.com/specvital/core/pkg/parser/strategies"
	"github.com/specvital/core/pkg/parser/strategies/gotesting"
	"github.com/specvital/core/pkg/parser/strategies/jest"
	"github.com/specvital/core/pkg/parser/strategies/playwright"
	"github.com/specvital/core/pkg/parser/strategies/vitest"
	"github.com/specvital/web/src/backend/common/clients/github"
	"golang.org/x/sync/errgroup"
)

// ErrRateLimitTooLow indicates GitHub API rate limit is insufficient for analysis.
var ErrRateLimitTooLow = errors.New("rate limit too low")

var strategiesOnce sync.Once

// InitializeParserStrategies registers all test file parsing strategies.
// Must be called once during application startup.
func InitializeParserStrategies() {
	strategiesOnce.Do(func() {
		gotesting.RegisterDefault()
		jest.RegisterDefault()
		playwright.RegisterDefault()
		vitest.RegisterDefault()
	})
}

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
}

type Service struct {
	githubClient *github.Client
}

func NewService(githubClient *github.Client) *Service {
	return &Service{
		githubClient: githubClient,
	}
}

func (s *Service) Analyze(ctx context.Context, owner, repo string) (*AnalysisResult, error) {
	rateLimit := s.githubClient.GetRateLimit()
	if rateLimit.Remaining > 0 && rateLimit.Remaining < minRateLimitRemaining {
		return nil, fmt.Errorf("%w: %d remaining", ErrRateLimitTooLow, rateLimit.Remaining)
	}

	commitInfo, err := s.githubClient.GetLatestCommit(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	files, err := s.githubClient.ListFiles(ctx, owner, repo)
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

	suites, err := s.parseTestFiles(ctx, owner, repo, testFiles)
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

func (s *Service) parseTestFiles(ctx context.Context, owner, repo string, files []github.FileInfo) ([]TestSuite, error) {
	if len(files) > maxFilesToProcess {
		slog.Warn("truncating test files", "total", len(files), "limit", maxFilesToProcess)
		files = files[:maxFilesToProcess]
	}

	registry := strategies.DefaultRegistry()

	var mu sync.Mutex
	suites := []TestSuite{}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentFetches)

	for _, file := range files {
		file := file
		g.Go(func() error {
			content, err := s.githubClient.GetFileContent(ctx, owner, repo, file.Path)
			if err != nil {
				slog.Warn("failed to get file content", "path", file.Path, "error", err)
				return nil
			}

			testFile, err := parseWithStrategies(ctx, registry, []byte(content), file.Path)
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

func parseWithStrategies(ctx context.Context, registry *strategies.Registry, source []byte, filename string) (*domain.TestFile, error) {
	strat := registry.FindStrategy(filename, source)
	if strat == nil {
		return nil, nil
	}
	return strat.Parse(ctx, source, filename)
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

func mapFramework(framework string) Framework {
	switch strings.ToLower(framework) {
	case "jest":
		return FrameworkJest
	case "vitest":
		return FrameworkVitest
	case "playwright":
		return FrameworkPlaywright
	case "gotesting", "go", "go-testing":
		return FrameworkGo
	default:
		ext := filepath.Ext(framework)
		if ext == ".go" {
			return FrameworkGo
		}
		return FrameworkJest
	}
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
	rateLimit := s.githubClient.GetRateLimit()
	if rateLimit.Remaining > 0 && rateLimit.Remaining < minRateLimitRemaining {
		return rateLimit.Remaining, fmt.Errorf("rate limit too low: %d remaining", rateLimit.Remaining)
	}
	return rateLimit.Remaining, nil
}

func (s *Service) GetRateLimit() github.RateLimitInfo {
	return s.githubClient.GetRateLimit()
}
