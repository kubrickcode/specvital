package vitest

import (
	"context"
	"regexp"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/detection/extraction"
	"github.com/specvital/core/pkg/parser/framework"
	"github.com/specvital/core/pkg/parser/framework/matchers"
	"github.com/specvital/core/pkg/parser/strategies/shared/jstest"
)

const frameworkName = "vitest"

func init() {
	framework.Register(NewDefinition())
}

func NewDefinition() *framework.Definition {
	return &framework.Definition{
		Name:      frameworkName,
		Languages: []domain.Language{domain.LanguageTypeScript, domain.LanguageJavaScript},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher("vitest", "vitest/"),
			matchers.NewConfigMatcher(
				"vitest.config.js",
				"vitest.config.ts",
				"vitest.config.mjs",
				"vitest.config.mts",
			),
		},
		ConfigParser: &VitestConfigParser{},
		Parser:       &VitestParser{},
		Priority:     framework.PrioritySpecialized,
	}
}

type VitestConfigParser struct{}

func (p *VitestConfigParser) Parse(ctx context.Context, configPath string, content []byte) (*framework.ConfigScope, error) {
	root := parseRoot(content)
	scope := framework.NewConfigScope(configPath, root)
	scope.Framework = frameworkName
	scope.GlobalsMode = parseGlobals(ctx, content)
	return scope, nil
}

type VitestParser struct{}

func (p *VitestParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	return jstest.Parse(ctx, source, filename, frameworkName)
}

var (
	configRootPattern    = regexp.MustCompile(`root\s*:\s*['"]([^'"]+)['"]`)
	configGlobalsPattern = regexp.MustCompile(`globals\s*:\s*true`)
)

func parseRoot(content []byte) string {
	if match := configRootPattern.FindSubmatch(content); match != nil {
		return string(match[1])
	}
	return ""
}

func parseGlobals(ctx context.Context, content []byte) bool {
	return extraction.MatchPatternExcludingComments(ctx, content, configGlobalsPattern)
}
