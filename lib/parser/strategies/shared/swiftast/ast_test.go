package swiftast

import (
	"context"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/swift"
)

func parseSwift(t *testing.T, content string) *sitter.Node {
	t.Helper()
	parser := sitter.NewParser()
	parser.SetLanguage(swift.GetLanguage())
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
			name:     "class with XCTestCase",
			content:  `class CalculatorTests: XCTestCase {}`,
			expected: "CalculatorTests",
		},
		{
			name:     "class with access modifier",
			content:  `public class PublicTests: XCTestCase {}`,
			expected: "PublicTests",
		},
		{
			name:     "final class",
			content:  `final class FinalTests: XCTestCase {}`,
			expected: "FinalTests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseSwift(t, tt.content)
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

func TestGetSuperTypes(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "single inheritance",
			content:  `class MyTests: XCTestCase {}`,
			expected: []string{"XCTestCase"},
		},
		{
			name:     "multiple protocols",
			content:  `class MyTests: XCTestCase, Codable {}`,
			expected: []string{"XCTestCase", "Codable"},
		},
		{
			name:     "no inheritance",
			content:  `class MyClass {}`,
			expected: nil,
		},
		{
			name:     "custom base class",
			content:  `class MyTests: BaseTestCase {}`,
			expected: []string{"BaseTestCase"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseSwift(t, tt.content)
			classNode := findClassDeclaration(root)
			if classNode == nil {
				t.Fatal("class node not found")
			}

			result := GetSuperTypes(classNode, []byte(tt.content))
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d super types, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected %q at index %d, got %q", tt.expected[i], i, v)
				}
			}
		})
	}
}

func TestIsXCTestCase(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "direct XCTestCase",
			content:  `class MyTests: XCTestCase {}`,
			expected: true,
		},
		{
			name:     "custom TestCase suffix",
			content:  `class MyTests: BaseTestCase {}`,
			expected: true,
		},
		{
			name:     "not a test class",
			content:  `class MyClass {}`,
			expected: false,
		},
		{
			name:     "extends other class",
			content:  `class MyClass: BaseClass {}`,
			expected: false,
		},
		{
			name:     "UseCase suffix - not a test",
			content:  `class LoginUseCase: BaseUseCase {}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseSwift(t, tt.content)
			classNode := findClassDeclaration(root)
			if classNode == nil {
				t.Fatal("class node not found")
			}

			result := IsXCTestCase(classNode, []byte(tt.content))
			if result != tt.expected {
				t.Errorf("IsXCTestCase: expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetFunctionName(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple function",
			content:  `func testExample() {}`,
			expected: "testExample",
		},
		{
			name:     "function with parameters",
			content:  `func testAdd(_ a: Int, _ b: Int) {}`,
			expected: "testAdd",
		},
		{
			name:     "async function",
			content:  `func testAsync() async {}`,
			expected: "testAsync",
		},
		{
			name:     "throws function",
			content:  `func testThrows() throws {}`,
			expected: "testThrows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseSwift(t, tt.content)
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

			result := GetFunctionName(funcNode, []byte(tt.content))
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsTestFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		expected bool
	}{
		{"valid test", "testExample", true},
		{"valid camelCase", "testMyFunction", true},
		{"missing uppercase", "testexample", false},
		{"too short", "test", false},
		{"not prefixed", "myTest", false},
		{"edge case length", "testA", true},
		{"empty string", "", false},
		{"only test prefix lowercase", "testing", false},
		{"setUp not a test", "setUp", false},
		{"tearDown not a test", "tearDown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTestFunction(tt.funcName)
			if result != tt.expected {
				t.Errorf("IsTestFunction(%q) = %v, want %v", tt.funcName, result, tt.expected)
			}
		})
	}
}

func TestIsSwiftTestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Filename patterns
		{"Test suffix", "MyTest.swift", true},
		{"Tests suffix", "MyTests.swift", true},
		{"non-test file", "MyClass.swift", false},
		{"java file", "MyTest.java", false},

		// Directory patterns
		{"in Tests dir", "Tests/MyFile.swift", true},
		{"in XCTests dir", "XCTests/MyFile.swift", true},
		{"in nested Tests dir", "project/Tests/Unit/MyFile.swift", true},

		// Combined
		{"Tests dir with test suffix", "Tests/MyTests.swift", true},
		{"Source dir non-test", "Sources/MyClass.swift", false},

		// Edge cases
		{"windows path", "Tests\\MyTests.swift", true},
		{"empty path", "", false},
		{"test in filename but not suffix", "TestHelper.swift", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSwiftTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsSwiftTestFile(%q): expected %v, got %v", tt.path, tt.expected, result)
			}
		})
	}
}

func TestGetClassBody(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectBody  bool
		expectFuncs int
	}{
		{
			name:        "class with methods",
			content:     `class MyTests: XCTestCase { func testOne() {} func testTwo() {} }`,
			expectBody:  true,
			expectFuncs: 2,
		},
		{
			name:        "empty class",
			content:     `class MyTests: XCTestCase {}`,
			expectBody:  true,
			expectFuncs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseSwift(t, tt.content)
			classNode := findClassDeclaration(root)
			if classNode == nil {
				t.Fatal("class node not found")
			}

			body := GetClassBody(classNode)
			if tt.expectBody && body == nil {
				t.Error("expected class body, got nil")
			}
			if !tt.expectBody && body != nil {
				t.Error("expected no class body, got non-nil")
			}
		})
	}
}
