package domain

// Location represents a source code location range.
type Location struct {
	// EndCol is the ending column (0-based).
	EndCol int `json:"endCol,omitempty"`
	// EndLine is the ending line number (1-based).
	EndLine int `json:"endLine"`
	// File is the file path.
	File string `json:"file"`
	// StartCol is the starting column (0-based).
	StartCol int `json:"startCol,omitempty"`
	// StartLine is the starting line number (1-based).
	StartLine int `json:"startLine"`
}
