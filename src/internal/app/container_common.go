package app

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const schedulerLockKey = "scheduler:auto-refresh:lock"

// ContainerConfig holds common configuration for dependency injection containers.
type ContainerConfig struct {
	EncryptionKey   string
	GeminiAPIKey    string
	GeminiPhase1Model string // optional: default gemini-2.5-flash
	GeminiPhase2Model string // optional: default gemini-2.5-flash-lite
	ParserVersion   string
	Pool            *pgxpool.Pool
}

// Validate checks that required common configuration fields are set.
func (c ContainerConfig) Validate() error {
	if c.Pool == nil {
		return fmt.Errorf("pool is required")
	}
	return nil
}

// ValidateAnalyzer checks that all analyzer-specific configuration is valid.
func (c ContainerConfig) ValidateAnalyzer() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.EncryptionKey == "" {
		return fmt.Errorf("encryption key is required")
	}
	if c.ParserVersion == "" {
		return fmt.Errorf("parser version is required")
	}
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("gemini API key is required")
	}
	return nil
}
