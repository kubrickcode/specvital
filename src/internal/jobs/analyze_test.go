package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/specvital/collector/internal/service"
	"github.com/specvital/collector/internal/service/mocks"
)

func TestNewAnalyzeHandler(t *testing.T) {
	mockSvc := &mocks.MockAnalysisService{}
	handler := NewAnalyzeHandler(mockSvc)
	if handler == nil {
		t.Error("expected handler, got nil")
	}
}

func TestAnalyzeHandler_ProcessTask(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		setupMock   func(*mocks.MockAnalysisService)
		wantErr     bool
		errContains string
	}{
		{
			name: "should process task successfully",
			payload: AnalyzePayload{
				Owner: "octocat",
				Repo:  "Hello-World",
			},
			setupMock: func(m *mocks.MockAnalysisService) {
				m.AnalyzeFunc = func(ctx context.Context, req service.AnalyzeRequest) error {
					if req.Owner != "octocat" {
						t.Errorf("expected owner 'octocat', got '%s'", req.Owner)
					}
					if req.Repo != "Hello-World" {
						t.Errorf("expected repo 'Hello-World', got '%s'", req.Repo)
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "should handle service error",
			payload: AnalyzePayload{
				Owner: "testowner",
				Repo:  "testrepo",
			},
			setupMock: func(m *mocks.MockAnalysisService) {
				m.AnalyzeFunc = func(ctx context.Context, req service.AnalyzeRequest) error {
					return errors.New("service error")
				}
			},
			wantErr:     true,
			errContains: "service error",
		},
		{
			name: "should propagate validation error from service",
			payload: AnalyzePayload{
				Owner: "",
				Repo:  "testrepo",
			},
			setupMock: func(m *mocks.MockAnalysisService) {
				m.AnalyzeFunc = func(ctx context.Context, req service.AnalyzeRequest) error {
					return service.ErrInvalidInput
				}
			},
			wantErr: true,
		},
		{
			name: "should propagate clone error from service",
			payload: AnalyzePayload{
				Owner: "nonexistent",
				Repo:  "repo",
			},
			setupMock: func(m *mocks.MockAnalysisService) {
				m.AnalyzeFunc = func(ctx context.Context, req service.AnalyzeRequest) error {
					return service.ErrCloneFailed
				}
			},
			wantErr: true,
		},
		{
			name: "should propagate save error from service",
			payload: AnalyzePayload{
				Owner: "owner",
				Repo:  "repo",
			},
			setupMock: func(m *mocks.MockAnalysisService) {
				m.AnalyzeFunc = func(ctx context.Context, req service.AnalyzeRequest) error {
					return service.ErrSaveFailed
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockSvc := &mocks.MockAnalysisService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}
			handler := NewAnalyzeHandler(mockSvc)

			payloadBytes, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			task := asynq.NewTask(TypeAnalyze, payloadBytes)

			// When
			err = handler.ProcessTask(context.Background(), task)

			// Then
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil {
					if !containsString(err.Error(), tt.errContains) {
						t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestAnalyzeHandler_ProcessTask_InvalidPayload(t *testing.T) {
	tests := []struct {
		name        string
		payload     []byte
		errContains string
	}{
		{
			name:        "malformed JSON",
			payload:     []byte(`{invalid json`),
			errContains: "unmarshal payload",
		},
		{
			name:        "invalid JSON type",
			payload:     []byte(`"just a string"`),
			errContains: "unmarshal payload",
		},
		{
			name:        "empty payload",
			payload:     []byte(``),
			errContains: "unmarshal payload",
		},
		{
			name:        "null payload",
			payload:     []byte(`null`),
			errContains: "",
		},
		{
			name:        "wrong field types",
			payload:     []byte(`{"owner": 123, "repo": true}`),
			errContains: "unmarshal payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockSvc := &mocks.MockAnalysisService{}
			handler := NewAnalyzeHandler(mockSvc)
			task := asynq.NewTask(TypeAnalyze, tt.payload)

			// When
			err := handler.ProcessTask(context.Background(), task)

			// Then
			if err == nil && tt.errContains != "" {
				t.Error("expected error, got nil")
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestAnalyzeHandler_ProcessTask_ServiceInvocation(t *testing.T) {
	t.Run("should pass correct parameters to service", func(t *testing.T) {
		// Given
		var capturedReq service.AnalyzeRequest
		mockSvc := &mocks.MockAnalysisService{
			AnalyzeFunc: func(ctx context.Context, req service.AnalyzeRequest) error {
				capturedReq = req
				return nil
			},
		}
		handler := NewAnalyzeHandler(mockSvc)

		payload := AnalyzePayload{
			Owner: "test-owner",
			Repo:  "test-repo",
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		task := asynq.NewTask(TypeAnalyze, payloadBytes)

		// When
		err = handler.ProcessTask(context.Background(), task)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedReq.Owner != "test-owner" {
			t.Errorf("expected owner 'test-owner', got '%s'", capturedReq.Owner)
		}
		if capturedReq.Repo != "test-repo" {
			t.Errorf("expected repo 'test-repo', got '%s'", capturedReq.Repo)
		}
	})

	t.Run("should not call service when payload unmarshal fails", func(t *testing.T) {
		// Given
		serviceCalled := false
		mockSvc := &mocks.MockAnalysisService{
			AnalyzeFunc: func(ctx context.Context, req service.AnalyzeRequest) error {
				serviceCalled = true
				return nil
			},
		}
		handler := NewAnalyzeHandler(mockSvc)
		task := asynq.NewTask(TypeAnalyze, []byte(`invalid json`))

		// When
		err := handler.ProcessTask(context.Background(), task)

		// Then
		if err == nil {
			t.Error("expected error, got nil")
		}
		if serviceCalled {
			t.Error("service should not be called when payload unmarshal fails")
		}
	})
}

func TestAnalyzeHandler_ProcessTask_ContextPropagation(t *testing.T) {
	t.Run("should propagate context to service", func(t *testing.T) {
		// Given
		type ctxKey string
		testKey := ctxKey("test-key")
		testValue := "test-value"

		var capturedCtx context.Context
		mockSvc := &mocks.MockAnalysisService{
			AnalyzeFunc: func(ctx context.Context, req service.AnalyzeRequest) error {
				capturedCtx = ctx
				return nil
			},
		}
		handler := NewAnalyzeHandler(mockSvc)

		payload := AnalyzePayload{
			Owner: "owner",
			Repo:  "repo",
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		task := asynq.NewTask(TypeAnalyze, payloadBytes)
		ctx := context.WithValue(context.Background(), testKey, testValue)

		// When
		err = handler.ProcessTask(ctx, task)

		// Then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedCtx == nil {
			t.Fatal("context was not propagated to service")
		}
		if capturedCtx.Value(testKey) != testValue {
			t.Errorf("expected context value '%s', got '%v'", testValue, capturedCtx.Value(testKey))
		}
	})
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
