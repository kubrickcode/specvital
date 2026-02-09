package buildinfo

import "testing"

func TestFormatVersionDisplay(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "pseudo-version format",
			version: "v1.5.1-0.20260112121406-deacdda09e17",
			want:    "v1.5.1 (deacdda)",
		},
		{
			name:    "simple semver",
			version: "v1.5.1",
			want:    "v1.5.1",
		},
		{
			name:    "unknown",
			version: "unknown",
			want:    "unknown",
		},
		{
			name:    "empty",
			version: "",
			want:    "",
		},
		{
			name:    "short pseudo-version",
			version: "v0.0.0-20260101000000-abc",
			want:    "v0.0.0 (abc)",
		},
		{
			name:    "full commit hash in pseudo-version",
			version: "v1.2.3-0.20260101000000-abcdef1234567890",
			want:    "v1.2.3 (abcdef1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatVersionDisplay(tt.version)
			if got != tt.want {
				t.Errorf("FormatVersionDisplay(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}

func TestExtractCoreVersion(t *testing.T) {
	t.Run("should return version string", func(t *testing.T) {
		version := ExtractCoreVersion()

		if version == "" {
			t.Error("expected non-empty version")
		}

		if version == "unknown" {
			t.Skip("skipping: test binary does not have specvital/core dependency")
		}

		if version[0] != 'v' {
			t.Errorf("expected version to start with 'v', got %q", version)
		}
	})
}
