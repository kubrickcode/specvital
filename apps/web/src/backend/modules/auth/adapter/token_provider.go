package adapter

import (
	"context"

	"github.com/specvital/web/src/backend/modules/auth/usecase"
)

type TokenProviderAdapter struct {
	usecase *usecase.GetUserGitHubTokenUseCase
}

func NewTokenProviderAdapter(uc *usecase.GetUserGitHubTokenUseCase) *TokenProviderAdapter {
	if uc == nil {
		panic("usecase is required")
	}
	return &TokenProviderAdapter{usecase: uc}
}

func (a *TokenProviderAdapter) GetUserGitHubToken(ctx context.Context, userID string) (string, error) {
	output, err := a.usecase.Execute(ctx, usecase.GetUserGitHubTokenInput{UserID: userID})
	if err != nil {
		return "", err
	}
	return output.Token, nil
}
