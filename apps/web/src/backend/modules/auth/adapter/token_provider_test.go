package adapter

import (
	"context"
	"testing"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/usecase"
)

type mockRepository struct {
	getOAuthAccountByUserIDFunc func(ctx context.Context, userID, provider string) (*entity.OAuthAccount, error)
}

func (m *mockRepository) CreateUser(_ context.Context, _ *entity.User) (string, error) {
	return "", nil
}

func (m *mockRepository) GetOAuthAccountByProvider(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
	return nil, nil
}

func (m *mockRepository) GetOAuthAccountByUserID(ctx context.Context, userID, provider string) (*entity.OAuthAccount, error) {
	if m.getOAuthAccountByUserIDFunc != nil {
		return m.getOAuthAccountByUserIDFunc(ctx, userID, provider)
	}
	return nil, nil
}

func (m *mockRepository) GetUserByID(_ context.Context, _ string) (*entity.User, error) {
	return nil, nil
}

func (m *mockRepository) UpdateLastLogin(_ context.Context, _ string) error {
	return nil
}

func (m *mockRepository) UpsertOAuthAccount(_ context.Context, _ *entity.OAuthAccount) (string, error) {
	return "", nil
}

type mockEncryptor struct {
	decryptFunc func(ciphertext string) (string, error)
}

func (m *mockEncryptor) Encrypt(_ string) (string, error) {
	return "", nil
}

func (m *mockEncryptor) Decrypt(ciphertext string) (string, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(ciphertext)
	}
	return ciphertext, nil
}

func TestTokenProviderAdapter_GetUserGitHubToken(t *testing.T) {
	encryptedToken := "encrypted-token"
	decryptedToken := "decrypted-token"

	repo := &mockRepository{
		getOAuthAccountByUserIDFunc: func(_ context.Context, userID, _ string) (*entity.OAuthAccount, error) {
			if userID == "user-123" {
				return &entity.OAuthAccount{
					AccessToken: &encryptedToken,
				}, nil
			}
			return nil, domain.ErrUserNotFound
		},
	}
	enc := &mockEncryptor{
		decryptFunc: func(ciphertext string) (string, error) {
			if ciphertext == encryptedToken {
				return decryptedToken, nil
			}
			return "", nil
		},
	}

	uc := usecase.NewGetUserGitHubTokenUseCase(enc, repo)
	adapter := NewTokenProviderAdapter(uc)

	token, err := adapter.GetUserGitHubToken(context.Background(), "user-123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token != decryptedToken {
		t.Errorf("expected token %q, got %q", decryptedToken, token)
	}
}

func TestNewTokenProviderAdapter_PanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got nil")
		}
	}()

	NewTokenProviderAdapter(nil)
}
