package middleware

import (
	"errors"
	"os"
	"testing"
)

func TestGetAllowedOrigins(t *testing.T) {
	tests := []struct {
		name       string
		envOrigins string
		envEnv     string
		want       []string
		wantErr    error
		setOrigins bool
		setEnv     bool
	}{
		{
			name:       "default origins when no env set",
			setOrigins: false,
			setEnv:     false,
			want:       []string{"http://localhost:5173"},
			wantErr:    nil,
		},
		{
			name:       "custom multiple origins",
			envOrigins: "https://example.com, https://api.example.com",
			setOrigins: true,
			setEnv:     false,
			want:       []string{"https://example.com", "https://api.example.com"},
			wantErr:    nil,
		},
		{
			name:       "single origin",
			envOrigins: "https://myapp.com",
			setOrigins: true,
			setEnv:     false,
			want:       []string{"https://myapp.com"},
			wantErr:    nil,
		},
		{
			name:       "error on missing origins in production",
			envOrigins: "",
			envEnv:     "production",
			setOrigins: false,
			setEnv:     true,
			want:       nil,
			wantErr:    ErrCORSOriginsNotSet,
		},
		{
			name:       "error on wildcard origin",
			envOrigins: "https://example.com,*",
			setOrigins: true,
			setEnv:     false,
			want:       nil,
			wantErr:    ErrWildcardOrigin,
		},
		{
			name:       "handles empty string between commas",
			envOrigins: "https://example.com,,https://api.example.com",
			setOrigins: true,
			setEnv:     false,
			want:       []string{"https://example.com", "https://api.example.com"},
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("CORS_ORIGINS")
			os.Unsetenv("ENV")
			t.Cleanup(func() {
				os.Unsetenv("CORS_ORIGINS")
				os.Unsetenv("ENV")
			})

			if tt.setOrigins {
				os.Setenv("CORS_ORIGINS", tt.envOrigins)
			}
			if tt.setEnv {
				os.Setenv("ENV", tt.envEnv)
			}

			got, err := GetAllowedOrigins()

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("expected %d origins, got %d", len(tt.want), len(got))
				return
			}

			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("origin[%d]: expected %s, got %s", i, want, got[i])
				}
			}
		})
	}
}

func TestCORS_ReturnsHandler(t *testing.T) {
	origins := []string{"http://localhost:5173"}

	handler := CORS(origins)
	if handler == nil {
		t.Error("expected non-nil handler")
	}
}
