package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
	"github.com/kubrickcode/specvital/packages/core/pkg/parser/tspool"
)

// RustExtractor extracts domain hints from Rust source code.
type RustExtractor struct{}

const (
	// use std::collections::HashMap;
	// use crate::models::User;
	// mod tests;
	rustUseQuery = `(use_declaration) @use`
	rustModQuery = `(mod_item name: (identifier) @mod)`

	// Method calls: obj.method(), Type::method()
	rustCallQuery = `
		(call_expression
			function: [
				(identifier) @call
				(scoped_identifier) @call
				(field_expression) @call
			]
		)
	`
)

func (e *RustExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	tree, err := tspool.Parse(ctx, domain.LanguageRust, source)
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

func (e *RustExtractor) extractImports(root *sitter.Node, source []byte) []string {
	seen := make(map[string]struct{})
	var imports []string

	// Extract use declarations
	useResults, err := tspool.QueryWithCache(root, source, domain.LanguageRust, rustUseQuery)
	if err == nil {
		for _, r := range useResults {
			if node, ok := r.Captures["use"]; ok {
				paths := extractRustUsePaths(node, source)
				for _, path := range paths {
					if path != "" {
						if _, exists := seen[path]; !exists {
							seen[path] = struct{}{}
							imports = append(imports, path)
						}
					}
				}
			}
		}
	}

	// Extract mod declarations
	modResults, err := tspool.QueryWithCache(root, source, domain.LanguageRust, rustModQuery)
	if err == nil {
		for _, r := range modResults {
			if node, ok := r.Captures["mod"]; ok {
				modName := getNodeText(node, source)
				if modName != "" {
					if _, exists := seen[modName]; !exists {
						seen[modName] = struct{}{}
						imports = append(imports, modName)
					}
				}
			}
		}
	}

	return imports
}

func (e *RustExtractor) extractCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageRust, rustCallQuery)
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

			// Convert :: to . for normalization
			call = strings.ReplaceAll(call, "::", ".")
			call = normalizeCall(call)
			if call == "" {
				continue
			}

			if ShouldFilterNoise(call) {
				continue
			}

			if isRustTestFrameworkCall(call) {
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

// extractRustUsePaths extracts paths from a use_declaration node.
// Handles: use std::collections::HashMap;
//
//	use crate::models::{User, Order};
//	use super::helpers;
//	use {anyhow::Context, bstr::ByteVec}; (multiline use list)
func extractRustUsePaths(node *sitter.Node, source []byte) []string {
	// Get the full use declaration text and extract the path
	text := getNodeText(node, source)
	if text == "" {
		return nil
	}

	// Remove "use " prefix and trailing ";"
	text = strings.TrimPrefix(text, "use ")
	text = strings.TrimSuffix(text, ";")
	text = strings.TrimSpace(text)

	// Handle direct use list: use { item1, item2 } or use {item1, item2}
	if strings.HasPrefix(text, "{") {
		return parseRustUseList(text)
	}

	// Handle use lists with base path: use crate::{a, b} -> crate
	if idx := strings.Index(text, "::{"); idx > 0 {
		text = text[:idx]
	}

	// Handle use as: use std::collections::HashMap as Map -> std::collections::HashMap
	if idx := strings.Index(text, " as "); idx > 0 {
		text = text[:idx]
	}

	// Handle wildcard: use std::* -> std
	text = strings.TrimSuffix(text, "::*")

	// Convert :: to / for consistency with other languages
	text = strings.ReplaceAll(text, "::", "/")

	return []string{text}
}

// parseRustUseList parses a use list like "{ item1::path, item2::path }" into individual paths.
// Handles nested braces: { a::{b, c}, d } -> ["a", "d"] (base paths only)
func parseRustUseList(text string) []string {
	// Remove surrounding braces
	text = strings.TrimPrefix(text, "{")
	text = strings.TrimSuffix(text, "}")
	text = strings.TrimSpace(text)

	if text == "" {
		return nil
	}

	// Split by comma, respecting nested braces
	items := splitRustUseItems(text)
	var paths []string

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		// Handle nested use list: item::{a, b} -> item (base path only)
		if idx := strings.Index(item, "::{"); idx > 0 {
			item = item[:idx]
		}

		// Handle alias: item::path as Alias -> item::path
		if idx := strings.Index(item, " as "); idx > 0 {
			item = item[:idx]
		}

		// Handle wildcard: item::* -> item
		item = strings.TrimSuffix(item, "::*")

		// Convert :: to / for consistency with other languages
		item = strings.ReplaceAll(item, "::", "/")

		if item != "" {
			paths = append(paths, item)
		}
	}

	return paths
}

// splitRustUseItems splits use list items by comma, respecting nested braces.
// Example: "a::{b, c}, d" -> ["a::{b, c}", "d"]
func splitRustUseItems(text string) []string {
	var items []string
	var current strings.Builder
	depth := 0

	for _, ch := range text {
		switch ch {
		case '{':
			depth++
			current.WriteRune(ch)
		case '}':
			depth--
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				items = append(items, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	// Add the last item
	if current.Len() > 0 {
		items = append(items, current.String())
	}

	return items
}

// rustTestFrameworkCalls contains patterns from Rust test frameworks
// that should be excluded from domain hints.
var rustTestFrameworkCalls = map[string]struct{}{
	// Standard test macros
	"assert":        {},
	"assert_eq":     {},
	"assert_ne":     {},
	"debug_assert":  {},
	"panic":         {},
	"unreachable":   {},
	"todo":          {},
	"unimplemented": {},
	// Common test utilities
	"println":  {},
	"print":    {},
	"eprintln": {},
	"eprint":   {},
	"dbg":      {},
	"format":   {},
	"vec":      {},
	// tokio-test
	"tokio.test": {},
	// proptest
	"proptest":       {},
	"prop_assert":    {},
	"prop_assert_eq": {},
	// Rust stdlib enums (domain noise in test files)
	"Ok":   {},
	"Err":  {},
	"Some": {},
	"None": {},
}

func isRustTestFrameworkCall(call string) bool {
	baseName := call
	if idx := strings.Index(call, "."); idx > 0 {
		baseName = call[:idx]
	}
	// Also check the full call for macros
	_, existsBase := rustTestFrameworkCalls[baseName]
	_, existsFull := rustTestFrameworkCalls[call]
	return existsBase || existsFull
}
