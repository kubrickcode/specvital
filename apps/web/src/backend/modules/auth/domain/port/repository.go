package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entity.User) (string, error)
	GetOAuthAccountByProvider(ctx context.Context, provider, providerUserID string) (*entity.OAuthAccount, error)
	GetOAuthAccountByUserID(ctx context.Context, userID, provider string) (*entity.OAuthAccount, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
	UpsertOAuthAccount(ctx context.Context, account *entity.OAuthAccount) (string, error)
}
