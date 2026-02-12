package prompt

import (
	"strings"
	"testing"

	"github.com/kubrickcode/specvital/apps/worker/internal/domain/specview"
)

func TestBuildPlacementUserPrompt(t *testing.T) {
	tests := []struct {
		name     string
		input    specview.PlacementInput
		contains []string
		notEmpty bool
	}{
		{
			name: "basic structure with new tests",
			input: specview.PlacementInput{
				Language: "Korean",
				ExistingStructure: &specview.Phase1Output{
					Domains: []specview.DomainGroup{
						{
							Name: "Authentication",
							Features: []specview.FeatureGroup{
								{Name: "Login"},
								{Name: "Password Management"},
							},
						},
						{
							Name: "Payment",
							Features: []specview.FeatureGroup{
								{Name: "Order"},
							},
						},
					},
				},
				NewTests: []specview.TestInfo{
					{Index: 0, Name: "test_2fa_setup"},
					{Index: 1, Name: "test_notification_send", SuitePath: "Notification > Email"},
				},
			},
			contains: []string{
				"Target Language: Korean",
				"<structure>",
				"D:Authentication",
				"F:Login",
				"F:Password Management",
				"D:Payment",
				"F:Order",
				"</structure>",
				"<new_tests>",
				"0|test_2fa_setup",
				"1|Notification > Email|test_notification_send",
				"</new_tests>",
				"Place 2 tests",
			},
			notEmpty: true,
		},
		{
			name: "empty new tests",
			input: specview.PlacementInput{
				Language: "English",
				ExistingStructure: &specview.Phase1Output{
					Domains: []specview.DomainGroup{
						{Name: "Auth", Features: []specview.FeatureGroup{{Name: "Login"}}},
					},
				},
				NewTests: []specview.TestInfo{},
			},
			contains: []string{
				"Place 0 tests",
			},
			notEmpty: true,
		},
		{
			name: "test without suite path",
			input: specview.PlacementInput{
				Language: "English",
				ExistingStructure: &specview.Phase1Output{
					Domains: []specview.DomainGroup{
						{Name: "User", Features: []specview.FeatureGroup{{Name: "Profile"}}},
					},
				},
				NewTests: []specview.TestInfo{
					{Index: 0, Name: "test_update_profile"},
				},
			},
			contains: []string{
				"0|test_update_profile",
			},
			notEmpty: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := BuildPlacementUserPrompt(tc.input)

			if tc.notEmpty && result == "" {
				t.Error("expected non-empty result")
			}

			for _, s := range tc.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected result to contain %q, got:\n%s", s, result)
				}
			}
		})
	}
}

func TestPlacementSystemPrompt(t *testing.T) {
	if PlacementSystemPrompt == "" {
		t.Error("PlacementSystemPrompt should not be empty")
	}

	requiredPhrases := []string{
		"Uncategorized",
		"placements",
		"test_index",
		"domain",
		"feature",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(PlacementSystemPrompt, phrase) {
			t.Errorf("PlacementSystemPrompt should contain %q", phrase)
		}
	}
}
