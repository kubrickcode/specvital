package port

import "github.com/specvital/web/src/backend/modules/auth/domain/entity"

// RefreshTokenResult contains the generated refresh token and its hash.
// The raw Token is sent to the client, while TokenHash is stored in the database.
type RefreshTokenResult struct {
	Token     string
	TokenHash string
}

type TokenManager interface {
	GenerateAccessToken(userID, login string) (string, error)
	GenerateRefreshToken() (*RefreshTokenResult, error)
	HashToken(token string) string
	ValidateAccessToken(tokenString string) (*entity.Claims, error)
}
