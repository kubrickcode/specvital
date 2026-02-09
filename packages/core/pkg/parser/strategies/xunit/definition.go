// Package xunit implements xUnit test framework support for C# test files.
package xunit

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
	"github.com/specvital/core/pkg/parser/strategies/shared/dotnetast"
)

const frameworkName = "xunit"

func init() {
	framework.Register(NewDefinition())
}

func NewDefinition() *framework.Definition {
	return &framework.Definition{
		Name:      frameworkName,
		Languages: []domain.Language{domain.LanguageCSharp},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher(
				"Xunit",
				"using Xunit",
			),
			&XUnitFileMatcher{},
			&XUnitContentMatcher{},
		},
		ConfigParser: nil,
		Parser:       &XUnitParser{},
		Priority:     framework.PriorityGeneric,
	}
}

// XUnitFileMatcher matches *Test.cs, *Tests.cs, Test*.cs, *Spec.cs, *Specs.cs files.
type XUnitFileMatcher struct{}

func (m *XUnitFileMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileName {
		return framework.NoMatch()
	}

	if dotnetast.IsCSharpTestFileName(signal.Value) {
		return framework.PartialMatch(20, "xUnit file naming convention")
	}

	return framework.NoMatch()
}

// XUnitContentMatcher matches xUnit specific patterns.
type XUnitContentMatcher struct{}

var xunitPatterns = []struct {
	pattern *regexp.Regexp
	desc    string
}{
	{regexp.MustCompile(`\[Fact\]`), "[Fact] attribute"},
	{regexp.MustCompile(`\[Theory\]`), "[Theory] attribute"},
	{regexp.MustCompile(`\[InlineData\(`), "[InlineData] attribute"},
	{regexp.MustCompile(`\[MemberData\(`), "[MemberData] attribute"},
	{regexp.MustCompile(`\[ClassData\(`), "[ClassData] attribute"},
	{regexp.MustCompile(`using\s+Xunit\s*;`), "using Xunit"},
	{regexp.MustCompile(`\[Skip\(`), "[Skip] attribute"},
}

func (m *XUnitContentMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileContent {
		return framework.NoMatch()
	}

	content, ok := signal.Context.([]byte)
	if !ok {
		content = []byte(signal.Value)
	}

	for _, p := range xunitPatterns {
		if p.pattern.Match(content) {
			return framework.PartialMatch(40, "Found xUnit pattern: "+p.desc)
		}
	}

	return framework.NoMatch()
}

// XUnitParser extracts test definitions from C# xUnit files.
type XUnitParser struct{}

func (p *XUnitParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	tree, err := parser.ParseWithPool(ctx, domain.LanguageCSharp, source)
	if err != nil {
		return nil, fmt.Errorf("xunit parser: failed to parse %s: %w", filename, err)
	}
	defer tree.Close()

	root := tree.RootNode()
	suites := parseTestClasses(root, source, filename)

	return &domain.TestFile{
		Path:      filename,
		Language:  domain.LanguageCSharp,
		Framework: frameworkName,
		Suites:    suites,
	}, nil
}

// maxNestedDepth limits recursion depth for nested class parsing.
// C# allows unlimited class nesting, but 20 levels provides a safe buffer
// (real-world maximum observed: 3 in FluentAssertions test suite).
const maxNestedDepth = 20

func getClassStatusAndModifier(attrLists []*sitter.Node, source []byte) (domain.TestStatus, string) {
	for _, attr := range dotnetast.GetAttributes(attrLists) {
		name := dotnetast.GetAttributeName(attr, source)
		if name == "Skip" || name == "SkipAttribute" {
			return domain.TestStatusSkipped, "[Skip]"
		}
	}
	return domain.TestStatusActive, ""
}

// getNamedParameterFromAttribute extracts a named parameter value from an attribute.
// Used for DisplayName from [Fact(DisplayName = "...")] or [Theory(DisplayName = "...")].
func getNamedParameterFromAttribute(attr *sitter.Node, source []byte, paramName string) string {
	argList := dotnetast.FindAttributeArgumentList(attr)
	if argList == nil {
		return ""
	}

	for i := 0; i < int(argList.ChildCount()); i++ {
		arg := argList.Child(i)
		if arg.Type() == dotnetast.NodeAttributeArgument {
			name, value := dotnetast.ParseAssignmentExpression(arg, source)
			if name == paramName {
				return value
			}
		}
	}
	return ""
}

// isSkipped checks if the attribute has Skip parameter set.
func isSkipped(attr *sitter.Node, source []byte) bool {
	argList := dotnetast.FindAttributeArgumentList(attr)
	if argList == nil {
		return false
	}

	for i := 0; i < int(argList.ChildCount()); i++ {
		arg := argList.Child(i)
		if arg.Type() == dotnetast.NodeAttributeArgument {
			name, _ := dotnetast.ParseAssignmentExpression(arg, source)
			if name == "Skip" {
				return true
			}
		}
	}
	return false
}

func parseTestClasses(root *sitter.Node, source []byte, filename string) []domain.TestSuite {
	var suites []domain.TestSuite

	parser.WalkTree(root, func(node *sitter.Node) bool {
		if node.Type() == dotnetast.NodeClassDeclaration {
			if suite := parseTestClassWithDepth(node, source, filename, 0); suite != nil {
				suites = append(suites, *suite)
			}
			return false
		}
		return true
	})

	return suites
}

func parseTestClassWithDepth(node *sitter.Node, source []byte, filename string, depth int) *domain.TestSuite {
	if depth > maxNestedDepth {
		return nil
	}

	className := dotnetast.GetClassName(node, source)
	if className == "" {
		return nil
	}

	attrLists := dotnetast.GetAttributeLists(node)
	classStatus, classModifier := getClassStatusAndModifier(attrLists, source)

	body := dotnetast.GetDeclarationList(node)
	if body == nil {
		return nil
	}

	var tests []domain.Test
	var nestedSuites []domain.TestSuite

	for _, child := range dotnetast.GetDeclarationChildren(body) {
		switch child.Type() {
		case dotnetast.NodeMethodDeclaration:
			tests = append(tests, parseTestMethod(child, source, filename, classStatus, classModifier)...)

		case dotnetast.NodeClassDeclaration:
			if nested := parseTestClassWithDepth(child, source, filename, depth+1); nested != nil {
				nestedSuites = append(nestedSuites, *nested)
			}
		}
	}

	if len(tests) == 0 && len(nestedSuites) == 0 {
		return nil
	}

	return &domain.TestSuite{
		Name:     className,
		Status:   classStatus,
		Modifier: classModifier,
		Location: parser.GetLocation(node, filename),
		Tests:    tests,
		Suites:   nestedSuites,
	}
}

func parseTestMethod(node *sitter.Node, source []byte, filename string, classStatus domain.TestStatus, classModifier string) []domain.Test {
	attrLists := dotnetast.GetAttributeLists(node)
	if len(attrLists) == 0 {
		return nil
	}

	methodName := dotnetast.GetMethodName(node, source)
	if methodName == "" {
		return nil
	}

	attributes := dotnetast.GetAttributes(attrLists)
	status := classStatus
	modifier := classModifier
	location := parser.GetLocation(node, filename)

	var tests []domain.Test
	hasFact := false
	hasTheory := false
	var displayName string
	var theorySkipped bool

	for _, attr := range attributes {
		name := dotnetast.GetAttributeName(attr, source)

		switch {
		case isFactAttribute(name):
			hasFact = true
			displayName = getNamedParameterFromAttribute(attr, source, "DisplayName")
			if isSkipped(attr, source) {
				status = domain.TestStatusSkipped
				modifier = "Skip"
			}

		case isTheoryAttribute(name):
			hasTheory = true
			displayName = getNamedParameterFromAttribute(attr, source, "DisplayName")
			if isSkipped(attr, source) {
				theorySkipped = true
			}

		case name == "InlineData" || name == "InlineDataAttribute":
			testStatus := status
			testModifier := modifier
			if theorySkipped {
				testStatus = domain.TestStatusSkipped
				testModifier = "Skip"
			}
			tests = append(tests, domain.Test{
				Name:     methodName,
				Status:   testStatus,
				Modifier: testModifier,
				Location: location,
			})
		}
	}

	// If [InlineData] attributes were found, return them
	if len(tests) > 0 {
		return tests
	}

	// [Fact] - count as single test
	if hasFact {
		testName := methodName
		if displayName != "" {
			testName = displayName
		}
		return []domain.Test{{
			Name:     testName,
			Status:   status,
			Modifier: modifier,
			Location: location,
		}}
	}

	// [Theory] without [InlineData] - count as single test
	// (includes [MemberData]/[ClassData] which expand at runtime)
	if hasTheory {
		testName := methodName
		if displayName != "" {
			testName = displayName
		}
		testStatus := status
		testModifier := modifier
		if theorySkipped {
			testStatus = domain.TestStatusSkipped
			testModifier = "Skip"
		}
		return []domain.Test{{
			Name:     testName,
			Status:   testStatus,
			Modifier: testModifier,
			Location: location,
		}}
	}

	return nil
}

// isFactAttribute checks if the attribute name represents a Fact-based test.
// xUnit custom test attributes must inherit from FactAttribute and follow
// the naming convention *Fact or *FactAttribute (e.g., UIFact, StaFact).
func isFactAttribute(name string) bool {
	return name == "Fact" || strings.HasSuffix(name, "Fact") ||
		name == "FactAttribute" || strings.HasSuffix(name, "FactAttribute")
}

// isTheoryAttribute checks if the attribute name represents a Theory-based test.
// xUnit custom test attributes must inherit from TheoryAttribute and follow
// the naming convention *Theory or *TheoryAttribute (e.g., UITheory).
func isTheoryAttribute(name string) bool {
	return name == "Theory" || strings.HasSuffix(name, "Theory") ||
		name == "TheoryAttribute" || strings.HasSuffix(name, "TheoryAttribute")
}
