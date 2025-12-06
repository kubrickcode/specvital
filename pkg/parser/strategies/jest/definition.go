package jest

import (
	"context"
	"regexp"

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
		{regexp.MustCompile(`\bjest\.fn\s*\(`), "jest.fn()"},
		{regexp.MustCompile(`\bjest\.mock\s*\(`), "jest.mock()"},
		{regexp.MustCompile(`\bjest\.spyOn\s*\(`), "jest.spyOn()"},
		{regexp.MustCompile(`\bjest\.useFakeTimers\s*\(`), "jest.useFakeTimers()"},
		{regexp.MustCompile(`\bjest\.clearAllMocks\s*\(`), "jest.clearAllMocks()"},
		{regexp.MustCompile(`\bjest\.resetAllMocks\s*\(`), "jest.resetAllMocks()"},
		{regexp.MustCompile(`\bjest\.restoreAllMocks\s*\(`), "jest.restoreAllMocks()"},
		{regexp.MustCompile(`\bjest\.setTimeout\s*\(`), "jest.setTimeout()"},
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
	root := parseRoot(content)
	scope := framework.NewConfigScope(configPath, root)
	scope.Framework = frameworkName
	scope.GlobalsMode = !parseInjectGlobalsFalse(content) // Jest defaults to true
	return scope, nil
}

type JestParser struct{}

func (p *JestParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	return jstest.Parse(ctx, source, filename, frameworkName)
}

var (
	configRootPattern         = regexp.MustCompile(`rootDir\s*:\s*['"]([^'"]+)['"]`)
	injectGlobalsFalsePattern = regexp.MustCompile(`injectGlobals\s*:\s*false`)
)

func parseRoot(content []byte) string {
	if match := configRootPattern.FindSubmatch(content); match != nil {
		return string(match[1])
	}
	return ""
}

func parseInjectGlobalsFalse(content []byte) bool {
	return injectGlobalsFalsePattern.Match(content)
}
