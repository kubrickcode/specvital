package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenUseCase struct {
	refreshTokenRepo port.RefreshTokenRepository
	repository       port.Repository
	tokenManager     port.TokenManager
}

func NewRefreshTokenUseCase(
	refreshTokenRepo port.RefreshTokenRepository,
	repository port.Repository,
	tokenManager port.TokenManager,
) *RefreshTokenUseCase {
	if refreshTokenRepo == nil {
		panic("refreshTokenRepo is required")
	}
	if repository == nil {
		panic("repository is required")
	}
	if tokenManager == nil {
		panic("tokenManager is required")
	}
	return &RefreshTokenUseCase{
		refreshTokenRepo: refreshTokenRepo,
		repository:       repository,
		tokenManager:     tokenManager,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	if input.RefreshToken == "" {
		return nil, domain.ErrRefreshTokenNotFound
	}

	tokenHash := uc.tokenManager.HashToken(input.RefreshToken)

	storedToken, err := uc.refreshTokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if err := uc.validateToken(ctx, storedToken); err != nil {
		return nil, err
	}

	user, err := uc.repository.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	newRefreshTokenResult, err := uc.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	now := time.Now()
	newStoredToken := &entity.RefreshToken{
		CreatedAt: now,
		ExpiresAt: now.Add(domain.RefreshTokenExpiry),
		FamilyID:  storedToken.FamilyID,
		Replaces:  &storedToken.ID,
		TokenHash: newRefreshTokenResult.TokenHash,
		UserID:    storedToken.UserID,
	}

	if _, err := uc.refreshTokenRepo.RotateToken(ctx, storedToken.ID, newStoredToken); err != nil {
		return nil, fmt.Errorf("rotate refresh token: %w", err)
	}

	accessToken, err := uc.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenResult.Token,
	}, nil
}

func (uc *RefreshTokenUseCase) validateToken(ctx context.Context, token *entity.RefreshToken) error {
	if token.IsRevoked() {
		if err := uc.refreshTokenRepo.RevokeFamily(ctx, token.FamilyID); err != nil {
			return fmt.Errorf("revoke token family: %w", err)
		}
		return domain.ErrTokenReuseDetected
	}

	if token.IsExpired() {
		return domain.ErrRefreshTokenExpired
	}

	return nil
}

func (uc *RefreshTokenUseCase) RevokeToken(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}

	tokenHash := uc.tokenManager.HashToken(rawToken)
	storedToken, err := uc.refreshTokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			return nil
		}
		return fmt.Errorf("get refresh token: %w", err)
	}

	if err := uc.refreshTokenRepo.Revoke(ctx, storedToken.ID); err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}
