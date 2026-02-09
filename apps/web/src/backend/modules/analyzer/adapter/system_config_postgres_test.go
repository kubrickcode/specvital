package adapter

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/specvital/web/src/backend/modules/analyzer/domain"
)

type mockSystemConfigQueries struct {
	value string
	err   error
}

func (m *mockSystemConfigQueries) GetSystemConfigValue(_ context.Context, _ string) (string, error) {
	return m.value, m.err
}

func TestSystemConfigPostgres_GetParserVersion(t *testing.T) {
	tests := []struct {
		name       string
		mockValue  string
		mockErr    error
		want       string
		wantErr    error
		wantErrMsg string
	}{
		{
			name:      "returns parser version",
			mockValue: "v1.5.1-0.20260112121406-deacdda09e17",
			mockErr:   nil,
			want:      "v1.5.1-0.20260112121406-deacdda09e17",
			wantErr:   nil,
		},
		{
			name:      "returns error when key not found",
			mockValue: "",
			mockErr:   pgx.ErrNoRows,
			want:      "",
			wantErr:   domain.ErrParserVersionNotConfigured,
		},
		{
			name:       "returns wrapped error on database failure",
			mockValue:  "",
			mockErr:    errors.New("connection failed"),
			want:       "",
			wantErrMsg: "get parser version from system_config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewSystemConfigPostgres(&mockSystemConfigQueries{
				value: tt.mockValue,
				err:   tt.mockErr,
			})

			got, err := adapter.GetParserVersion(context.Background())

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetParserVersion() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetParserVersion() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErrMsg != "" {
				if err == nil {
					t.Errorf("GetParserVersion() error = nil, wantErrMsg %q", tt.wantErrMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("GetParserVersion() error = %v, should contain %q", err, tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("GetParserVersion() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("GetParserVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
