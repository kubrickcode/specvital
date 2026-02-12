package domain_hints

import (
	"context"
	"strings"
	"testing"

	"github.com/kubrickcode/specvital/lib/parser/domain"
)

func TestRustExtractor_Extract_UseStatements(t *testing.T) {
	source := []byte(`
use std::collections::HashMap;
use crate::models::User;
use super::helpers;
use tokio::sync::mpsc;
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"std/collections/HashMap": true,
		"crate/models/User":       true,
		"super/helpers":           true,
		"tokio/sync/mpsc":         true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestRustExtractor_Extract_UseList(t *testing.T) {
	source := []byte(`
use std::collections::{HashMap, HashSet};
use crate::{models, services};
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Should extract the base path before the list
	expectedImports := map[string]bool{
		"std/collections": true,
		"crate":           true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestRustExtractor_Extract_UseWildcard(t *testing.T) {
	source := []byte(`
use std::prelude::*;
use crate::models::*;
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"std/prelude":  true,
		"crate/models": true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestRustExtractor_Extract_UseAlias(t *testing.T) {
	source := []byte(`
use std::collections::HashMap as Map;
use crate::models::User as UserModel;
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"std/collections/HashMap": true,
		"crate/models/User":       true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestRustExtractor_Extract_ModDeclarations(t *testing.T) {
	source := []byte(`
mod tests;
mod helpers;
pub mod utils;
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"tests":   true,
		"helpers": true,
		"utils":   true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected mod %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestRustExtractor_Extract_MethodCalls(t *testing.T) {
	source := []byte(`
use std::collections::HashMap;

fn test_service() {
    user_service.create(user);
    PaymentGateway::process(payment);
    notification_service.send_email(user);
}
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"user_service.create":             true,
		"PaymentGateway.process":          true,
		"notification_service.send_email": true,
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	for call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q to be included, got %v", call, hints.Calls)
		}
	}
}

func TestRustExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`// empty file`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestRustExtractor_Extract_TestFrameworkCalls(t *testing.T) {
	source := []byte(`
use crate::services::payment;

#[test]
fn test_payment() {
    assert_eq!(result, expected);
    assert!(condition);
    println!("debug output");

    payment_service.process(order);
    Result::Ok(value);
}
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"assert_eq", "assert", "println"}
	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected test framework call %q to be excluded, got %v", call, hints.Calls)
		}
	}

	// Domain calls should be included
	if !callSet["payment_service.process"] {
		t.Errorf("expected payment_service.process call, got %v", hints.Calls)
	}
}

func TestRustExtractor_Extract_CargoTestFile(t *testing.T) {
	source := []byte(`
use crate::models::Order;
use crate::services::payment::PaymentGateway;
use tokio::sync::mpsc;

mod test_helpers;

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_payment_processing() {
        let gateway = PaymentGateway::new();
        let order = Order::create(100);

        gateway.process(order).await;
        notification_service.send_confirmation(order.id);
    }
}
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Verify imports
	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	expectedImports := []string{
		"crate/models/Order",
		"crate/services/payment/PaymentGateway",
		"tokio/sync/mpsc",
		"test_helpers",
	}
	for _, imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q, got %v", imp, hints.Imports)
		}
	}

	// Verify calls
	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	expectedCalls := []string{
		"PaymentGateway.new",
		"Order.create",
		"gateway.process",
		"notification_service.send_confirmation",
	}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q, got %v", call, hints.Calls)
		}
	}
}

func TestRustExtractor_Extract_Deduplication(t *testing.T) {
	source := []byte(`
use std::collections::HashMap;
use std::collections::HashMap;

fn test() {
    user_service.create(1);
    user_service.create(2);
}
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Count occurrences
	importCount := 0
	for _, imp := range hints.Imports {
		if imp == "std/collections/HashMap" {
			importCount++
		}
	}
	if importCount != 1 {
		t.Errorf("expected 'std/collections/HashMap' to appear once, got %d times in %v", importCount, hints.Imports)
	}

	callCount := 0
	for _, call := range hints.Calls {
		if call == "user_service.create" {
			callCount++
		}
	}
	if callCount != 1 {
		t.Errorf("expected 'user_service.create' to appear once, got %d times in %v", callCount, hints.Calls)
	}
}

func TestGetExtractor_Rust(t *testing.T) {
	ext := GetExtractor(domain.LanguageRust)
	if ext == nil {
		t.Error("expected extractor for Rust, got nil")
	}

	_, ok := ext.(*RustExtractor)
	if !ok {
		t.Errorf("expected RustExtractor, got %T", ext)
	}
}

func TestRustExtractor_Extract_StdlibEnumFiltering(t *testing.T) {
	source := []byte(`
#[test]
fn test_result_handling() {
	let result = compute();
	match result {
		Ok(value) => {
			process(value);
		}
		Err(error) => {
			handle_error(error);
		}
	}

	let maybe = find();
	match maybe {
		Some(item) => use_item(item),
		None => set_default(),
	}
}
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	t.Run("stdlib enums filtered", func(t *testing.T) {
		callSet := make(map[string]bool)
		for _, call := range hints.Calls {
			callSet[call] = true
		}

		// Stdlib enum constructors should be excluded
		excluded := []string{"Ok", "Err", "Some", "None"}
		for _, call := range excluded {
			if callSet[call] {
				t.Errorf("expected stdlib enum %q to be excluded, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		callSet := make(map[string]bool)
		for _, call := range hints.Calls {
			callSet[call] = true
		}

		// Domain calls should be included
		expectedCalls := []string{"process", "handle_error", "use_item", "set_default"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestRustExtractor_Extract_NoiseFilter(t *testing.T) {
	source := []byte(`
fn test_decimal_literals() {
    // Decimal literals should not be extracted as calls
    let x = (0.5).abs();
    let y = (1.0).to_string();
    let z = (123.456).floor();

    // Real domain calls should be included
    math_service.calculate(x);
    number_utils.format(y);
}
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("decimal literal noise filtered", func(t *testing.T) {
		// (0.5).abs() should not produce "(0." or similar noise
		for call := range callSet {
			if strings.HasPrefix(call, "(") {
				t.Errorf("decimal literal noise %q should be filtered, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"math_service.calculate", "number_utils.format"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestParseRustUseList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"simple list",
			"{ anyhow::Context, bstr::ByteVec }",
			[]string{"anyhow/Context", "bstr/ByteVec"},
		},
		{
			"multiline list",
			"{\n    grep_matcher,\n    bstr,\n}",
			[]string{"grep_matcher", "bstr"},
		},
		{
			"list with alias",
			"{ anyhow::Context as Ctx, bstr::ByteVec }",
			[]string{"anyhow/Context", "bstr/ByteVec"},
		},
		{
			"single item",
			"{ single::item }",
			[]string{"single/item"},
		},
		{
			"nested path",
			"{ std::collections::HashMap, tokio::sync::mpsc }",
			[]string{"std/collections/HashMap", "tokio/sync/mpsc"},
		},
		{
			"with trailing comma",
			"{ item1, item2, }",
			[]string{"item1", "item2"},
		},
		{
			"nested braces",
			"{ regex_automata::{PatternSet, meta::Regex}, bstr::ByteSlice }",
			[]string{"regex_automata", "bstr/ByteSlice"},
		},
		{
			"deeply nested multiline",
			"{\n    aho_corasick::AhoCorasick,\n    regex_automata::{\n        PatternSet,\n        meta::Regex,\n    },\n}",
			[]string{"aho_corasick/AhoCorasick", "regex_automata"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRustUseList(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("parseRustUseList(%q) = %v (len=%d), want %v (len=%d)",
					tt.input, got, len(got), tt.want, len(tt.want))
				return
			}
			for i, g := range got {
				if g != tt.want[i] {
					t.Errorf("parseRustUseList(%q)[%d] = %q, want %q", tt.input, i, g, tt.want[i])
				}
			}
		})
	}
}

func TestRustExtractor_Extract_DirectUseList(t *testing.T) {
	source := []byte(`
use {anyhow::Context, bstr::ByteVec};
use {
    grep_matcher::LineTerminator,
    grep_searcher::SearcherBuilder,
};
`)

	extractor := &RustExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"anyhow/Context":                true,
		"bstr/ByteVec":                  true,
		"grep_matcher/LineTerminator":   true,
		"grep_searcher/SearcherBuilder": true,
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	for imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected import %q to be included, got %v", imp, hints.Imports)
		}
	}

	// Verify no noise patterns (imports starting with "{")
	for _, imp := range hints.Imports {
		if strings.HasPrefix(imp, "{") {
			t.Errorf("found noise pattern starting with '{': %q", imp)
		}
	}
}
