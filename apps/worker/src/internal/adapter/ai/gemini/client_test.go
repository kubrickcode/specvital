package gemini

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "empty API key",
			config:  Config{},
			wantErr: true,
		},
		{
			name: "valid config with API key",
			config: Config{
				APIKey: "test-api-key",
			},
			wantErr: false,
		},
		{
			name: "valid config with all fields",
			config: Config{
				APIKey:      "test-api-key",
				Phase1Model: "gemini-2.5-flash",
				Phase2Model: "gemini-2.5-flash-lite",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultModels(t *testing.T) {
	if defaultPhase1Model != "gemini-2.5-flash" {
		t.Errorf("expected default phase1 model to be gemini-2.5-flash, got %s", defaultPhase1Model)
	}
	if defaultPhase2Model != "gemini-2.5-flash-lite" {
		t.Errorf("expected default phase2 model to be gemini-2.5-flash-lite, got %s", defaultPhase2Model)
	}
}

func TestDefaultSeed(t *testing.T) {
	if defaultSeed != 42 {
		t.Errorf("expected default seed to be 42, got %d", defaultSeed)
	}
}
