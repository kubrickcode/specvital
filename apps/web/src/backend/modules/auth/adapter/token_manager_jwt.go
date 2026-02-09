package adapter

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/golang-jwt/jwt/v5"

	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

const (
	Issuer            = "specvital"
	RefreshTokenBytes = 32
)

type JWTTokenManager struct {
	accessExpiry time.Duration
	secret       []byte
}

var _ port.TokenManager = (*JWTTokenManager)(nil)

func NewJWTTokenManager(secret string) (*JWTTokenManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if len(secret) < 32 {
		return nil, errors.New("jwt secret must be at least 32 characters")
	}
	return &JWTTokenManager{
		accessExpiry: domain.AccessTokenExpiry,
		secret:       []byte(secret),
	}, nil
}

func (m *JWTTokenManager) GenerateAccessToken(userID, login string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}
	if login == "" {
		return "", errors.New("login is required")
	}

	now := time.Now()
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiry)),
		},
		Login: login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTTokenManager) GenerateRefreshToken() (*port.RefreshTokenResult, error) {
	tokenBytes := make([]byte, RefreshTokenBytes)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, errors.Wrap(err, "generate random bytes")
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	tokenHash := m.HashToken(token)

	return &port.RefreshTokenResult{
		Token:     token,
		TokenHash: tokenHash,
	}, nil
}

func (m *JWTTokenManager) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (m *JWTTokenManager) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrapf(domain.ErrInvalidToken, "unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}
		return nil, errors.Wrapf(domain.ErrInvalidToken, "token validation failed: %v", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	if claims.Issuer != Issuer {
		return nil, errors.Wrapf(domain.ErrInvalidToken, "invalid issuer: expected %s", Issuer)
	}

	if claims.ExpiresAt == nil || claims.IssuedAt == nil {
		return nil, errors.Wrap(domain.ErrInvalidToken, "missing required time claims")
	}

	return &entity.Claims{
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
		Issuer:    claims.Issuer,
		Login:     claims.Login,
		Subject:   claims.Subject,
	}, nil
}

type jwtClaims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func (c *jwtClaims) UserID() string {
	return c.Subject
}
