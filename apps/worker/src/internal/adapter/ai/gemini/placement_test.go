package gemini

import (
	"strings"
	"testing"

	"github.com/kubrickcode/specvital/apps/worker/src/internal/domain/specview"
)

func TestParsePlacementResponse(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		expectedCount int
		wantErr       bool
		errContains   string
		validate      func(*testing.T, *specview.PlacementOutput)
	}{
		{
			name: "valid response with two placements",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "Authentication", "feature": "Login"},
					{"test_index": 1, "domain": "Payment", "feature": "Order"}
				]
			}`,
			expectedCount: 2,
			wantErr:       false,
			validate: func(t *testing.T, output *specview.PlacementOutput) {
				if len(output.Placements) != 2 {
					t.Errorf("expected 2 placements, got %d", len(output.Placements))
				}
				if output.Placements[0].DomainName != "Authentication" {
					t.Errorf("expected domain 'Authentication', got %q", output.Placements[0].DomainName)
				}
				if output.Placements[0].FeatureName != "Login" {
					t.Errorf("expected feature 'Login', got %q", output.Placements[0].FeatureName)
				}
				if output.Placements[1].TestIndex != 1 {
					t.Errorf("expected test_index 1, got %d", output.Placements[1].TestIndex)
				}
			},
		},
		{
			name: "uncategorized placement",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "Uncategorized", "feature": "Uncategorized"}
				]
			}`,
			expectedCount: 1,
			wantErr:       false,
			validate: func(t *testing.T, output *specview.PlacementOutput) {
				if output.Placements[0].DomainName != "Uncategorized" {
					t.Errorf("expected domain 'Uncategorized', got %q", output.Placements[0].DomainName)
				}
			},
		},
		{
			name:          "invalid json",
			json:          `{invalid}`,
			expectedCount: 1,
			wantErr:       true,
			errContains:   "json unmarshal",
		},
		{
			name: "invalid test_index negative",
			json: `{
				"placements": [
					{"test_index": -1, "domain": "Auth", "feature": "Login"}
				]
			}`,
			expectedCount: 1,
			wantErr:       true,
			errContains:   "invalid test_index",
		},
		{
			name: "invalid test_index exceeds count",
			json: `{
				"placements": [
					{"test_index": 5, "domain": "Auth", "feature": "Login"}
				]
			}`,
			expectedCount: 2,
			wantErr:       true,
			errContains:   "invalid test_index",
		},
		{
			name: "duplicate test_index",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "Auth", "feature": "Login"},
					{"test_index": 0, "domain": "Payment", "feature": "Order"}
				]
			}`,
			expectedCount: 2,
			wantErr:       true,
			errContains:   "duplicate test_index",
		},
		{
			name: "empty domain",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "", "feature": "Login"}
				]
			}`,
			expectedCount: 1,
			wantErr:       true,
			errContains:   "empty domain or feature",
		},
		{
			name: "empty feature",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "Auth", "feature": ""}
				]
			}`,
			expectedCount: 1,
			wantErr:       true,
			errContains:   "empty domain or feature",
		},
		{
			name: "count mismatch fewer",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "Auth", "feature": "Login"}
				]
			}`,
			expectedCount: 2,
			wantErr:       true,
			errContains:   "placement count mismatch",
		},
		{
			name: "count mismatch more with invalid index",
			json: `{
				"placements": [
					{"test_index": 0, "domain": "Auth", "feature": "Login"},
					{"test_index": 1, "domain": "Payment", "feature": "Order"},
					{"test_index": 2, "domain": "User", "feature": "Profile"}
				]
			}`,
			expectedCount: 2,
			wantErr:       true,
			errContains:   "invalid test_index",
		},
		{
			name:          "empty placements when expected none",
			json:          `{"placements": []}`,
			expectedCount: 0,
			wantErr:       false,
			validate: func(t *testing.T, output *specview.PlacementOutput) {
				if len(output.Placements) != 0 {
					t.Errorf("expected 0 placements, got %d", len(output.Placements))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := parsePlacementResponse(tc.json, tc.expectedCount)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("expected error to contain %q, got %q", tc.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.validate != nil {
				tc.validate(t, output)
			}
		})
	}
}
