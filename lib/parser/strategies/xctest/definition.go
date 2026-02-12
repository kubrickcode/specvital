// Package xctest implements XCTest framework support for Swift test files.
package xctest

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/lib/parser"
	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/framework"
	"github.com/kubrickcode/specvital/lib/parser/framework/matchers"
	"github.com/kubrickcode/specvital/lib/parser/strategies/shared/swiftast"
)

const frameworkName = framework.FrameworkXCTest

// Detection confidence scores (aligned with 4-stage detection system).
const (
	confidenceFileName = 20
	confidenceContent  = 40
)

func init() {
	framework.Register(NewDefinition())
}

func NewDefinition() *framework.Definition {
	return &framework.Definition{
		Name:      frameworkName,
		Languages: []domain.Language{domain.LanguageSwift},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher(
				"import XCTest",
				"@testable import",
			),
			&XCTestFileMatcher{},
			&XCTestContentMatcher{},
		},
		ConfigParser: nil,
		Parser:       &XCTestParser{},
		Priority:     framework.PriorityGeneric,
	}
}

// XCTestFileMatcher matches *Test.swift and *Tests.swift files.
type XCTestFileMatcher struct{}

func (m *XCTestFileMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileName {
		return framework.NoMatch()
	}

	if swiftast.IsSwiftTestFile(signal.Value) {
		return framework.PartialMatch(confidenceFileName, "XCTest file naming convention")
	}

	return framework.NoMatch()
}

// XCTestContentMatcher matches XCTest-specific patterns.
type XCTestContentMatcher struct{}

var xctestPatterns = []struct {
	pattern *regexp.Regexp
	desc    string
}{
	{regexp.MustCompile(`class\s+\w+\s*:\s*XCTestCase`), "extends XCTestCase"},
	{regexp.MustCompile(`\bfunc\s+test[A-Z]\w*\s*\(`), "test method"},
	{regexp.MustCompile(`\bXCTAssert`), "XCT assertion"},
	{regexp.MustCompile(`\bXCTFail\b`), "XCTFail"},
	{regexp.MustCompile(`\bXCTExpect`), "XCTExpect"},
	{regexp.MustCompile(`\bXCTSkip\b`), "XCTSkip"},
	{regexp.MustCompile(`import\s+XCTest`), "XCTest import"},
}

func (m *XCTestContentMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileContent {
		return framework.NoMatch()
	}

	content, ok := signal.Context.([]byte)
	if !ok {
		content = []byte(signal.Value)
	}

	for _, p := range xctestPatterns {
		if p.pattern.Match(content) {
			return framework.PartialMatch(confidenceContent, "Found XCTest pattern: "+p.desc)
		}
	}

	return framework.NoMatch()
}

// XCTestParser extracts test definitions from Swift XCTest files.
type XCTestParser struct{}

func (p *XCTestParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	tree, err := parser.ParseWithPool(ctx, domain.LanguageSwift, source)
	if err != nil {
		return nil, fmt.Errorf("xctest parser: failed to parse %s: %w", filename, err)
	}
	defer tree.Close()

	root := tree.RootNode()
	suites := parseTestClasses(root, source, filename)

	return &domain.TestFile{
		Path:      filename,
		Language:  domain.LanguageSwift,
		Framework: frameworkName,
		Suites:    suites,
	}, nil
}

func parseTestClasses(root *sitter.Node, source []byte, filename string) []domain.TestSuite {
	var suites []domain.TestSuite

	parser.WalkTree(root, func(node *sitter.Node) bool {
		if node.Type() == swiftast.NodeClassDeclaration {
			if suite := parseTestClass(node, source, filename); suite != nil {
				suites = append(suites, *suite)
			}
			return false
		}
		return true
	})

	return suites
}

func parseTestClass(node *sitter.Node, source []byte, filename string) *domain.TestSuite {
	className := swiftast.GetClassName(node, source)
	if className == "" {
		return nil
	}

	if !swiftast.IsXCTestCase(node, source) {
		return nil
	}

	suite := &domain.TestSuite{
		Name:     className,
		Status:   domain.TestStatusActive,
		Location: parser.GetLocation(node, filename),
	}

	body := swiftast.GetClassBody(node)
	if body == nil {
		return nil
	}

	parseTestMethods(body, source, filename, suite)

	if len(suite.Tests) == 0 {
		return nil
	}

	return suite
}

func parseTestMethods(body *sitter.Node, source []byte, filename string, suite *domain.TestSuite) {
	parser.WalkTree(body, func(node *sitter.Node) bool {
		if node.Type() == swiftast.NodeFunctionDeclaration {
			if test := parseTestMethod(node, source, filename); test != nil {
				suite.Tests = append(suite.Tests, *test)
			}
			return false
		}
		return true
	})
}

func parseTestMethod(node *sitter.Node, source []byte, filename string) *domain.Test {
	funcName := swiftast.GetFunctionName(node, source)
	if funcName == "" {
		return nil
	}

	if !swiftast.IsTestFunction(funcName) {
		return nil
	}

	status := domain.TestStatusActive
	modifier := ""

	// Check for throws XCTSkip pattern in function body
	if hasXCTSkip(node, source) {
		status = domain.TestStatusSkipped
		modifier = "XCTSkip"
	}

	// Check for async keyword
	if isAsyncFunction(node, source) {
		modifier = appendModifier(modifier, "async")
	}

	return &domain.Test{
		Name:     funcName,
		Status:   status,
		Modifier: modifier,
		Location: parser.GetLocation(node, filename),
	}
}

func hasXCTSkip(node *sitter.Node, source []byte) bool {
	content := node.Content(source)
	return strings.Contains(content, "XCTSkip") || strings.Contains(content, "throw XCTSkip")
}

func isAsyncFunction(node *sitter.Node, source []byte) bool {
	content := node.Content(source)
	return strings.Contains(content, "async ")
}

func appendModifier(existing, new string) string {
	if existing == "" {
		return new
	}
	return existing + ", " + new
}
