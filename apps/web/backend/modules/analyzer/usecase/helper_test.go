package usecase

import "testing"

func TestFormatParserVersion(t *testing.T) {
	tests := []struct {
		name    string
		version *string
		want    string
	}{
		{
			name:    "nil input returns empty string",
			version: nil,
			want:    "",
		},
		{
			name:    "empty string returns empty string",
			version: ptr(""),
			want:    "",
		},
		{
			name:    "standard version returns as-is",
			version: ptr("v1.5.1"),
			want:    "v1.5.1",
		},
		{
			name:    "pseudo-version with patch version",
			version: ptr("v1.5.1-0.20260112121406-deacdda09e17"),
			want:    "v1.5.1 (deacdda)",
		},
		{
			name:    "pseudo-version with zero version",
			version: ptr("v0.0.0-20260112121406-abc1234567890"),
			want:    "v0.0.0 (abc1234)",
		},
		{
			name:    "pseudo-version with long commit hash truncated to 7 chars",
			version: ptr("v2.0.0-0.20250101000000-abcdefghijklmn"),
			want:    "v2.0.0 (abcdefg)",
		},
		{
			name:    "pre-release version without pseudo-version format",
			version: ptr("v1.0.0-beta.1"),
			want:    "v1.0.0-beta.1",
		},
		{
			name:    "version with single hyphen",
			version: ptr("v1.0.0-rc1"),
			want:    "v1.0.0-rc1",
		},
		{
			name:    "commit hash exactly 7 chars",
			version: ptr("v1.0.0-0.20260101000000-abcdefg"),
			want:    "v1.0.0 (abcdefg)",
		},
		{
			name:    "commit hash shorter than 7 chars",
			version: ptr("v1.0.0-0.20260101000000-abc"),
			want:    "v1.0.0 (abc)",
		},
		{
			name:    "version with extra hyphens returns as-is",
			version: ptr("v1.0.0-alpha-beta-gamma"),
			want:    "v1.0.0-alpha-beta-gamma",
		},
		{
			name:    "format without v prefix still parses pseudo-version",
			version: ptr("1.0.0-0.20260101000000-abc"),
			want:    "1.0.0 (abc)",
		},
		{
			name:    "invalid timestamp-like second part returns as-is",
			version: ptr("v1.0.0-99999999999999-abc"),
			want:    "v1.0.0 (abc)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatParserVersion(tt.version)
			if got != tt.want {
				t.Errorf("FormatParserVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}
