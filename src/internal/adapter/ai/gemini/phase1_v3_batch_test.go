package gemini

import (
	"testing"
)

func TestParseV3BatchResponse_Success(t *testing.T) {
	jsonStr := `[{"d": "Authentication", "f": "Login"}, {"d": "Payment", "f": "Checkout"}]`

	results, err := parseV3BatchResponse(jsonStr)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Domain != "Authentication" {
		t.Errorf("expected domain 'Authentication', got %q", results[0].Domain)
	}
	if results[0].Feature != "Login" {
		t.Errorf("expected feature 'Login', got %q", results[0].Feature)
	}
	if results[1].Domain != "Payment" {
		t.Errorf("expected domain 'Payment', got %q", results[1].Domain)
	}
	if results[1].Feature != "Checkout" {
		t.Errorf("expected feature 'Checkout', got %q", results[1].Feature)
	}
}

func TestParseV3BatchResponse_EmptyArray(t *testing.T) {
	jsonStr := `[]`

	results, err := parseV3BatchResponse(jsonStr)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestParseV3BatchResponse_EmptyString(t *testing.T) {
	_, err := parseV3BatchResponse("")

	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

func TestParseV3BatchResponse_NullArray(t *testing.T) {
	_, err := parseV3BatchResponse("null")

	if err == nil {
		t.Fatal("expected error for null response")
	}
}

func TestParseV3BatchResponse_ObjectInsteadOfArray(t *testing.T) {
	_, err := parseV3BatchResponse(`{"invalid": "not an array"}`)

	if err == nil {
		t.Fatal("expected error for JSON object instead of array")
	}
}

func TestParseV3BatchResponse_MalformedJSON(t *testing.T) {
	_, err := parseV3BatchResponse(`[{"d": "Auth", "f": "Login"`)

	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestParseV3BatchResponse_SingleItem(t *testing.T) {
	jsonStr := `[{"d": "Uncategorized", "f": "General"}]`

	results, err := parseV3BatchResponse(jsonStr)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Domain != "Uncategorized" {
		t.Errorf("expected domain 'Uncategorized', got %q", results[0].Domain)
	}
}

func TestValidateV3BatchCount_Match(t *testing.T) {
	results := []v3BatchResult{
		{Domain: "Auth", Feature: "Login"},
		{Domain: "Payment", Feature: "Checkout"},
	}

	err := validateV3BatchCount(results, 2)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateV3BatchCount_TooFew(t *testing.T) {
	results := []v3BatchResult{
		{Domain: "Auth", Feature: "Login"},
	}

	err := validateV3BatchCount(results, 3)

	if err == nil {
		t.Fatal("expected error for count mismatch")
	}
}

func TestValidateV3BatchCount_TooMany(t *testing.T) {
	results := []v3BatchResult{
		{Domain: "Auth", Feature: "Login"},
		{Domain: "Payment", Feature: "Checkout"},
		{Domain: "User", Feature: "Profile"},
	}

	err := validateV3BatchCount(results, 2)

	if err == nil {
		t.Fatal("expected error for count mismatch")
	}
}

func TestValidateV3BatchCount_EmptyExpected(t *testing.T) {
	results := []v3BatchResult{}

	err := validateV3BatchCount(results, 0)

	if err != nil {
		t.Errorf("unexpected error for empty expected: %v", err)
	}
}

func TestValidateV3BatchCount_EmptyGot(t *testing.T) {
	results := []v3BatchResult{}

	err := validateV3BatchCount(results, 2)

	if err == nil {
		t.Fatal("expected error for count mismatch")
	}
}

func TestParseV3BatchResponse_EmptyDomain(t *testing.T) {
	jsonStr := `[{"d": "", "f": "Login"}]`

	_, err := parseV3BatchResponse(jsonStr)

	if err == nil {
		t.Fatal("expected error for empty domain")
	}
}

func TestParseV3BatchResponse_EmptyFeature(t *testing.T) {
	jsonStr := `[{"d": "Auth", "f": ""}]`

	_, err := parseV3BatchResponse(jsonStr)

	if err == nil {
		t.Fatal("expected error for empty feature")
	}
}

func TestParseV3BatchResponse_EmptyFieldAtSecondIndex(t *testing.T) {
	jsonStr := `[{"d": "Auth", "f": "Login"}, {"d": "", "f": "Checkout"}]`

	_, err := parseV3BatchResponse(jsonStr)

	if err == nil {
		t.Fatal("expected error for empty domain at index 1")
	}
}
