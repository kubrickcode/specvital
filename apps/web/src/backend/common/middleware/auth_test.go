package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authadapter "github.com/specvital/web/src/backend/modules/auth/adapter"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
)

const testSecret = "test-secret-that-is-at-least-32-characters-long"

func TestRequireAuth_ValidToken(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	token, _ := jwtManager.GenerateAccessToken("user-123", "testuser")

	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	var capturedUserID string
	handler := m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if capturedUserID != "user-123" {
		t.Errorf("expected userID user-123, got %s", capturedUserID)
	}
}

func TestRequireAuth_MissingCookie(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	handler := m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("expected Content-Type application/problem+json, got %s", ct)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	handler := m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "invalid-token"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("expected Content-Type application/problem+json, got %s", ct)
	}
}

func TestOptionalAuth_ValidToken(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	token, _ := jwtManager.GenerateAccessToken("user-123", "testuser")

	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	var capturedUserID string
	handler := m.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if capturedUserID != "user-123" {
		t.Errorf("expected userID user-123, got %s", capturedUserID)
	}
}

func TestOptionalAuth_MissingCookie(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	var capturedUserID string
	handler := m.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if capturedUserID != "" {
		t.Errorf("expected empty userID, got %s", capturedUserID)
	}
}

func TestOptionalAuth_InvalidToken(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	var capturedUserID string
	handler := m.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "invalid-token"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if capturedUserID != "" {
		t.Errorf("expected empty userID, got %s", capturedUserID)
	}
}

func TestOptionalAuth_NonAPIRoute_SkipsAuth(t *testing.T) {
	jwtManager, _ := authadapter.NewJWTTokenManager(testSecret)
	token, _ := jwtManager.GenerateAccessToken("user-123", "testuser")

	m := NewAuthMiddleware(jwtManager, "auth_token", "refresh_token")

	var capturedUserID string
	handler := m.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	// Non-API route should skip auth parsing entirely
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	// Claims should NOT be set for non-API routes
	if capturedUserID != "" {
		t.Errorf("expected empty userID for non-API route, got %s", capturedUserID)
	}
}

func TestGetClaims_NoClaims(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	claims := GetClaims(req.Context())

	if claims != nil {
		t.Error("expected nil claims")
	}
}

func TestWithClaims(t *testing.T) {
	claims := &entity.Claims{
		ExpiresAt: time.Now().Add(time.Hour),
		IssuedAt:  time.Now(),
		Issuer:    "specvital",
		Login:     "testuser",
		Subject:   "user-123",
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := WithClaims(req.Context(), claims)

	retrieved := GetClaims(ctx)
	if retrieved == nil {
		t.Fatal("expected claims to be set")
	}
	if retrieved.UserID() != "user-123" {
		t.Errorf("expected userID user-123, got %s", retrieved.UserID())
	}
	if retrieved.Login != "testuser" {
		t.Errorf("expected login testuser, got %s", retrieved.Login)
	}
}
