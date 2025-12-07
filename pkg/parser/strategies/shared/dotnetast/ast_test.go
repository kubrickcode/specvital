package dotnetast

import (
	"context"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
)

func parseCS(t *testing.T, source string) *sitter.Node {
	t.Helper()
	parser := sitter.NewParser()
	parser.SetLanguage(csharp.GetLanguage())
	tree, err := parser.ParseCtx(context.Background(), nil, []byte(source))
	if err != nil {
		t.Fatalf("failed to parse C# source: %v", err)
	}
	return tree.RootNode()
}

func TestGetClassName(t *testing.T) {
	source := `public class MyTestClass { }`
	root := parseCS(t, source)

	var className string
	walkTree(root, func(n *sitter.Node) bool {
		if n.Type() == NodeClassDeclaration {
			className = GetClassName(n, []byte(source))
			return false
		}
		return true
	})

	if className != "MyTestClass" {
		t.Errorf("expected 'MyTestClass', got '%s'", className)
	}
}

func TestGetMethodName(t *testing.T) {
	source := `public class C { public void TestMethod() { } }`
	root := parseCS(t, source)

	var methodName string
	walkTree(root, func(n *sitter.Node) bool {
		if n.Type() == NodeMethodDeclaration {
			methodName = GetMethodName(n, []byte(source))
			return false
		}
		return true
	})

	if methodName != "TestMethod" {
		t.Errorf("expected 'TestMethod', got '%s'", methodName)
	}
}

func TestGetAttributeName(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "simple attribute",
			source:   `public class C { [Fact] public void Test() { } }`,
			expected: "Fact",
		},
		{
			name:     "qualified attribute",
			source:   `public class C { [Xunit.Fact] public void Test() { } }`,
			expected: "Fact",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := parseCS(t, tt.source)

			var attrName string
			walkTree(root, func(n *sitter.Node) bool {
				if n.Type() == NodeAttribute {
					attrName = GetAttributeName(n, []byte(tt.source))
					return false
				}
				return true
			})

			if attrName != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, attrName)
			}
		})
	}
}

func TestHasAttribute(t *testing.T) {
	source := `public class C { [Fact] [Theory] public void Test() { } }`
	root := parseCS(t, source)

	var attrLists []*sitter.Node
	walkTree(root, func(n *sitter.Node) bool {
		if n.Type() == NodeMethodDeclaration {
			attrLists = GetAttributeLists(n)
			return false
		}
		return true
	})

	if !HasAttribute(attrLists, []byte(source), "Fact") {
		t.Error("expected to find Fact attribute")
	}
	if !HasAttribute(attrLists, []byte(source), "Theory") {
		t.Error("expected to find Theory attribute")
	}
	if HasAttribute(attrLists, []byte(source), "Skip") {
		t.Error("should not find Skip attribute")
	}
}

func TestGetDeclarationList(t *testing.T) {
	source := `public class C { public void M() { } }`
	root := parseCS(t, source)

	var body *sitter.Node
	walkTree(root, func(n *sitter.Node) bool {
		if n.Type() == NodeClassDeclaration {
			body = GetDeclarationList(n)
			return false
		}
		return true
	})

	if body == nil {
		t.Error("expected non-nil declaration list")
	}
	if body.Type() != NodeDeclarationList {
		t.Errorf("expected declaration_list, got '%s'", body.Type())
	}
}

func walkTree(node *sitter.Node, visitor func(*sitter.Node) bool) {
	if !visitor(node) {
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		walkTree(node.Child(i), visitor)
	}
}
