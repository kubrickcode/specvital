package kotlinast

import (
	"context"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/kotlin"
)

func parseKotlin(t *testing.T, content string) *sitter.Node {
	t.Helper()
	parser := sitter.NewParser()
	parser.SetLanguage(kotlin.GetLanguage())
	tree, err := parser.ParseCtx(context.Background(), nil, []byte(content))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	return tree.RootNode()
}

func findClassDeclaration(root *sitter.Node) *sitter.Node {
	for i := 0; i < int(root.ChildCount()); i++ {
		child := root.Child(i)
		if child.Type() == NodeClassDeclaration {
			return child
		}
	}
	return nil
}

func TestGetClassName(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple class",
			content:  `class MyTest {}`,
			expected: "MyTest",
		},
		{
			name:     "class with FunSpec",
			content:  `class CalculatorTest : FunSpec({})`,
			expected: "CalculatorTest",
		},
		{
			name:     "class with annotation",
			content:  `@Disabled class SkippedTest : StringSpec({})`,
			expected: "SkippedTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseKotlin(t, tt.content)
			classNode := findClassDeclaration(root)
			if classNode == nil {
				t.Fatal("class node not found")
			}

			result := GetClassName(classNode, []byte(tt.content))
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsKotestSpec(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expected      bool
		expectedStyle string
	}{
		{
			name:          "FunSpec",
			content:       `class MyTest : FunSpec({})`,
			expected:      true,
			expectedStyle: SpecFunSpec,
		},
		{
			name:          "StringSpec",
			content:       `class MyTest : StringSpec({})`,
			expected:      true,
			expectedStyle: SpecStringSpec,
		},
		{
			name:          "BehaviorSpec",
			content:       `class MyTest : BehaviorSpec({})`,
			expected:      true,
			expectedStyle: SpecBehaviorSpec,
		},
		{
			name:          "DescribeSpec",
			content:       `class MyTest : DescribeSpec({})`,
			expected:      true,
			expectedStyle: SpecDescribeSpec,
		},
		{
			name:          "AnnotationSpec",
			content:       `class MyTest : AnnotationSpec() {}`,
			expected:      true,
			expectedStyle: SpecAnnotationSpec,
		},
		{
			name:          "not a spec",
			content:       `class MyClass {}`,
			expected:      false,
			expectedStyle: "",
		},
		{
			name:          "extends other class",
			content:       `class MyClass : BaseClass() {}`,
			expected:      false,
			expectedStyle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseKotlin(t, tt.content)
			classNode := findClassDeclaration(root)
			if classNode == nil {
				t.Fatal("class node not found")
			}

			isSpec, style := IsKotestSpec(classNode, []byte(tt.content))
			if isSpec != tt.expected {
				t.Errorf("IsKotestSpec: expected %v, got %v", tt.expected, isSpec)
			}
			if style != tt.expectedStyle {
				t.Errorf("style: expected %q, got %q", tt.expectedStyle, style)
			}
		})
	}
}

func TestExtractStringContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple string",
			content:  `val x = "hello world"`,
			expected: "hello world",
		},
		{
			name:     "empty string",
			content:  `val x = ""`,
			expected: "",
		},
	}

	parser := sitter.NewParser()
	parser.SetLanguage(kotlin.GetLanguage())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := parser.ParseCtx(context.Background(), nil, []byte(tt.content))
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			root := tree.RootNode()
			var stringNode *sitter.Node

			var findString func(n *sitter.Node)
			findString = func(n *sitter.Node) {
				if n.Type() == NodeStringLiteral || n.Type() == NodeLineStringLiteral {
					stringNode = n
					return
				}
				for i := 0; i < int(n.ChildCount()); i++ {
					findString(n.Child(i))
					if stringNode != nil {
						return
					}
				}
			}
			findString(root)

			if stringNode == nil {
				t.Fatal("string node not found")
			}

			result := ExtractStringContent(stringNode, []byte(tt.content))
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsKotlinTestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Filename patterns
		{"Test suffix", "UserTest.kt", true},
		{"Tests suffix", "UserTests.kt", true},
		{"Spec suffix", "UserSpec.kt", true},
		{"Test prefix", "TestUser.kt", true},
		{"kts script test", "UserSpec.kts", true},
		{"non-test file", "User.kt", false},
		{"java file", "UserTest.java", false},

		// Directory patterns
		{"in test dir", "src/test/kotlin/User.kt", true},
		{"in tests dir", "project/tests/User.kt", true},
		{"test dir mid path", "project/test/User.kt", true},
		{"src/test dir", "src/test/User.kt", true},

		// Combined
		{"test dir with test suffix", "src/test/kotlin/UserTest.kt", true},
		{"main dir non-test", "src/main/kotlin/User.kt", false},

		// Edge cases
		{"windows path", "src\\test\\kotlin\\User.kt", true},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsKotlinTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsKotlinTestFile(%q): expected %v, got %v", tt.path, tt.expected, result)
			}
		})
	}
}

func TestHasAnnotation(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		annotation string
		expected   bool
	}{
		{
			name:       "has Test annotation",
			content:    `@Test fun myTest() {}`,
			annotation: "Test",
			expected:   true,
		},
		{
			name:       "has Disabled annotation",
			content:    `@Disabled @Test fun myTest() {}`,
			annotation: "Disabled",
			expected:   true,
		},
		{
			name:       "no matching annotation",
			content:    `@Test fun myTest() {}`,
			annotation: "Disabled",
			expected:   false,
		},
		{
			name:       "no annotations",
			content:    `fun myTest() {}`,
			annotation: "Test",
			expected:   false,
		},
	}

	parser := sitter.NewParser()
	parser.SetLanguage(kotlin.GetLanguage())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := parser.ParseCtx(context.Background(), nil, []byte(tt.content))
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			root := tree.RootNode()
			var funcNode *sitter.Node

			var findFunc func(n *sitter.Node)
			findFunc = func(n *sitter.Node) {
				if n.Type() == NodeFunctionDeclaration {
					funcNode = n
					return
				}
				for i := 0; i < int(n.ChildCount()); i++ {
					findFunc(n.Child(i))
					if funcNode != nil {
						return
					}
				}
			}
			findFunc(root)

			if funcNode == nil {
				t.Fatal("function node not found")
			}

			modifiers := GetModifiers(funcNode)
			result := HasAnnotation(modifiers, []byte(tt.content), tt.annotation)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetCallExpressionName(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple call",
			content:  `fun x() { test("name") }`,
			expected: "test",
		},
		{
			name:     "println call",
			content:  `fun x() { println("hello") }`,
			expected: "println",
		},
	}

	parser := sitter.NewParser()
	parser.SetLanguage(kotlin.GetLanguage())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := parser.ParseCtx(context.Background(), nil, []byte(tt.content))
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			root := tree.RootNode()
			var callNode *sitter.Node

			var findCall func(n *sitter.Node)
			findCall = func(n *sitter.Node) {
				if n.Type() == NodeCallExpression {
					callNode = n
					return
				}
				for i := 0; i < int(n.ChildCount()); i++ {
					findCall(n.Child(i))
					if callNode != nil {
						return
					}
				}
			}
			findCall(root)

			if callNode == nil {
				t.Fatal("call expression not found")
			}

			result := GetCallExpressionName(callNode, []byte(tt.content))
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
