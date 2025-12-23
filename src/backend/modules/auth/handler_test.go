package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/domain"
)

type mockService struct {
	authResult      *AuthResult
	currentUser     *domain.User
	err             error
	githubToken     string
	initiateAuthURL string
}

func (m *mockService) GetCurrentUser(_ context.Context, _ string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.currentUser, nil
}

func (m *mockService) GetUserGitHubToken(_ context.Context, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.githubToken, nil
}

func (m *mockService) HandleOAuthCallback(_ context.Context, _, _ string) (*AuthResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.authResult, nil
}

func (m *mockService) InitiateOAuth(_ context.Context) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.initiateAuthURL, nil
}

var _ Service = (*mockService)(nil)

func setupTestRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	apiHandlers := api.NewAPIHandlers(&mockAnalyzerHandler{}, handler, NewMockBookmarkHandler(), &mockRepositoryHandler{})
	strictHandler := api.NewStrictHandler(apiHandlers, nil)
	api.HandlerFromMux(strictHandler, r)
	return r
}

type mockAnalyzerHandler struct{}

func (m *mockAnalyzerHandler) AnalyzeRepository(_ context.Context, _ api.AnalyzeRepositoryRequestObject) (api.AnalyzeRepositoryResponseObject, error) {
	return nil, nil
}

func (m *mockAnalyzerHandler) GetAnalysisStatus(_ context.Context, _ api.GetAnalysisStatusRequestObject) (api.GetAnalysisStatusResponseObject, error) {
	return nil, nil
}

type mockRepositoryHandler struct{}

func (m *mockRepositoryHandler) GetRecentRepositories(_ context.Context, _ api.GetRecentRepositoriesRequestObject) (api.GetRecentRepositoriesResponseObject, error) {
	return nil, nil
}

func (m *mockRepositoryHandler) GetRepositoryStats(_ context.Context, _ api.GetRepositoryStatsRequestObject) (api.GetRepositoryStatsResponseObject, error) {
	return nil, nil
}

func (m *mockRepositoryHandler) GetUpdateStatus(_ context.Context, _ api.GetUpdateStatusRequestObject) (api.GetUpdateStatusResponseObject, error) {
	return nil, nil
}

func (m *mockRepositoryHandler) ReanalyzeRepository(_ context.Context, _ api.ReanalyzeRepositoryRequestObject) (api.ReanalyzeRepositoryResponseObject, error) {
	return nil, nil
}

func TestAuthLogin_Success(t *testing.T) {
	svc := &mockService{
		initiateAuthURL: "https://github.com/login/oauth/authorize?client_id=test",
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestAuthLogin_ServiceError(t *testing.T) {
	svc := &mockService{
		err: domain.ErrInvalidState,
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestAuthCallback_Success(t *testing.T) {
	svc := &mockService{
		authResult: &AuthResult{
			Token: "test-jwt-token",
			User: &domain.User{
				AvatarURL: "https://github.com/avatar.png",
				ID:        "user-123",
				Username:  "testuser",
			},
		},
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback?code=test-code&state=test-state", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("expected status 302, got %d", rec.Code)
	}

	cookie := rec.Header().Get("Set-Cookie")
	if cookie == "" {
		t.Error("expected Set-Cookie header to be set")
	}

	location := rec.Header().Get("Location")
	if location != "http://localhost:5173" {
		t.Errorf("expected Location header to be http://localhost:5173, got %s", location)
	}
}

func TestAuthCallback_InvalidState(t *testing.T) {
	svc := &mockService{
		err: domain.ErrInvalidState,
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback?code=test-code&state=invalid-state", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestAuthCallback_InvalidCode(t *testing.T) {
	svc := &mockService{
		err: domain.ErrInvalidOAuthCode,
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback?code=invalid-code&state=test-state", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestAuthCallback_MissingParameters(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"missing code", "/api/auth/callback?state=test-state"},
		{"missing state", "/api/auth/callback?code=test-code"},
		{"both missing", "/api/auth/callback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewHandler(&HandlerConfig{
				CookieSecure: false,
				FrontendURL:  "http://localhost:5173",
				Logger:       logger.New(),
				Service:      &mockService{},
			})
			if err != nil {
				t.Fatalf("failed to create handler: %v", err)
			}
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rec.Code)
			}
		})
	}
}

func TestAuthLogout_Success(t *testing.T) {
	svc := &mockService{}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	cookie := rec.Header().Get("Set-Cookie")
	if cookie == "" {
		t.Error("expected Set-Cookie header to be set for logout")
	}
	if !strings.Contains(cookie, "Max-Age=") {
		t.Error("logout cookie should contain Max-Age for expiration")
	}
}

func TestAuthMe_Unauthenticated(t *testing.T) {
	svc := &mockService{}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMe_UserNotFound(t *testing.T) {
	svc := &mockService{
		err: domain.ErrUserNotFound,
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	// Use direct handler call with claims in context
	claims := &domain.Claims{}
	claims.Subject = "user-123"
	ctx := middleware.WithClaims(context.Background(), claims)
	req := api.AuthMeRequestObject{}

	resp, _ := handler.AuthMe(ctx, req)
	if _, ok := resp.(api.AuthMe401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

func TestAuthMe_Success(t *testing.T) {
	svc := &mockService{
		currentUser: &domain.User{
			AvatarURL: "https://github.com/avatar.png",
			ID:        "user-123",
			Username:  "testuser",
		},
	}
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      svc,
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	// Use direct handler call with claims in context
	claims := &domain.Claims{}
	claims.Subject = "user-123"
	ctx := middleware.WithClaims(context.Background(), claims)
	req := api.AuthMeRequestObject{}

	resp, _ := handler.AuthMe(ctx, req)
	if _, ok := resp.(api.AuthMe200JSONResponse); !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}
}

func TestBuildAuthCookie(t *testing.T) {
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: true,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      &mockService{},
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	cookie := handler.buildAuthCookie("test-token")
	if cookie == "" {
		t.Error("expected cookie string to be non-empty")
	}

	// Check cookie attributes
	tests := []struct {
		name     string
		contains string
	}{
		{"HttpOnly", "HttpOnly"},
		{"SameSite", "SameSite=Lax"},
		{"Secure", "Secure"},
		{"Path", "Path=/"},
	}

	for _, tc := range tests {
		if !strings.Contains(cookie, tc.contains) {
			t.Errorf("cookie should contain %s, got: %s", tc.contains, cookie)
		}
	}
}

func TestBuildLogoutCookie(t *testing.T) {
	handler, err := NewHandler(&HandlerConfig{
		CookieSecure: false,
		FrontendURL:  "http://localhost:5173",
		Logger:       logger.New(),
		Service:      &mockService{},
	})
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	cookie := handler.BuildLogoutCookie()
	if cookie == "" {
		t.Error("expected cookie string to be non-empty")
	}

	// Logout cookie should have Max-Age=0 or negative
	if !strings.Contains(cookie, "Max-Age=") {
		t.Error("logout cookie should contain Max-Age")
	}
}
