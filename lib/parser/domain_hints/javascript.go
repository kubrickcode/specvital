package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/tspool"
)

// JavaScriptExtractor extracts domain hints from JavaScript/TypeScript source code.
type JavaScriptExtractor struct {
	lang domain.Language
}

const (
	// ES6 imports: import x from 'y', import { x } from 'y', import 'y'
	jsImportQuery = `
		(import_statement
			source: (string) @import
		)
	`

	// CommonJS: require('x'), require("x")
	// Note: predicate #eq? is not guaranteed to work in all tree-sitter bindings,
	// so we also filter by function name in the extraction code.
	jsRequireQuery = `
		(call_expression
			function: (identifier) @func
			arguments: (arguments (string) @import)
		)
	`

	// Function calls: obj.method(), func()
	jsCallQuery = `
		(call_expression
			function: [
				(identifier) @call
				(member_expression) @call
			]
		)
	`
)

func (e *JavaScriptExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	tree, err := tspool.Parse(ctx, e.lang, source)
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

// extractImports extracts both ES6 and CommonJS imports.
// Uses best-effort extraction: query errors are ignored to allow partial results
// when one import style fails (e.g., ES6 query fails but CommonJS succeeds).
func (e *JavaScriptExtractor) extractImports(root *sitter.Node, source []byte) []string {
	seen := make(map[string]struct{})
	var imports []string

	// ES6 imports (errors ignored for best-effort extraction)
	es6Results, err := tspool.QueryWithCache(root, source, e.lang, jsImportQuery)
	if err == nil {
		for _, r := range es6Results {
			if node, ok := r.Captures["import"]; ok {
				path := trimJSQuotes(getNodeText(node, source))
				if ShouldFilterImportNoise(path) || isTypeOnlyImportNode(node) {
					continue
				}
				if _, exists := seen[path]; !exists {
					seen[path] = struct{}{}
					imports = append(imports, path)
				}
			}
		}
	}

	// CommonJS require (errors ignored for best-effort extraction)
	requireResults, err := tspool.QueryWithCache(root, source, e.lang, jsRequireQuery)
	if err == nil {
		for _, r := range requireResults {
			// Manually filter: only accept require() calls
			funcNode, hasFuncNode := r.Captures["func"]
			if !hasFuncNode || getNodeText(funcNode, source) != "require" {
				continue
			}
			if node, ok := r.Captures["import"]; ok {
				path := trimJSQuotes(getNodeText(node, source))
				if ShouldFilterImportNoise(path) {
					continue
				}
				if _, exists := seen[path]; !exists {
					seen[path] = struct{}{}
					imports = append(imports, path)
				}
			}
		}
	}

	return imports
}

func (e *JavaScriptExtractor) extractCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, e.lang, jsCallQuery)
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	calls := make([]string, 0, len(results))

	for _, r := range results {
		if node, ok := r.Captures["call"]; ok {
			call := getNodeText(node, source)
			if call == "" {
				continue
			}
			// Skip require() calls as they're already captured in imports
			if call == "require" {
				continue
			}
			// Normalize: remove whitespace and limit to 2 segments
			call = normalizeCall(call)
			if call == "" {
				continue
			}
			// Filter universal noise patterns (before more specific checks)
			if ShouldFilterNoise(call) {
				continue
			}
			// Skip test framework calls
			if isTestFrameworkCall(call) {
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

func trimJSQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	first := s[0]
	last := s[len(s)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') || (first == '`' && last == '`') {
		return s[1 : len(s)-1]
	}
	return s
}

func isTypeOnlyImportNode(node *sitter.Node) bool {
	parent := node.Parent()
	if parent == nil {
		return false
	}
	// Check if import_statement has "type" keyword
	if parent.Type() == "import_statement" {
		for i := 0; i < int(parent.ChildCount()); i++ {
			child := parent.Child(i)
			if child != nil && child.Type() == "type" {
				return true
			}
		}
	}
	return false
}

// testFrameworkCalls contains base function names from common test frameworks
// that should be excluded from domain hints. These calls don't provide
// meaningful domain classification signals. Method calls like test.describe()
// or expect().toBe() are filtered by checking the base name before the dot.
var testFrameworkCalls = map[string]struct{}{
	"describe": {}, "it": {}, "test": {}, "expect": {},
	"beforeEach": {}, "afterEach": {}, "beforeAll": {}, "afterAll": {},
	"vi": {}, "jest": {}, "cy": {},
	"fn": {}, // Standalone mock function (jest.fn(), vi.fn() result)
}

func isTestFrameworkCall(call string) bool {
	baseName := call
	if idx := strings.Index(call, "."); idx > 0 {
		baseName = call[:idx]
	}
	_, exists := testFrameworkCalls[baseName]
	return exists
}
