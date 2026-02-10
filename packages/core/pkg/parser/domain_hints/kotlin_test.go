package domain_hints

import (
	"context"
	"strings"
	"testing"

	"github.com/kubrickcode/specvital/packages/core/pkg/domain"
)

func TestKotlinExtractor_Extract_Imports(t *testing.T) {
	source := []byte(`
package kotest

import io.kotest.core.spec.style.StringSpec
import io.kotest.matchers.shouldBe
import com.example.service.UserService
import org.junit.jupiter.api.Test

class KotestSpec : StringSpec({
})
`)

	extractor := &KotlinExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"io.kotest.core.spec.style.StringSpec": true,
		"io.kotest.matchers.shouldBe":          true,
		"com.example.service.UserService":      true,
		"org.junit.jupiter.api.Test":           true,
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

func TestKotlinExtractor_Extract_Calls(t *testing.T) {
	source := []byte(`
package com.example

import io.kotest.core.spec.style.FunSpec

class CalculatorTest : FunSpec({
    test("add two numbers") {
        val calculator = Calculator()
        val result = calculator.add(1, 2)
        userService.findById(123)
        paymentGateway.process(order)
    }
})
`)

	extractor := &KotlinExtractor{}
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

func TestKotlinExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`// empty file`)

	extractor := &KotlinExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestKotlinExtractor_Extract_TestFrameworkCalls(t *testing.T) {
	source := []byte(`
package com.example

import io.kotest.core.spec.style.FunSpec
import io.kotest.matchers.shouldBe

class CalculatorTest : FunSpec({
    test("add two numbers") {
        val result = calculator.add(1, 2)
        result shouldBe 3
        userService.validate(user)
    }
})
`)

	extractor := &KotlinExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"shouldBe", "test"}
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

func TestKotlinExtractor_Extract_KotestSpec(t *testing.T) {
	source := []byte(`
package kotest

import io.kotest.core.spec.style.StringSpec
import io.kotest.matchers.shouldBe
import com.example.service.PaymentService

class PaymentSpec : StringSpec({
    "payment should be processed" {
        val service = PaymentService()
        val result = service.process(order)
        stripe.confirm(result.intentId)
        result.status shouldBe "success"
    }
})
`)

	extractor := &KotlinExtractor{}
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
		"io.kotest.core.spec.style.StringSpec",
		"io.kotest.matchers.shouldBe",
		"com.example.service.PaymentService",
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

	expectedCalls := []string{"service.process", "stripe.confirm"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q, got %v", call, hints.Calls)
		}
	}
}

func TestKotlinExtractor_Extract_ChainedCalls(t *testing.T) {
	source := []byte(`
package com.example

class Test {
    fun test() {
        // Long chains should be normalized to 2 segments
        client.api.users.findAll()
        response.data.items.first().value
    }
}
`)

	extractor := &KotlinExtractor{}
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

func TestGetExtractor_Kotlin(t *testing.T) {
	ext := GetExtractor(domain.LanguageKotlin)
	if ext == nil {
		t.Error("expected extractor for Kotlin, got nil")
	}

	_, ok := ext.(*KotlinExtractor)
	if !ok {
		t.Errorf("expected KotlinExtractor, got %T", ext)
	}
}

func TestKotlinExtractor_Extract_NoiseFilter(t *testing.T) {
	source := []byte(`
package com.example

class DecimalTest {
    fun test() {
        // Decimal literals should not be extracted as calls
        val x = (0.5).toInt()
        val y = (1.0).toString()
        val z = (123.456).roundToInt()

        // Real domain calls should be included
        mathService.calculate(x)
        numberUtils.format(y)
    }
}
`)

	extractor := &KotlinExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("decimal literal noise filtered", func(t *testing.T) {
		// (0.5).toInt() should not produce "(0." or similar noise
		for call := range callSet {
			if strings.HasPrefix(call, "(") {
				t.Errorf("decimal literal noise %q should be filtered, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"mathService.calculate", "numberUtils.format"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestKotlinExtractor_Extract_StdlibFilter(t *testing.T) {
	source := []byte(`
package com.example

class CollectionTest {
    fun test() {
        // Kotlin stdlib functions should be filtered (no domain signal)
        val list = listOf(1, 2, 3)
        val set = setOf("a", "b")
        val map = mapOf("key" to "value")
        val empty = emptyList<String>()
        val pair = Pair("first", "second")
        error("something went wrong")
        require(list.isNotEmpty())
        check(set.size > 0)

        // Domain calls should be included
        userService.findAll()
        orderRepository.save(order)
    }
}
`)

	extractor := &KotlinExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("stdlib functions filtered", func(t *testing.T) {
		excludedCalls := []string{
			"listOf", "setOf", "mapOf", "emptyList", "Pair", "error", "require", "check",
		}
		for _, call := range excludedCalls {
			if callSet[call] {
				t.Errorf("expected stdlib call %q to be excluded, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"userService.findAll", "orderRepository.save"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}
