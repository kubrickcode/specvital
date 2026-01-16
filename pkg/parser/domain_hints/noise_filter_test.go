package domain_hints

import "testing"

func TestShouldFilterImportNoise(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		filter bool
	}{
		// Empty string
		{"empty string", "", true},

		// Bare relative markers (noise - no domain signal)
		{"bare dot", ".", true},
		{"bare double dot", "..", true},

		// Relative paths with content (valid - has domain signal)
		{"relative parent with path", "../utils", false},
		{"relative current with path", "./helper", false},
		{"deep relative path", "../../config/settings", false},

		// Normal imports (valid)
		{"npm package", "lodash", false},
		{"scoped package", "@types/node", false},
		{"node builtin", "node:path", false},
		{"absolute path", "/absolute/path", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldFilterImportNoise(tt.path)
			if got != tt.filter {
				t.Errorf("ShouldFilterImportNoise(%q) = %v, want %v", tt.path, got, tt.filter)
			}
		})
	}
}

func TestShouldFilterNoise(t *testing.T) {
	tests := []struct {
		name   string
		call   string
		filter bool
	}{
		// Empty string
		{"empty string", "", true},

		// Bracket patterns (spread array)
		{"starts with bracket", "[.", true},
		{"starts with bracket long", "[...arr]", true},
		{"spreadsheet identifier edge case", "[.forEach", true},

		// Parenthesis patterns (decimal literal bug)
		{"starts with paren decimal", "(0.", true},
		{"starts with paren decimal 1", "(1.", true},
		{"starts with paren complex", "(0.5).toFixed", true},

		// String literal method calls (parser artifact)
		{"double quote string method", `"str".encode`, true},
		{"single quote string method", `'str'.upper`, true},
		{"unicode string method", `"ööö".encode`, true},
		{"url string literal", `"http://example.com".format`, true},

		// Function arguments leaked (parser artifact)
		{"function with kwarg", `func(key="value")`, true},
		{"method with kwarg", `requests.Request(method="GET")`, true},
		{"contains equals", `config=value`, true},

		// URL patterns leaked (parser artifact)
		{"url in call http", `requests.Request("GET","http://example`, true},
		{"url in call https", `fetch("https://api.com/endpoint")`, true},
		{"normal call without url", `requests.get`, false},

		// Cheerio/jQuery selector
		{"dollar singleton", "$", true},
		{"dollar with method", "$.ajax", false},

		// Generic callback variable name
		{"fn callback", "fn", true},
		{"fn prefixed", "fnCallback", false},

		// Dot-prefix patterns (C++ :: scope operator conversion artifact)
		// e.g., "::testing::" → ".testing" or "::CreateEvent" → ".CreateEvent"
		{"dot prefix testing", ".testing", true},
		{"dot prefix std", ".std", true},
		{"dot prefix CreateEvent", ".CreateEvent", true},
		{"dot prefix InterlockedIncrement", ".InterlockedIncrement", true},
		{"dot prefix foo", ".foo", true},
		{"dot prefix method", ".foo.bar", true},
		// Valid calls with dots (not at start)
		{"dotted normal", "std.string", false},
		{"dotted method", "userService.create", false},

		// JavaScript inline comments leaked into call
		{"inline comment simple", "res.json()//comment", true},
		{"inline comment with text", "res.json()//Byspec,theruntimecanonly", true},
		{"inline comment multiple", "func()//first//second", true},
		{"no comment", "res.json()", false},

		// Short standalone calls (1-2 chars) - all filtered (no domain signal)
		{"single char lowercase", "a", true},
		{"single char uppercase", "A", true},
		{"single char digit", "1", true},
		{"single char underscore", "_", true},
		{"single char space", " ", true},
		{"single char dot", ".", true},
		{"single char bracket", "[", true},
		{"two char callback", "cb", true},
		{"two char temp var", "xy", true},
		{"two char uppercase", "AB", true},

		// Short calls with dot are preserved (package.method pattern)
		{"two char with dot", "io.Reader", false},
		{"two char pkg call", "os.Exit", false},

		// Normal identifiers
		{"normal identifier", "doSomething", false},
		{"dotted identifier", "service.method", false},
		{"three char identifier", "abc", false},
		{"three char call", "fmt", false},

		// Unbalanced parentheses (parser artifact from method chaining)
		{"unbalanced paren go", "json.NewDecoder(w", true},
		{"unbalanced paren ts", "expect(mockChromeStorage.session", true},
		{"unbalanced paren nested", "expect(badge.style", true},
		{"unbalanced paren complex", "expect(dropdown?.classList", true},

		// Balanced parentheses (normal calls - should pass)
		{"balanced paren simple", "json.Marshal(data)", false},
		{"balanced paren method chain", "expect(result).toEqual", false},
		{"balanced paren empty", "func()", false},
		{"balanced paren nested", "outer(inner(x))", false},

		// Two-letter module names without dots (NOISE - no domain signal)
		// Files: tRPC v10.45.2 analysis
		//   - packages/tests/server/interop/websockets.test.ts: import WebSocket from 'ws'
		//   - packages/tests/server/websockets.test.ts: import WebSocket from 'ws'
		{"ws module import", "ws", true},
		// Files: tRPC v10.45.2 analysis
		//   - packages/tests/server/react/formData.test.tsx: import * as fs from 'fs'
		{"fs module import", "fs", true},
		// Other common node builtins
		{"os module", "os", true},
		{"io module", "io", true},

		// Two-letter calls WITH dots (domain signal preserved)
		{"ws.Server call", "ws.Server", false},
		{"fs.read call", "fs.readFile", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldFilterNoise(tt.call)
			if got != tt.filter {
				t.Errorf("ShouldFilterNoise(%q) = %v, want %v", tt.call, got, tt.filter)
			}
		})
	}
}
