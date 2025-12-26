package domain

import "errors"

var (
	ErrInstallationNotFound  = errors.New("github app installation not found")
	ErrInstallationSuspended = errors.New("github app installation is suspended")
	ErrInvalidPrivateKey     = errors.New("invalid github app private key")
	ErrInvalidWebhookPayload = errors.New("invalid webhook payload")
	ErrMissingAppID          = errors.New("github app id is required")
	ErrMissingSignature      = errors.New("missing webhook signature")
	ErrSignatureVerifyFailed = errors.New("webhook signature verification failed")
	ErrTokenGenerationFailed = errors.New("failed to generate installation token")
	ErrWeakWebhookSecret     = errors.New("webhook secret must be at least 20 characters")
)
