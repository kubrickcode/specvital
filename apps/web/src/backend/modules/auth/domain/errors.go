package domain

import "github.com/cockroachdb/errors"

var (
	ErrAccessDenied         = errors.New("access denied by user")
	ErrDevLoginDisabled     = errors.New("dev login is disabled")
	ErrInvalidCode          = errors.New("invalid authorization code")
	ErrInvalidGitHubToken   = errors.New("invalid or expired github access token")
	ErrInvalidOAuthCode     = errors.New("invalid oauth code")
	ErrInvalidState         = errors.New("invalid oauth state")
	ErrInvalidToken         = errors.New("invalid token")
	ErrNetworkFailure       = errors.New("network communication failed")
	ErrNoGitHubToken        = errors.New("user has no github access token")
	ErrRateLimited          = errors.New("github api rate limit exceeded")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired         = errors.New("token expired")
	ErrTokenReuseDetected   = errors.New("refresh token reuse detected")
	ErrUserNotFound         = errors.New("user not found")
)
