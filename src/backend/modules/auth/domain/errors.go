package domain

import "github.com/cockroachdb/errors"

var (
	ErrCodebaseNotFound = errors.New("codebase not found")
	ErrInvalidOAuthCode = errors.New("invalid oauth code")
	ErrInvalidState     = errors.New("invalid oauth state")
	ErrInvalidToken     = errors.New("invalid token")
	ErrNoGitHubToken    = errors.New("user has no github access token")
	ErrTokenExpired     = errors.New("token expired")
	ErrUserNotFound     = errors.New("user not found")
)
