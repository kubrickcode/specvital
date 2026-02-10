// Package framework provides a unified framework definition system for test framework detection and parsing.
// It replaces the dual registry system (matchers + strategies) with a single unified Definition type.
package framework

import (
	"context"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
)

// Definition represents a unified test framework specification that combines:
// - Framework identity (Name, Languages)
// - Detection components (Matchers)
// - Configuration parsing (ConfigParser)
// - Test file parsing (Parser)
// - Priority for detection ordering
//
// Each framework (Jest, Vitest, Playwright, Go testing) provides a Definition
// that is registered with the global Registry.
type Definition struct {
	// Name is the framework identifier (e.g., "jest", "vitest", "playwright").
	Name string

	// Languages specifies which programming languages this framework supports.
	Languages []domain.Language

	// Matchers are detection rules that determine if a test file uses this framework.
	// Multiple matchers can be registered for different detection strategies
	// (import statements, config files, file content patterns, etc.).
	Matchers []Matcher

	// ConfigParser extracts configuration information from framework config files.
	// Used to determine settings like globals mode, test patterns, etc.
	// May be nil if the framework doesn't have configuration files.
	ConfigParser ConfigParser

	// Parser extracts test definitions from source code files.
	// Produces a domain.TestFile with all test suites and test cases.
	Parser Parser

	// Priority determines detection order when multiple frameworks could match.
	// Higher priority frameworks are checked first.
	// Use PriorityGeneric (100), PriorityE2E (150), or PrioritySpecialized (200).
	Priority int
}

// Matcher defines the interface for framework detection rules.
// Matchers analyze different signals (imports, config files, content patterns)
// to determine if a test file belongs to a specific framework.
type Matcher interface {
	// Match evaluates a signal and returns a confidence score.
	// Returns MatchResult with confidence (0-100) and supporting evidence.
	Match(ctx context.Context, signal Signal) MatchResult
}

// Signal represents a detection signal that matchers can evaluate.
type Signal struct {
	// Type indicates what kind of signal this is (import, config file, etc.).
	Type SignalType

	// Value contains the signal data (import path, file name, content, etc.).
	Value string

	// Context provides additional metadata specific to the signal type.
	// For example, file content for content-based matching.
	Context interface{}
}

// SignalType categorizes different kinds of detection signals.
type SignalType int

const (
	// SignalImport represents an import statement (e.g., "import { test } from 'vitest'").
	SignalImport SignalType = iota

	// SignalConfigFile represents a config file name (e.g., "jest.config.js").
	SignalConfigFile

	// SignalFileContent represents file content patterns (e.g., "test.describe(").
	SignalFileContent

	// SignalFileName represents a test file name pattern (e.g., "*.test.ts").
	SignalFileName
)

// MatchResult contains the outcome of a matcher evaluation.
type MatchResult struct {
	// Confidence is a score from 0-100 indicating how certain the match is.
	// 0 = no match, 100 = definite match.
	Confidence int

	// Evidence is a list of specific indicators that support this match.
	// For example: ["import '@jest/globals'", "jest.config.js found"].
	Evidence []string

	// Negative indicates this is a definite non-match.
	// If true, this framework should be excluded from consideration.
	// For example, Vitest with globals:false should not match files without imports.
	Negative bool
}

// ConfigParser extracts configuration information from framework config files.
type ConfigParser interface {
	// Parse reads and interprets a framework configuration file.
	// Returns ConfigScope containing parsed settings like test patterns, globals mode, etc.
	// Returns error if the config file cannot be parsed.
	Parse(ctx context.Context, configPath string, content []byte) (*ConfigScope, error)
}

// Parser extracts test definitions from source code files.
type Parser interface {
	// Parse analyzes source code and extracts test suites and test cases.
	// Returns a domain.TestFile containing all discovered tests.
	// Returns error if the file cannot be parsed or doesn't contain valid tests.
	Parse(ctx context.Context, source []byte, filename string) (*domain.TestFile, error)
}

// NoMatch returns a MatchResult indicating no match was found.
func NoMatch() MatchResult {
	return MatchResult{Confidence: 0}
}

// DefiniteMatch returns a MatchResult indicating a definite match with evidence.
func DefiniteMatch(evidence ...string) MatchResult {
	return MatchResult{
		Confidence: 100,
		Evidence:   evidence,
	}
}

// PartialMatch returns a MatchResult with a specific confidence level and evidence.
func PartialMatch(confidence int, evidence ...string) MatchResult {
	return MatchResult{
		Confidence: confidence,
		Evidence:   evidence,
	}
}

// NegativeMatch returns a MatchResult indicating this framework definitely does not match.
func NegativeMatch(evidence ...string) MatchResult {
	return MatchResult{
		Confidence: 0,
		Evidence:   evidence,
		Negative:   true,
	}
}
