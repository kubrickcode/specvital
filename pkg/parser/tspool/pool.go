// Package tspool provides pooled tree-sitter parsers for concurrent parsing.
//
// This package centralizes parser pooling logic to:
//   - Enable reuse across different parser components
//   - Reduce parser allocation overhead via sync.Pool
//   - Ensure thread-safe parser sharing
//
// Separate pools are maintained per language (Go, JavaScript, TypeScript)
// to avoid language switching overhead.
//
// Thread-safety: Parsers returned by Get are NOT safe for concurrent use.
// Each goroutine must Get its own parser or use the Parse helper.
package tspool

import (
	"context"
	"fmt"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"

	"github.com/specvital/core/pkg/domain"
)

// MaxTreeDepth is the maximum recursion depth when walking AST trees.
const MaxTreeDepth = 1000

var (
	goLang   *sitter.Language
	javaLang *sitter.Language
	jsLang   *sitter.Language
	pyLang   *sitter.Language
	tsLang   *sitter.Language

	langOnce sync.Once
)

func initLanguages() {
	langOnce.Do(func() {
		goLang = golang.GetLanguage()
		javaLang = java.GetLanguage()
		jsLang = javascript.GetLanguage()
		pyLang = python.GetLanguage()
		tsLang = typescript.GetLanguage()
	})
}

// GetLanguage returns the tree-sitter language for the given domain language.
func GetLanguage(lang domain.Language) *sitter.Language {
	initLanguages()
	switch lang {
	case domain.LanguageGo:
		return goLang
	case domain.LanguageJava:
		return javaLang
	case domain.LanguageJavaScript:
		return jsLang
	case domain.LanguagePython:
		return pyLang
	default:
		return tsLang
	}
}

var (
	goParserPool   sync.Pool
	javaParserPool sync.Pool
	jsParserPool   sync.Pool
	pyParserPool   sync.Pool
	tsParserPool   sync.Pool
)

func getParserPool(lang domain.Language) *sync.Pool {
	switch lang {
	case domain.LanguageGo:
		return &goParserPool
	case domain.LanguageJava:
		return &javaParserPool
	case domain.LanguageJavaScript:
		return &jsParserPool
	case domain.LanguagePython:
		return &pyParserPool
	default:
		return &tsParserPool
	}
}

// Get returns a pooled parser for the given language.
// The returned parser is NOT safe for concurrent use.
// Use Put to return the parser after use.
func Get(lang domain.Language) *sitter.Parser {
	pool := getParserPool(lang)

	if p := pool.Get(); p != nil {
		if parser, ok := p.(*sitter.Parser); ok {
			return parser
		}
	}

	initLanguages()
	parser := sitter.NewParser()
	parser.SetLanguage(GetLanguage(lang))
	return parser
}

// Put returns a parser to the pool.
func Put(lang domain.Language, parser *sitter.Parser) {
	if parser == nil {
		return
	}
	pool := getParserPool(lang)
	pool.Put(parser)
}

// Parse parses source using a pooled parser.
// Caller MUST call tree.Close() to free resources.
// The parser is automatically returned to the pool after parsing.
func Parse(ctx context.Context, lang domain.Language, source []byte) (*sitter.Tree, error) {
	parser := Get(lang)
	defer Put(lang, parser)

	tree, err := parser.ParseCtx(ctx, nil, source)
	if err != nil {
		return nil, fmt.Errorf("parse %s failed: %w", lang, err)
	}

	return tree, nil
}
