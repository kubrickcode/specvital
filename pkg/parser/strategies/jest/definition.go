package jest

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/framework"
	"github.com/specvital/core/pkg/parser/framework/matchers"
	"github.com/specvital/core/pkg/parser/strategies/shared/jstest"
)

func init() {
	framework.Register(NewDefinition())
}

func NewDefinition() *framework.Definition {
	return &framework.Definition{
		Name:      frameworkName,
		Languages: []domain.Language{domain.LanguageTypeScript, domain.LanguageJavaScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("@jest/globals", "@jest/", "jest"),
			matchers.NewConfigMatcher(
				"jest.config.js",
				"jest.config.ts",
				"jest.config.mjs",
				"jest.config.cjs",
				"jest.config.json",
			),
			&JestContentMatcher{},
		},
		ConfigParser: &JestConfigParser{},
		Parser:       &JestParser{},
		Priority:     framework.PriorityGeneric,
	}
}

// JestContentMatcher matches jest-specific patterns (jest.fn, jest.mock, etc.).
type JestContentMatcher struct{}

func (m *JestContentMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileContent {
		return framework.NoMatch()
	}

	content, ok := signal.Context.([]byte)
	if !ok {
		content = []byte(signal.Value)
	}

	jestPatterns := []struct {
		pattern *regexp.Regexp
		desc    string
	}{
		{regexp.MustCompile(`\bjest\.advanceTimersByTime\s*\(`), "jest.advanceTimersByTime()"},
		{regexp.MustCompile(`\bjest\.clearAllMocks\s*\(`), "jest.clearAllMocks()"},
		{regexp.MustCompile(`\bjest\.fn\s*\(`), "jest.fn()"},
		{regexp.MustCompile(`\bjest\.isolateModules\s*\(`), "jest.isolateModules()"},
		{regexp.MustCompile(`\bjest\.mock\s*\(`), "jest.mock()"},
		{regexp.MustCompile(`\bjest\.resetAllMocks\s*\(`), "jest.resetAllMocks()"},
		{regexp.MustCompile(`\bjest\.resetModules\s*\(`), "jest.resetModules()"},
		{regexp.MustCompile(`\bjest\.restoreAllMocks\s*\(`), "jest.restoreAllMocks()"},
		{regexp.MustCompile(`\bjest\.runAllTimers\s*\(`), "jest.runAllTimers()"},
		{regexp.MustCompile(`\bjest\.runOnlyPendingTimers\s*\(`), "jest.runOnlyPendingTimers()"},
		{regexp.MustCompile(`\bjest\.setTimeout\s*\(`), "jest.setTimeout()"},
		{regexp.MustCompile(`\bjest\.spyOn\s*\(`), "jest.spyOn()"},
		{regexp.MustCompile(`\bjest\.useFakeTimers\s*\(`), "jest.useFakeTimers()"},
		{regexp.MustCompile(`\bjest\.useRealTimers\s*\(`), "jest.useRealTimers()"},
	}

	var evidence []string
	for _, p := range jestPatterns {
		if p.pattern.Match(content) {
			evidence = append(evidence, "Found Jest-specific pattern: "+p.desc)
			return framework.PartialMatch(40, evidence...)
		}
	}

	return framework.NoMatch()
}

type JestConfigParser struct{}

func (p *JestConfigParser) Parse(ctx context.Context, configPath string, content []byte) (*framework.ConfigScope, error) {
	rootDir := parseRootDir(content)
	scope := framework.NewConfigScope(configPath, rootDir)
	scope.Framework = frameworkName
	scope.GlobalsMode = !parseInjectGlobalsFalse(content) // Jest defaults to true

	configDir := filepath.Dir(configPath)
	roots := parseRoots(content, configDir, rootDir)
	if len(roots) > 0 {
		scope.Roots = roots
	}

	return scope, nil
}

type JestParser struct{}

func (p *JestParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	return jstest.Parse(ctx, source, filename, frameworkName)
}

var (
	configRootDirPattern      = regexp.MustCompile(`rootDir\s*:\s*['"]([^'"]+)['"]`)
	configRootsPattern        = regexp.MustCompile(`roots\s*:\s*\[([^\]]+)\]`)
	configRootItemPattern     = regexp.MustCompile(`['"]([^'"]+)['"]`)
	injectGlobalsFalsePattern = regexp.MustCompile(`injectGlobals\s*:\s*false`)
)

func parseRootDir(content []byte) string {
	if match := configRootDirPattern.FindSubmatch(content); match != nil {
		return string(match[1])
	}
	return ""
}

func parseRoots(content []byte, configDir string, rootDir string) []string {
	match := configRootsPattern.FindSubmatch(content)
	if match == nil {
		return nil
	}

	rootsContent := match[1]
	items := configRootItemPattern.FindAllSubmatch(rootsContent, -1)
	if len(items) == 0 {
		return nil
	}

	resolvedRootDir := configDir
	if rootDir != "" {
		resolvedRootDir = filepath.Clean(filepath.Join(configDir, rootDir))
	}

	var roots []string
	for _, item := range items {
		root := string(item[1])
		hadRootDirPlaceholder := strings.Contains(root, "<rootDir>")
		root = strings.ReplaceAll(root, "<rootDir>", resolvedRootDir)

		if !filepath.IsAbs(root) && !hadRootDirPlaceholder {
			root = filepath.Join(configDir, root)
		}

		root = filepath.Clean(root)
		roots = append(roots, root)
	}

	return roots
}

func parseInjectGlobalsFalse(content []byte) bool {
	return injectGlobalsFalsePattern.Match(content)
}
