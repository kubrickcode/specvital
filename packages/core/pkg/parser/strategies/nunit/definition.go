// Package nunit implements NUnit test framework support for C# test files.
package nunit

import (
	"context"
	"fmt"
	"regexp"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/framework/matchers"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/strategies/shared/dotnetast"
)

const frameworkName = "nunit"

func init() {
	framework.Register(NewDefinition())
}

func NewDefinition() *framework.Definition {
	return &framework.Definition{
		Name:      frameworkName,
		Languages: []domain.Language{domain.LanguageCSharp},
		Matchers: []framework.Matcher{
			matchers.NewImportMatcher(
				"NUnit.Framework",
				"using NUnit.Framework",
			),
			&NUnitFileMatcher{},
			&NUnitContentMatcher{},
		},
		ConfigParser: nil,
		Parser:       &NUnitParser{},
		Priority:     framework.PriorityGeneric,
	}
}

// NUnitFileMatcher matches *Test.cs, *Tests.cs, Test*.cs, *Spec.cs, *Specs.cs files.
type NUnitFileMatcher struct{}

func (m *NUnitFileMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileName {
		return framework.NoMatch()
	}

	if dotnetast.IsCSharpTestFileName(signal.Value) {
		return framework.PartialMatch(20, "NUnit file naming convention")
	}

	return framework.NoMatch()
}

// NUnitContentMatcher matches NUnit specific patterns.
type NUnitContentMatcher struct{}

var nunitPatterns = []struct {
	pattern *regexp.Regexp
	desc    string
}{
	{regexp.MustCompile(`\[Test\]`), "[Test] attribute"},
	{regexp.MustCompile(`\[TestCase\(`), "[TestCase] attribute"},
	{regexp.MustCompile(`\[TestFixture\]`), "[TestFixture] attribute"},
	{regexp.MustCompile(`\[TestCaseSource\(`), "[TestCaseSource] attribute"},
	{regexp.MustCompile(`\[Theory\]`), "[Theory] attribute"},
	{regexp.MustCompile(`using\s+NUnit\.Framework\s*;`), "using NUnit.Framework"},
	{regexp.MustCompile(`\[Ignore\(`), "[Ignore] attribute"},
}

func (m *NUnitContentMatcher) Match(ctx context.Context, signal framework.Signal) framework.MatchResult {
	if signal.Type != framework.SignalFileContent {
		return framework.NoMatch()
	}

	content, ok := signal.Context.([]byte)
	if !ok {
		content = []byte(signal.Value)
	}

	for _, p := range nunitPatterns {
		if p.pattern.Match(content) {
			return framework.PartialMatch(40, "Found NUnit pattern: "+p.desc)
		}
	}

	return framework.NoMatch()
}

// NUnitParser extracts test definitions from C# NUnit files.
type NUnitParser struct{}

func (p *NUnitParser) Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error) {
	tree, err := parser.ParseWithPool(ctx, domain.LanguageCSharp, source)
	if err != nil {
		return nil, fmt.Errorf("nunit parser: failed to parse %s: %w", filename, err)
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
// (real-world maximum observed: 3 in NUnit's own test suite).
const maxNestedDepth = 20

func getClassStatusAndModifier(attrLists []*sitter.Node, source []byte) (domain.TestStatus, string) {
	for _, attr := range dotnetast.GetAttributes(attrLists) {
		name := dotnetast.GetAttributeName(attr, source)
		if name == "Ignore" || name == "IgnoreAttribute" {
			return domain.TestStatusSkipped, "[Ignore]"
		}
	}
	return domain.TestStatusActive, ""
}

// getNamedParameterFromAttribute extracts a named parameter value from an attribute.
// Used for Description from [Test(Description = "...")] or TestName from [TestCase(TestName = "...")].
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

	// Check for [Ignore] attribute (reuses class-level check)
	if methodStatus, methodModifier := getClassStatusAndModifier(attrLists, source); methodStatus == domain.TestStatusSkipped {
		status = methodStatus
		modifier = methodModifier
	}

	location := parser.GetLocation(node, filename)
	var tests []domain.Test
	hasSimpleTest := false
	hasTestCaseSource := false
	var testDescription string

	for _, attr := range attributes {
		name := dotnetast.GetAttributeName(attr, source)

		switch name {
		case "Test", "TestAttribute", "Theory", "TheoryAttribute":
			hasSimpleTest = true
			testDescription = getNamedParameterFromAttribute(attr, source, "Description")

		case "TestCase", "TestCaseAttribute":
			testName := getNamedParameterFromAttribute(attr, source, "TestName")
			if testName == "" {
				testName = methodName
			}
			tests = append(tests, domain.Test{
				Name:     testName,
				Status:   status,
				Modifier: modifier,
				Location: location,
			})

		case "TestCaseSource", "TestCaseSourceAttribute":
			hasTestCaseSource = true
		}
	}

	// If [TestCase] attributes were found, return them
	if len(tests) > 0 {
		return tests
	}

	// [Test], [Theory], or [TestCaseSource] - count as single test
	if hasSimpleTest || hasTestCaseSource {
		testName := methodName
		if testDescription != "" {
			testName = testDescription
		}
		return []domain.Test{{
			Name:     testName,
			Status:   status,
			Modifier: modifier,
			Location: location,
		}}
	}

	return nil
}
