package usecase

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type GetUserGitHubTokenInput struct {
	UserID string
}

type GetUserGitHubTokenOutput struct {
	Token string
}

type GetUserGitHubTokenUseCase struct {
	encryptor  port.Encryptor
	repository port.Repository
}

func NewGetUserGitHubTokenUseCase(
	encryptor port.Encryptor,
	repository port.Repository,
) *GetUserGitHubTokenUseCase {
	if encryptor == nil {
		panic("encryptor is required")
	}
	if repository == nil {
		panic("repository is required")
	}
	return &GetUserGitHubTokenUseCase{
		encryptor:  encryptor,
		repository: repository,
	}
}

func (uc *GetUserGitHubTokenUseCase) Execute(ctx context.Context, input GetUserGitHubTokenInput) (*GetUserGitHubTokenOutput, error) {
	account, err := uc.repository.GetOAuthAccountByUserID(ctx, input.UserID, entity.ProviderGitHub)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get oauth account by user id: %w", err)
	}

	if account.AccessToken == nil {
		return nil, domain.ErrNoGitHubToken
	}

	decrypted, err := uc.encryptor.Decrypt(*account.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("decrypt token: %w", err)
	}

	return &GetUserGitHubTokenOutput{Token: decrypted}, nil
}
