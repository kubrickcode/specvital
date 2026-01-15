package domain_hints

import "testing"

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

		// Single character calls - all filtered (no domain signal)
		{"single char lowercase", "a", true},
		{"single char uppercase", "A", true},
		{"single char digit", "1", true},
		{"single char underscore", "_", true},
		{"single char space", " ", true},
		{"single char dot", ".", true},
		{"single char bracket", "[", true},

		// Normal identifiers
		{"normal identifier", "doSomething", false},
		{"dotted identifier", "service.method", false},

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
