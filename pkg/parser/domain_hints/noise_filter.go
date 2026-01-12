package domain_hints

// ShouldFilterNoise filters out universal noise patterns from domain hints.
// Removes: empty strings, malformed identifiers, and framework-specific noise.
func ShouldFilterNoise(call string) bool {
	if call == "" {
		return true
	}

	// Malformed patterns: starts with "[" (e.g., "[." from spread array handling)
	if call[0] == '[' {
		return true
	}

	// Malformed patterns: starts with "(" (e.g., "(0.", "(1." from decimal literals)
	if call[0] == '(' {
		return true
	}

	// Cheerio/jQuery selector: single "$" is noise
	if call == "$" {
		return true
	}

	// Single character not matching valid identifier pattern
	if len(call) == 1 {
		return !IsValidIdentifierChar(rune(call[0]))
	}

	return false
}

// IsValidIdentifierChar checks if a rune is a valid identifier character.
// Valid: A-Z, a-z, 0-9, underscore
func IsValidIdentifierChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_'
}
