package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
	subscriptionport "github.com/specvital/web/src/backend/modules/subscription/domain/port"
)

type HandleOAuthCallbackInput struct {
	Code  string
	State string
}

type HandleOAuthCallbackOutput struct {
	AccessToken  string
	RefreshToken string
	User         *entity.User
}

type HandleOAuthCallbackUseCase struct {
	encryptor        port.Encryptor
	oauthClient      port.OAuthClient
	refreshTokenRepo port.RefreshTokenRepository
	repository       port.Repository
	stateStore       port.StateStore
	subscriber       subscriptionport.Subscriber
	tokenManager     port.TokenManager
}

func NewHandleOAuthCallbackUseCase(
	encryptor port.Encryptor,
	oauthClient port.OAuthClient,
	refreshTokenRepo port.RefreshTokenRepository,
	repository port.Repository,
	stateStore port.StateStore,
	subscriber subscriptionport.Subscriber,
	tokenManager port.TokenManager,
) *HandleOAuthCallbackUseCase {
	if encryptor == nil {
		panic("encryptor is required")
	}
	if oauthClient == nil {
		panic("oauthClient is required")
	}
	if refreshTokenRepo == nil {
		panic("refreshTokenRepo is required")
	}
	if repository == nil {
		panic("repository is required")
	}
	if stateStore == nil {
		panic("stateStore is required")
	}
	if subscriber == nil {
		panic("subscriber is required")
	}
	if tokenManager == nil {
		panic("tokenManager is required")
	}
	return &HandleOAuthCallbackUseCase{
		encryptor:        encryptor,
		oauthClient:      oauthClient,
		refreshTokenRepo: refreshTokenRepo,
		repository:       repository,
		stateStore:       stateStore,
		subscriber:       subscriber,
		tokenManager:     tokenManager,
	}
}

func (uc *HandleOAuthCallbackUseCase) Execute(ctx context.Context, input HandleOAuthCallbackInput) (*HandleOAuthCallbackOutput, error) {
	if err := uc.stateStore.Validate(ctx, input.State); err != nil {
		if errors.Is(err, domain.ErrInvalidState) {
			return nil, err
		}
		return nil, fmt.Errorf("validate state: %w", err)
	}

	accessToken, err := uc.oauthClient.ExchangeCode(ctx, input.Code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	oauthUser, err := uc.oauthClient.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get oauth user info: %w", err)
	}

	encryptedToken, err := uc.encryptor.Encrypt(accessToken)
	if err != nil {
		return nil, fmt.Errorf("encrypt token: %w", err)
	}

	user, err := uc.findOrCreateUser(ctx, oauthUser, encryptedToken)
	if err != nil {
		return nil, fmt.Errorf("find or create user: %w", err)
	}

	if err := uc.repository.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("update last login: %w", err)
	}

	jwtToken, err := uc.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
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

	return &HandleOAuthCallbackOutput{
		AccessToken:  jwtToken,
		RefreshToken: refreshTokenResult.Token,
		User:         user,
	}, nil
}

func (uc *HandleOAuthCallbackUseCase) findOrCreateUser(ctx context.Context, oauthUser *entity.OAuthUserInfo, encryptedToken string) (*entity.User, error) {
	account, err := uc.repository.GetOAuthAccountByProvider(ctx, entity.ProviderGitHub, oauthUser.ExternalID)
	if err == nil {
		if err := uc.updateOAuthAccount(ctx, account, oauthUser, encryptedToken); err != nil {
			return nil, err
		}
		return uc.repository.GetUserByID(ctx, account.UserID)
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		return uc.createNewUser(ctx, oauthUser, encryptedToken)
	}

	return nil, fmt.Errorf("get oauth account by provider: %w", err)
}

func (uc *HandleOAuthCallbackUseCase) updateOAuthAccount(ctx context.Context, account *entity.OAuthAccount, oauthUser *entity.OAuthUserInfo, encryptedToken string) error {
	account.AccessToken = &encryptedToken
	account.ProviderUsername = &oauthUser.Username

	if _, err := uc.repository.UpsertOAuthAccount(ctx, account); err != nil {
		return fmt.Errorf("update oauth account: %w", err)
	}
	return nil
}

// createNewUser creates a user, OAuth account, and assigns default subscription.
// Note: These operations are not wrapped in a transaction. If AssignDefaultPlan fails,
// the user exists without a subscription. This is handled via backfill mechanism
// (GetUsersWithoutActiveSubscription query) to ensure eventual consistency.
func (uc *HandleOAuthCallbackUseCase) createNewUser(ctx context.Context, oauthUser *entity.OAuthUserInfo, encryptedToken string) (*entity.User, error) {
	newUser := &entity.User{
		AvatarURL: oauthUser.AvatarURL,
		Email:     oauthUser.Email,
		Username:  oauthUser.Username,
	}

	userID, err := uc.repository.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	oauthAccount := &entity.OAuthAccount{
		AccessToken:      &encryptedToken,
		Provider:         entity.ProviderGitHub,
		ProviderUserID:   oauthUser.ExternalID,
		ProviderUsername: &oauthUser.Username,
		UserID:           userID,
	}

	if _, err := uc.repository.UpsertOAuthAccount(ctx, oauthAccount); err != nil {
		return nil, fmt.Errorf("create oauth account: %w", err)
	}

	if err := uc.subscriber.AssignDefaultPlan(ctx, userID); err != nil {
		return nil, fmt.Errorf("assign default plan: %w", err)
	}

	newUser.ID = userID
	return newUser, nil
}
