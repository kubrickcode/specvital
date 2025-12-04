// Package domain defines the core types for test file representation.
package domain

// Language represents a programming language.
type Language string

// Supported languages for test file parsing.
const (
	LanguageTypeScript Language = "typescript"
	LanguageJavaScript Language = "javascript"
	LanguageGo         Language = "go"
)
