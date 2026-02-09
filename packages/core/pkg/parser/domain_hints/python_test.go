package domain_hints

import (
	"context"
	"strings"
	"testing"

	"github.com/specvital/core/pkg/domain"
)

func TestPythonExtractor_Extract_Imports(t *testing.T) {
	source := []byte(`
import pytest
import os.path
from stripe import PaymentIntent
from myapp.services import UserService
from datetime import datetime, timedelta
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedImports := map[string]bool{
		"pytest":         true,
		"os.path":        true,
		"stripe":         true,
		"myapp.services": true,
		"datetime":       true,
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

func TestPythonExtractor_Extract_Calls(t *testing.T) {
	source := []byte(`
import pytest

def test_create_payment():
    payment_service.create_intent("amount")
    user_repo.find_by_id(1)
    result = order_service.process(order)
    do_something()
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	expectedCalls := map[string]bool{
		"payment_service.create_intent": true,
		"user_repo.find_by_id":          true,
		"order_service.process":         true,
		"do_something":                  true,
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

func TestPythonExtractor_Extract_EmptyFile(t *testing.T) {
	source := []byte(`# empty file`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints != nil {
		t.Errorf("expected nil for empty file, got %+v", hints)
	}
}

func TestPythonExtractor_Extract_TestFrameworkCalls(t *testing.T) {
	source := []byte(`
import pytest

@pytest.fixture
def mock_user():
    return {"id": 1}

def test_with_fixture(mock_user):
    pytest.mark.skip("reason")
    auth_service.validate(mock_user)
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// Test framework calls should be excluded
	excludedCalls := []string{"pytest.fixture", "pytest.mark"}
	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected test framework call %q to be excluded", call)
		}
	}

	// Domain calls should be included
	if !callSet["auth_service.validate"] {
		t.Errorf("expected auth_service.validate call, got %v", hints.Calls)
	}
}

func TestPythonExtractor_Extract_ChainedCalls(t *testing.T) {
	source := []byte(`
import pytest

def test_chained():
    # Long chains should be normalized to 2 segments
    client.api.users.create().json()
    response.data.items.first().value
`)

	extractor := &PythonExtractor{}
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

func TestPythonExtractor_Extract_RelativeImports(t *testing.T) {
	source := []byte(`
from . import utils
from .. import parent
from .models import User
from ..services import PaymentService
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	// Relative imports should capture the relative path
	expectedImports := []string{".models", "..services"}
	for _, imp := range expectedImports {
		if !importSet[imp] {
			t.Errorf("expected relative import %q to be included, got %v", imp, hints.Imports)
		}
	}
}

func TestPythonExtractor_Extract_PytestFile(t *testing.T) {
	source := []byte(`
import pytest
from stripe import PaymentIntent
from myapp.models import Order
from myapp.services.payment import PaymentService

class TestPaymentFlow:
    @pytest.fixture
    def payment_service(self):
        return PaymentService()

    def test_create_payment(self, payment_service):
        order = Order(amount=100)
        result = payment_service.create(order)
        stripe_api.confirm_intent(result.intent_id)
        assert result.status == "pending"
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	// Verify imports
	importSet := make(map[string]bool)
	for _, imp := range hints.Imports {
		importSet[imp] = true
	}

	expectedImports := []string{"pytest", "stripe", "myapp.models", "myapp.services.payment"}
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

	expectedCalls := []string{"payment_service.create", "stripe_api.confirm_intent"}
	for _, call := range expectedCalls {
		if !callSet[call] {
			t.Errorf("expected call %q, got %v", call, hints.Calls)
		}
	}
}

func TestPythonExtractor_Extract_UnittestCalls(t *testing.T) {
	source := []byte(`
import unittest

class TestPayment(unittest.TestCase):
    def setUp(self):
        self.client = PaymentClient()

    def test_payment(self):
        result = self.client.create_payment(100)
        self.assertEqual(result.status, "success")
        self.assertTrue(result.confirmed)
        payment_service.validate(result)
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	// self.* calls should be excluded
	excludedCalls := []string{"self.assertEqual", "self.assertTrue", "self.client"}
	for _, call := range excludedCalls {
		if callSet[call] {
			t.Errorf("expected self.* call %q to be excluded", call)
		}
	}

	// Domain calls should be included
	expectedDomainCalls := []string{"PaymentClient", "payment_service.validate"}
	for _, call := range expectedDomainCalls {
		if !callSet[call] {
			t.Errorf("expected %q call, got %v", call, hints.Calls)
		}
	}
}

func TestGetExtractor_Python(t *testing.T) {
	ext := GetExtractor(domain.LanguagePython)
	if ext == nil {
		t.Error("expected extractor for Python, got nil")
	}

	_, ok := ext.(*PythonExtractor)
	if !ok {
		t.Errorf("expected PythonExtractor, got %T", ext)
	}
}

func TestPythonExtractor_Extract_NoiseFilter(t *testing.T) {
	source := []byte(`
import pytest

def test_noise_patterns():
    # Bracket patterns should be filtered
    items = [1, 2, 3]
    result = [...items]

    # Dollar sign patterns (rare in Python but test anyway)
    $ = None

    # Real domain calls should be included
    user_service.create(data)
    payment_gateway.process(amount)
`)

	extractor := &PythonExtractor{}
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
			if len(call) > 0 && (call[0] == '[' || call[0] == '(' || call == "$") {
				t.Errorf("noise pattern %q should be filtered, got %v", call, hints.Calls)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"user_service.create", "payment_gateway.process"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestPythonExtractor_Extract_StringLiteralMethods(t *testing.T) {
	source := []byte(`
import pytest

def test_string_methods():
    # String literal method calls - should be filtered
    result = "ööö".encode("utf-8")
    url = "http://example.com".format()
    data = 'hello'.upper()

    # Normal method calls - should be included
    user_service.create(data)
    response.json()
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("string literal methods filtered", func(t *testing.T) {
		for call := range callSet {
			if len(call) > 0 && (call[0] == '"' || call[0] == '\'') {
				t.Errorf("string literal method %q should be filtered", call)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"user_service.create", "response.json"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}

func TestPythonExtractor_Extract_FunctionArgumentsFiltered(t *testing.T) {
	source := []byte(`
import requests

def test_request():
    # Function calls with arguments leaked - should be filtered
    requests.Request(method="GET", url="http://example.com")
    response = session.get("http://api.example.com", headers=headers)

    # Normal calls - should be included
    client.send(data)
    parser.parse(content)
`)

	extractor := &PythonExtractor{}
	hints := extractor.Extract(context.Background(), source)

	if hints == nil {
		t.Fatal("expected hints, got nil")
	}

	callSet := make(map[string]bool)
	for _, call := range hints.Calls {
		callSet[call] = true
	}

	t.Run("arguments not included in calls", func(t *testing.T) {
		for call := range callSet {
			if strings.Contains(call, "=") {
				t.Errorf("call %q should not contain arguments (=)", call)
			}
		}
	})

	t.Run("domain calls included", func(t *testing.T) {
		expectedCalls := []string{"client.send", "parser.parse"}
		for _, call := range expectedCalls {
			if !callSet[call] {
				t.Errorf("expected domain call %q, got %v", call, hints.Calls)
			}
		}
	})
}
