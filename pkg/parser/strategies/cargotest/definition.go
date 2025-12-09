// Package cargotest implements Rust cargo test framework support.
package cargotest

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser"
	"github.com/specvital/core/pkg/parser/framework"
	"github.com/specvital/core/pkg/parser/framework/matchers"
)

const frameworkName = "cargo-test"

// Tree-sitter node types for Rust
const (
	nodeAttributeItem = "attribute_item"
	nodeAttribute     = "attribute"
	nodeFunctionItem  = "function_item"
	nodeModItem       = "mod_item"
	nodeIdentifier    = "identifier"
	nodeMetaItem      = "meta_item"
)

func init() {
	framework.Register(NewDefinition())
}

func NewDefinition() *framework.Definition {
	return &framework.Definition{
		Name:      frameworkName,
		Languages: []domain.Language{domain.LanguageRust},
		Matchers: []framework.Matcher{
			&CargoTestFileMatcher{},
			matchers.NewConfigMatcher("Cargo.toml"),
			&CargoTestContentMatcher{},
		},
		ConfigParser: nil,
		Parser:       &CargoTestParser{},
		Priority:     framework.PriorityGeneric,
	}
}

// CargoTestFileMatcher matches Rust test files.
type CargoTestFileMatcher struct{}

func (m *CargoTestFileMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileName {
		return framework.NoMatch()
	}

	filename := signal.Value

	if strings.HasSuffix(filename, "_test.rs") {
		return framework.PartialMatch(20, "Rust test file naming: *_test.rs")
	}

	if strings.HasSuffix(filename, ".rs") {
		if strings.Contains(filename, "/tests/") || strings.HasPrefix(filename, "tests/") {
			return framework.PartialMatch(20, "Rust test directory: tests/*.rs")
		}
	}

	return framework.NoMatch()
}

// CargoTestContentMatcher matches #[test] and #[cfg(test)] patterns.
type CargoTestContentMatcher struct{}

var cargoTestPatterns = []struct {
	pattern *regexp.Regexp
	desc    string
}{
	{regexp.MustCompile(`#\[test\]`), "#[test] attribute"},
	{regexp.MustCompile(`#\[cfg\(test\)\]`), "#[cfg(test)] attribute"},
	{regexp.MustCompile(`#\[ignore\]`), "#[ignore] attribute"},
	{regexp.MustCompile(`#\[should_panic`), "#[should_panic] attribute"},
}

func (m *CargoTestContentMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileContent {
		return framework.NoMatch()
	}

	content, ok := signal.Context.([]byte)
	if !ok {
		content = []byte(signal.Value)
	}

	for _, p := range cargoTestPatterns {
		if p.pattern.Match(content) {
			return framework.PartialMatch(40, "Found Rust pattern: "+p.desc)
		}
	}

	return framework.NoMatch()
}

// CargoTestParser extracts test definitions from Rust source files.
type CargoTestParser struct{}

func (p *CargoTestParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	tree, err := parser.ParseWithPool(ctx, domain.LanguageRust, source)
	if err != nil {
		return nil, fmt.Errorf("cargo-test parser: failed to parse %s: %w", filename, err)
	}
	defer tree.Close()

	root := tree.RootNode()
	file := &domain.TestFile{
		Path:      filename,
		Language:  domain.LanguageRust,
		Framework: frameworkName,
	}

	// Use WalkTree for depth-protected traversal (prevents stack overflow)
	parseRustAST(root, source, filename, file)
	return file, nil
}

// parseRustAST traverses the AST using depth-protected WalkTree.
// It handles test modules and test functions at the top level and within #[cfg(test)] modules.
func parseRustAST(root *sitter.Node, source []byte, filename string, file *domain.TestFile) {
	// Track test modules by node start byte position to associate tests with their parent suite
	testModules := make(map[uint32]*domain.TestSuite)

	parser.WalkTree(root, func(node *sitter.Node) bool {
		switch node.Type() {
		case nodeModItem:
			if isTestModule(node, source) {
				name := extractModuleName(node, source)
				if name != "" {
					suite := &domain.TestSuite{
						Name:     name,
						Status:   domain.TestStatusActive,
						Location: parser.GetLocation(node, filename),
					}
					// Store suite keyed by start byte position
					testModules[node.StartByte()] = suite
				}
			}
			return true // Continue into children

		case nodeFunctionItem:
			attrs := collectAttributes(node, source)
			if !attrs.isTest {
				return false // Skip non-test functions
			}

			name := extractFunctionName(node, source)
			if name == "" {
				return false
			}

			test := buildTest(name, attrs, node, filename)

			// Find parent test module, if any
			parentSuite := findParentTestSuite(node, testModules)
			if parentSuite != nil {
				parentSuite.Tests = append(parentSuite.Tests, test)
			} else {
				file.Tests = append(file.Tests, test)
			}
			return false // No need to traverse into function body
		}

		return true // Continue traversal for other node types
	})

	// Add non-empty test suites to file
	for _, suite := range testModules {
		if len(suite.Tests) > 0 || len(suite.Suites) > 0 {
			file.Suites = append(file.Suites, *suite)
		}
	}
}

// findParentTestSuite finds the nearest ancestor test module for a node.
func findParentTestSuite(node *sitter.Node, testModules map[uint32]*domain.TestSuite) *domain.TestSuite {
	current := node.Parent()
	for current != nil {
		if suite, ok := testModules[current.StartByte()]; ok {
			return suite
		}
		current = current.Parent()
	}
	return nil
}

// buildTest creates a Test from function attributes.
func buildTest(name string, attrs testAttributes, node *sitter.Node, filename string) domain.Test {
	status := domain.TestStatusActive
	modifier := ""

	if attrs.isIgnore {
		status = domain.TestStatusSkipped
		modifier = "#[ignore]"
	}

	if attrs.shouldPanic != "" {
		if modifier != "" {
			modifier += " " + attrs.shouldPanic
		} else {
			modifier = attrs.shouldPanic
		}
	}

	return domain.Test{
		Name:     name,
		Status:   status,
		Modifier: modifier,
		Location: parser.GetLocation(node, filename),
	}
}

type testAttributes struct {
	isTest       bool
	isIgnore     bool
	shouldPanic  string // Full attribute text (e.g., "#[should_panic(expected = \"...\")]")
}

// getPrecedingAttributes returns attribute_item nodes immediately preceding the given node.
func getPrecedingAttributes(node *sitter.Node) []*sitter.Node {
	parent := node.Parent()
	if parent == nil {
		return nil
	}

	nodeIndex := -1
	for i := 0; i < int(parent.ChildCount()); i++ {
		if parent.Child(i) == node {
			nodeIndex = i
			break
		}
	}

	if nodeIndex == -1 {
		return nil
	}

	var attrs []*sitter.Node
	for i := nodeIndex - 1; i >= 0; i-- {
		child := parent.Child(i)
		if child.Type() != nodeAttributeItem {
			break
		}
		attrs = append(attrs, child)
	}

	return attrs
}

func collectAttributes(funcNode *sitter.Node, source []byte) testAttributes {
	attrs := testAttributes{}

	for _, attrNode := range getPrecedingAttributes(funcNode) {
		attrName := extractAttributeName(attrNode, source)
		switch attrName {
		case "test":
			attrs.isTest = true
		case "ignore":
			attrs.isIgnore = true
		case "should_panic":
			attrs.shouldPanic = parser.GetNodeText(attrNode, source)
		}
	}

	return attrs
}

func extractAttributeName(attrItem *sitter.Node, source []byte) string {
	attr := parser.FindChildByType(attrItem, nodeAttribute)
	if attr == nil {
		return ""
	}

	ident := parser.FindChildByType(attr, nodeIdentifier)
	if ident != nil {
		return parser.GetNodeText(ident, source)
	}

	// Handle complex attributes like #[cfg(test)] where identifier is nested in meta_item
	meta := parser.FindChildByType(attr, nodeMetaItem)
	if meta != nil {
		ident = parser.FindChildByType(meta, nodeIdentifier)
		if ident != nil {
			return parser.GetNodeText(ident, source)
		}
	}

	return ""
}

func extractFunctionName(funcNode *sitter.Node, source []byte) string {
	name := funcNode.ChildByFieldName("name")
	if name == nil {
		return ""
	}
	return parser.GetNodeText(name, source)
}

func extractModuleName(modNode *sitter.Node, source []byte) string {
	name := modNode.ChildByFieldName("name")
	if name == nil {
		return ""
	}
	return parser.GetNodeText(name, source)
}

func isTestModule(modNode *sitter.Node, source []byte) bool {
	for _, attrNode := range getPrecedingAttributes(modNode) {
		attrName := extractAttributeName(attrNode, source)
		if attrName == "cfg" {
			attrText := parser.GetNodeText(attrNode, source)
			if strings.Contains(attrText, "cfg(test)") {
				return true
			}
		}
	}

	name := extractModuleName(modNode, source)
	return name == "tests"
}
