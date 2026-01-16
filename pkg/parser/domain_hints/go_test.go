package domain_hints

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/domain"
)

func TestGoExtractor_Extract(t *testing.T) {
	source := []byte(`package order

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"myapp/repository"
	"myapp/services/inventory"
)

func TestCreateOrder(t *testing.T) {
	mockCart := Cart{Items: []Item{{ID: 1, Qty: 2}}}

	t.Run("should create order from cart", func(t *testing.T) {
		result, err := orderService.CreateFromCart(mockCart)
		assert.NoError(t, err)
		assert.Equal(t, "pending", result.Status)
	})
}
`)

	extractor := &GoExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	t.Run("imports", func(t *testing.T) {
		// "testing" is filtered as Go stdlib
		expectedImports := []string{
			"github.com/stretchr/testify/assert",
			"myapp/repository",
			"myapp/services/inventory",
		}
		if len(hints.Imports) != len(expectedImports) {
			t.Errorf("imports count: got %d, want %d", len(hints.Imports), len(expectedImports))
		}
		for i, expected := range expectedImports {
			if i >= len(hints.Imports) {
				break
			}
			if hints.Imports[i] != expected {
				t.Errorf("imports[%d]: got %q, want %q", i, hints.Imports[i], expected)
			}
		}
	})
}

func TestGoExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`package empty`)

	extractor := &GoExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestGoExtractor_Extract_Calls(t *testing.T) {
	source := []byte(`package test

import "testing"

func TestSomething(t *testing.T) {
	authService.ValidateToken("token")
	userRepo.FindByID(1)
	result, err := orderService.Create(order)
	doSomething()
}
`)

	extractor := &GoExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"authService.ValidateToken": true,
		"userRepo.FindByID":         true,
		"orderService.Create":       true,
		"doSomething":               true,
	}

	for _, call := range hints.Calls {
		delete(expectedCalls, call)
	}

	if len(expectedCalls) > 0 {
		t.Errorf("missing calls: %v, got: %v", expectedCalls, hints.Calls)
	}
}

func TestGetExtractor(t *testing.T) {
	tests := []struct {
		lang    domain.Language
		wantNil bool
	}{
		{domain.LanguageGo, false},
		{domain.LanguageJavaScript, false},
		{domain.LanguageTypeScript, false},
		{domain.LanguageTSX, false},
		{domain.LanguagePython, false},
		{domain.LanguageJava, false},
		{domain.LanguageKotlin, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			ext := GetExtractor(tt.lang)
			if tt.wantNil && ext != nil {
				t.Errorf("expected nil extractor for %s", tt.lang)
			}
			if !tt.wantNil && ext == nil {
				t.Errorf("expected extractor for %s, got nil", tt.lang)
			}
		})
	}
}


func TestGoExtractor_Extract_StdlibFiltering(t *testing.T) {
	source := []byte(`package test

import (
	"io"
	"os"
	"fs"
	"fmt"
	"context"
	"time"
	"encoding/json"
	"net/http"
	"testing"
	"github.com/stretchr/testify/require"
	"myapp/domain/order"
	"myapp/services/payment"
)

func TestSomething(t *testing.T) {
	require.NoError(t, nil)
}
`)

	extractor := &GoExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	t.Run("stdlib imports filtered", func(t *testing.T) {
		stdlibPackages := []string{"io", "os", "fs", "fmt", "context", "time", "encoding/json", "net/http", "testing"}
		for _, stdlib := range stdlibPackages {
			for _, imp := range hints.Imports {
				if imp == stdlib {
					t.Errorf("stdlib import %q should be filtered", stdlib)
				}
			}
		}
	})

	t.Run("domain imports kept", func(t *testing.T) {
		domainImports := map[string]bool{
			"github.com/stretchr/testify/require": false,
			"myapp/domain/order":                  false,
			"myapp/services/payment":              false,
		}
		for _, imp := range hints.Imports {
			if _, exists := domainImports[imp]; exists {
				domainImports[imp] = true
			}
		}
		for imp, found := range domainImports {
			if !found {
				t.Errorf("domain import %q should be kept, got imports: %v", imp, hints.Imports)
			}
		}
	})
}

func TestIsGoStdlibImport(t *testing.T) {
	tests := []struct {
		importPath string
		wantFilter bool
	}{
		// Direct stdlib packages
		{"io", true},
		{"os", true},
		{"fs", true},
		{"fmt", true},
		{"context", true},
		{"time", true},
		{"testing", true},
		{"errors", true},
		{"strings", true},
		{"bytes", true},
		// Nested stdlib packages
		{"encoding/json", true},
		{"encoding/xml", true},
		{"encoding/gob", true},
		{"net/http", true},
		{"net/url", true},
		{"crypto/sha256", true},
		{"io/fs", true},
		{"io/ioutil", true},
		// Non-stdlib (should NOT be filtered)
		{"github.com/stretchr/testify/require", false},
		{"github.com/stretchr/testify/assert", false},
		{"myapp/domain/order", false},
		{"myapp/services/payment", false},
		{"golang.org/x/sync/errgroup", false},
		{"google.golang.org/grpc", false},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			got := isGoStdlibImport(tt.importPath)
			if got != tt.wantFilter {
				t.Errorf("isGoStdlibImport(%q) = %v, want %v", tt.importPath, got, tt.wantFilter)
			}
		})
	}
}

func TestGoExtractor_Extract_NoiseFiltering(t *testing.T) {
	source := []byte(`package test

import "testing"

func TestSpread(t *testing.T) {
	// This would produce "[." pattern if not filtered
	result := []int{1, 2}
	expanded := append([]int{}, result...)
	doSomething()
}
`)

	extractor := &GoExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	t.Run("no malformed identifiers", func(t *testing.T) {
		for _, call := range hints.Calls {
			if len(call) > 0 && call[0] == '[' {
				t.Errorf("found malformed call: %q", call)
			}
		}
	})

	t.Run("no empty strings", func(t *testing.T) {
		for _, call := range hints.Calls {
			if call == "" {
				t.Error("found empty call")
			}
		}
	})
}
