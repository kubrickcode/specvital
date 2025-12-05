package detection

import (
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/specvital/core/pkg/parser/detection/matchers"
)

// ProjectContext provides project-level metadata for framework detection.
// Enables source-agnostic detection for both local filesystem and remote sources (e.g., GitHub API).
//
// Thread-safety:
//   - Construction: Use [ProjectContextBuilder] for safe construction. Builder methods are NOT
//     thread-safe; complete all builder operations in a single goroutine before calling Build().
//   - After construction: ProjectContext is immutable. Concurrent reads are safe.
//   - Do NOT call AddConfigFile or SetConfigContent after passing to Scan().
type ProjectContext struct {
	ConfigFiles    []string
	ConfigContents map[string]*matchers.ConfigInfo
}

func NewProjectContext() *ProjectContext {
	return &ProjectContext{
		ConfigFiles:    []string{},
		ConfigContents: make(map[string]*matchers.ConfigInfo),
	}
}

// AddConfigFile adds a config file path. NOT thread-safe; use only during construction.
func (pc *ProjectContext) AddConfigFile(path string) {
	pc.ConfigFiles = append(pc.ConfigFiles, path)
}

// SetConfigContent sets parsed config info for a path. NOT thread-safe; use only during construction.
// Requires pc to be initialized via NewProjectContext() or NewProjectContextBuilder().
func (pc *ProjectContext) SetConfigContent(path string, info *matchers.ConfigInfo) {
	pc.ConfigContents[path] = info
}

// FindApplicableConfig returns the nearest config for a source file path.
// Traverses up directory tree; root-level configs act as fallbacks.
// Returns nil if no config found.
func (pc *ProjectContext) FindApplicableConfig(filePath string) *matchers.ConfigInfo {
	_, info := pc.findNearestConfig(filePath, nil)
	return info
}

// FindConfigPath returns the nearest config path for a source file with optional framework filter.
func (pc *ProjectContext) FindConfigPath(filePath, framework string) string {
	var filter func(*matchers.ConfigInfo) bool
	if framework != "" {
		filter = func(info *matchers.ConfigInfo) bool {
			return info.Framework == framework
		}
	}
	path, _ := pc.findNearestConfig(filePath, filter)
	return path
}

func (pc *ProjectContext) findNearestConfig(filePath string, filter func(*matchers.ConfigInfo) bool) (string, *matchers.ConfigInfo) {
	if pc == nil || len(pc.ConfigContents) == 0 {
		return "", nil
	}

	filePath = filepath.ToSlash(filePath)
	fileDir := filepath.Dir(filePath)

	type candidate struct {
		depth int
		info  *matchers.ConfigInfo
		path  string
	}

	var candidates []candidate
	for path, info := range pc.ConfigContents {
		if info == nil {
			continue
		}
		if filter != nil && !filter(info) {
			continue
		}

		normalizedPath := filepath.ToSlash(path)
		configDir := filepath.Dir(normalizedPath)
		if configDir == "." {
			configDir = ""
		}

		if configDir == "" || strings.HasPrefix(fileDir, configDir) {
			depth := strings.Count(configDir, "/")
			candidates = append(candidates, candidate{path: path, info: info, depth: depth})
		}
	}

	if len(candidates) == 0 {
		return "", nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].depth > candidates[j].depth
	})

	return candidates[0].path, candidates[0].info
}

func (pc *ProjectContext) HasConfigFile(pattern string) bool {
	if pc == nil {
		return false
	}

	return slices.ContainsFunc(pc.ConfigFiles, func(path string) bool {
		return filepath.Base(path) == pattern
	})
}
