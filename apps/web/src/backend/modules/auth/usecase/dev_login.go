package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

const (
	DefaultTestUserID       = "00000000-0000-0000-0000-000000000001"
	DefaultTestUserUsername = "test-user"
)

type DevLoginInput struct {
	UserID   string
	Username string
}

type DevLoginOutput struct {
	AccessToken  string
	RefreshToken string
	User         *entity.User
}

type DevLoginUseCase struct {
	isEnabled        bool
	refreshTokenRepo port.RefreshTokenRepository
	repository       port.Repository
	tokenManager     port.TokenManager
}

func NewDevLoginUseCase(
	isEnabled bool,
	refreshTokenRepo port.RefreshTokenRepository,
	repository port.Repository,
	tokenManager port.TokenManager,
) *DevLoginUseCase {
	if refreshTokenRepo == nil {
		panic("refreshTokenRepo is required")
	}
	if repository == nil {
		panic("repository is required")
	}
	if tokenManager == nil {
		panic("tokenManager is required")
	}
	return &DevLoginUseCase{
		isEnabled:        isEnabled,
		refreshTokenRepo: refreshTokenRepo,
		repository:       repository,
		tokenManager:     tokenManager,
	}
}

func (uc *DevLoginUseCase) Execute(ctx context.Context, input DevLoginInput) (*DevLoginOutput, error) {
	if !uc.isEnabled {
		return nil, domain.ErrDevLoginDisabled
	}

	user, err := uc.findOrCreateTestUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("find or create test user: %w", err)
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenResult, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	now := time.Now()
	refreshToken := &entity.RefreshToken{
		CreatedAt: now,
		ExpiresAt: now.Add(domain.RefreshTokenExpiry),
		FamilyID:  uuid.New().String(),
		TokenHash: refreshTokenResult.TokenHash,
		UserID:    user.ID,
	}

	if _, err := uc.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	return &DevLoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenResult.Token,
		User:         user,
	}, nil
}

func (uc *DevLoginUseCase) findOrCreateTestUser(ctx context.Context, input DevLoginInput) (*entity.User, error) {
	userID := input.UserID
	if userID == "" {
		userID = DefaultTestUserID
	}

	user, err := uc.repository.GetUserByID(ctx, userID)
	if err == nil {
		return user, nil
	}

	if err != domain.ErrUserNotFound {
		return nil, err
	}

	username := input.Username
	if username == "" {
		username = DefaultTestUserUsername
	}

	newUser := &entity.User{
		ID:        userID,
		AvatarURL: "https://avatars.githubusercontent.com/u/0?v=4",
		Username:  username,
	}

	createdID, err := uc.repository.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	newUser.ID = createdID
	return newUser, nil
}
