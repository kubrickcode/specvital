package handler

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
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
	"github.com/specvital/web/src/backend/modules/auth/usecase"
	"github.com/specvital/web/src/backend/modules/user"
)

type mockRepository struct {
	createUserFunc                func(ctx context.Context, user *entity.User) (string, error)
	getOAuthAccountByProviderFunc func(ctx context.Context, provider, externalID string) (*entity.OAuthAccount, error)
	getUserByIDFunc               func(ctx context.Context, id string) (*entity.User, error)
}

func (m *mockRepository) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	return "new-user-id", nil
}

func (m *mockRepository) GetOAuthAccountByProvider(ctx context.Context, provider, externalID string) (*entity.OAuthAccount, error) {
	if m.getOAuthAccountByProviderFunc != nil {
		return m.getOAuthAccountByProviderFunc(ctx, provider, externalID)
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockRepository) GetOAuthAccountByUserID(_ context.Context, _, _ string) (*entity.OAuthAccount, error) {
	return nil, nil
}

func (m *mockRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockRepository) UpdateLastLogin(_ context.Context, _ string) error {
	return nil
}

func (m *mockRepository) UpsertOAuthAccount(_ context.Context, _ *entity.OAuthAccount) (string, error) {
	return "", nil
}

type mockOAuthClient struct {
	exchangeCodeFunc    func(ctx context.Context, code string) (string, error)
	generateAuthURLFunc func(state string) (string, error)
	getUserInfoFunc     func(ctx context.Context, accessToken string) (*entity.OAuthUserInfo, error)
}

func (m *mockOAuthClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	if m.exchangeCodeFunc != nil {
		return m.exchangeCodeFunc(ctx, code)
	}
	return "access-token", nil
}

func (m *mockOAuthClient) GenerateAuthURL(state string) (string, error) {
	if m.generateAuthURLFunc != nil {
		return m.generateAuthURLFunc(state)
	}
	return "https://github.com/login/oauth/authorize?state=" + state, nil
}

func (m *mockOAuthClient) GetUserInfo(ctx context.Context, accessToken string) (*entity.OAuthUserInfo, error) {
	if m.getUserInfoFunc != nil {
		return m.getUserInfoFunc(ctx, accessToken)
	}
	return &entity.OAuthUserInfo{
		ExternalID: "12345",
		Username:   "testuser",
		AvatarURL:  "https://github.com/avatar.png",
	}, nil
}

type mockStateStore struct {
	createFunc   func(ctx context.Context) (string, error)
	validateFunc func(ctx context.Context, state string) error
}

func (m *mockStateStore) Create(ctx context.Context) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx)
	}
	return "test-state", nil
}

func (m *mockStateStore) Validate(ctx context.Context, state string) error {
	if m.validateFunc != nil {
		return m.validateFunc(ctx, state)
	}
	return nil
}

func (m *mockStateStore) Close() error {
	return nil
}

type mockEncryptor struct {
	encryptFunc func(plaintext string) (string, error)
	decryptFunc func(ciphertext string) (string, error)
}

func (m *mockEncryptor) Encrypt(plaintext string) (string, error) {
	if m.encryptFunc != nil {
		return m.encryptFunc(plaintext)
	}
	return "encrypted-" + plaintext, nil
}

func (m *mockEncryptor) Decrypt(ciphertext string) (string, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(ciphertext)
	}
	return ciphertext, nil
}

type mockTokenManager struct {
	generateFunc func(userID, login string) (string, error)
	validateFunc func(tokenString string) (*entity.Claims, error)
}

func (m *mockTokenManager) Generate(userID, login string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(userID, login)
	}
	return "jwt-token-" + userID, nil
}

func (m *mockTokenManager) Validate(tokenString string) (*entity.Claims, error) {
	if m.validateFunc != nil {
		return m.validateFunc(tokenString)
	}
	return nil, nil
}

func createTestHandler(
	getCurrentUserUC *usecase.GetCurrentUserUseCase,
	handleOAuthCallbackUC *usecase.HandleOAuthCallbackUseCase,
	initiateOAuthUC *usecase.InitiateOAuthUseCase,
) *Handler {
	handler, _ := NewHandler(&HandlerConfig{
		CookieSecure:        false,
		FrontendURL:         "http://localhost:5173",
		GetCurrentUser:      getCurrentUserUC,
		HandleOAuthCallback: handleOAuthCallbackUC,
		InitiateOAuth:       initiateOAuthUC,
		Logger:              logger.New(),
	})
	return handler
}

func setupTestRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	apiHandlers := api.NewAPIHandlers(&mockAnalyzerHandler{}, user.NewMockHandler(), handler, user.NewMockHandler(), &mockGitHubHandler{}, &mockRepositoryHandler{}, nil)
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

type mockGitHubHandler struct{}

func (m *mockGitHubHandler) GetOrganizationRepositories(_ context.Context, _ api.GetOrganizationRepositoriesRequestObject) (api.GetOrganizationRepositoriesResponseObject, error) {
	return nil, nil
}

func (m *mockGitHubHandler) GetUserGitHubOrganizations(_ context.Context, _ api.GetUserGitHubOrganizationsRequestObject) (api.GetUserGitHubOrganizationsResponseObject, error) {
	return nil, nil
}

func (m *mockGitHubHandler) GetUserGitHubRepositories(_ context.Context, _ api.GetUserGitHubRepositoriesRequestObject) (api.GetUserGitHubRepositoriesResponseObject, error) {
	return nil, nil
}

func createDefaultUseCases() (*usecase.GetCurrentUserUseCase, *usecase.HandleOAuthCallbackUseCase, *usecase.InitiateOAuthUseCase, *mockRepository, *mockOAuthClient, *mockStateStore) {
	repo := &mockRepository{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{}
	encryptor := &mockEncryptor{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, repo, stateStore, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)

	return getCurrentUser, handleOAuthCallback, initiateOAuth, repo, oauthClient, stateStore
}

func TestAuthLogin_Success(t *testing.T) {
	getCurrentUser, handleOAuthCallback, initiateOAuth, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
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
	repo := &mockRepository{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{
		createFunc: func(_ context.Context) (string, error) {
			return "", domain.ErrInvalidState
		},
	}
	encryptor := &mockEncryptor{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, repo, stateStore, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestAuthCallback_Success(t *testing.T) {
	repo := &mockRepository{
		getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
			return &entity.User{
				ID:        "user-123",
				Username:  "testuser",
				AvatarURL: "https://github.com/avatar.png",
			}, nil
		},
	}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{}
	encryptor := &mockEncryptor{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, repo, stateStore, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
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
	repo := &mockRepository{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{
		validateFunc: func(_ context.Context, _ string) error {
			return domain.ErrInvalidState
		},
	}
	encryptor := &mockEncryptor{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, repo, stateStore, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback?code=test-code&state=invalid-state", nil)
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
			getCurrentUser, handleOAuthCallback, initiateOAuth, _, _, _ := createDefaultUseCases()
			handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
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
	getCurrentUser, handleOAuthCallback, initiateOAuth, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
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
	getCurrentUser, handleOAuthCallback, initiateOAuth, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMe_UserNotFound(t *testing.T) {
	repo := &mockRepository{
		getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, repo, &mockStateStore{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)

	claims := &entity.Claims{}
	claims.Subject = "user-123"
	ctx := middleware.WithClaims(context.Background(), claims)
	req := api.AuthMeRequestObject{}

	resp, _ := handler.AuthMe(ctx, req)
	if _, ok := resp.(api.AuthMe401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

func TestAuthMe_Success(t *testing.T) {
	repo := &mockRepository{
		getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
			return &entity.User{
				ID:        "user-123",
				Username:  "testuser",
				AvatarURL: "https://github.com/avatar.png",
			}, nil
		},
	}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, repo, &mockStateStore{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)

	claims := &entity.Claims{}
	claims.Subject = "user-123"
	ctx := middleware.WithClaims(context.Background(), claims)
	req := api.AuthMeRequestObject{}

	resp, _ := handler.AuthMe(ctx, req)
	if _, ok := resp.(api.AuthMe200JSONResponse); !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}
}

func TestBuildAuthCookie(t *testing.T) {
	getCurrentUser, handleOAuthCallback, initiateOAuth, _, _, _ := createDefaultUseCases()
	handler, _ := NewHandler(&HandlerConfig{
		CookieSecure:        true,
		FrontendURL:         "http://localhost:5173",
		GetCurrentUser:      getCurrentUser,
		HandleOAuthCallback: handleOAuthCallback,
		InitiateOAuth:       initiateOAuth,
		Logger:              logger.New(),
	})

	cookie := handler.buildAuthCookie("test-token")
	if cookie == "" {
		t.Error("expected cookie string to be non-empty")
	}

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
	getCurrentUser, handleOAuthCallback, initiateOAuth, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth)

	cookie := handler.BuildLogoutCookie()
	if cookie == "" {
		t.Error("expected cookie string to be non-empty")
	}

	if !strings.Contains(cookie, "Max-Age=") {
		t.Error("logout cookie should contain Max-Age")
	}
}

var _ port.Repository = (*mockRepository)(nil)
var _ port.OAuthClient = (*mockOAuthClient)(nil)
var _ port.StateStore = (*mockStateStore)(nil)
var _ port.Encryptor = (*mockEncryptor)(nil)
var _ port.TokenManager = (*mockTokenManager)(nil)
