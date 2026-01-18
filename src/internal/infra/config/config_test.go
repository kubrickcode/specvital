package config

import (
	"os"
	"testing"
)

func TestLoadQueueConfig_Defaults(t *testing.T) {
	clearQueueEnvVars(t)

	cfg := loadQueueConfig()

	if cfg.Analyzer.Priority != defaultAnalyzerPriorityWorkers {
		t.Errorf("Analyzer.Priority = %d, want %d", cfg.Analyzer.Priority, defaultAnalyzerPriorityWorkers)
	}
	if cfg.Analyzer.Default != defaultAnalyzerDefaultWorkers {
		t.Errorf("Analyzer.Default = %d, want %d", cfg.Analyzer.Default, defaultAnalyzerDefaultWorkers)
	}
	if cfg.Analyzer.Scheduled != defaultAnalyzerScheduledWorkers {
		t.Errorf("Analyzer.Scheduled = %d, want %d", cfg.Analyzer.Scheduled, defaultAnalyzerScheduledWorkers)
	}
	if cfg.Specgen.Priority != defaultSpecgenPriorityWorkers {
		t.Errorf("Specgen.Priority = %d, want %d", cfg.Specgen.Priority, defaultSpecgenPriorityWorkers)
	}
	if cfg.Specgen.Default != defaultSpecgenDefaultWorkers {
		t.Errorf("Specgen.Default = %d, want %d", cfg.Specgen.Default, defaultSpecgenDefaultWorkers)
	}
	if cfg.Specgen.Scheduled != defaultSpecgenScheduledWorkers {
		t.Errorf("Specgen.Scheduled = %d, want %d", cfg.Specgen.Scheduled, defaultSpecgenScheduledWorkers)
	}
}

func TestLoadQueueConfig_EnvOverride(t *testing.T) {
	clearQueueEnvVars(t)
	t.Setenv("ANALYZER_QUEUE_PRIORITY_WORKERS", "10")
	t.Setenv("ANALYZER_QUEUE_DEFAULT_WORKERS", "7")
	t.Setenv("SPECGEN_QUEUE_SCHEDULED_WORKERS", "4")

	cfg := loadQueueConfig()

	if cfg.Analyzer.Priority != 10 {
		t.Errorf("Analyzer.Priority = %d, want 10", cfg.Analyzer.Priority)
	}
	if cfg.Analyzer.Default != 7 {
		t.Errorf("Analyzer.Default = %d, want 7", cfg.Analyzer.Default)
	}
	if cfg.Specgen.Scheduled != 4 {
		t.Errorf("Specgen.Scheduled = %d, want 4", cfg.Specgen.Scheduled)
	}
	// Non-overridden values should use defaults
	if cfg.Analyzer.Scheduled != defaultAnalyzerScheduledWorkers {
		t.Errorf("Analyzer.Scheduled = %d, want %d", cfg.Analyzer.Scheduled, defaultAnalyzerScheduledWorkers)
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue int
		want         int
	}{
		{"empty uses default", "", 5, 5},
		{"valid number", "10", 5, 10},
		{"invalid string uses default", "abc", 5, 5},
		{"negative uses default", "-1", 5, 5},
		{"zero uses default", "0", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_ENV_INT"
			if tt.envValue != "" {
				t.Setenv(key, tt.envValue)
			} else {
				os.Unsetenv(key)
			}

			got := getEnvInt(key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvInt(%q, %d) = %d, want %d", tt.envValue, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func clearQueueEnvVars(t *testing.T) {
	t.Helper()
	envVars := []string{
		"ANALYZER_QUEUE_PRIORITY_WORKERS",
		"ANALYZER_QUEUE_DEFAULT_WORKERS",
		"ANALYZER_QUEUE_SCHEDULED_WORKERS",
		"SPECGEN_QUEUE_PRIORITY_WORKERS",
		"SPECGEN_QUEUE_DEFAULT_WORKERS",
		"SPECGEN_QUEUE_SCHEDULED_WORKERS",
	}
	for _, env := range envVars {
		os.Unsetenv(env)
	}
}
