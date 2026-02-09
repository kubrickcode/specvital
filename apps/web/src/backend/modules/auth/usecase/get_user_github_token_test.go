package usecase

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

type mockEncryptor struct {
	encryptFunc func(plaintext string) (string, error)
	decryptFunc func(ciphertext string) (string, error)
}

func (m *mockEncryptor) Encrypt(plaintext string) (string, error) {
	if m.encryptFunc != nil {
		return m.encryptFunc(plaintext)
	}
	return plaintext, nil
}

func (m *mockEncryptor) Decrypt(ciphertext string) (string, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(ciphertext)
	}
	return ciphertext, nil
}

func TestGetUserGitHubTokenUseCase_Execute(t *testing.T) {
	encryptedToken := "encrypted-token"
	decryptedToken := "decrypted-token"

	tests := []struct {
		name      string
		input     GetUserGitHubTokenInput
		setupRepo func() *mockRepository
		setupEnc  func() *mockEncryptor
		wantErr   error
		wantToken string
	}{
		{
			name:  "should return decrypted token",
			input: GetUserGitHubTokenInput{UserID: "user-123"},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getOAuthAccountByUserIDFunc: func(_ context.Context, userID, provider string) (*entity.OAuthAccount, error) {
						if userID == "user-123" && provider == entity.ProviderGitHub {
							return &entity.OAuthAccount{
								AccessToken: &encryptedToken,
							}, nil
						}
						return nil, domain.ErrUserNotFound
					},
				}
			},
			setupEnc: func() *mockEncryptor {
				return &mockEncryptor{
					decryptFunc: func(ciphertext string) (string, error) {
						if ciphertext == encryptedToken {
							return decryptedToken, nil
						}
						return "", nil
					},
				}
			},
			wantErr:   nil,
			wantToken: decryptedToken,
		},
		{
			name:  "should return error when user not found",
			input: GetUserGitHubTokenInput{UserID: "unknown"},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getOAuthAccountByUserIDFunc: func(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
						return nil, domain.ErrUserNotFound
					},
				}
			},
			setupEnc: func() *mockEncryptor {
				return &mockEncryptor{}
			},
			wantErr:   domain.ErrUserNotFound,
			wantToken: "",
		},
		{
			name:  "should return error when no github token",
			input: GetUserGitHubTokenInput{UserID: "user-no-token"},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getOAuthAccountByUserIDFunc: func(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
						return &entity.OAuthAccount{
							AccessToken: nil,
						}, nil
					},
				}
			},
			setupEnc: func() *mockEncryptor {
				return &mockEncryptor{}
			},
			wantErr:   domain.ErrNoGitHubToken,
			wantToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			enc := tt.setupEnc()
			uc := NewGetUserGitHubTokenUseCase(enc, repo)

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

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if output.Token != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, output.Token)
			}
		})
	}
}
