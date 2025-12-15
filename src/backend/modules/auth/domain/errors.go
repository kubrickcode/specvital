package domain

import "github.com/cockroachdb/errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidOAuthCode = errors.New("invalid oauth code")
	ErrInvalidState     = errors.New("invalid oauth state")
)
