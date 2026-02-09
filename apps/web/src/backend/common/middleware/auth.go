package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/specvital/web/src/backend/common/httputil"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
)

type contextKey string

const (
	claimsKey       contextKey = "claims"
	clientIPKey     contextKey = "client_ip"
	refreshTokenKey contextKey = "refresh_token"
)

func GetClaims(ctx context.Context) *entity.Claims {
	claims, _ := ctx.Value(claimsKey).(*entity.Claims)
	return claims
}

func GetRefreshToken(ctx context.Context) string {
	token, _ := ctx.Value(refreshTokenKey).(string)
	return token
}

func GetClientIP(ctx context.Context) string {
	ip, _ := ctx.Value(clientIPKey).(string)
	return ip
}

func GetUserID(ctx context.Context) string {
	claims := GetClaims(ctx)
	if claims == nil {
		return ""
	}
	return claims.UserID()
}

func WithClaims(ctx context.Context, claims *entity.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func withClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

func WithRefreshToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, refreshTokenKey, token)
}

type AuthMiddleware struct {
	accessCookieName  string
	refreshCookieName string
	tokenManager      port.TokenManager
}

func NewAuthMiddleware(tokenManager port.TokenManager, accessCookieName, refreshCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		accessCookieName:  accessCookieName,
		refreshCookieName: refreshCookieName,
		tokenManager:      tokenManager,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if refreshCookie, err := r.Cookie(m.refreshCookieName); err == nil {
			ctx = WithRefreshToken(ctx, refreshCookie.Value)
		}

		claims, err := m.extractClaims(r)
		if err != nil {
			writeUnauthorized(w, "authentication required")
			return
		}

		ctx = WithClaims(ctx, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAPIRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		ctx = withClientIP(ctx, httputil.ExtractClientIP(r))

		claims, err := m.extractClaims(r)
		if err == nil && claims != nil {
			ctx = WithClaims(ctx, claims)
		}

		if refreshCookie, err := r.Cookie(m.refreshCookieName); err == nil {
			ctx = WithRefreshToken(ctx, refreshCookie.Value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isAPIRoute checks if the path should have authentication parsing.
// Currently uses simple prefix matching. If routing becomes more complex,
// consider using a route matcher or configuration-based approach.
func isAPIRoute(path string) bool {
	return strings.HasPrefix(path, "/api")
}

func (m *AuthMiddleware) extractClaims(r *http.Request) (*entity.Claims, error) {
	cookie, err := r.Cookie(m.accessCookieName)
	if err != nil {
		return nil, err
	}

	claims, err := m.tokenManager.ValidateAccessToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func writeUnauthorized(w http.ResponseWriter, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(api.NewUnauthorized(detail))
}
