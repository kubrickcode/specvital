package domain_hints

import (
	"context"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/kubrickcode/specvital/lib/parser/domain"
	"github.com/kubrickcode/specvital/lib/parser/tspool"
)

// CppExtractor extracts domain hints from C++ source code.
type CppExtractor struct{}

const (
	// #include <iostream>, #include "myheader.h"
	cppIncludeQuery = `(preproc_include path: (_) @include)`

	// Function/method calls
	cppCallQuery = `
		(call_expression
			function: [
				(identifier) @call
				(qualified_identifier) @call
				(field_expression) @call
			]
		)
	`
)

func (e *CppExtractor) Extract(ctx context.Context, source []byte) *domain.DomainHints {
	tree, err := tspool.Parse(ctx, domain.LanguageCpp, source)
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

func (e *CppExtractor) extractImports(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageCpp, cppIncludeQuery)
	if err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	var imports []string

	for _, r := range results {
		if node, ok := r.Captures["include"]; ok {
			path := extractCppIncludePath(node, source)
			if path == "" {
				continue
			}

			// Filter C++ stdlib headers (no domain classification value)
			if isCppStdlibImport(path) {
				continue
			}

			if _, exists := seen[path]; exists {
				continue
			}
			seen[path] = struct{}{}
			imports = append(imports, path)
		}
	}

	return imports
}

func (e *CppExtractor) extractCalls(root *sitter.Node, source []byte) []string {
	results, err := tspool.QueryWithCache(root, source, domain.LanguageCpp, cppCallQuery)
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
			// Handle -> operator for pointer access
			call = strings.ReplaceAll(call, "->", ".")
			call = normalizeCall(call)
			if call == "" {
				continue
			}

			if ShouldFilterNoise(call) {
				continue
			}

			if isCppTestFrameworkCall(call) {
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

// extractCppIncludePath extracts the path from an include directive.
// Handles: #include <iostream>
//
//	#include "myheader.h"
//	#include <gtest/gtest.h>
func extractCppIncludePath(node *sitter.Node, source []byte) string {
	text := getNodeText(node, source)
	if text == "" {
		return ""
	}

	text = strings.TrimSpace(text)

	// Remove angle brackets or quotes
	if len(text) >= 2 {
		if (text[0] == '<' && text[len(text)-1] == '>') ||
			(text[0] == '"' && text[len(text)-1] == '"') {
			text = text[1 : len(text)-1]
		}
	}

	return text
}

// cppTestFrameworkCalls contains patterns from C++ test frameworks
// that should be excluded from domain hints.
var cppTestFrameworkCalls = map[string]struct{}{
	// Google Test
	"EXPECT_TRUE":              {},
	"EXPECT_FALSE":             {},
	"EXPECT_EQ":                {},
	"EXPECT_NE":                {},
	"EXPECT_LT":                {},
	"EXPECT_LE":                {},
	"EXPECT_GT":                {},
	"EXPECT_GE":                {},
	"EXPECT_STREQ":             {},
	"EXPECT_STRNE":             {},
	"EXPECT_THROW":             {},
	"EXPECT_NO_THROW":          {},
	"EXPECT_DEATH":             {},
	"ASSERT_TRUE":              {},
	"ASSERT_FALSE":             {},
	"ASSERT_EQ":                {},
	"ASSERT_NE":                {},
	"ASSERT_LT":                {},
	"ASSERT_LE":                {},
	"ASSERT_GT":                {},
	"ASSERT_GE":                {},
	"ASSERT_STREQ":             {},
	"ASSERT_STRNE":             {},
	"ASSERT_THROW":             {},
	"ASSERT_NO_THROW":          {},
	"ASSERT_DEATH":             {},
	"TEST":                     {},
	"TEST_F":                   {},
	"TEST_P":                   {},
	"TYPED_TEST":               {},
	"TYPED_TEST_SUITE":         {},
	"INSTANTIATE_TEST_SUITE_P": {},
	// Catch2
	"REQUIRE":         {},
	"REQUIRE_FALSE":   {},
	"REQUIRE_THROWS":  {},
	"REQUIRE_NOTHROW": {},
	"CHECK":           {},
	"CHECK_FALSE":     {},
	"CHECK_THROWS":    {},
	"CHECK_NOTHROW":   {},
	"SECTION":         {},
	"TEST_CASE":       {},
	"SCENARIO":        {},
	"GIVEN":           {},
	"WHEN":            {},
	"THEN":            {},
	// Common utilities
	"std.cout": {},
	"std.cerr": {},
	"std.endl": {},
	"printf":   {},
	"fprintf":  {},
	"cout":     {},
	"cerr":     {},
}

func isCppTestFrameworkCall(call string) bool {
	baseName := call
	if idx := strings.Index(call, "."); idx > 0 {
		baseName = call[:idx]
	}
	_, existsBase := cppTestFrameworkCalls[baseName]
	_, existsFull := cppTestFrameworkCalls[call]
	return existsBase || existsFull
}

// cppStdlibHeaders contains C++ standard library headers that provide
// no domain classification signal and should be filtered from imports.
// These are universal language primitives without domain-specific meaning.
var cppStdlibHeaders = map[string]struct{}{
	// C++ STL containers
	"array":         {},
	"deque":         {},
	"forward_list":  {},
	"list":          {},
	"map":           {},
	"queue":         {},
	"set":           {},
	"span":          {},
	"stack":         {},
	"unordered_map": {},
	"unordered_set": {},
	"vector":        {},
	// C++ STL algorithms and utilities
	"algorithm":        {},
	"any":              {},
	"bitset":           {},
	"charconv":         {}, // C++17
	"chrono":           {},
	"compare":          {},
	"complex":          {},
	"concepts":         {}, // C++20
	"coroutine":        {}, // C++20
	"exception":        {},
	"execution":        {}, // C++17
	"expected":         {},
	"filesystem":       {}, // C++17
	"format":           {}, // C++20
	"functional":       {},
	"initializer_list": {},
	"iterator":         {},
	"limits":           {},
	"locale":           {},
	"memory":           {},
	"numbers":          {}, // C++20
	"numeric":          {},
	"optional":         {},
	"random":           {},
	"ranges":           {},
	"ratio":            {},
	"regex":            {},
	"source_location":  {}, // C++20
	"string":           {},
	"string_view":      {},
	"tuple":            {},
	"type_traits":      {},
	"typeinfo":         {},
	"utility":          {},
	"valarray":         {},
	"variant":          {},
	// C++ I/O streams
	"fstream":   {},
	"iomanip":   {},
	"ios":       {},
	"iosfwd":    {},
	"iostream":  {},
	"istream":   {},
	"ostream":   {},
	"sstream":   {},
	"streambuf": {},
	// C++ threading and synchronization
	"atomic":             {},
	"barrier":            {},
	"condition_variable": {},
	"future":             {},
	"latch":              {},
	"mutex":              {},
	"semaphore":          {},
	"shared_mutex":       {},
	"stop_token":         {},
	"thread":             {},
	// C++ error handling
	"stdexcept":    {},
	"system_error": {},
	// C++ memory management
	"new":              {},
	"memory_resource":  {},
	"scoped_allocator": {},
	// C compatibility headers
	"cassert":   {},
	"cctype":    {},
	"cerrno":    {},
	"cfenv":     {},
	"cfloat":    {},
	"cinttypes": {},
	"climits":   {},
	"clocale":   {},
	"cmath":     {},
	"csetjmp":   {},
	"csignal":   {},
	"cstdarg":   {},
	"cstddef":   {},
	"cstdint":   {},
	"cstdio":    {},
	"cstdlib":   {},
	"cstring":   {},
	"ctime":     {},
	"cuchar":    {},
	"cwchar":    {},
	"cwctype":   {},
	// C headers (legacy)
	"assert.h":   {},
	"ctype.h":    {},
	"errno.h":    {},
	"fenv.h":     {},
	"float.h":    {},
	"inttypes.h": {},
	"limits.h":   {},
	"locale.h":   {},
	"math.h":     {},
	"setjmp.h":   {},
	"signal.h":   {},
	"stdarg.h":   {},
	"stddef.h":   {},
	"stdint.h":   {},
	"stdio.h":    {},
	"stdlib.h":   {},
	"string.h":   {},
	"time.h":     {},
	"uchar.h":    {},
	"wchar.h":    {},
	"wctype.h":   {},
	// Platform-specific (POSIX)
	"dirent.h":    {},
	"dlfcn.h":     {},
	"fcntl.h":     {},
	"fnmatch.h":   {},
	"glob.h":      {},
	"grp.h":       {},
	"poll.h":      {},
	"pthread.h":   {},
	"pwd.h":       {},
	"sched.h":     {},
	"semaphore.h": {},
	"strings.h":   {},
	"syslog.h":    {},
	"termios.h":   {},
	"unistd.h":    {},
	// Platform-specific (Windows)
	"windows.h": {},
	"direct.h":  {},
	"io.h":      {},
	"objbase.h": {},
}

// isCppStdlibImport checks if the import path is a C++ standard library header.
func isCppStdlibImport(importPath string) bool {
	// Exact match
	if _, exists := cppStdlibHeaders[importPath]; exists {
		return true
	}

	// Prefix match for POSIX system directories
	// These are universal OS primitives without domain-specific signal
	for _, prefix := range cppStdlibPrefixes {
		if strings.HasPrefix(importPath, prefix) {
			return true
		}
	}

	return false
}

// cppStdlibPrefixes contains directory prefixes for system headers
// that should be filtered regardless of specific filename.
var cppStdlibPrefixes = []string{
	"sys/",     // POSIX system headers (sys/socket.h, sys/mman.h, etc.)
	"netinet/", // Network headers (netinet/in.h, netinet/tcp.h, etc.)
	"arpa/",    // ARPA headers (arpa/inet.h, etc.)
	"linux/",   // Linux-specific headers
	"bits/",    // glibc implementation details
}
