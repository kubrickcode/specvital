package detection

import (
	"context"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/detection/config"
	"github.com/specvital/core/pkg/parser/detection/matchers"
)

var langExtMap = map[string]domain.Language{
	".cjs": domain.LanguageJavaScript,
	".cts": domain.LanguageTypeScript,
	".go":  domain.LanguageGo,
	".js":  domain.LanguageJavaScript,
	".jsx": domain.LanguageJavaScript,
	".mjs": domain.LanguageJavaScript,
	".mts": domain.LanguageTypeScript,
	".ts":  domain.LanguageTypeScript,
	".tsx": domain.LanguageTypeScript,
}

type Detector struct {
	matcherRegistry *matchers.Registry
	scopeResolver   *config.Resolver
}

func NewDetector(matcherRegistry *matchers.Registry, scopeResolver *config.Resolver) *Detector {
	return &Detector{
		matcherRegistry: matcherRegistry,
		scopeResolver:   scopeResolver,
	}
}

// Detect performs hierarchical framework detection (backward compatible).
func (d *Detector) Detect(ctx context.Context, filePath string, content []byte) Result {
	return d.DetectWithContext(ctx, filePath, content, nil)
}

// DetectWithContext performs hierarchical framework detection.
// Hierarchy: ProjectContext (globals mode) → Import → ScopeConfig → Unknown
func (d *Detector) DetectWithContext(ctx context.Context, filePath string, content []byte, projectCtx *ProjectContext) Result {
	if result, ok := d.detectFromProjectContext(filePath, projectCtx); ok {
		return result
	}

	if result, ok := d.detectFromImports(ctx, filePath, content); ok {
		return result
	}

	if result, ok := d.detectFromScopeConfig(filePath); ok {
		return result
	}

	return Unknown()
}

func (d *Detector) detectFromProjectContext(filePath string, projectCtx *ProjectContext) (Result, bool) {
	if projectCtx == nil {
		return Result{}, false
	}

	configInfo := projectCtx.FindApplicableConfig(filePath)
	if configInfo == nil || !configInfo.GlobalsMode {
		return Result{}, false
	}

	configPath := projectCtx.FindConfigPath(filePath, configInfo.Framework)
	return FromProjectContext(configInfo.Framework, configPath), true
}

func (d *Detector) detectFromImports(ctx context.Context, filePath string, content []byte) (Result, bool) {
	lang := detectLanguage(filePath)
	if lang == "" {
		return Result{}, false
	}

	compatibleMatchers := d.matcherRegistry.FindByLanguage(lang)
	if len(compatibleMatchers) == 0 {
		return Result{}, false
	}

	sorted := sortedByPriority(compatibleMatchers)

	imports := sorted[0].ExtractImports(ctx, content)
	if len(imports) == 0 {
		return Result{}, false
	}

	if matcher := findMatchingMatcher(sorted, imports); matcher != nil {
		return FromImport(matcher.Name()), true
	}
	return Result{}, false
}

func (d *Detector) detectFromScopeConfig(filePath string) (Result, bool) {
	if d.scopeResolver == nil {
		return Result{}, false
	}

	sorted := sortedByPriority(d.matcherRegistry.All())

	for _, matcher := range sorted {
		patterns := matcher.ConfigPatterns()
		if len(patterns) == 0 {
			continue
		}
		if configPath, found := d.scopeResolver.ResolveConfig(filePath, patterns); found {
			return FromScopeConfig(matcher.Name(), configPath), true
		}
	}
	return Result{}, false
}

func sortedByPriority(matcherList []matchers.Matcher) []matchers.Matcher {
	sorted := make([]matchers.Matcher, len(matcherList))
	copy(sorted, matcherList)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority() > sorted[j].Priority()
	})

	return sorted
}

func findMatchingMatcher(matcherList []matchers.Matcher, imports []string) matchers.Matcher {
	for _, matcher := range matcherList {
		if slices.ContainsFunc(imports, matcher.MatchImport) {
			return matcher
		}
	}
	return nil
}

func detectLanguage(filePath string) domain.Language {
	ext := strings.ToLower(filepath.Ext(filePath))
	return langExtMap[ext]
}
