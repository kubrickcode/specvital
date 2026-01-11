package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/tspool"
)

// PythonExtractor extracts domain hints from Python source code.
type PythonExtractor struct{}

const (
	// import x, import x.y.z
	pyImportQuery = `(import_statement (dotted_name) @import)`

	// from x import y, from x.y import z
	pyFromImportQuery = `
		(import_from_statement
			module_name: (dotted_name) @module
		)
		(import_from_statement
			module_name: (relative_import) @module
		)
	`

	// Function calls: func(), obj.method()
	pyCallQuery = `
		(call
			function: [
				(identifier) @call
				(attribute) @call
			]
		)
	`
)

func (e *PythonExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	tree, err := tspool.Parse(ctx, domain.LanguagePython, source)
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

func (e *PythonExtractor) extractImports(root *sitter.Node, source []byte) []string {
	seen := make(map[string]struct{})
	var imports []string

	// import x, import x.y.z
	importResults, err := tspool.QueryWithCache(root, source, domain.LanguagePython, pyImportQuery)
	if err == nil {
		for _, r := range importResults {
			if node, ok := r.Captures["import"]; ok {
				modPath := getNodeText(node, source)
				if modPath != "" {
					if _, exists := seen[modPath]; !exists {
						seen[modPath] = struct{}{}
						imports = append(imports, modPath)
					}
				}
			}
		}
	}

	// from x import y
	fromResults, err := tspool.QueryWithCache(root, source, domain.LanguagePython, pyFromImportQuery)
	if err == nil {
		for _, r := range fromResults {
			if node, ok := r.Captures["module"]; ok {
				modPath := getNodeText(node, source)
				modPath = strings.TrimSpace(modPath)
				// Skip bare relative markers (., ..) as they provide no domain classification value.
				// Keep relative imports with module paths (.models, ..services).
				if modPath != "" && modPath != "." && modPath != ".." {
					if _, exists := seen[modPath]; !exists {
						seen[modPath] = struct{}{}
						imports = append(imports, modPath)
					}
				}
			}
		}
	}

	return imports
}

func (e *PythonExtractor) extractCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguagePython, pyCallQuery)
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
			// Normalize: remove whitespace and limit to 2 segments
			call = normalizeCall(call)
			if call == "" {
				continue
			}
			// Skip test framework calls
			if isPythonTestFrameworkCall(call) {
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

// pythonTestFrameworkCalls contains base function names from Python test frameworks
// that should be excluded from domain hints.
var pythonTestFrameworkCalls = map[string]struct{}{
	// pytest markers and decorators
	"pytest": {}, "test": {}, "fixture": {}, "mark": {}, "parametrize": {},
	"skip": {}, "skipif": {}, "xfail": {},
	// pytest setup/teardown
	"setup": {}, "teardown": {},
	"setup_method": {}, "teardown_method": {},
	"setup_class": {}, "teardown_class": {},
	"setup_module": {}, "teardown_module": {},
	// pytest built-in fixtures
	"raises": {}, "monkeypatch": {}, "caplog": {}, "capsys": {}, "tmpdir": {},
	"request": {}, "pytestconfig": {}, "tmp_path": {},
	// unittest
	"unittest": {}, "setUp": {}, "tearDown": {},
	"setUpClass": {}, "tearDownClass": {},
	"setUpModule": {}, "tearDownModule": {},
	// mock
	"mock": {}, "patch": {}, "Mock": {}, "MagicMock": {},
	// self (unittest method calls like self.assertEqual)
	"self": {},
}

func isPythonTestFrameworkCall(call string) bool {
	baseName := call
	if idx := strings.Index(call, "."); idx > 0 {
		baseName = call[:idx]
	}
	_, exists := pythonTestFrameworkCalls[baseName]
	return exists
}
