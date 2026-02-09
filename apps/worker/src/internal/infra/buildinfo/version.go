package buildinfo

import (
	"runtime/debug"
	"strings"
)

const coreModulePath = "github.com/specvital/core"

// ExtractCoreVersion extracts the version of specvital/core from build info.
// Returns the module version string (e.g., "v1.5.1-0.20260112121406-deacdda09e17")
// or "unknown" if the module is not found.
func ExtractCoreVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	for _, dep := range info.Deps {
		if dep.Path == coreModulePath {
			if dep.Replace != nil {
				return dep.Replace.Version
			}
			return dep.Version
		}
	}

	return "unknown"
}

// FormatVersionDisplay formats a pseudo-version for user display.
// Input: "v1.5.1-0.20260112121406-deacdda09e17"
// Output: "v1.5.1 (deacdda)"
func FormatVersionDisplay(version string) string {
	if version == "unknown" || version == "" {
		return version
	}

	parts := strings.Split(version, "-")
	if len(parts) < 3 {
		return version
	}

	semver := parts[0]
	commitHash := parts[len(parts)-1]

	if len(commitHash) >= 7 {
		commitHash = commitHash[:7]
	}

	return semver + " (" + commitHash + ")"
}
