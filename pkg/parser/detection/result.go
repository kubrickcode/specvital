package detection

import "fmt"

// Confidence thresholds for detection quality
const (
	// ConfidenceDefinite indicates definite framework match (71-100).
	// No fallback detection needed.
	ConfidenceDefinite = 71

	// ConfidenceModerate indicates likely framework match (31-70).
	// Should warn user if conflicts occur.
	ConfidenceModerate = 31

	// ConfidenceWeak indicates weak framework match (1-30).
	// Should try fallback detection strategies.
	ConfidenceWeak = 1
)

// Result represents the outcome of framework detection for a test file.
type Result struct {
	// Framework is the detected framework name (e.g., "jest", "vitest").
	// Empty string if no framework detected.
	Framework string

	// Confidence is the total confidence score (0-100).
	// Higher scores indicate more certain detection.
	Confidence int

	// Evidence is a list of detection signals that contributed to this result.
	Evidence []Evidence

	// Scope is the config scope that applies to this file (if scope-based detection succeeded).
	// May be nil if no config scope applies.
	Scope interface{} // framework.ConfigScope, but avoid import cycle
}

// Evidence represents a single detection signal.
type Evidence struct {
	// Source identifies where this evidence came from.
	// Values: "config-scope", "globals-mode", "import", "content", "filename"
	Source string

	// Description provides human-readable details about this evidence.
	Description string

	// Confidence is the confidence points contributed by this evidence.
	Confidence int

	// Negative indicates this evidence rules out the framework.
	// For example: file imports from a different framework.
	Negative bool
}

func (r Result) IsDefinite() bool {
	return r.Confidence >= ConfidenceDefinite
}

func (r Result) IsModerate() bool {
	return r.Confidence >= ConfidenceModerate && r.Confidence < ConfidenceDefinite
}

func (r Result) IsWeak() bool {
	return r.Confidence >= ConfidenceWeak && r.Confidence < ConfidenceModerate
}

func (r Result) ConfidenceLevel() string {
	switch {
	case r.Confidence >= ConfidenceDefinite:
		return "definite"
	case r.Confidence >= ConfidenceModerate:
		return "moderate"
	case r.Confidence >= ConfidenceWeak:
		return "weak"
	default:
		return "none"
	}
}

func (r Result) String() string {
	if r.Framework == "" {
		return "no framework detected"
	}
	return fmt.Sprintf("%s (confidence: %d/%s)", r.Framework, r.Confidence, r.ConfidenceLevel())
}

func Unknown() Result {
	return Result{
		Framework:  "",
		Confidence: 0,
		Evidence:   nil,
	}
}

func (r *Result) AddEvidence(ev Evidence) {
	if r.Evidence == nil {
		r.Evidence = make([]Evidence, 0, 4)
	}
	r.Evidence = append(r.Evidence, ev)
}
