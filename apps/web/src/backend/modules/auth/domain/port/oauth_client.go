package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

type OAuthClient interface {
	ExchangeCode(ctx context.Context, code string) (string, error)
	GenerateAuthURL(state string) (string, error)
	GetUserInfo(ctx context.Context, accessToken string) (*entity.OAuthUserInfo, error)
}
