package usecase

import (
	"context"
	"testing"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

func TestDevLoginUseCase_Execute(t *testing.T) {
	t.Run("should return error when disabled", func(t *testing.T) {
		uc := NewDevLoginUseCase(
			false,
			&mockRefreshTokenRepository{},
			&mockRepository{},
			&mockTokenManager{},
		)

		_, err := uc.Execute(context.Background(), DevLoginInput{})
		if err != domain.ErrDevLoginDisabled {
			t.Errorf("expected ErrDevLoginDisabled, got %v", err)
		}
	})

	t.Run("should create new user when not exists", func(t *testing.T) {
		var createdUser *entity.User
		repo := &mockRepository{
			getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
				return nil, domain.ErrUserNotFound
			},
			createUserFunc: func(_ context.Context, user *entity.User) (string, error) {
				createdUser = user
				return user.ID, nil
			},
		}

		uc := NewDevLoginUseCase(
			true,
			&mockRefreshTokenRepository{},
			repo,
			&mockTokenManager{},
		)

		output, err := uc.Execute(context.Background(), DevLoginInput{
			Username: "custom-user",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if createdUser == nil {
			t.Fatal("expected user to be created")
		}
		if createdUser.Username != "custom-user" {
			t.Errorf("expected username 'custom-user', got %s", createdUser.Username)
		}
		if output.AccessToken == "" {
			t.Error("expected access token")
		}
		if output.RefreshToken == "" {
			t.Error("expected refresh token")
		}
	})

	t.Run("should use existing user when found", func(t *testing.T) {
		existingUser := &entity.User{
			ID:        "existing-user-id",
			Username:  "existing-user",
			AvatarURL: "https://example.com/avatar.png",
		}
		repo := &mockRepository{
			getUserByIDFunc: func(_ context.Context, id string) (*entity.User, error) {
				if id == "existing-user-id" {
					return existingUser, nil
				}
				return nil, domain.ErrUserNotFound
			},
		}

		uc := NewDevLoginUseCase(
			true,
			&mockRefreshTokenRepository{},
			repo,
			&mockTokenManager{},
		)

		output, err := uc.Execute(context.Background(), DevLoginInput{
			UserID: "existing-user-id",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.User.ID != "existing-user-id" {
			t.Errorf("expected user ID 'existing-user-id', got %s", output.User.ID)
		}
		if output.User.Username != "existing-user" {
			t.Errorf("expected username 'existing-user', got %s", output.User.Username)
		}
	})

	t.Run("should use default values when no input provided", func(t *testing.T) {
		var createdUser *entity.User
		repo := &mockRepository{
			getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
				return nil, domain.ErrUserNotFound
			},
			createUserFunc: func(_ context.Context, user *entity.User) (string, error) {
				createdUser = user
				return user.ID, nil
			},
		}

		uc := NewDevLoginUseCase(
			true,
			&mockRefreshTokenRepository{},
			repo,
			&mockTokenManager{},
		)

		_, err := uc.Execute(context.Background(), DevLoginInput{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if createdUser.ID != DefaultTestUserID {
			t.Errorf("expected default user ID, got %s", createdUser.ID)
		}
		if createdUser.Username != DefaultTestUserUsername {
			t.Errorf("expected default username, got %s", createdUser.Username)
		}
	})
}
