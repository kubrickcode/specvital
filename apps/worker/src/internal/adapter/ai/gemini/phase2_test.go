package gemini

import (
	"context"
	"testing"

	"github.com/kubrickcode/specvital/apps/worker/src/internal/domain/specview"
)

func TestParsePhase2Response_ValidJSON(t *testing.T) {
	jsonStr := `{
		"conversions": [
			{
				"index": 0,
				"description": "사용자가 올바른 자격 증명으로 로그인할 수 있다",
				"confidence": 0.92
			},
			{
				"index": 1,
				"description": "잘못된 비밀번호로 로그인 시 오류를 반환한다",
				"confidence": 0.88
			}
		]
	}`

	// Index mapping: 0→5, 1→7 (simulates original test indices)
	indexMapping := []int{5, 7}

	output, err := parsePhase2Response(jsonStr, indexMapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(output.Behaviors) != 2 {
		t.Errorf("expected 2 behaviors, got %d", len(output.Behaviors))
	}

	// Verify index mapping: AI's 0 → original 5
	behavior := output.Behaviors[0]
	if behavior.TestIndex != 5 {
		t.Errorf("expected test index 5 (mapped from 0), got %d", behavior.TestIndex)
	}
	if behavior.Confidence != 0.92 {
		t.Errorf("expected confidence 0.92, got %f", behavior.Confidence)
	}
	if behavior.Description == "" {
		t.Error("expected non-empty description")
	}

	// Verify second behavior: AI's 1 → original 7
	if output.Behaviors[1].TestIndex != 7 {
		t.Errorf("expected test index 7 (mapped from 1), got %d", output.Behaviors[1].TestIndex)
	}
}

func TestParsePhase2Response_InvalidJSON(t *testing.T) {
	_, err := parsePhase2Response("not json", nil)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParsePhase2Response_EmptyConversions(t *testing.T) {
	jsonStr := `{"conversions": []}`

	output, err := parsePhase2Response(jsonStr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(output.Behaviors) != 0 {
		t.Errorf("expected 0 behaviors, got %d", len(output.Behaviors))
	}
}

func TestValidatePhase2Output_ValidOutput(t *testing.T) {
	input := specview.Phase2Input{
		DomainContext: "Auth",
		FeatureName:   "Login",
		Tests: []specview.TestForConversion{
			{Index: 0, Name: "Test1"},
			{Index: 1, Name: "Test2"},
		},
	}

	output := &specview.Phase2Output{
		Behaviors: []specview.BehaviorSpec{
			{TestIndex: 0, Description: "Description 1", Confidence: 0.9},
			{TestIndex: 1, Description: "Description 2", Confidence: 0.85},
		},
	}

	err := validatePhase2Output(context.Background(), output, input)
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidatePhase2Output_NilOutput(t *testing.T) {
	input := specview.Phase2Input{}

	err := validatePhase2Output(context.Background(), nil, input)
	if err == nil {
		t.Error("expected error for nil output")
	}
}

func TestValidatePhase2Output_EmptyBehaviors(t *testing.T) {
	input := specview.Phase2Input{}
	output := &specview.Phase2Output{Behaviors: []specview.BehaviorSpec{}}

	err := validatePhase2Output(context.Background(), output, input)
	if err == nil {
		t.Error("expected error for empty behaviors")
	}
}

func TestValidatePhase2Output_EmptyDescription(t *testing.T) {
	input := specview.Phase2Input{
		Tests: []specview.TestForConversion{
			{Index: 0, Name: "Test1"},
		},
	}

	output := &specview.Phase2Output{
		Behaviors: []specview.BehaviorSpec{
			{TestIndex: 0, Description: "", Confidence: 0.9}, // Empty description
		},
	}

	err := validatePhase2Output(context.Background(), output, input)
	if err == nil {
		t.Error("expected error for empty description")
	}
}

func TestValidatePhase2Output_UnexpectedTestIndex(t *testing.T) {
	input := specview.Phase2Input{
		Tests: []specview.TestForConversion{
			{Index: 0, Name: "Test1"},
		},
	}

	output := &specview.Phase2Output{
		Behaviors: []specview.BehaviorSpec{
			{TestIndex: 0, Description: "Valid", Confidence: 0.9},
			{TestIndex: 999, Description: "Invalid index", Confidence: 0.9}, // 999 doesn't exist
		},
	}

	err := validatePhase2Output(context.Background(), output, input)
	if err == nil {
		t.Error("expected error for unexpected test index")
	}
}

func TestValidatePhase2Output_PartialCoverage(t *testing.T) {
	// Partial coverage should log warning but not fail
	input := specview.Phase2Input{
		Tests: []specview.TestForConversion{
			{Index: 0, Name: "Test1"},
			{Index: 1, Name: "Test2"},
			{Index: 2, Name: "Test3"},
		},
	}

	output := &specview.Phase2Output{
		Behaviors: []specview.BehaviorSpec{
			{TestIndex: 0, Description: "Description 1", Confidence: 0.9},
			// Missing index 1 and 2
		},
	}

	// Should not error, just warn
	err := validatePhase2Output(context.Background(), output, input)
	if err != nil {
		t.Errorf("partial coverage should not error: %v", err)
	}
}
