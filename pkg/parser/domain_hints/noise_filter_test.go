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

		// Cheerio/jQuery selector
		{"dollar singleton", "$", true},
		{"dollar with method", "$.ajax", false},

		// Generic callback variable name
		{"fn callback", "fn", true},
		{"fn prefixed", "fnCallback", false},

		// Single character identifiers
		{"single valid identifier", "a", false},
		{"single valid identifier upper", "A", false},
		{"single valid identifier digit", "1", false},
		{"single valid identifier underscore", "_", false},
		{"single invalid identifier space", " ", true},
		{"single invalid identifier dot", ".", true},
		{"single invalid identifier bracket", "[", true},

		// Normal identifiers
		{"normal identifier", "doSomething", false},
		{"dotted identifier", "service.method", false},
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

func TestIsValidIdentifierChar(t *testing.T) {
	tests := []struct {
		name  string
		char  rune
		valid bool
	}{
		{"lowercase", 'a', true},
		{"uppercase", 'Z', true},
		{"digit", '5', true},
		{"underscore", '_', true},
		{"dollar", '$', false},
		{"space", ' ', false},
		{"dot", '.', false},
		{"bracket", '[', false},
		{"paren", '(', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidIdentifierChar(tt.char)
			if got != tt.valid {
				t.Errorf("IsValidIdentifierChar(%q) = %v, want %v", tt.char, got, tt.valid)
			}
		})
	}
}
