package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/strategies/shared/kotlinast"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/tspool"
)

// KotlinExtractor extracts domain hints from Kotlin source code.
type KotlinExtractor struct{}

const (
	// import x.y.z, import x.y.*
	kotlinImportQuery = `(import_header) @import`

	// Function calls: obj.method(), function()
	kotlinCallQuery = `(call_expression) @call`
)

func (e *KotlinExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	// Sanitize source to handle NULL bytes
	source = kotlinast.SanitizeSource(source)

	tree, err := tspool.Parse(ctx, domain.LanguageKotlin, source)
	if err != nil {
		return nil
	}
	defer tree.Close()

	root := tree.RootNode()

	hints := &domain.DomainHints{
		Imports: e.extractImports(root, source),
		Calls:   e.extractCalls(root, source),
	}

	if len(hints.Imports) == 0 && len(hints.Calls) == 0 {
		return nil
	}

	return hints
}

func (e *KotlinExtractor) extractImports(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageKotlin, kotlinImportQuery)
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	var imports []string

	for _, r := range results {
		if node, ok := r.Captures["import"]; ok {
			importPath := extractKotlinImportPath(node, source)
			if importPath != "" {
				if _, exists := seen[importPath]; !exists {
					seen[importPath] = struct{}{}
					imports = append(imports, importPath)
				}
			}
		}
	}

	return imports
}

// extractKotlinImportPath extracts the import path from an import_header node.
// Handles: import x.y.z, import x.y.*
func extractKotlinImportPath(node *sitter.Node, source []byte) string {
	// Look for identifier chain
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()

		// identifier: the import path
		if childType == "identifier" {
			return child.Content(source)
		}
	}

	return ""
}

func (e *KotlinExtractor) extractCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageKotlin, kotlinCallQuery)
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	calls := make([]string, 0, len(results))

	for _, r := range results {
		if node, ok := r.Captures["call"]; ok {
			call := extractKotlinCallExpression(node, source)
			if call == "" {
				continue
			}
			// Normalize to 2 segments
			call = normalizeCall(call)
			if call == "" {
				continue
			}
			// Filter noise patterns
			if ShouldFilterNoise(call) {
				continue
			}
			// Skip test framework calls
			if isKotlinTestFrameworkCall(call) {
				continue
			}
			if _, exists := seen[call]; exists {
				continue
			}
			seen[call] = struct{}{}
			calls = append(calls, call)
		}
	}

	return calls
}

// extractKotlinCallExpression extracts the call expression.
// Handles: obj.method(), function(), obj.field.method()
func extractKotlinCallExpression(node *sitter.Node, source []byte) string {
	// call_expression structure: (navigation_expression | simple_identifier) call_suffix
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()

		// navigation_expression: obj.method
		if childType == "navigation_expression" {
			return child.Content(source)
		}
		// simple_identifier: function
		if childType == "simple_identifier" {
			return child.Content(source)
		}
	}

	return ""
}

// kotlinTestFrameworkCalls contains base names from Kotlin test frameworks
// and stdlib functions that should be excluded from domain hints.
var kotlinTestFrameworkCalls = map[string]struct{}{
	// Kotest matchers
	"shouldBe": {}, "shouldNotBe": {}, "shouldThrow": {}, "shouldNotThrow": {},
	"shouldBeNull": {}, "shouldNotBeNull": {}, "shouldContain": {},
	"shouldHaveSize": {}, "shouldBeEmpty": {}, "shouldNotBeEmpty": {},
	// Kotest spec DSL
	"describe": {}, "context": {}, "it": {}, "should": {}, "test": {},
	"feature": {}, "scenario": {}, "given": {}, "when": {}, "then": {},
	"expect": {}, "xdescribe": {}, "xit": {}, "xtest": {},
	// JUnit assertions (Kotlin style)
	"assertEquals": {}, "assertNotEquals": {}, "assertTrue": {}, "assertFalse": {},
	"assertNull": {}, "assertNotNull": {}, "assertThrows": {}, "assertDoesNotThrow": {},
	"Assertions": {},
	// Mockk
	"mockk": {}, "every": {}, "verify": {}, "slot": {}, "spyk": {},
	"confirmVerified": {}, "coEvery": {}, "coVerify": {},
	// Kotlin stdlib - collection factory functions (no domain signal)
	"listOf": {}, "mutableListOf": {}, "setOf": {}, "mutableSetOf": {},
	"mapOf": {}, "mutableMapOf": {}, "arrayOf": {}, "arrayOfNulls": {},
	// Kotlin stdlib - empty collections (no domain signal)
	"emptyList": {}, "emptySet": {}, "emptyMap": {}, "emptyArray": {},
	// Kotlin stdlib - utility functions (no domain signal)
	"Pair": {}, "Triple": {},
	// Kotlin stdlib - exception/validation functions (no domain signal)
	"error": {}, "require": {}, "requireNotNull": {}, "check": {}, "checkNotNull": {},
}

func isKotlinTestFrameworkCall(call string) bool {
	baseName := call
	if idx := strings.Index(call, "."); idx > 0 {
		baseName = call[:idx]
	}
	_, exists := kotlinTestFrameworkCalls[baseName]
	return exists
}
