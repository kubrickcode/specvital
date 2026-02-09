package port

import (
	"context"

	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) (string, error)
	GetByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeFamily(ctx context.Context, familyID string) error
	RevokeUserTokens(ctx context.Context, userID string) error
	RotateToken(ctx context.Context, oldTokenID string, newToken *entity.RefreshToken) (string, error)
}
