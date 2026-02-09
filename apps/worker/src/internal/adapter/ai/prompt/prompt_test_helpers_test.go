package prompt

import "strings"

// estimateTokenCount provides a rough token estimate for testing purposes.
// Uses ~1.3 tokens per word heuristic which is approximate for Gemini.
// Real token count may vary by ±20%. For production, use actual tokenizer.
func estimateTokenCount(text string) int {
	words := strings.Fields(text)
	return len(words) * 4 / 3
}
