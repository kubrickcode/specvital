// Package detection provides hierarchical test framework detection.
package detection

type Confidence int

const (
	ConfidenceUnknown Confidence = iota
	ConfidenceLow                // scope config (config file exists in directory tree)
	ConfidenceMedium             // import (explicit import statement)
	ConfidenceHigh               // project context (config parsed with globals mode)
)

type Source string

const (
	SourceUnknown        Source = "unknown"
	SourceImport         Source = "import"
	SourceScopeConfig    Source = "scope_config"
	SourceProjectContext Source = "project_context"
	SourcePragma         Source = "pragma"
)

const FrameworkUnknown = "unknown"

type Result struct {
	ConfigPath string
	Confidence Confidence
	Framework  string
	Source     Source
}

func (r Result) IsUnknown() bool {
	return r.Framework == "" || r.Framework == FrameworkUnknown || r.Confidence == ConfidenceUnknown
}

func Unknown() Result {
	return Result{
		Confidence: ConfidenceUnknown,
		Framework:  FrameworkUnknown,
		Source:     SourceUnknown,
	}
}

func FromImport(framework string) Result {
	return Result{
		Confidence: ConfidenceMedium,
		Framework:  framework,
		Source:     SourceImport,
	}
}

func FromScopeConfig(framework, configPath string) Result {
	return Result{
		ConfigPath: configPath,
		Confidence: ConfidenceLow,
		Framework:  framework,
		Source:     SourceScopeConfig,
	}
}

func FromProjectContext(framework, configPath string) Result {
	return Result{
		ConfigPath: configPath,
		Confidence: ConfidenceHigh,
		Framework:  framework,
		Source:     SourceProjectContext,
	}
}
