package usecase

import (
	"context"
	"fmt"

	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type InitiateOAuthInput struct{}

type InitiateOAuthOutput struct {
	AuthURL string
}

type InitiateOAuthUseCase struct {
	oauthClient port.OAuthClient
	stateStore  port.StateStore
}

func NewInitiateOAuthUseCase(oauthClient port.OAuthClient, stateStore port.StateStore) *InitiateOAuthUseCase {
	if oauthClient == nil {
		panic("oauthClient is required")
	}
	if stateStore == nil {
		panic("stateStore is required")
	}
	return &InitiateOAuthUseCase{
		oauthClient: oauthClient,
		stateStore:  stateStore,
	}
}

func (uc *InitiateOAuthUseCase) Execute(ctx context.Context, _ InitiateOAuthInput) (*InitiateOAuthOutput, error) {
	state, err := uc.stateStore.Create(ctx)
	if err != nil {
		return nil, fmt.Errorf("create state: %w", err)
	}

	authURL, err := uc.oauthClient.GenerateAuthURL(state)
	if err != nil {
		return nil, fmt.Errorf("generate auth url: %w", err)
	}

	return &InitiateOAuthOutput{AuthURL: authURL}, nil
}
