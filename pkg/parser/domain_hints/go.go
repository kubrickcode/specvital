package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser/tspool"
)

// GoExtractor extracts domain hints from Go source code.
type GoExtractor struct{}

const (
	goImportQuery = `(import_spec path: (interpreted_string_literal) @import)`

	goCallQuery = `
		(function_declaration
			body: (block
				(expression_statement
					(call_expression
						function: [
							(identifier) @call
							(selector_expression) @call
						]
					)
				)
			)
		)
		(function_declaration
			body: (block
				(short_var_declaration
					right: (expression_list
						(call_expression
							function: [
								(identifier) @call
								(selector_expression) @call
							]
						)
					)
				)
			)
		)
	`
)

func (e *GoExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	tree, err := tspool.Parse(ctx, domain.LanguageGo, source)
	if err != nil {
		return nil
	}
	defer tree.Close()

	root := tree.RootNode()

	hints := &domain.DomainHints{
		Imports: extractGoImports(root, source),
		Calls:   extractGoCalls(root, source),
	}

	if len(hints.Imports) == 0 && len(hints.Calls) == 0 {
		return nil
	}

	return hints
}

func extractGoImports(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageGo, goImportQuery)
	if err != nil {
		return nil
	}

	imports := make([]string, 0, len(results))
	for _, r := range results {
		if node, ok := r.Captures["import"]; ok {
			path := trimQuotes(getNodeText(node, source))
			if path == "" {
				continue
			}
			if isGoStdlibImport(path) {
				continue
			}
			imports = append(imports, path)
		}
	}

	return imports
}

// goStdlibImports contains Go standard library packages that provide
// no domain classification signal and should be filtered from imports.
var goStdlibImports = map[string]struct{}{
	// I/O packages (most common noise from Grafana analysis)
	"io":        {},
	"io/fs":     {},
	"io/ioutil": {},
	"os":        {},
	"fs":        {},
	// Formatting and conversion
	"fmt":     {},
	"strconv": {},
	// Time and context (common utility packages)
	"time":    {},
	"context": {},
	// Basic data manipulation
	"bytes":   {},
	"strings": {},
	"unicode": {},
	// Reflection and runtime
	"reflect": {},
	"runtime": {},
	// Encoding (common utility)
	"encoding":       {},
	"encoding/json":  {},
	"encoding/xml":   {},
	"encoding/csv":   {},
	"encoding/base64": {},
	// Path manipulation
	"path":          {},
	"path/filepath": {},
	// Sync primitives
	"sync":       {},
	"sync/atomic": {},
	// Error handling
	"errors": {},
	// Sorting
	"sort": {},
	// Math
	"math":      {},
	"math/rand": {},
	// Logging
	"log": {},
	// Regex
	"regexp": {},
	// Bufio
	"bufio": {},
	// Testing (Go test framework - provides no domain signal)
	"testing": {},
}

// isGoStdlibImport checks if the import path is a Go standard library package.
func isGoStdlibImport(importPath string) bool {
	// Check exact match first
	if _, exists := goStdlibImports[importPath]; exists {
		return true
	}

	// Check if it's a sub-package of stdlib (e.g., "encoding/gob")
	if idx := strings.Index(importPath, "/"); idx > 0 {
		basePkg := importPath[:idx]
		if basePkg == "encoding" || basePkg == "crypto" || basePkg == "archive" ||
			basePkg == "compress" || basePkg == "container" || basePkg == "database" ||
			basePkg == "debug" || basePkg == "go" || basePkg == "hash" ||
			basePkg == "html" || basePkg == "image" || basePkg == "index" ||
			basePkg == "internal" || basePkg == "log" || basePkg == "math" ||
			basePkg == "mime" || basePkg == "net" || basePkg == "os" ||
			basePkg == "plugin" || basePkg == "runtime" || basePkg == "sync" ||
			basePkg == "syscall" || basePkg == "testing" || basePkg == "text" ||
			basePkg == "time" || basePkg == "unicode" {
			return true
		}
	}

	return false
}

func extractGoCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageGo, goCallQuery)
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
			// Normalize: remove whitespace and limit to 2 segments
			call = normalizeCall(call)
			if call == "" {
				continue
			}
			// Filter noise patterns
			if ShouldFilterNoise(call) {
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

func trimQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '`' && s[len(s)-1] == '`') {
		return s[1 : len(s)-1]
	}
	return s
}

// normalizeCall normalizes a function call for domain hints:
// 1. Removes whitespace (newlines, extra spaces)
// 2. Limits to first 2 segments (e.g., "a.b.c.d" -> "a.b")
// This reduces token count while preserving meaningful domain context.
func normalizeCall(call string) string {
	// Remove all whitespace (newlines, tabs, spaces)
	var result strings.Builder
	for _, r := range call {
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' {
			result.WriteRune(r)
		}
	}
	call = result.String()

	// Limit to 2 segments
	parts := strings.SplitN(call, ".", 3)
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return call
}

func getNodeText(node *sitter.Node, source []byte) (result string) {
	start := node.StartByte()
	end := node.EndByte()
	sourceLen := uint32(len(source))

	if start > sourceLen || end > sourceLen {
		return ""
	}

	defer func() {
		if r := recover(); r != nil {
			result = ""
		}
	}()

	return node.Content(source)
}
