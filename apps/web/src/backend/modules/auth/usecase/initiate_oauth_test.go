package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

type mockOAuthClient struct {
	exchangeCodeFunc    func(ctx context.Context, code string) (string, error)
	generateAuthURLFunc func(state string) (string, error)
	getUserInfoFunc     func(ctx context.Context, accessToken string) (*entity.OAuthUserInfo, error)
}

func (m *mockOAuthClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	if m.exchangeCodeFunc != nil {
		return m.exchangeCodeFunc(ctx, code)
	}
	return "", nil
}

func (m *mockOAuthClient) GenerateAuthURL(state string) (string, error) {
	if m.generateAuthURLFunc != nil {
		return m.generateAuthURLFunc(state)
	}
	return "", nil
}

func (m *mockOAuthClient) GetUserInfo(ctx context.Context, accessToken string) (*entity.OAuthUserInfo, error) {
	if m.getUserInfoFunc != nil {
		return m.getUserInfoFunc(ctx, accessToken)
	}
	return nil, nil
}

type mockStateStore struct {
	createFunc   func(ctx context.Context) (string, error)
	validateFunc func(ctx context.Context, state string) error
	closeFunc    func() error
}

func (m *mockStateStore) Create(ctx context.Context) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx)
	}
	return "", nil
}

func (m *mockStateStore) Validate(ctx context.Context, state string) error {
	if m.validateFunc != nil {
		return m.validateFunc(ctx, state)
	}
	return nil
}

func (m *mockStateStore) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestInitiateOAuthUseCase_Execute(t *testing.T) {
	tests := []struct {
		name            string
		setupOAuth      func() *mockOAuthClient
		setupStateStore func() *mockStateStore
		wantErr         bool
		wantAuthURL     string
	}{
		{
			name: "should return auth URL",
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{
					generateAuthURLFunc: func(state string) (string, error) {
						return "https://github.com/login/oauth/authorize?state=" + state, nil
					},
				}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					createFunc: func(_ context.Context) (string, error) {
						return "test-state", nil
					},
				}
			},
			wantErr:     false,
			wantAuthURL: "https://github.com/login/oauth/authorize?state=test-state",
		},
		{
			name: "should return error when state creation fails",
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					createFunc: func(_ context.Context) (string, error) {
						return "", errors.New("state creation failed")
					},
				}
			},
			wantErr:     true,
			wantAuthURL: "",
		},
		{
			name: "should return error when auth URL generation fails",
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{
					generateAuthURLFunc: func(_ string) (string, error) {
						return "", errors.New("URL generation failed")
					},
				}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					createFunc: func(_ context.Context) (string, error) {
						return "test-state", nil
					},
				}
			},
			wantErr:     true,
			wantAuthURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oauth := tt.setupOAuth()
			stateStore := tt.setupStateStore()
			uc := NewInitiateOAuthUseCase(oauth, stateStore)

			output, err := uc.Execute(context.Background(), InitiateOAuthInput{})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output.AuthURL != tt.wantAuthURL {
				t.Errorf("expected auth URL %q, got %q", tt.wantAuthURL, output.AuthURL)
			}
		})
	}
}
