package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/specvital/web/src/backend/common/logger"
	"github.com/specvital/web/src/backend/common/middleware"
	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/auth/adapter/mapper"
	"github.com/specvital/web/src/backend/modules/auth/domain"
	"github.com/specvital/web/src/backend/modules/auth/domain/entity"
	"github.com/specvital/web/src/backend/modules/auth/usecase"
)

const (
	AccessCookieName  = "auth_token"
	RefreshCookieName = "refresh_token"
	CookiePath        = "/"
)

var (
	AccessCookieMaxAge  = int(domain.AccessTokenExpiry.Seconds())
	RefreshCookieMaxAge = int(domain.RefreshTokenExpiry.Seconds())
)

type Handler struct {
	cookieDomain        string
	cookieSecure        bool
	devLogin            *usecase.DevLoginUseCase
	frontendURL         string
	getCurrentUser      *usecase.GetCurrentUserUseCase
	handleOAuthCallback *usecase.HandleOAuthCallbackUseCase
	initiateOAuth       *usecase.InitiateOAuthUseCase
	logger              *logger.Logger
	refreshToken        *usecase.RefreshTokenUseCase
}

type HandlerConfig struct {
	CookieDomain        string
	CookieSecure        bool
	DevLogin            *usecase.DevLoginUseCase
	FrontendURL         string
	GetCurrentUser      *usecase.GetCurrentUserUseCase
	HandleOAuthCallback *usecase.HandleOAuthCallbackUseCase
	InitiateOAuth       *usecase.InitiateOAuthUseCase
	Logger              *logger.Logger
	RefreshToken        *usecase.RefreshTokenUseCase
}

var _ api.AuthHandlers = (*Handler)(nil)

func NewHandler(cfg *HandlerConfig) (*Handler, error) {
	if cfg == nil {
		return nil, errors.New("handler config is required")
	}
	if cfg.FrontendURL == "" {
		return nil, errors.New("frontend URL is required")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.GetCurrentUser == nil {
		return nil, errors.New("GetCurrentUser usecase is required")
	}
	if cfg.HandleOAuthCallback == nil {
		return nil, errors.New("HandleOAuthCallback usecase is required")
	}
	if cfg.InitiateOAuth == nil {
		return nil, errors.New("InitiateOAuth usecase is required")
	}
	if cfg.RefreshToken == nil {
		return nil, errors.New("RefreshToken usecase is required")
	}
	return &Handler{
		cookieDomain:        cfg.CookieDomain,
		cookieSecure:        cfg.CookieSecure,
		devLogin:            cfg.DevLogin,
		frontendURL:         cfg.FrontendURL,
		getCurrentUser:      cfg.GetCurrentUser,
		handleOAuthCallback: cfg.HandleOAuthCallback,
		initiateOAuth:       cfg.InitiateOAuth,
		logger:              cfg.Logger,
		refreshToken:        cfg.RefreshToken,
	}, nil
}

func (h *Handler) AuthLogin(ctx context.Context, _ api.AuthLoginRequestObject) (api.AuthLoginResponseObject, error) {
	output, err := h.initiateOAuth.Execute(ctx, usecase.InitiateOAuthInput{})
	if err != nil {
		h.logger.Error(ctx, "failed to initiate oauth", "error", err)
		return api.AuthLogin500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to initiate login"),
		}, nil
	}

	return api.AuthLogin200JSONResponse(mapper.ToLoginResponse(output.AuthURL)), nil
}

func (h *Handler) AuthCallback(ctx context.Context, request api.AuthCallbackRequestObject) (api.AuthCallbackResponseObject, error) {
	code := request.Params.Code
	state := request.Params.State

	if code == "" || state == "" {
		return api.AuthCallback400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("missing code or state parameter"),
		}, nil
	}

	result, err := h.handleOAuthCallback.Execute(ctx, usecase.HandleOAuthCallbackInput{
		Code:  code,
		State: state,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidState) {
			return api.AuthCallback400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid or expired state parameter"),
			}, nil
		}
		if errors.Is(err, domain.ErrInvalidOAuthCode) {
			return api.AuthCallback400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: api.NewBadRequest("invalid authorization code"),
			}, nil
		}

		h.logger.Error(ctx, "oauth callback failed", "error", err)
		return api.AuthCallback500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("authentication failed"),
		}, nil
	}

	accessCookie := h.buildAccessCookie(result.AccessToken)
	refreshCookie := h.buildRefreshCookie(result.RefreshToken)

	return authCallbackResponse{
		accessCookie:  accessCookie,
		frontendURL:   h.frontendURL,
		refreshCookie: refreshCookie,
	}, nil
}

func (h *Handler) AuthLogout(ctx context.Context, _ api.AuthLogoutRequestObject) (api.AuthLogoutResponseObject, error) {
	refreshTokenValue := middleware.GetRefreshToken(ctx)
	if refreshTokenValue != "" {
		if err := h.refreshToken.RevokeToken(ctx, refreshTokenValue); err != nil {
			h.logger.Warn(ctx, "failed to revoke refresh token", "error", err)
		}
	}

	return logoutResponse{
		accessCookie:  h.buildExpiredCookie(AccessCookieName),
		refreshCookie: h.buildExpiredCookie(RefreshCookieName),
		success:       true,
	}, nil
}

func (h *Handler) AuthRefresh(ctx context.Context, _ api.AuthRefreshRequestObject) (api.AuthRefreshResponseObject, error) {
	refreshTokenValue := middleware.GetRefreshToken(ctx)
	if refreshTokenValue == "" {
		return api.AuthRefresh401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("refresh token required"),
		}, nil
	}

	output, err := h.refreshToken.Execute(ctx, usecase.RefreshTokenInput{
		RefreshToken: refreshTokenValue,
	})
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) ||
			errors.Is(err, domain.ErrRefreshTokenExpired) ||
			errors.Is(err, domain.ErrTokenReuseDetected) ||
			errors.Is(err, domain.ErrUserNotFound) {
			return api.AuthRefresh401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("invalid or expired refresh token"),
			}, nil
		}

		h.logger.Error(ctx, "failed to refresh token", "error", err)
		return api.AuthRefresh500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("token refresh failed"),
		}, nil
	}

	return refreshResponse{
		accessCookie:  h.buildAccessCookie(output.AccessToken),
		refreshCookie: h.buildRefreshCookie(output.RefreshToken),
		success:       true,
	}, nil
}

func (h *Handler) AuthMe(ctx context.Context, _ api.AuthMeRequestObject) (api.AuthMeResponseObject, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return api.AuthMe401ApplicationProblemPlusJSONResponse{
			UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("authentication required"),
		}, nil
	}

	output, err := h.getCurrentUser.Execute(ctx, usecase.GetCurrentUserInput{UserID: userID})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return api.AuthMe401ApplicationProblemPlusJSONResponse{
				UnauthorizedApplicationProblemPlusJSONResponse: api.NewUnauthorized("user not found"),
			}, nil
		}

		h.logger.Error(ctx, "failed to get current user", "error", err, "userID", userID)
		return api.AuthMe500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("failed to retrieve user information"),
		}, nil
	}

	return api.AuthMe200JSONResponse(mapper.ToUserInfo(output.User)), nil
}

func (h *Handler) AuthDevLogin(ctx context.Context, request api.AuthDevLoginRequestObject) (api.AuthDevLoginResponseObject, error) {
	if h.devLogin == nil {
		return api.AuthDevLogin403ApplicationProblemPlusJSONResponse{
			ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("dev login is disabled"),
		}, nil
	}

	var userID, username string
	if request.Body != nil {
		if request.Body.UserID != nil {
			userID = *request.Body.UserID
		}
		if request.Body.Username != nil {
			username = *request.Body.Username
		}
	}

	output, err := h.devLogin.Execute(ctx, usecase.DevLoginInput{
		UserID:   userID,
		Username: username,
	})
	if err != nil {
		if errors.Is(err, domain.ErrDevLoginDisabled) {
			return api.AuthDevLogin403ApplicationProblemPlusJSONResponse{
				ForbiddenApplicationProblemPlusJSONResponse: api.NewForbidden("dev login is disabled in production"),
			}, nil
		}

		h.logger.Error(ctx, "dev login failed", "error", err)
		return api.AuthDevLogin500ApplicationProblemPlusJSONResponse{
			InternalErrorApplicationProblemPlusJSONResponse: api.NewInternalError("dev login failed"),
		}, nil
	}

	return devLoginResponse{
		accessCookie:  h.buildAccessCookie(output.AccessToken),
		cookieSecure:  h.cookieSecure,
		refreshCookie: h.buildRefreshCookie(output.RefreshToken),
		user:          output.User,
	}, nil
}

func (h *Handler) buildAccessCookie(token string) string {
	cookie := &http.Cookie{
		Domain:   h.cookieDomain,
		HttpOnly: true,
		MaxAge:   AccessCookieMaxAge,
		Name:     AccessCookieName,
		Path:     CookiePath,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cookieSecure,
		Value:    token,
	}
	return cookie.String()
}

func (h *Handler) buildRefreshCookie(token string) string {
	cookie := &http.Cookie{
		Domain:   h.cookieDomain,
		HttpOnly: true,
		MaxAge:   RefreshCookieMaxAge,
		Name:     RefreshCookieName,
		Path:     CookiePath,
		SameSite: http.SameSiteStrictMode,
		Secure:   h.cookieSecure,
		Value:    token,
	}
	return cookie.String()
}

func (h *Handler) buildExpiredCookie(name string) string {
	sameSite := http.SameSiteLaxMode
	if name == RefreshCookieName {
		sameSite = http.SameSiteStrictMode
	}

	cookie := &http.Cookie{
		Domain:   h.cookieDomain,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		MaxAge:   -1,
		Name:     name,
		Path:     CookiePath,
		SameSite: sameSite,
		Secure:   h.cookieSecure,
		Value:    "",
	}
	return cookie.String()
}

type authCallbackResponse struct {
	accessCookie  string
	frontendURL   string
	refreshCookie string
}

func (r authCallbackResponse) VisitAuthCallbackResponse(w http.ResponseWriter) error {
	w.Header().Add("Set-Cookie", r.accessCookie)
	w.Header().Add("Set-Cookie", r.refreshCookie)
	w.Header().Set("Location", r.frontendURL)
	w.WriteHeader(http.StatusFound)
	return nil
}

type logoutResponse struct {
	accessCookie  string
	refreshCookie string
	success       bool
}

func (r logoutResponse) VisitAuthLogoutResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", r.accessCookie)
	w.Header().Add("Set-Cookie", r.refreshCookie)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(api.LogoutResponse{Success: r.success})
}

type refreshResponse struct {
	accessCookie  string
	refreshCookie string
	success       bool
}

func (r refreshResponse) VisitAuthRefreshResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", r.accessCookie)
	w.Header().Add("Set-Cookie", r.refreshCookie)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(api.RefreshResponse{Success: r.success})
}

type devLoginResponse struct {
	accessCookie  string
	cookieSecure  bool
	refreshCookie string
	user          *entity.User
}

func (r devLoginResponse) VisitAuthDevLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", r.accessCookie)
	w.Header().Add("Set-Cookie", r.refreshCookie)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(api.DevLoginResponse{
		Success: true,
		User:    mapper.ToUserInfo(r.user),
	})
}
