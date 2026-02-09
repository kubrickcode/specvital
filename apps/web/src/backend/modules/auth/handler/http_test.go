package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/domain/port"
	"github.com/specvital/web/src/backend/modules/auth/usecase"
	specviewhandler "github.com/specvital/web/src/backend/modules/spec-view/handler"
	subscriptionport "github.com/specvital/web/src/backend/modules/subscription/domain/port"
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
	generateAccessTokenFunc  func(userID, login string) (string, error)
	generateRefreshTokenFunc func() (*port.RefreshTokenResult, error)
	hashTokenFunc            func(token string) string
	validateAccessTokenFunc  func(tokenString string) (*entity.Claims, error)
}

func (m *mockTokenManager) GenerateAccessToken(userID, login string) (string, error) {
	if m.generateAccessTokenFunc != nil {
		return m.generateAccessTokenFunc(userID, login)
	}
	return "access-token-" + userID, nil
}

func (m *mockTokenManager) GenerateRefreshToken() (*port.RefreshTokenResult, error) {
	if m.generateRefreshTokenFunc != nil {
		return m.generateRefreshTokenFunc()
	}
	return &port.RefreshTokenResult{
		Token:     "mock-refresh-token",
		TokenHash: "mock-refresh-hash",
	}, nil
}

func (m *mockTokenManager) HashToken(token string) string {
	if m.hashTokenFunc != nil {
		return m.hashTokenFunc(token)
	}
	return "hashed-" + token
}

func (m *mockTokenManager) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	if m.validateAccessTokenFunc != nil {
		return m.validateAccessTokenFunc(tokenString)
	}
	return nil, nil
}

type mockRefreshTokenRepo struct {
	createFunc          func(ctx context.Context, token *entity.RefreshToken) (string, error)
	getByHashFunc       func(ctx context.Context, hash string) (*entity.RefreshToken, error)
	revokeFunc          func(ctx context.Context, id string) error
	revokeFamilyFunc    func(ctx context.Context, familyID string) error
	revokeUserTokenFunc func(ctx context.Context, userID string) error
	rotateTokenFunc     func(ctx context.Context, oldTokenID string, newToken *entity.RefreshToken) (string, error)
}

func (m *mockRefreshTokenRepo) Create(ctx context.Context, token *entity.RefreshToken) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, token)
	}
	return "refresh-token-id", nil
}

func (m *mockRefreshTokenRepo) GetByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	if m.getByHashFunc != nil {
		return m.getByHashFunc(ctx, hash)
	}
	return nil, nil
}

func (m *mockRefreshTokenRepo) Revoke(ctx context.Context, id string) error {
	if m.revokeFunc != nil {
		return m.revokeFunc(ctx, id)
	}
	return nil
}

func (m *mockRefreshTokenRepo) RevokeFamily(ctx context.Context, familyID string) error {
	if m.revokeFamilyFunc != nil {
		return m.revokeFamilyFunc(ctx, familyID)
	}
	return nil
}

func (m *mockRefreshTokenRepo) RevokeUserTokens(ctx context.Context, userID string) error {
	if m.revokeUserTokenFunc != nil {
		return m.revokeUserTokenFunc(ctx, userID)
	}
	return nil
}

func (m *mockRefreshTokenRepo) RotateToken(ctx context.Context, oldTokenID string, newToken *entity.RefreshToken) (string, error) {
	if m.rotateTokenFunc != nil {
		return m.rotateTokenFunc(ctx, oldTokenID, newToken)
	}
	return "new-token-id", nil
}

type mockSubscriber struct {
	assignDefaultPlanFunc func(ctx context.Context, userID string) error
}

var _ subscriptionport.Subscriber = (*mockSubscriber)(nil)

func (m *mockSubscriber) AssignDefaultPlan(ctx context.Context, userID string) error {
	if m.assignDefaultPlanFunc != nil {
		return m.assignDefaultPlanFunc(ctx, userID)
	}
	return nil
}

func createTestHandler(
	getCurrentUserUC *usecase.GetCurrentUserUseCase,
	handleOAuthCallbackUC *usecase.HandleOAuthCallbackUseCase,
	initiateOAuthUC *usecase.InitiateOAuthUseCase,
	refreshTokenUC *usecase.RefreshTokenUseCase,
) *Handler {
	handler, _ := NewHandler(&HandlerConfig{
		CookieSecure:        false,
		FrontendURL:         "http://localhost:5173",
		GetCurrentUser:      getCurrentUserUC,
		HandleOAuthCallback: handleOAuthCallbackUC,
		InitiateOAuth:       initiateOAuthUC,
		Logger:              logger.New(),
		RefreshToken:        refreshTokenUC,
	})
	return handler
}

func setupTestRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	apiHandlers := api.NewAPIHandlers(&mockAnalyzerHandler{}, user.NewMockHandler(), handler, user.NewMockHandler(), &mockGitHubHandler{}, &mockGitHubAppHandler{}, &mockPricingHandler{}, &mockRepositoryHandler{}, specviewhandler.NewMockHandler(), &mockSubscriptionHandler{}, &mockUsageHandler{}, user.NewMockHandler(), nil)
	strictHandler := api.NewStrictHandler(apiHandlers, nil)
	api.HandlerFromMux(strictHandler, r)
	return r
}

type mockAnalyzerHandler struct{}

func (m *mockAnalyzerHandler) AnalyzeRepository(_ context.Context, _ api.AnalyzeRepositoryRequestObject) (api.AnalyzeRepositoryResponseObject, error) {
	return nil, nil
}

func (m *mockAnalyzerHandler) GetAnalysisHistory(_ context.Context, _ api.GetAnalysisHistoryRequestObject) (api.GetAnalysisHistoryResponseObject, error) {
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

type mockGitHubAppHandler struct{}

func (m *mockGitHubAppHandler) GetGitHubAppInstallURL(_ context.Context, _ api.GetGitHubAppInstallURLRequestObject) (api.GetGitHubAppInstallURLResponseObject, error) {
	return api.GetGitHubAppInstallURL200JSONResponse{InstallURL: ""}, nil
}

func (m *mockGitHubAppHandler) GetUserGitHubAppInstallations(_ context.Context, _ api.GetUserGitHubAppInstallationsRequestObject) (api.GetUserGitHubAppInstallationsResponseObject, error) {
	return api.GetUserGitHubAppInstallations200JSONResponse{Data: []api.GitHubAppInstallation{}}, nil
}

type mockPricingHandler struct{}

func (m *mockPricingHandler) GetPricing(_ context.Context, _ api.GetPricingRequestObject) (api.GetPricingResponseObject, error) {
	return api.GetPricing200JSONResponse{Data: []api.PricingPlan{}}, nil
}

type mockSubscriptionHandler struct{}

func (m *mockSubscriptionHandler) GetUserSubscription(_ context.Context, _ api.GetUserSubscriptionRequestObject) (api.GetUserSubscriptionResponseObject, error) {
	return api.GetUserSubscription200JSONResponse{}, nil
}

type mockUsageHandler struct{}

func (m *mockUsageHandler) CheckQuota(_ context.Context, _ api.CheckQuotaRequestObject) (api.CheckQuotaResponseObject, error) {
	return api.CheckQuota200JSONResponse{}, nil
}

func (m *mockUsageHandler) GetCurrentUsage(_ context.Context, _ api.GetCurrentUsageRequestObject) (api.GetCurrentUsageResponseObject, error) {
	return api.GetCurrentUsage200JSONResponse{}, nil
}

func createDefaultUseCases() (*usecase.GetCurrentUserUseCase, *usecase.HandleOAuthCallbackUseCase, *usecase.InitiateOAuthUseCase, *usecase.RefreshTokenUseCase, *mockRepository, *mockOAuthClient, *mockStateStore, *mockRefreshTokenRepo) {
	repo := &mockRepository{}
	refreshTokenRepo := &mockRefreshTokenRepo{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{}
	encryptor := &mockEncryptor{}
	subscriber := &mockSubscriber{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, refreshTokenRepo, repo, stateStore, subscriber, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)
	refreshToken := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, tokenManager)

	return getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, repo, oauthClient, stateStore, refreshTokenRepo
}

func TestAuthLogin_Success(t *testing.T) {
	getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
	refreshTokenRepo := &mockRefreshTokenRepo{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{
		createFunc: func(_ context.Context) (string, error) {
			return "", domain.ErrInvalidState
		},
	}
	encryptor := &mockEncryptor{}
	subscriber := &mockSubscriber{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, refreshTokenRepo, repo, stateStore, subscriber, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)
	refreshToken := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, tokenManager)

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
	refreshTokenRepo := &mockRefreshTokenRepo{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{}
	encryptor := &mockEncryptor{}
	subscriber := &mockSubscriber{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, refreshTokenRepo, repo, stateStore, subscriber, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)
	refreshToken := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, tokenManager)

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
	refreshTokenRepo := &mockRefreshTokenRepo{}
	oauthClient := &mockOAuthClient{}
	stateStore := &mockStateStore{
		validateFunc: func(_ context.Context, _ string) error {
			return domain.ErrInvalidState
		},
	}
	encryptor := &mockEncryptor{}
	subscriber := &mockSubscriber{}
	tokenManager := &mockTokenManager{}

	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(encryptor, oauthClient, refreshTokenRepo, repo, stateStore, subscriber, tokenManager)
	initiateOAuth := usecase.NewInitiateOAuthUseCase(oauthClient, stateStore)
	refreshToken := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, tokenManager)

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
			getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
			handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
	getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
	getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)
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
	refreshTokenRepo := &mockRefreshTokenRepo{}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, refreshTokenRepo, repo, &mockStateStore{}, &mockSubscriber{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})
	refreshToken := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, &mockTokenManager{})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)

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
	refreshTokenRepo := &mockRefreshTokenRepo{}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, refreshTokenRepo, repo, &mockStateStore{}, &mockSubscriber{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})
	refreshToken := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, &mockTokenManager{})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)

	claims := &entity.Claims{}
	claims.Subject = "user-123"
	ctx := middleware.WithClaims(context.Background(), claims)
	req := api.AuthMeRequestObject{}

	resp, _ := handler.AuthMe(ctx, req)
	if _, ok := resp.(api.AuthMe200JSONResponse); !ok {
		t.Errorf("expected 200 response, got %T", resp)
	}
}

func TestBuildAccessCookie(t *testing.T) {
	getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
	handler, _ := NewHandler(&HandlerConfig{
		CookieSecure:        true,
		FrontendURL:         "http://localhost:5173",
		GetCurrentUser:      getCurrentUser,
		HandleOAuthCallback: handleOAuthCallback,
		InitiateOAuth:       initiateOAuth,
		Logger:              logger.New(),
		RefreshToken:        refreshToken,
	})

	cookie := handler.buildAccessCookie("test-token")
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

func TestBuildExpiredCookie(t *testing.T) {
	getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)

	cookie := handler.buildExpiredCookie(AccessCookieName)
	if cookie == "" {
		t.Error("expected cookie string to be non-empty")
	}

	if !strings.Contains(cookie, "Max-Age=") {
		t.Error("logout cookie should contain Max-Age")
	}
}

func TestAuthRefresh_Success(t *testing.T) {
	now := time.Now()
	repo := &mockRepository{
		getUserByIDFunc: func(_ context.Context, _ string) (*entity.User, error) {
			return &entity.User{ID: "user-123", Username: "testuser"}, nil
		},
	}
	refreshTokenRepo := &mockRefreshTokenRepo{
		getByHashFunc: func(_ context.Context, _ string) (*entity.RefreshToken, error) {
			return &entity.RefreshToken{
				ID:        "token-123",
				UserID:    "user-123",
				FamilyID:  "family-123",
				TokenHash: "hashed-valid-refresh-token",
				CreatedAt: now.Add(-time.Hour),
				ExpiresAt: now.Add(6 * 24 * time.Hour),
			}, nil
		},
	}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, refreshTokenRepo, repo, &mockStateStore{}, &mockSubscriber{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})
	refreshTokenUC := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, &mockTokenManager{
		hashTokenFunc: func(token string) string {
			return "hashed-" + token
		},
	})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshTokenUC)

	ctx := middleware.WithRefreshToken(context.Background(), "valid-refresh-token")
	req := api.AuthRefreshRequestObject{}

	resp, err := handler.AuthRefresh(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(refreshResponse); !ok {
		t.Errorf("expected refreshResponse, got %T", resp)
	}
}

func TestAuthRefresh_MissingToken(t *testing.T) {
	getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken, _, _, _, _ := createDefaultUseCases()
	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshToken)

	ctx := context.Background()
	req := api.AuthRefreshRequestObject{}

	resp, err := handler.AuthRefresh(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.AuthRefresh401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

func TestAuthRefresh_ExpiredToken(t *testing.T) {
	repo := &mockRepository{}
	refreshTokenRepo := &mockRefreshTokenRepo{
		getByHashFunc: func(_ context.Context, _ string) (*entity.RefreshToken, error) {
			return nil, domain.ErrRefreshTokenExpired
		},
	}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, refreshTokenRepo, repo, &mockStateStore{}, &mockSubscriber{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})
	refreshTokenUC := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, &mockTokenManager{
		hashTokenFunc: func(token string) string {
			return "hashed-" + token
		},
	})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshTokenUC)

	ctx := middleware.WithRefreshToken(context.Background(), "expired-token")
	req := api.AuthRefreshRequestObject{}

	resp, err := handler.AuthRefresh(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.AuthRefresh401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

func TestAuthRefresh_TokenReuseDetected(t *testing.T) {
	repo := &mockRepository{}
	refreshTokenRepo := &mockRefreshTokenRepo{
		getByHashFunc: func(_ context.Context, _ string) (*entity.RefreshToken, error) {
			return nil, domain.ErrTokenReuseDetected
		},
	}
	getCurrentUser := usecase.NewGetCurrentUserUseCase(repo)
	handleOAuthCallback := usecase.NewHandleOAuthCallbackUseCase(&mockEncryptor{}, &mockOAuthClient{}, refreshTokenRepo, repo, &mockStateStore{}, &mockSubscriber{}, &mockTokenManager{})
	initiateOAuth := usecase.NewInitiateOAuthUseCase(&mockOAuthClient{}, &mockStateStore{})
	refreshTokenUC := usecase.NewRefreshTokenUseCase(refreshTokenRepo, repo, &mockTokenManager{
		hashTokenFunc: func(token string) string {
			return "hashed-" + token
		},
	})

	handler := createTestHandler(getCurrentUser, handleOAuthCallback, initiateOAuth, refreshTokenUC)

	ctx := middleware.WithRefreshToken(context.Background(), "reused-token")
	req := api.AuthRefreshRequestObject{}

	resp, err := handler.AuthRefresh(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := resp.(api.AuthRefresh401ApplicationProblemPlusJSONResponse); !ok {
		t.Errorf("expected 401 response, got %T", resp)
	}
}

var _ port.Repository = (*mockRepository)(nil)
var _ port.OAuthClient = (*mockOAuthClient)(nil)
var _ port.StateStore = (*mockStateStore)(nil)
var _ port.Encryptor = (*mockEncryptor)(nil)
var _ port.TokenManager = (*mockTokenManager)(nil)
