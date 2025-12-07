// Package domain defines the core types for test file representation.
package domain

// Language represents a programming language.
type Language string

// Supported languages for test file parsing.
const (
	LanguageGo         Language = "go"
	LanguageJavaScript Language = "javascript"
	LanguagePython     Language = "python"
	LanguageTypeScript Language = "typescript"
)
