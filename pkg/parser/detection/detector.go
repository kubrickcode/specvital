package detection

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/detection/extraction"
	"github.com/specvital/core/pkg/parser/framework"
)

// Detector performs framework detection using early-return approach.
// Detection priority (highest to lowest):
// 1. Import statements - explicit developer intent (immediate return)
// 2. Config scope - project-level configuration
// 3. Content patterns - framework-specific code patterns
//
// The first successful match at any level immediately returns.
type Detector struct {
	registry     *framework.Registry
	projectScope *framework.AggregatedProjectScope
}

// NewDetector creates a new detector.
func NewDetector(registry *framework.Registry) *Detector {
	return &Detector{
		registry:     registry,
		projectScope: nil,
	}
}

// SetProjectScope configures the detector with project-wide configuration context.
func (d *Detector) SetProjectScope(scope *framework.AggregatedProjectScope) {
	d.projectScope = scope
}

// Detect performs framework detection on a test file.
// Uses early-return: first match wins based on priority.
func (d *Detector) Detect(ctx context.Context, filePath string, content []byte) Result {
	lang := detectLanguage(filePath)
	if lang == "" {
		return Unknown()
	}

	// Go test files are detected by naming convention (*_test.go)
	if lang == domain.LanguageGo {
		if strings.HasSuffix(filepath.Base(filePath), "_test.go") {
			return Confirmed("go-testing", SourceContentPattern)
		}
		return Unknown()
	}

	frameworks := d.registry.FindByLanguage(lang)
	if len(frameworks) == 0 {
		return Unknown()
	}

	if fw := d.detectFromImport(ctx, lang, content, frameworks); fw != "" {
		return Confirmed(fw, SourceImport)
	}

	if result := d.detectFromScope(filePath, lang); result.Framework != "" {
		return result
	}

	if fw := d.detectFromContent(ctx, content, frameworks); fw != "" {
		return Confirmed(fw, SourceContentPattern)
	}

	return Unknown()
}

// detectFromImport checks for framework-specific import statements.
// Returns framework name if found, empty string otherwise.
func (d *Detector) detectFromImport(ctx context.Context, lang domain.Language, content []byte, frameworks []*framework.Definition) string {
	var imports []string

	switch lang {
	case domain.LanguageTypeScript, domain.LanguageJavaScript:
		imports = extraction.ExtractJSImports(ctx, content)
	case domain.LanguageGo:
		// Go uses testing package directly, not detected via imports
		return ""
	case domain.LanguageJava:
		imports = extraction.ExtractJavaImports(ctx, content)
	case domain.LanguagePython:
		imports = extraction.ExtractPythonImports(ctx, content)
	case domain.LanguageCSharp:
		imports = extraction.ExtractCSharpUsings(ctx, content)
	}

	if len(imports) == 0 {
		return ""
	}

	for _, fw := range frameworks {
		for _, matcher := range fw.Matchers {
			for _, imp := range imports {
				signal := framework.Signal{
					Type:  framework.SignalImport,
					Value: imp,
				}

				mr := matcher.Match(ctx, signal)
				if mr.Confidence > 0 && !mr.Negative {
					return fw.Name
				}
			}
		}
	}

	return ""
}

// detectFromScope checks if file is within a config scope.
// Returns Result with framework and scope if found.
func (d *Detector) detectFromScope(filePath string, lang domain.Language) Result {
	if d.projectScope == nil {
		return Unknown()
	}

	type scopeMatch struct {
		path  string
		scope *framework.ConfigScope
		depth int
	}

	var matches []scopeMatch
	for path, scope := range d.projectScope.Configs {
		def := d.registry.Find(scope.Framework)
		if def == nil {
			continue
		}

		// Check language compatibility
		langCompatible := false
		for _, l := range def.Languages {
			if l == lang {
				langCompatible = true
				break
			}
		}
		if !langCompatible {
			continue
		}

		if scope.Contains(filePath) {
			matches = append(matches, scopeMatch{
				path:  path,
				scope: scope,
				depth: scope.Depth(),
			})
		}
	}

	if len(matches) == 0 {
		return Unknown()
	}

	best := matches[0]
	for _, m := range matches[1:] {
		if m.depth > best.depth {
			best = m
		}
	}

	return ConfirmedWithScope(best.scope.Framework, best.scope)
}

// detectFromContent checks for framework-specific content patterns.
// Returns framework name if found, empty string otherwise.
func (d *Detector) detectFromContent(ctx context.Context, content []byte, frameworks []*framework.Definition) string {
	for _, fw := range frameworks {
		for _, matcher := range fw.Matchers {
			signal := framework.Signal{
				Type:    framework.SignalFileContent,
				Value:   "",
				Context: content,
			}

			mr := matcher.Match(ctx, signal)
			if mr.Confidence > 0 && !mr.Negative {
				return fw.Name
			}
		}
	}

	return ""
}

func detectLanguage(filePath string) domain.Language {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".ts", ".tsx":
		return domain.LanguageTypeScript
	case ".js", ".jsx", ".mjs", ".cjs":
		return domain.LanguageJavaScript
	case ".go":
		return domain.LanguageGo
	case ".java":
		return domain.LanguageJava
	case ".py":
		return domain.LanguagePython
	case ".cs":
		return domain.LanguageCSharp
	default:
		return ""
	}
}
