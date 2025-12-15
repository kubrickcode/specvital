package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/specvital/web/src/backend/modules/auth/domain"
)

const testSecret = "test-secret-key-must-be-at-least-32-chars"

func TestNewManager(t *testing.T) {
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
			_, err := NewManager(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GenerateAndValidate(t *testing.T) {
	m, err := NewManager(testSecret)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	userID := "user-123"
	login := "testuser"

	token, err := m.Generate(userID, login)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if token == "" {
		t.Error("Generate() returned empty token")
	}

	claims, err := m.Validate(token)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
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
}

func TestManager_ValidateInvalidToken(t *testing.T) {
	m, _ := NewManager(testSecret)

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
			_, err := m.Validate(tt.token)
			if err == nil {
				t.Error("Validate() expected error, got nil")
			}
		})
	}
}

func TestManager_ValidateExpiredToken(t *testing.T) {
	m, _ := NewManager(testSecret)

	// Create expired token
	claims := domain.Claims{
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

	_, err := m.Validate(tokenString)
	if err != domain.ErrTokenExpired {
		t.Errorf("Validate() error = %v, want %v", err, domain.ErrTokenExpired)
	}
}

func generateTokenWithSecret(secret string) string {
	claims := domain.Claims{
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
