package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
	subscriptionport "github.com/specvital/web/src/backend/modules/subscription/domain/port"
)

type mockTokenManager struct {
	generateAccessTokenFunc  func(userID, login string) (string, error)
	generateRefreshTokenFunc func() (*port.RefreshTokenResult, error)
	hashTokenFunc            func(token string) string
	validateAccessTokenFunc  func(tokenString string) (*entity.Claims, error)
}

func (m *mockTokenManager) GenerateAccessToken(userID, login string) (string, error) {
	if m.generateAccessTokenFunc != nil {
		return m.generateAccessTokenFunc(userID, login)
	}
	return "mock-access-token", nil
}

func (m *mockTokenManager) GenerateRefreshToken() (*port.RefreshTokenResult, error) {
	if m.generateRefreshTokenFunc != nil {
		return m.generateRefreshTokenFunc()
	}
	return &port.RefreshTokenResult{
		Token:     "mock-refresh-token",
		TokenHash: "mock-token-hash",
	}, nil
}

func (m *mockTokenManager) HashToken(token string) string {
	if m.hashTokenFunc != nil {
		return m.hashTokenFunc(token)
	}
	return "hashed-" + token
}

func (m *mockTokenManager) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	if m.validateAccessTokenFunc != nil {
		return m.validateAccessTokenFunc(tokenString)
	}
	return nil, nil
}

type mockRefreshTokenRepo struct {
	createFunc          func(ctx context.Context, token *entity.RefreshToken) (string, error)
	getByHashFunc       func(ctx context.Context, hash string) (*entity.RefreshToken, error)
	revokeFunc          func(ctx context.Context, id string) error
	revokeFamilyFunc    func(ctx context.Context, familyID string) error
	revokeUserTokenFunc func(ctx context.Context, userID string) error
	rotateTokenFunc     func(ctx context.Context, oldTokenID string, newToken *entity.RefreshToken) (string, error)
}

func (m *mockRefreshTokenRepo) Create(ctx context.Context, token *entity.RefreshToken) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, token)
	}
	return "refresh-token-id", nil
}

func (m *mockRefreshTokenRepo) GetByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	if m.getByHashFunc != nil {
		return m.getByHashFunc(ctx, hash)
	}
	return nil, nil
}

func (m *mockRefreshTokenRepo) Revoke(ctx context.Context, id string) error {
	if m.revokeFunc != nil {
		return m.revokeFunc(ctx, id)
	}
	return nil
}

func (m *mockRefreshTokenRepo) RevokeFamily(ctx context.Context, familyID string) error {
	if m.revokeFamilyFunc != nil {
		return m.revokeFamilyFunc(ctx, familyID)
	}
	return nil
}

func (m *mockRefreshTokenRepo) RevokeUserTokens(ctx context.Context, userID string) error {
	if m.revokeUserTokenFunc != nil {
		return m.revokeUserTokenFunc(ctx, userID)
	}
	return nil
}

func (m *mockRefreshTokenRepo) RotateToken(ctx context.Context, oldTokenID string, newToken *entity.RefreshToken) (string, error) {
	if m.rotateTokenFunc != nil {
		return m.rotateTokenFunc(ctx, oldTokenID, newToken)
	}
	return "new-token-id", nil
}

type mockSubscriber struct {
	assignDefaultPlanFunc func(ctx context.Context, userID string) error
}

var _ subscriptionport.Subscriber = (*mockSubscriber)(nil)

func (m *mockSubscriber) AssignDefaultPlan(ctx context.Context, userID string) error {
	if m.assignDefaultPlanFunc != nil {
		return m.assignDefaultPlanFunc(ctx, userID)
	}
	return nil
}

func TestHandleOAuthCallbackUseCase_Execute(t *testing.T) {
	tests := []struct {
		name                  string
		input                 HandleOAuthCallbackInput
		setupEncryptor        func() *mockEncryptor
		setupOAuth            func() *mockOAuthClient
		setupRefreshTokenRepo func() *mockRefreshTokenRepo
		setupRepo             func() *mockRepository
		setupStateStore       func() *mockStateStore
		setupSubscriber       func() *mockSubscriber
		setupToken            func() *mockTokenManager
		wantErr               error
		wantAnyErr            bool
		wantAccessToken       string
	}{
		{
			name: "should handle callback for new user",
			input: HandleOAuthCallbackInput{
				Code:  "valid-code",
				State: "valid-state",
			},
			setupEncryptor: func() *mockEncryptor {
				return &mockEncryptor{
					encryptFunc: func(plaintext string) (string, error) {
						return "encrypted-" + plaintext, nil
					},
				}
			},
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{
					exchangeCodeFunc: func(_ context.Context, code string) (string, error) {
						if code == "valid-code" {
							return "access-token", nil
						}
						return "", errors.New("invalid code")
					},
					getUserInfoFunc: func(_ context.Context, _ string) (*entity.OAuthUserInfo, error) {
						return &entity.OAuthUserInfo{
							ExternalID: "12345",
							Username:   "testuser",
							AvatarURL:  "https://github.com/avatar",
						}, nil
					},
				}
			},
			setupRefreshTokenRepo: func() *mockRefreshTokenRepo {
				return &mockRefreshTokenRepo{}
			},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getOAuthAccountByProviderFunc: func(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
						return nil, domain.ErrUserNotFound
					},
					createUserFunc: func(_ context.Context, _ *entity.User) (string, error) {
						return "new-user-id", nil
					},
					upsertOAuthAccountFunc: func(_ context.Context, _ *entity.OAuthAccount) (string, error) {
						return "oauth-id", nil
					},
					updateLastLoginFunc: func(_ context.Context, _ string) error {
						return nil
					},
				}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					validateFunc: func(_ context.Context, state string) error {
						if state == "valid-state" {
							return nil
						}
						return domain.ErrInvalidState
					},
				}
			},
			setupSubscriber: func() *mockSubscriber {
				return &mockSubscriber{}
			},
			setupToken: func() *mockTokenManager {
				return &mockTokenManager{
					generateAccessTokenFunc: func(userID, login string) (string, error) {
						return "jwt-token-" + userID, nil
					},
				}
			},
			wantErr:         nil,
			wantAccessToken: "jwt-token-new-user-id",
		},
		{
			name: "should handle callback for existing user",
			input: HandleOAuthCallbackInput{
				Code:  "valid-code",
				State: "valid-state",
			},
			setupEncryptor: func() *mockEncryptor {
				return &mockEncryptor{
					encryptFunc: func(plaintext string) (string, error) {
						return "encrypted-" + plaintext, nil
					},
				}
			},
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{
					exchangeCodeFunc: func(_ context.Context, _ string) (string, error) {
						return "access-token", nil
					},
					getUserInfoFunc: func(_ context.Context, _ string) (*entity.OAuthUserInfo, error) {
						return &entity.OAuthUserInfo{
							ExternalID: "12345",
							Username:   "existinguser",
							AvatarURL:  "https://github.com/avatar",
						}, nil
					},
				}
			},
			setupRefreshTokenRepo: func() *mockRefreshTokenRepo {
				return &mockRefreshTokenRepo{}
			},
			setupRepo: func() *mockRepository {
				token := "old-token"
				return &mockRepository{
					getOAuthAccountByProviderFunc: func(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
						return &entity.OAuthAccount{
							ID:          "oauth-id",
							UserID:      "existing-user-id",
							AccessToken: &token,
						}, nil
					},
					getUserByIDFunc: func(_ context.Context, id string) (*entity.User, error) {
						if id == "existing-user-id" {
							return &entity.User{
								ID:       "existing-user-id",
								Username: "existinguser",
							}, nil
						}
						return nil, domain.ErrUserNotFound
					},
					upsertOAuthAccountFunc: func(_ context.Context, _ *entity.OAuthAccount) (string, error) {
						return "oauth-id", nil
					},
					updateLastLoginFunc: func(_ context.Context, _ string) error {
						return nil
					},
				}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					validateFunc: func(_ context.Context, _ string) error {
						return nil
					},
				}
			},
			setupSubscriber: func() *mockSubscriber {
				return &mockSubscriber{}
			},
			setupToken: func() *mockTokenManager {
				return &mockTokenManager{
					generateAccessTokenFunc: func(userID, _ string) (string, error) {
						return "jwt-token-" + userID, nil
					},
				}
			},
			wantErr:         nil,
			wantAccessToken: "jwt-token-existing-user-id",
		},
		{
			name: "should return error for invalid state",
			input: HandleOAuthCallbackInput{
				Code:  "valid-code",
				State: "invalid-state",
			},
			setupEncryptor: func() *mockEncryptor {
				return &mockEncryptor{}
			},
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{}
			},
			setupRefreshTokenRepo: func() *mockRefreshTokenRepo {
				return &mockRefreshTokenRepo{}
			},
			setupRepo: func() *mockRepository {
				return &mockRepository{}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					validateFunc: func(_ context.Context, _ string) error {
						return domain.ErrInvalidState
					},
				}
			},
			setupSubscriber: func() *mockSubscriber {
				return &mockSubscriber{}
			},
			setupToken: func() *mockTokenManager {
				return &mockTokenManager{}
			},
			wantErr:         domain.ErrInvalidState,
			wantAccessToken: "",
		},
		{
			name: "should return error when assign default plan fails for new user",
			input: HandleOAuthCallbackInput{
				Code:  "valid-code",
				State: "valid-state",
			},
			setupEncryptor: func() *mockEncryptor {
				return &mockEncryptor{
					encryptFunc: func(plaintext string) (string, error) {
						return "encrypted-" + plaintext, nil
					},
				}
			},
			setupOAuth: func() *mockOAuthClient {
				return &mockOAuthClient{
					exchangeCodeFunc: func(_ context.Context, code string) (string, error) {
						return "access-token", nil
					},
					getUserInfoFunc: func(_ context.Context, _ string) (*entity.OAuthUserInfo, error) {
						return &entity.OAuthUserInfo{
							ExternalID: "12345",
							Username:   "newuser",
							AvatarURL:  "https://github.com/avatar",
						}, nil
					},
				}
			},
			setupRefreshTokenRepo: func() *mockRefreshTokenRepo {
				return &mockRefreshTokenRepo{}
			},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getOAuthAccountByProviderFunc: func(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
						return nil, domain.ErrUserNotFound
					},
					createUserFunc: func(_ context.Context, _ *entity.User) (string, error) {
						return "new-user-id", nil
					},
					upsertOAuthAccountFunc: func(_ context.Context, _ *entity.OAuthAccount) (string, error) {
						return "oauth-id", nil
					},
				}
			},
			setupStateStore: func() *mockStateStore {
				return &mockStateStore{
					validateFunc: func(_ context.Context, _ string) error {
						return nil
					},
				}
			},
			setupSubscriber: func() *mockSubscriber {
				return &mockSubscriber{
					assignDefaultPlanFunc: func(_ context.Context, _ string) error {
						return errors.New("subscription service unavailable")
					},
				}
			},
			setupToken: func() *mockTokenManager {
				return &mockTokenManager{}
			},
			wantAnyErr:      true,
			wantAccessToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.setupEncryptor()
			oauth := tt.setupOAuth()
			refreshTokenRepo := tt.setupRefreshTokenRepo()
			repo := tt.setupRepo()
			stateStore := tt.setupStateStore()
			subscriber := tt.setupSubscriber()
			tokenMgr := tt.setupToken()

			uc := NewHandleOAuthCallbackUseCase(enc, oauth, refreshTokenRepo, repo, stateStore, subscriber, tokenMgr)

			output, err := uc.Execute(context.Background(), tt.input)

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

			if tt.wantAnyErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output.AccessToken != tt.wantAccessToken {
				t.Errorf("expected access token %q, got %q", tt.wantAccessToken, output.AccessToken)
			}

			if output.RefreshToken == "" {
				t.Error("expected refresh token, got empty")
			}

			if output.User == nil {
				t.Error("expected user, got nil")
			}
		})
	}
}
