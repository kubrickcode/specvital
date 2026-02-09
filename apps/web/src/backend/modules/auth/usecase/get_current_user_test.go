package usecase

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

type mockRepository struct {
	getUserByIDFunc               func(ctx context.Context, id string) (*entity.User, error)
	createUserFunc                func(ctx context.Context, user *entity.User) (string, error)
	getOAuthAccountByProviderFunc func(ctx context.Context, provider, providerUserID string) (*entity.OAuthAccount, error)
	getOAuthAccountByUserIDFunc   func(ctx context.Context, userID, provider string) (*entity.OAuthAccount, error)
	updateLastLoginFunc           func(ctx context.Context, userID string) error
	upsertOAuthAccountFunc        func(ctx context.Context, account *entity.OAuthAccount) (string, error)
}

func (m *mockRepository) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	return "", nil
}

func (m *mockRepository) GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*entity.OAuthAccount, error) {
	if m.getOAuthAccountByProviderFunc != nil {
		return m.getOAuthAccountByProviderFunc(ctx, provider, providerUserID)
	}
	return nil, nil
}

func (m *mockRepository) GetOAuthAccountByUserID(ctx context.Context, userID, provider string) (*entity.OAuthAccount, error) {
	if m.getOAuthAccountByUserIDFunc != nil {
		return m.getOAuthAccountByUserIDFunc(ctx, userID, provider)
	}
	return nil, nil
}

func (m *mockRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	if m.updateLastLoginFunc != nil {
		return m.updateLastLoginFunc(ctx, userID)
	}
	return nil
}

func (m *mockRepository) UpsertOAuthAccount(ctx context.Context, account *entity.OAuthAccount) (string, error) {
	if m.upsertOAuthAccountFunc != nil {
		return m.upsertOAuthAccountFunc(ctx, account)
	}
	return "", nil
}

func TestGetCurrentUserUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     GetCurrentUserInput
		setupRepo func() *mockRepository
		wantErr   error
		wantUser  bool
	}{
		{
			name:  "should return user when found",
			input: GetCurrentUserInput{UserID: "user-123"},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getUserByIDFunc: func(_ context.Context, id string) (*entity.User, error) {
						if id == "user-123" {
							return &entity.User{
								ID:       "user-123",
								Username: "testuser",
							}, nil
						}
						return nil, domain.ErrUserNotFound
					},
				}
			},
			wantErr:  nil,
			wantUser: true,
		},
		{
			name:  "should return error when user not found",
			input: GetCurrentUserInput{UserID: "unknown"},
			setupRepo: func() *mockRepository {
				return &mockRepository{
					getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
						return nil, domain.ErrUserNotFound
					},
				}
			},
			wantErr:  domain.ErrUserNotFound,
			wantUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			uc := NewGetCurrentUserUseCase(repo)

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

			if tt.wantUser && output.User == nil {
				t.Error("expected user, got nil")
			}
		})
	}
}
