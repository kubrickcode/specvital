package domain_hints

import (
	"context"
	"testing"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
)

func TestJavaExtractor_Extract_Imports(t *testing.T) {
	source := []byte(`
package com.example.project;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.params.ParameterizedTest;
import static org.junit.jupiter.api.Assertions.assertEquals;
import com.example.service.*;

class CalculatorTests {
}
`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"org.junit.jupiter.api.Test":                    true,
		"org.junit.jupiter.api.DisplayName":             true,
		"org.junit.jupiter.params.ParameterizedTest":    true,
		"org.junit.jupiter.api.Assertions.assertEquals": true,
		"com.example.service.*":                         true,
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

func TestJavaExtractor_Extract_Calls(t *testing.T) {
	source := []byte(`
package com.example.project;

import org.junit.jupiter.api.Test;

class CalculatorTests {
	@Test
	void testAdd() {
		Calculator calculator = new Calculator();
		int result = calculator.add(1, 2);
		userService.findById(123);
		paymentGateway.process(order);
	}
}
`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"calculator.add":         true,
		"userService.findById":   true,
		"paymentGateway.process": true,
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

func TestJavaExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`// empty file`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestJavaExtractor_Extract_TestFrameworkCalls(t *testing.T) {
	source := []byte(`
package com.example.project;

import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;

class CalculatorTests {
	@Test
	void testAdd() {
		Calculator calculator = new Calculator();
		assertEquals(2, calculator.add(1, 1));
		assertTrue(calculator.isPositive(5));
		assertThrows(Exception.class, () -> calculator.divide(1, 0));
		userService.validate(user);
	}
}
`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"assertEquals", "assertTrue", "assertThrows"}
	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected test framework call %q to be excluded", call)
		}
	}

	// Domain calls should be included
	if !callSet["calculator.add"] {
		t.Errorf("expected calculator.add call, got %v", hints.Calls)
	}
	if !callSet["userService.validate"] {
		t.Errorf("expected userService.validate call, got %v", hints.Calls)
	}
}

func TestJavaExtractor_Extract_StaticImport(t *testing.T) {
	source := []byte(`
package com.example;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.hamcrest.MatcherAssert.assertThat;

class Test {}
`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	if !importSet["org.junit.jupiter.api.Assertions.assertEquals"] {
		t.Errorf("expected static import, got %v", hints.Imports)
	}
	if !importSet["org.hamcrest.MatcherAssert.assertThat"] {
		t.Errorf("expected static import, got %v", hints.Imports)
	}
}

func TestJavaExtractor_Extract_ChainedCalls(t *testing.T) {
	source := []byte(`
package com.example;

class Test {
	void test() {
		// Long chains should be normalized to 2 segments
		client.api.users.findAll();
		response.data.items.first().getValue();
	}
}
`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Should be normalized to 2 segments
	expectedCalls := []string{"client.api", "response.data"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected %q call (2-segment normalized), got %v", call, hints.Calls)
		}
	}
}

func TestGetExtractor_Java(t *testing.T) {
	ext := GetExtractor(domain.LanguageJava)
	if ext == nil {
		t.Error("expected extractor for Java, got nil")
	}

	_, ok := ext.(*JavaExtractor)
	if !ok {
		t.Errorf("expected JavaExtractor, got %T", ext)
	}
}

func TestJavaExtractor_Extract_NoiseFilter(t *testing.T) {
	source := []byte(`
package com.example;

class Test {
	void test() {
		// Decimal literal patterns should be filtered
		double x = (0.5).doubleValue();
		double y = (1.0).toString();

		// Real domain calls should be included
		userService.create(data);
		paymentGateway.process(amount);
	}
}
`)

	extractor := &JavaExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("noise patterns filtered", func(t *testing.T) {
		for call := range callSet {
			if len(call) > 0 && (call[0] == '[' || call[0] == '(') {
				t.Errorf("noise pattern %q should be filtered, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"userService.create", "paymentGateway.process"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}
