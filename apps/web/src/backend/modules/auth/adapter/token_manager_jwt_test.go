package adapter

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/specvital/web/src/backend/modules/auth/domain"
)

const testSecret = "test-secret-key-must-be-at-least-32-chars"

func TestNewJWTTokenManager(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "valid secret",
			secret:  testSecret,
			wantErr: false,
		},
		{
			name:    "empty secret",
			secret:  "",
			wantErr: true,
		},
		{
			name:    "short secret",
			secret:  "short",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewJWTTokenManager(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJWTTokenManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWTTokenManager_GenerateAccessToken(t *testing.T) {
	m, err := NewJWTTokenManager(testSecret)
	if err != nil {
		t.Fatalf("NewJWTTokenManager() error = %v", err)
	}

	userID := "user-123"
	login := "testuser"

	token, err := m.GenerateAccessToken(userID, login)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateAccessToken() returned empty token")
	}

	claims, err := m.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID() != userID {
		t.Errorf("claims.UserID() = %v, want %v", claims.UserID(), userID)
	}
	if claims.Login != login {
		t.Errorf("claims.Login = %v, want %v", claims.Login, login)
	}
	if claims.Issuer != Issuer {
		t.Errorf("claims.Issuer = %v, want %v", claims.Issuer, Issuer)
	}

	expectedExpiry := time.Now().Add(domain.AccessTokenExpiry)
	if claims.ExpiresAt.Before(expectedExpiry.Add(-time.Minute)) || claims.ExpiresAt.After(expectedExpiry.Add(time.Minute)) {
		t.Errorf("claims.ExpiresAt = %v, expected around %v", claims.ExpiresAt, expectedExpiry)
	}
}

func TestJWTTokenManager_GenerateAccessToken_EmptyInputs(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	tests := []struct {
		login  string
		name   string
		userID string
	}{
		{name: "empty userID", userID: "", login: "testuser"},
		{name: "empty login", userID: "user-123", login: ""},
		{name: "both empty", userID: "", login: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.GenerateAccessToken(tt.userID, tt.login)
			if err == nil {
				t.Errorf("GenerateAccessToken() should return error for %s", tt.name)
			}
		})
	}
}

func TestJWTTokenManager_ValidateAccessToken_InvalidToken(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "invalid token format",
			token:   "invalid.token.here",
			wantErr: domain.ErrInvalidToken,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: domain.ErrInvalidToken,
		},
		{
			name:    "wrong signature",
			token:   generateTokenWithSecret("wrong-secret-key-must-be-32-chars-long"),
			wantErr: domain.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.ValidateAccessToken(tt.token)
			if err == nil {
				t.Error("ValidateAccessToken() expected error, got nil")
			}
		})
	}
}

func TestJWTTokenManager_ValidateAccessToken_ExpiredToken(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   "user-123",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
		Login: "testuser",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	_, err := m.ValidateAccessToken(tokenString)
	if err != domain.ErrTokenExpired {
		t.Errorf("ValidateAccessToken() error = %v, want %v", err, domain.ErrTokenExpired)
	}
}

func TestJWTTokenManager_ValidateAccessToken_InvalidIssuer(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "other-issuer",
			Subject:   "user-123",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Login: "testuser",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	_, err := m.ValidateAccessToken(tokenString)
	if err == nil {
		t.Error("ValidateAccessToken() should reject token with invalid issuer")
	}
}

func TestJWTTokenManager_GenerateRefreshToken(t *testing.T) {
	m, err := NewJWTTokenManager(testSecret)
	if err != nil {
		t.Fatalf("NewJWTTokenManager() error = %v", err)
	}

	result, err := m.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if result.Token == "" {
		t.Error("GenerateRefreshToken() returned empty token")
	}

	if result.TokenHash == "" {
		t.Error("GenerateRefreshToken() returned empty hash")
	}

	if result.Token == result.TokenHash {
		t.Error("token and hash should be different")
	}

	recomputedHash := m.HashToken(result.Token)
	if recomputedHash != result.TokenHash {
		t.Errorf("hash mismatch: got %v, want %v", recomputedHash, result.TokenHash)
	}
}

func TestJWTTokenManager_GenerateRefreshToken_Uniqueness(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	tokens := make(map[string]bool)
	hashes := make(map[string]bool)

	for i := 0; i < 100; i++ {
		result, err := m.GenerateRefreshToken()
		if err != nil {
			t.Fatalf("GenerateRefreshToken() error = %v", err)
		}

		if tokens[result.Token] {
			t.Errorf("duplicate token generated: %s", result.Token)
		}
		tokens[result.Token] = true

		if hashes[result.TokenHash] {
			t.Errorf("duplicate hash generated: %s", result.TokenHash)
		}
		hashes[result.TokenHash] = true
	}
}

func TestJWTTokenManager_HashToken(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	token := "test-refresh-token-value"
	hash1 := m.HashToken(token)
	hash2 := m.HashToken(token)

	if hash1 != hash2 {
		t.Errorf("same token should produce same hash: %v != %v", hash1, hash2)
	}

	differentHash := m.HashToken("different-token")
	if hash1 == differentHash {
		t.Error("different tokens should produce different hashes")
	}
}

func TestJWTTokenManager_HashToken_Empty(t *testing.T) {
	m, _ := NewJWTTokenManager(testSecret)

	hash := m.HashToken("")
	if hash == "" {
		t.Error("empty token should still produce a hash")
	}
}

func generateTokenWithSecret(secret string) string {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   "user-123",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Login: "testuser",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
