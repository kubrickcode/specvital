// Package swiftast provides shared Swift AST traversal utilities for test framework parsers.
package swiftast

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// Swift AST node types.
const (
	NodeClassDeclaration     = "class_declaration"
	NodeFunctionDeclaration  = "function_declaration"
	NodeImportDeclaration    = "import_declaration"
	NodeIdentifier           = "simple_identifier"
	NodeTypeIdentifier       = "type_identifier"
	NodeInheritanceSpecifier = "inheritance_specifier"
	NodeClassBody            = "class_body"
	NodeFunctionBody         = "function_body"
	NodeStatements           = "statements"
	NodeAttribute            = "attribute"
	NodeModifiers            = "modifiers"
	NodeUserType             = "user_type"
)

// GetClassName extracts the class name from a class_declaration node.
func GetClassName(node *sitter.Node, source []byte) string {
	// Try field name first
	nameNode := node.ChildByFieldName("name")
	if nameNode != nil {
		return nameNode.Content(source)
	}

	// Fallback: look for identifier child
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == NodeIdentifier || child.Type() == NodeTypeIdentifier {
			return child.Content(source)
		}
	}
	return ""
}

// GetSuperTypes returns the inherited types from a class_declaration.
func GetSuperTypes(node *sitter.Node, source []byte) []string {
	var types []string

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == NodeInheritanceSpecifier || child.Type() == "inheritance_clause" {
			for j := 0; j < int(child.ChildCount()); j++ {
				typeChild := child.Child(j)
				typeName := extractTypeName(typeChild, source)
				if typeName != "" {
					types = append(types, typeName)
				}
			}
		}
	}
	return types
}

func extractTypeName(node *sitter.Node, source []byte) string {
	if node == nil {
		return ""
	}

	switch node.Type() {
	case NodeTypeIdentifier, NodeIdentifier:
		return node.Content(source)
	case NodeUserType:
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if name := extractTypeName(child, source); name != "" {
				return name
			}
		}
		return node.Content(source)
	}

	// Recurse into children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if name := extractTypeName(child, source); name != "" {
			return name
		}
	}
	return ""
}

// IsXCTestCase checks if a class inherits from XCTestCase or a likely test base class.
// Detects:
// - Direct inheritance: class MyTests: XCTestCase
// - Indirect inheritance: class MyTests: BaseTestCase (common pattern)
// - Class names ending with "Tests" or "Test" inheriting from *TestCase
func IsXCTestCase(node *sitter.Node, source []byte) bool {
	supers := GetSuperTypes(node, source)
	for _, s := range supers {
		if s == "XCTestCase" {
			return true
		}
		// Common pattern: BaseTestCase, CustomTestCase, etc.
		if strings.HasSuffix(s, "TestCase") {
			return true
		}
	}
	return false
}

// GetClassBody returns the class body node from a class_declaration.
func GetClassBody(node *sitter.Node) *sitter.Node {
	return node.ChildByFieldName("body")
}

// GetFunctionName extracts the function name from a function_declaration node.
func GetFunctionName(node *sitter.Node, source []byte) string {
	// Try field name first
	nameNode := node.ChildByFieldName("name")
	if nameNode != nil {
		return nameNode.Content(source)
	}

	// Look for identifier children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == NodeIdentifier {
			return child.Content(source)
		}
	}
	return ""
}

// IsTestFunction checks if a function is a test method (starts with "test").
func IsTestFunction(funcName string) bool {
	return strings.HasPrefix(funcName, "test") && len(funcName) > 4 && funcName[4] >= 'A' && funcName[4] <= 'Z'
}

// IsSwiftTestFile checks if the path matches Swift test file patterns.
func IsSwiftTestFile(path string) bool {
	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	base := normalizedPath
	if idx := strings.LastIndex(normalizedPath, "/"); idx >= 0 {
		base = normalizedPath[idx+1:]
	}

	if !strings.HasSuffix(base, ".swift") {
		return false
	}

	name := strings.TrimSuffix(base, ".swift")

	// Swift test naming conventions: *Tests.swift, *Test.swift
	if strings.HasSuffix(name, "Test") || strings.HasSuffix(name, "Tests") {
		return true
	}

	// Directory-based detection
	if strings.Contains(normalizedPath, "/Tests/") ||
		strings.Contains(normalizedPath, "/XCTests/") ||
		strings.Contains(normalizedPath, "Tests/") {
		return true
	}

	return false
}
