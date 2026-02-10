package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/strategies/shared/javaast"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/tspool"
)

// JavaExtractor extracts domain hints from Java source code.
type JavaExtractor struct{}

const (
	// import x.y.z; import x.y.*; import static x.y.z;
	javaImportQuery = `
		(import_declaration) @import
	`

	// Function calls: obj.method(), ClassName.staticMethod()
	javaCallQuery = `
		(method_invocation) @call
	`
)

func (e *JavaExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	// Sanitize source to handle NULL bytes
	source = javaast.SanitizeSource(source)

	tree, err := tspool.Parse(ctx, domain.LanguageJava, source)
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

func (e *JavaExtractor) extractImports(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageJava, javaImportQuery)
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	var imports []string

	for _, r := range results {
		if node, ok := r.Captures["import"]; ok {
			importPath := extractJavaImportPath(node, source)
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

// extractJavaImportPath extracts the import path from an import_declaration node.
// Handles: import x.y.z; import x.y.*; import static x.y.z;
func extractJavaImportPath(node *sitter.Node, source []byte) string {
	var scopedIdentifier string
	hasAsterisk := false

	// Skip "import", "static", and ";" - extract the scoped_identifier
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		childType := child.Type()

		// scoped_identifier: x.y.z
		if childType == "scoped_identifier" {
			scopedIdentifier = child.Content(source)
		}
		// asterisk: for wildcard imports (x.y.*)
		if childType == "asterisk" {
			hasAsterisk = true
		}
		// identifier: single-word import (rare, but possible)
		if childType == "identifier" && scopedIdentifier == "" {
			scopedIdentifier = child.Content(source)
		}
	}

	if scopedIdentifier == "" {
		return ""
	}

	if hasAsterisk {
		return scopedIdentifier + ".*"
	}
	return scopedIdentifier
}

func (e *JavaExtractor) extractCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageJava, javaCallQuery)
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	calls := make([]string, 0, len(results))

	for _, r := range results {
		if node, ok := r.Captures["call"]; ok {
			call := extractJavaMethodCall(node, source)
			if call == "" {
				continue
			}
			// Normalize to 2 segments
			call = normalizeCall(call)
			if call == "" {
				continue
			}
			// Filter universal noise patterns (before more specific checks)
			if ShouldFilterNoise(call) {
				continue
			}
			// Skip test framework calls
			if isJavaTestFrameworkCall(call) {
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

// extractJavaMethodCall extracts the method call expression.
// Handles: obj.method(), ClassName.staticMethod(), method()
func extractJavaMethodCall(node *sitter.Node, source []byte) string {
	// method_invocation structure: object.name(arguments)
	// or: name(arguments)

	var parts []string

	// Extract object (if exists)
	objectNode := node.ChildByFieldName("object")
	if objectNode != nil {
		parts = append(parts, objectNode.Content(source))
	}

	// Extract method name
	nameNode := node.ChildByFieldName("name")
	if nameNode != nil {
		parts = append(parts, nameNode.Content(source))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ".")
}

// javaTestFrameworkCalls contains base names from Java test frameworks
// that should be excluded from domain hints.
var javaTestFrameworkCalls = map[string]struct{}{
	// JUnit assertions
	"assertEquals": {}, "assertNotEquals": {}, "assertTrue": {}, "assertFalse": {},
	"assertNull": {}, "assertNotNull": {}, "assertSame": {}, "assertNotSame": {},
	"assertArrayEquals": {}, "assertThrows": {}, "assertDoesNotThrow": {},
	"assertAll": {}, "assertTimeout": {}, "assertTimeoutPreemptively": {},
	"fail": {}, "assumeTrue": {}, "assumeFalse": {},
	// Assertions class prefix
	"Assertions": {},
	// Hamcrest
	"assertThat": {}, "is": {}, "equalTo": {}, "hasSize": {}, "contains": {},
	"containsString": {}, "startsWith": {}, "endsWith": {},
	"MatcherAssert": {},
	// Mockito
	"mock": {}, "spy": {}, "when": {}, "verify": {}, "doReturn": {},
	"doThrow": {}, "doNothing": {}, "times": {}, "never": {}, "any": {},
	"eq": {}, "anyString": {}, "anyInt": {}, "anyLong": {},
	"Mockito": {},
	// AssertJ
	"isEqualTo": {}, "isNotNull": {},
	// Java Object methods (inherited from java.lang.Object, no domain signal)
	"getClass": {}, "toString": {}, "hashCode": {}, "equals": {}, "clone": {},
	"getClass()": {}, "toString()": {}, "hashCode()": {},
}

func isJavaTestFrameworkCall(call string) bool {
	baseName := call
	if idx := strings.Index(call, "."); idx > 0 {
		baseName = call[:idx]
	}
	_, exists := javaTestFrameworkCalls[baseName]
	return exists
}
