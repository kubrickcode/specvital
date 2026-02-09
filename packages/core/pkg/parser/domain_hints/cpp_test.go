package domain_hints

import (
	"context"
	"testing"

	"github.com/specvital/core/pkg/domain"
)

func TestCppExtractor_Extract_IncludeStatements(t *testing.T) {
	source := []byte(`
#include <iostream>
#include <vector>
#include "myheader.h"
#include <gtest/gtest.h>
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Note: iostream and vector are stdlib headers and should be filtered
	expectedImports := map[string]bool{
		"myheader.h":    true,
		"gtest/gtest.h": true,
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

	// Verify stdlib headers are filtered
	filteredHeaders := []string{"iostream", "vector"}
	for _, stdlib := range filteredHeaders {
		if importSet[stdlib] {
			t.Errorf("stdlib header %q should be filtered", stdlib)
		}
	}
}

func TestCppExtractor_Extract_SystemHeaders(t *testing.T) {
	// Test that stdlib-only source returns nil (no domain-relevant imports)
	source := []byte(`
#include <string>
#include <memory>
#include <algorithm>
#include <functional>
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	// All these are stdlib headers, so hints should be nil (no domain signal)
	if hints != nil {
		t.Errorf("expected nil for stdlib-only file, got imports=%v, calls=%v", hints.Imports, hints.Calls)
	}
}

func TestCppExtractor_Extract_MixedHeaders(t *testing.T) {
	// Test mix of stdlib and domain headers
	source := []byte(`
#include <string>
#include <memory>
#include "domain/service.h"
#include <boost/asio.hpp>
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	// Stdlib should be filtered
	stdlibHeaders := []string{"string", "memory"}
	for _, stdlib := range stdlibHeaders {
		if importSet[stdlib] {
			t.Errorf("stdlib header %q should be filtered", stdlib)
		}
	}

	// Domain headers should be included
	domainHeaders := []string{"domain/service.h", "boost/asio.hpp"}
	for _, domain := range domainHeaders {
		if !importSet[domain] {
			t.Errorf("domain header %q should be included, got %v", domain, hints.Imports)
		}
	}
}

func TestCppExtractor_Extract_LocalHeaders(t *testing.T) {
	source := []byte(`
#include "services/payment.h"
#include "models/user.h"
#include "../common/utils.h"
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"services/payment.h": true,
		"models/user.h":      true,
		"../common/utils.h":  true,
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

func TestCppExtractor_Extract_MethodCalls(t *testing.T) {
	source := []byte(`
#include <iostream>

void testFunction() {
    userService.create(user);
    PaymentGateway::process(payment);
    notificationService->sendEmail(user);
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"userService.create":            true,
		"PaymentGateway.process":        true,
		"notificationService.sendEmail": true,
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

func TestCppExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`// empty file`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestCppExtractor_Extract_TestFrameworkCalls(t *testing.T) {
	source := []byte(`
#include <gtest/gtest.h>

TEST(PaymentTest, ProcessPayment) {
    EXPECT_EQ(result, expected);
    ASSERT_TRUE(condition);

    paymentService.process(order);
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"EXPECT_EQ", "ASSERT_TRUE", "TEST"}
	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected test framework call %q to be excluded, got %v", call, hints.Calls)
		}
	}

	// Domain calls should be included
	if !callSet["paymentService.process"] {
		t.Errorf("expected paymentService.process call, got %v", hints.Calls)
	}
}

func TestCppExtractor_Extract_GTestFile(t *testing.T) {
	source := []byte(`
#include <gtest/gtest.h>
#include "services/payment.h"
#include "models/order.h"

class PaymentTest : public ::testing::Test {
protected:
    void SetUp() override {
        gateway = std::make_unique<PaymentGateway>();
    }

    std::unique_ptr<PaymentGateway> gateway;
};

TEST_F(PaymentTest, ProcessPayment) {
    Order order(100);

    gateway->process(order);
    notificationService->sendConfirmation(order.id);

    EXPECT_TRUE(gateway->isComplete());
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Verify imports
	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	expectedImports := []string{"gtest/gtest.h", "services/payment.h", "models/order.h"}
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

	expectedCalls := []string{"gateway.process", "notificationService.sendConfirmation"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q, got %v", call, hints.Calls)
		}
	}
}

func TestCppExtractor_Extract_Catch2File(t *testing.T) {
	source := []byte(`
#include <catch2/catch_test_macros.hpp>
#include "services/user.h"

TEST_CASE("User creation", "[user]") {
    SECTION("valid user") {
        userService.create(validData);
        repository.save(user);

        REQUIRE(user.isValid());
    }
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Verify imports
	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	expectedImports := []string{"catch2/catch_test_macros.hpp", "services/user.h"}
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

	expectedCalls := []string{"userService.create", "repository.save"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q, got %v", call, hints.Calls)
		}
	}

	// REQUIRE should be excluded
	if callSet["REQUIRE"] {
		t.Errorf("expected REQUIRE to be excluded, got %v", hints.Calls)
	}
}

func TestCppExtractor_Extract_Deduplication(t *testing.T) {
	source := []byte(`
#include "myheader.h"
#include "myheader.h"

void test() {
    userService.create(1);
    userService.create(2);
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Count occurrences of domain header (not stdlib)
	importCount := 0
	for _, imp := range hints.Imports {
		if imp == "myheader.h" {
			importCount++
		}
	}
	if importCount != 1 {
		t.Errorf("expected 'myheader.h' to appear once, got %d times in %v", importCount, hints.Imports)
	}

	callCount := 0
	for _, call := range hints.Calls {
		if call == "userService.create" {
			callCount++
		}
	}
	if callCount != 1 {
		t.Errorf("expected 'userService.create' to appear once, got %d times in %v", callCount, hints.Calls)
	}
}

func TestGetExtractor_Cpp(t *testing.T) {
	ext := GetExtractor(domain.LanguageCpp)
	if ext == nil {
		t.Error("expected extractor for C++, got nil")
	}

	_, ok := ext.(*CppExtractor)
	if !ok {
		t.Errorf("expected CppExtractor, got %T", ext)
	}
}

func TestCppExtractor_Extract_NamespacedCalls(t *testing.T) {
	source := []byte(`
#include <vector>

void test() {
    std::vector<int> v;
    MyNamespace::Service::getInstance();
    payment::gateway::process(order);
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Should be normalized to 2 segments
	expectedCalls := []string{"MyNamespace.Service", "payment.gateway"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected %q call (2-segment normalized), got %v", call, hints.Calls)
		}
	}
}

func TestCppExtractor_Extract_StdlibImportsFiltered(t *testing.T) {
	source := []byte(`
#include <iostream>
#include <vector>
#include <string>
#include <memory>
#include <algorithm>
#include "services/payment.h"
#include "models/user.h"
#include <gtest/gtest.h>
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	t.Run("stdlib imports filtered", func(t *testing.T) {
		stdlibHeaders := []string{"iostream", "vector", "string", "memory", "algorithm"}
		for _, stdlib := range stdlibHeaders {
			if importSet[stdlib] {
				t.Errorf("stdlib import %q should be filtered", stdlib)
			}
		}
	})

	t.Run("domain imports included", func(t *testing.T) {
		domainImports := []string{"services/payment.h", "models/user.h", "gtest/gtest.h"}
		for _, domain := range domainImports {
			if !importSet[domain] {
				t.Errorf("domain import %q should be included, got %v", domain, hints.Imports)
			}
		}
	})
}

func TestIsCppStdlibImport(t *testing.T) {
	tests := []struct {
		importPath string
		wantFilter bool
	}{
		// C++ STL containers
		{"vector", true},
		{"map", true},
		{"set", true},
		{"unordered_map", true},
		{"deque", true},
		// C++ STL algorithms and utilities
		{"algorithm", true},
		{"memory", true},
		{"string", true},
		{"functional", true},
		{"chrono", true},
		{"optional", true},
		// C++ I/O streams
		{"iostream", true},
		{"fstream", true},
		{"sstream", true},
		// C compatibility headers
		{"cstdlib", true},
		{"cstring", true},
		{"cmath", true},
		// C headers (legacy)
		{"stdio.h", true},
		{"stdlib.h", true},
		{"string.h", true},
		// Platform-specific (POSIX) - exact match
		{"unistd.h", true},
		{"pthread.h", true},
		// Platform-specific (POSIX) - prefix match
		{"sys/stat.h", true},
		{"sys/socket.h", true},
		{"sys/mman.h", true},
		{"netinet/in.h", true},
		{"netinet/tcp.h", true},
		{"arpa/inet.h", true},
		{"linux/limits.h", true},
		// Platform-specific (Windows)
		{"windows.h", true},
		// Non-stdlib (should NOT be filtered)
		{"gtest/gtest.h", false},
		{"gmock/gmock.h", false},
		{"services/payment.h", false},
		{"models/user.h", false},
		{"myheader.h", false},
		{"boost/asio.hpp", false},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			got := isCppStdlibImport(tt.importPath)
			if got != tt.wantFilter {
				t.Errorf("isCppStdlibImport(%q) = %v, want %v", tt.importPath, got, tt.wantFilter)
			}
		})
	}
}

func TestCppExtractor_Extract_DotPrefixCallsFiltered(t *testing.T) {
	// This tests that dot-prefix patterns from :: scope operator are filtered
	// by the noise_filter.go universal filter
	source := []byte(`
#include <iostream>

void test() {
    // These would generate .testing and .CreateEvent after :: conversion
    userService.create(user);
    PaymentGateway::process(payment);
}
`)

	extractor := &CppExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Verify no dot-prefix patterns leaked through
	for call := range callSet {
		if len(call) > 0 && call[0] == '.' {
			t.Errorf("dot-prefix pattern %q should be filtered", call)
		}
	}

	// Domain calls should be included
	expectedCalls := []string{"userService.create", "PaymentGateway.process"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected domain call %q, got %v", call, hints.Calls)
		}
	}
}

func TestCppExtractor_Extract_NoiseFilter(t *testing.T) {
	source := []byte(`
#include <iostream>

void test() {
    // Decimal literal patterns should be filtered
    auto x = (0.5f);
    auto y = (1.0).toString();

    // Real domain calls should be included
    userService.create(user);
    PaymentGateway::process(payment);
}
`)

	extractor := &CppExtractor{}
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
		expectedCalls := []string{"userService.create", "PaymentGateway.process"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}
