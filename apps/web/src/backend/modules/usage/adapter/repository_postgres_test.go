package adapter

import (
	"testing"

	"github.com/specvital/web/src/backend/modules/usage/domain/entity"
)

func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventType entity.EventType
		want      bool
	}{
		{
			name:      "specview is valid",
			eventType: entity.EventTypeSpecview,
			want:      true,
		},
		{
			name:      "analysis is valid",
			eventType: entity.EventTypeAnalysis,
			want:      true,
		},
		{
			name:      "empty string is invalid",
			eventType: "",
			want:      false,
		},
		{
			name:      "unknown type is invalid",
			eventType: "unknown",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventType.IsValid(); got != tt.want {
				t.Errorf("EventType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUuidFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			input:   "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uuidFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("uuidFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
